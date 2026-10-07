// MedCredix — API demo tanpa backend (untuk deploy statis di Vercel).
//
// Menggantikan backend Golang: seluruh endpoint /api/* dilayani di browser dari
// snapshot data demo `~/data/demo-db.json` (hasil seeder + scoring engine Go asli).
// - Kontrak respons identik dengan backend: data langsung dikembalikan, error dilempar sebagai Error.
// - Filter, urutan, RBAC 10 peran, isolasi tenant, dan jejak audit meniru handler Go.
// - Perubahan data (PATCH/POST) tersimpan di memori selama tab terbuka; muat ulang = reset ke data awal.
// - Skor (Health/BDCS/Trust/CCS dst.) adalah hasil hitung statis dari snapshot — tidak dihitung ulang.

type Row = Record<string, any>

class ApiErr extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

// ---------- RBAC (port dari backend/internal/auth) ----------
const PERMISSIONS = [
  'VIEW', 'CREATE', 'EDIT', 'DELETE', 'VERIFY', 'APPROVE', 'SIGN',
  'EXPORT', 'PRINT', 'ASSIGN', 'RESOLVE', 'REOPEN', 'CONFIGURE',
]
const ROLE_LABELS: Record<string, string> = {
  super_admin: 'SPKD Super Admin',
  admin_rs: 'Admin RS',
  hospital_mgmt: 'Manajemen RS',
  credential_officer: 'Petugas Kredensial',
  medical_committee: 'Komite Medis',
  hr_sdm: 'HR / SDM',
  auditor: 'Auditor',
  claims_team: 'Tim Klaim',
  it_simrs: 'IT / SIMRS',
  viewer: 'Viewer',
}
const ROLE_PERMISSIONS: Record<string, string[]> = {
  super_admin: PERMISSIONS,
  admin_rs: PERMISSIONS,
  hospital_mgmt: ['VIEW', 'EXPORT', 'PRINT', 'ASSIGN'],
  credential_officer: ['VIEW', 'CREATE', 'EDIT', 'VERIFY', 'ASSIGN', 'EXPORT', 'PRINT'],
  medical_committee: ['VIEW', 'VERIFY', 'APPROVE', 'SIGN', 'REOPEN', 'RESOLVE', 'EXPORT'],
  hr_sdm: ['VIEW', 'CREATE', 'EDIT', 'EXPORT'],
  auditor: ['VIEW', 'EXPORT', 'PRINT'],
  claims_team: ['VIEW', 'CREATE', 'EDIT', 'ASSIGN', 'RESOLVE', 'EXPORT'],
  it_simrs: ['VIEW', 'CREATE', 'EDIT', 'CONFIGURE', 'EXPORT'],
  viewer: ['VIEW'],
}

// ---------- Database in-memory ----------
let DB: any = null

async function db() {
  if (!DB) {
    const mod: any = await import('~/data/demo-db.json')
    DB = structuredClone(mod.default ?? mod)
  }
  return DB
}

const T = (name: string): Row[] => DB.tables[name] ?? (DB.tables[name] = [])

// ---------- Util ----------
function nowIso(): string {
  const d = new Date(Date.now() + 7 * 3600_000)
  return d.toISOString().slice(0, 19) + '+07:00'
}

function newId(): string {
  return 'c' + Date.now().toString(16) + Math.random().toString(16).slice(2, 8)
}

function parseDate(s?: string | null): string | null {
  if (!s) return null
  if (/^\d{4}-\d{2}-\d{2}$/.test(s)) return `${s}T00:00:00Z`
  const t = new Date(s)
  return isNaN(t.getTime()) ? null : t.toISOString()
}

const pad3 = (n: number) => String(n).padStart(3, '0')

/** LIKE '%term%' ala SQLite (case-insensitive ASCII). */
function like(v: any, term: string): boolean {
  if (v === null || v === undefined) return false
  return String(v).toLowerCase().includes(term.toLowerCase())
}

/** Pembanding urutan: null lebih kecil (ala SQLite), string biner. */
function cmp(a: any, b: any): number {
  if (a === b) return 0
  if (a === null || a === undefined) return -1
  if (b === null || b === undefined) return 1
  return a < b ? -1 : 1
}

function sortBy<T>(rows: T[], ...keys: [keyof T | string, 'asc' | 'desc'][]): T[] {
  return rows
    .map((r, i) => [r, i] as [T, number])
    .sort((x, y) => {
      for (const [k, dir] of keys) {
        const c = cmp((x[0] as any)[k], (y[0] as any)[k])
        if (c !== 0) return dir === 'asc' ? c : -c
      }
      return x[1] - y[1]
    })
    .map((p) => p[0])
}

const inH = (name: string, hid: string) => T(name).filter((r) => r.hospitalId === hid)

const trimFloat = (v: number) => (Number.isInteger(v) ? String(v) : v.toFixed(1))
const round1 = (v: number) => Math.round(v * 10) / 10
const roundHalf = (v: number) => Math.floor(v * 10 + 0.5) / 10

// ---------- Sesi ----------
function sessionOf(u: Row) {
  return {
    id: u.id,
    name: u.name,
    email: u.email,
    role: u.role,
    roleLabel: ROLE_LABELS[u.role] ?? u.role,
    title: u.title ?? null,
    hospitalId: u.hospitalId ?? null,
    permissions: ROLE_PERMISSIONS[u.role] ?? ['VIEW'],
  }
}
type Session = ReturnType<typeof sessionOf>

function getUser(userId: string): Session {
  const u = userId ? T('User').find((x) => x.id === userId) : null
  if (!u || !u.active) throw new ApiErr(401, 'Sesi tidak valid. Silakan masuk kembali.')
  return sessionOf(u)
}

function resolveHospital(user: Session, requested?: string | null): string {
  if (user.role === 'super_admin') return requested || 'h_nusantara'
  if (!user.hospitalId) throw new ApiErr(403, 'Akun Anda tidak terikat ke rumah sakit mana pun.')
  if (requested && requested !== user.hospitalId) {
    throw new ApiErr(403, 'Akses lintas rumah sakit tidak diizinkan (isolasi tenant).')
  }
  return user.hospitalId
}

function requirePerm(user: Session, perm: string) {
  if (!user.permissions.includes(perm)) {
    throw new ApiErr(403, `Peran ${user.roleLabel} tidak memiliki izin ${perm}.`)
  }
}

function logAudit(o: {
  hospitalId: string | null
  actor: { name: string; role: string } | null
  action: string
  entityType: string
  entityId?: string | null
  description: string
  before?: string | null
  after?: string | null
}) {
  T('AuditLog').push({
    id: newId(),
    hospitalId: o.hospitalId,
    actorName: o.actor?.name ?? 'MedCredix Engine',
    actorRole: o.actor?.role ?? 'system',
    action: o.action,
    entityType: o.entityType,
    entityId: o.entityId ?? null,
    description: o.description,
    before: o.before ?? null,
    after: o.after ?? null,
    ip: '10.20.0.1',
    correlationId: `corr-${Date.now().toString(16)}-${Math.floor(Math.random() * 0xffff).toString(16).padStart(4, '0')}`,
    createdAt: nowIso(),
  })
}

const userBriefs = (hid: string) =>
  T('User').filter((u) => u.hospitalId === hid).map((u) => ({ id: u.id, name: u.name, role: u.role }))

// ---------- Handler ----------
type Ctx = { user: Session; q: URLSearchParams; body: any }
type Handler = (c: Ctx) => any

const OPEN_ACTION = ['Baru', 'Ditugaskan', 'Dalam Pengerjaan', 'Menunggu Verifikasi']

const handlers: Record<string, Handler> = {}

// ----- bootstrap & dashboard -----
handlers['GET /bootstrap'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const hospital = T('Hospital').find((h) => h.id === hid)
  if (!hospital) throw new ApiErr(404, 'Data tidak ditemukan.')
  const intel = DB.intel[hid]
  const notifications = sortBy(inH('Notification', hid), ['createdAt', 'desc']).slice(0, 8)
  const badges = {
    warnings: inH('EarlyWarning', hid).filter((w) => w.status === 'Baru').length,
    risks: inH('Risk', hid).filter((r) => r.status === 'Terbuka').length,
    actions: inH('Action', hid).filter((a) => OPEN_ACTION.includes(a.status)).length,
  }
  return {
    user,
    hospital,
    hospitals:
      user.role === 'super_admin'
        ? T('Hospital').map((h) => ({ id: h.id, name: h.name, city: h.city }))
        : null,
    intel: {
      healthScore: intel.healthScore,
      healthStatus: intel.healthStatus,
      bcds: intel.bcds,
      kelompok: intel.kelompok,
      trust: intel.trust,
      ccs: intel.ccs.score,
      cpcs: intel.cpcs.score,
      ur: intel.ur,
      domains: intel.domains.map((d: Row) => ({ key: d.key, label: d.label, score: d.score, weight: d.weight, group: d.group })),
      trend: intel.trend,
      prereq: intel.prereq,
    },
    badges,
    notifications,
    alertRules: inH('AlertRule', hid),
    users: userBriefs(hid),
  }
}

handlers['GET /dashboard'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const intel = DB.intel[hid]
  const risks = sortBy(inH('Risk', hid).filter((r) => r.status !== 'Ditutup'), ['level', 'desc'], ['detectedAt', 'desc'])
  const warnings = sortBy(
    inH('EarlyWarning', hid).filter((w) => ['Baru', 'Diketahui', 'Ditugaskan'].includes(w.status)),
    ['severity', 'desc'], ['detectedAt', 'desc'],
  )
  const actions = sortBy(inH('Action', hid).filter((a) => OPEN_ACTION.includes(a.status)), ['priority', 'desc'], ['dueDate', 'asc'])
  const evidence = sortBy(inH('EvidenceItem', hid), ['timestamp', 'desc']).slice(0, 6)
  const resolved = sortBy(
    inH('Action', hid).filter((a) => ['Selesai', 'Ditutup'].includes(a.status) && a.resolvedAt),
    ['resolvedAt', 'desc'],
  ).slice(0, 6)
  const pathwayWorst = sortBy(inH('ClinicalPathway', hid), ['avgAdherence', 'asc']).slice(0, 3)

  const timeline: Row[] = []
  for (const t of intel.trend) {
    timeline.push({
      period: t.period,
      type: 'Skor',
      label: `${trimFloat(t.health)} · BDCS ${trimFloat(t.bcds)} · Trust ${trimFloat(t.trust)}`,
      delta: null,
    })
  }
  for (const a of resolved) {
    timeline.push({ period: String(a.resolvedAt).slice(0, 7), type: 'Perbaikan', label: `${a.code} — ${a.title}`, delta: null })
  }
  const tl = sortBy(timeline, ['period', 'asc'])
  const count = (lv: string) => risks.filter((r) => r.level === lv).length
  return {
    intel, risks, warnings, actions, topActions: actions.slice(0, 5), evidence, pathwayWorst,
    timeline: tl,
    riskSummary: { tinggi: count('Tinggi'), sedang: count('Sedang'), rendah: count('Rendah') },
  }
}

// ----- scores (WHY panel) -----
const DOMAIN_LABELS: Record<string, string> = {
  workforce: 'Clinical Workforce', credentialing: 'Credentialing Readiness',
  digital: 'Digital Maturity', medication: 'Medication Management',
  diagnostic: 'Diagnostic Readiness', interop: 'Interoperability',
  claim: 'Claim Integrity', pathway: 'Clinical Pathway',
  safety: 'Patient Safety', trust: 'Data Trust',
}

handlers['GET /scores'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const domainKey = q.get('domain') || 'health'
  const intel = DB.intel[hid]

  if (domainKey === 'health') {
    return {
      kind: 'health',
      formula:
        'Hospital Health Score = Σ(skor domain × bobot domain). Bobot: Claim Integrity 16%, Credentialing 14%, Patient Safety 14%, Clinical Workforce 12%, Clinical Pathway 10%, Digital Maturity 10%, Medication 8%, Diagnostik 8%, Interop 8%.',
      formulaVersion: DB.formulaVersion,
      domains: intel.domains.map((d: Row) => ({
        key: d.key, label: d.label, score: d.score, weight: d.weight, contribution: round1(d.score * d.weight),
      })),
      healthScore: intel.healthScore, bcds: intel.bcds, trust: intel.trust, kelompok: intel.kelompok,
      stageTable: DB.stageTable,
      risks: inH('Risk', hid).filter((r) => r.status !== 'Ditutup').slice(0, 6),
      evidence: inH('EvidenceItem', hid).filter((e) => e.refType === 'Skor').slice(0, 10),
      actions: inH('Action', hid).filter((a) => ['Baru', 'Ditugaskan', 'Dalam Pengerjaan'].includes(a.status)).slice(0, 6),
      trend: intel.trend,
      recommendations: DB.healthRecs[hid] ?? [],
    }
  }

  if (domainKey === 'bcds') {
    return {
      kind: 'bcds',
      formula:
        'BPJS Credential Readiness (BDCS) = (Kelompok A × 25%) + (B × 30%) + (C × 35%) + (D × 10%), lalu dikalikan Trust Multiplier — ditampilkan transparan.',
      kelompok: intel.kelompok, bcds: intel.bcds, trust: intel.trust,
      trustBands: DB.trustBands,
      stage: intel.bcds.stage, lulus: intel.bcds.lulus, prereqOk: intel.bcds.prereqOk,
      trend: intel.trend, formulaVersion: DB.formulaVersion,
    }
  }

  const domain = intel.domains.find((d: Row) => d.key === domainKey)
  if (!domain) throw new ApiErr(404, 'Domain tidak dikenal.')
  const risks = inH('Risk', hid).filter((r) => r.status !== 'Ditutup')
  const evidence = inH('EvidenceItem', hid).filter((e) => e.refId === `DOMAIN:${domainKey}`)
  const dl = DOMAIN_LABELS[domainKey]
  const relatedRisks = risks.filter(
    (r) => r.domain === dl || String(r.domain).toLowerCase().includes(domainKey.toLowerCase()),
  )
  const relatedActions = inH('Action', hid).filter(
    (a) => a.sourceRef && (a.sourceRef === `DOMAIN:${domainKey}` || relatedRisks.some((r) => r.code === a.sourceRef)),
  )
  return {
    kind: 'domain', domain, formulaVersion: DB.formulaVersion,
    evidence,
    risks: relatedRisks.length ? relatedRisks : null,
    actions: relatedActions.length ? relatedActions : null,
    recommendations: DB.domainRecs[hid]?.[domainKey] ?? [],
    trend: intel.trend, healthScore: intel.healthScore,
  }
}

// ----- risks & warnings -----
handlers['GET /risks'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  let rows = inH('Risk', hid)
  if (q.get('level')) rows = rows.filter((r) => r.level === q.get('level'))
  if (q.get('category')) rows = rows.filter((r) => r.category === q.get('category'))
  if (q.get('status')) rows = rows.filter((r) => r.status === q.get('status'))
  rows = sortBy(rows, ['detectedAt', 'desc'])
  return { risks: rows, evidence: inH('EvidenceItem', hid).filter((e) => e.refType === 'Risiko'), total: rows.length }
}

handlers['PATCH /risks'] = ({ user, q, body }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const risk = inH('Risk', hid).find((r) => r.id === body?.id)
  if (!risk) throw new ApiErr(404, 'Risiko tidak ditemukan.')
  if (body.action === 'ack' || body.action === 'assign') requirePerm(user, 'ASSIGN')
  else requirePerm(user, 'EDIT')
  const before = risk.status
  const upd: Row = {}
  if (body.status) upd.status = body.status
  else if (body.action === 'ack') upd.status = 'Ditinjau'
  if (body.owner) upd.suggestedOwner = body.owner
  if (body.deadline) {
    const d = parseDate(body.deadline)
    if (d) upd.deadline = d
  }
  if (!Object.keys(upd).length) return risk
  Object.assign(risk, upd)
  logAudit({
    hospitalId: hid, actor: user, action: 'update', entityType: 'Risiko', entityId: risk.id,
    description: `${risk.code}: ${body.action}${body.owner ? ' → ' + body.owner : ''}`,
    before, after: risk.status,
  })
  return risk
}

handlers['GET /warnings'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  let rows = inH('EarlyWarning', hid)
  if (q.get('severity')) rows = rows.filter((w) => w.severity === q.get('severity'))
  if (q.get('status')) rows = rows.filter((w) => w.status === q.get('status'))
  rows = sortBy(rows, ['status', 'asc'], ['detectedAt', 'desc'])
  return { warnings: rows, rules: sortBy(inH('AlertRule', hid), ['name', 'asc']), total: rows.length }
}

handlers['PATCH /warnings'] = ({ user, q, body }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const w = inH('EarlyWarning', hid).find((x) => x.id === body?.id)
  if (!w) throw new ApiErr(404, 'Peringatan tidak ditemukan.')
  if (body.action === 'resolve') requirePerm(user, 'RESOLVE')
  else requirePerm(user, 'ASSIGN')
  const before = w.status
  const statusMap: Record<string, string> = { ack: 'Diketahui', assign: 'Ditugaskan', resolve: 'Teratasi' }
  w.status = statusMap[body.action] ?? w.status
  if (body.action === 'ack' || body.action === 'assign') {
    w.acknowledgedBy = user.name
    w.acknowledgedAt = nowIso()
  }
  if (body.action === 'assign') {
    const assignee = body.assignee || w.acknowledgedBy || w.assignedTo || ''
    if (assignee) w.assignedTo = assignee
  }
  if (body.dueDate) {
    const d = parseDate(body.dueDate)
    if (d) w.dueDate = d
  }
  logAudit({
    hospitalId: hid, actor: user, action: 'update', entityType: 'Peringatan', entityId: w.id,
    description: `${w.code}: ${body.action}`, before, after: w.status,
  })
  return w
}

// ----- actions -----
handlers['GET /actions'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const id = q.get('id')
  if (id) {
    const a = inH('Action', hid).find((x) => x.id === id)
    if (!a) throw new ApiErr(404, 'Aksi tidak ditemukan.')
    const events = sortBy(T('ActionEvent').filter((e) => e.actionId === id), ['createdAt', 'asc'])
    return { action: a, events, users: userBriefs(hid) }
  }
  let rows = inH('Action', hid)
  if (q.get('status')) rows = rows.filter((a) => a.status === q.get('status'))
  if (q.get('priority')) rows = rows.filter((a) => a.priority === q.get('priority'))
  const s = q.get('q')
  if (s) rows = rows.filter((a) => like(a.title, s) || like(a.code, s) || like(a.owner, s))
  rows = sortBy(rows, ['createdAt', 'desc'])
  return {
    actions: rows,
    stats: {
      total: rows.length,
      open: rows.filter((a) => OPEN_ACTION.includes(a.status)).length,
      high: rows.filter((a) => a.priority === 'Tinggi' && a.status !== 'Selesai' && a.status !== 'Ditutup').length,
      resolved: rows.filter((a) => a.status === 'Selesai' || a.status === 'Ditutup').length,
    },
  }
}

handlers['POST /actions'] = ({ user, q, body }) => {
  if (!body) throw new ApiErr(400, 'Payload tidak valid.')
  const hid = resolveHospital(user, q.get('hospitalId'))
  requirePerm(user, 'CREATE')
  const code = 'ACT-' + pad3(inH('Action', hid).length + 1)
  const sourceType = body.sourceType || 'Manual'
  const now = nowIso()
  const created: Row = {
    id: newId(), hospitalId: hid, code, title: body.title ?? '', description: body.description ?? '',
    priority: body.priority || 'Sedang', status: 'Baru', sourceType, sourceRef: body.sourceRef ?? null,
    owner: body.owner || 'Belum ditugaskan', ownerRole: body.ownerRole ?? '',
    dueDate: parseDate(body.dueDate), expectedOutcome: body.expectedOutcome ?? null,
    progressNote: null, resolutionNote: null, resolvedAt: null, createdAt: now, updatedAt: now,
  }
  T('Action').push(created)
  T('ActionEvent').push({
    id: newId(), actionId: created.id, type: 'created', fromValue: null, toValue: 'Baru', byName: user.name,
    note: `Aksi dibuat dari ${sourceType}${body.sourceRef ? ` (${body.sourceRef})` : ''}`, createdAt: now,
  })
  logAudit({
    hospitalId: hid, actor: user, action: 'create', entityType: 'Aksi', entityId: created.id,
    description: `Membuat aksi ${code}: ${body.title}`,
  })
  return created
}

handlers['PATCH /actions'] = ({ user, q, body }) => {
  if (!body) throw new ApiErr(400, 'Payload tidak valid.')
  const hid = resolveHospital(user, q.get('hospitalId'))
  const a = inH('Action', hid).find((x) => x.id === body.id)
  if (!a) throw new ApiErr(404, 'Aksi tidak ditemukan.')
  const now = nowIso()

  if (body.action === 'assign') {
    requirePerm(user, 'ASSIGN')
    const prevOwner = a.owner
    const owner = (body.owner ?? '').trim() ? body.owner : a.owner
    const ownerRole = (body.ownerRole ?? '').trim() ? body.ownerRole : a.ownerRole
    if (a.status === 'Baru') a.status = 'Ditugaskan'
    a.owner = owner
    a.ownerRole = ownerRole
    if (body.dueDate) {
      const d = parseDate(body.dueDate)
      if (d) a.dueDate = d
    }
    a.updatedAt = now
    T('ActionEvent').push({
      id: newId(), actionId: a.id, type: 'assigned', fromValue: prevOwner, toValue: owner,
      byName: user.name, note: body.note ?? null, createdAt: now,
    })
    logAudit({
      hospitalId: hid, actor: user, action: 'update', entityType: 'Aksi', entityId: a.id,
      description: `${a.code} ditugaskan ke ${owner}`, before: prevOwner, after: owner,
    })
    return a
  }

  if (body.action === 'status') {
    requirePerm(user, 'EDIT')
    const prev = a.status
    a.status = body.status
    if (body.note !== undefined && body.note !== null) a.progressNote = body.note
    if (body.status === 'Selesai' && body.resolutionNote !== undefined && body.resolutionNote !== null) {
      a.resolutionNote = body.resolutionNote
    }
    if (body.status === 'Selesai' || body.status === 'Ditutup') a.resolvedAt = now
    a.updatedAt = now
    T('ActionEvent').push({
      id: newId(), actionId: a.id, type: 'status', fromValue: prev, toValue: body.status,
      byName: user.name, note: body.note ?? null, createdAt: now,
    })
    logAudit({
      hospitalId: hid, actor: user, action: 'update', entityType: 'Aksi', entityId: a.id,
      description: `${a.code}: ${prev} → ${body.status}`, before: prev, after: body.status,
    })
    return a
  }

  if (body.action === 'comment') {
    const ev = {
      id: newId(), actionId: a.id, type: 'comment', fromValue: null, toValue: null,
      byName: user.name, note: body.note ?? '', createdAt: now,
    }
    T('ActionEvent').push(ev)
    return ev
  }
  throw new ApiErr(400, 'Aksi tidak dikenal.')
}

// ----- praktisi, kredensial, prasyarat -----
handlers['GET /practitioners'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const id = q.get('id')
  if (id) {
    const p = inH('Practitioner', hid).find((x) => x.id === id)
    if (!p) throw new ApiErr(404, 'Praktisi tidak ditemukan.')
    const credentials = sortBy(T('PractitionerCredential').filter((c) => c.practitionerId === id), ['type', 'asc'])
    const ids = new Set(credentials.map((c) => c.id))
    const documents = inH('CredentialDocument', hid).filter((d) => ids.has(d.credentialId))
    const standard = p.specialty ? inH('SpecialtyStandard', hid).find((s) => s.specialty === p.specialty) ?? null : null
    return { practitioner: p, credentials, documents, standard }
  }
  let rows = inH('Practitioner', hid)
  const s = q.get('q')
  if (s) rows = rows.filter((p) => like(p.fullName, s) || like(p.specialty, s) || like(p.sipNumber, s) || like(p.strNumber, s))
  if (q.get('prof')) rows = rows.filter((p) => p.profession === q.get('prof'))
  if (q.get('status')) rows = rows.filter((p) => p.status === q.get('status'))
  if (q.get('dept')) rows = rows.filter((p) => like(p.department, q.get('dept')!))
  rows = sortBy(rows, ['fullName', 'asc'])
  const depts: string[] = []
  for (const p of rows) if (!depts.includes(p.department)) depts.push(p.department)
  return { practitioners: rows, depts, total: rows.length }
}

handlers['GET /credentials'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  let rows = inH('PractitionerCredential', hid)
  if (q.get('status')) rows = rows.filter((c) => c.status === q.get('status'))
  if (q.get('type')) rows = rows.filter((c) => c.type === q.get('type'))
  const s = q.get('q')
  if (s) rows = rows.filter((c) => like(c.name, s) || like(c.number, s))
  rows = sortBy(rows, ['expiryDate', 'asc'])
  const pmap = new Map(inH('Practitioner', hid).map((p) => [p.id, p]))
  const credentials = rows.map((c) => {
    const p = pmap.get(c.practitionerId)
    return { ...c, practitionerName: p?.fullName ?? '—', profession: p?.profession ?? '—', specialty: p?.specialty ?? null }
  })
  return { credentials, total: credentials.length }
}

handlers['PATCH /credentials'] = ({ user, q, body }) => {
  if (!body) throw new ApiErr(400, 'Payload tidak valid.')
  const hid = resolveHospital(user, q.get('hospitalId'))
  const cred = inH('PractitionerCredential', hid).find((c) => c.id === body.id)
  if (!cred) throw new ApiErr(404, 'Kredensial tidak ditemukan.')
  const transitions: Record<string, [string, string]> = {
    submit: ['Diajukan', 'EDIT'],
    verify: ['Verifikasi', 'VERIFY'],
    review: ['Tinjauan Komite', 'VERIFY'],
    approve: ['Disetujui', 'APPROVE'],
    activate: ['Aktif', 'APPROVE'],
    suspend: ['Ditangguhkan', 'APPROVE'],
    reactivate: ['Aktif', 'APPROVE'],
  }
  const t = transitions[body.action]
  if (!t) throw new ApiErr(400, 'Aksi tidak dikenal.')
  requirePerm(user, t[1])
  const prev = cred.status
  const bump =
    ['Disetujui', 'Aktif', 'Segera Berakhir'].includes(cred.status) && body.action !== 'reactivate'
  cred.status = t[0]
  if (bump) cred.version = cred.version + 1
  const now = nowIso()
  if (body.action === 'verify' || body.action === 'review') {
    cred.verifiedById = user.id
    cred.verifiedAt = now
  }
  if (['approve', 'activate', 'reactivate', 'suspend'].includes(body.action)) {
    cred.approvedById = user.id
    cred.approvedAt = now
  }
  cred.updatedAt = now
  const auditAction =
    body.action === 'approve' || body.action === 'activate' ? 'approve'
      : body.action === 'verify' || body.action === 'review' ? 'verify' : 'update'
  logAudit({
    hospitalId: hid, actor: user, action: auditAction, entityType: 'Kredensial', entityId: cred.id,
    description: `${body.action} → ${t[0]} (${cred.name})`, before: prev, after: t[0],
  })
  return cred
}

handlers['GET /requirements'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  return { requirements: sortBy(inH('Requirement', hid), ['code', 'asc']) }
}

handlers['PATCH /requirements'] = ({ user, q, body }) => {
  if (!body) throw new ApiErr(400, 'Payload tidak valid.')
  const hid = resolveHospital(user, q.get('hospitalId'))
  requirePerm(user, body.action === 'verify' ? 'VERIFY' : 'EDIT')
  const r = inH('Requirement', hid).find((x) => x.id === body.id)
  if (!r) throw new ApiErr(404, 'Prasyarat tidak ditemukan.')
  const prev = r.status
  let status = r.status
  if (body.action === 'submit') status = 'Menunggu Verifikasi'
  else if (body.action === 'verify') status = 'Valid'
  else if (body.action === 'invalidate') status = 'Belum Dilengkapi'
  r.status = status
  r.updatedAt = nowIso()
  if (body.docNumber !== undefined && body.docNumber !== null) r.docNumber = body.docNumber
  if (body.issuedAt) { const d = parseDate(body.issuedAt); if (d) r.issuedAt = d }
  if (body.docExpiry) { const d = parseDate(body.docExpiry); if (d) r.docExpiry = d }
  if (body.fileName !== undefined && body.fileName !== null) r.fileName = body.fileName
  if (body.notes !== undefined && body.notes !== null) r.notes = body.notes
  logAudit({
    hospitalId: hid, actor: user, action: body.action === 'verify' ? 'verify' : 'update',
    entityType: 'Prasyarat', entityId: r.id, description: `${r.code} → ${status}`, before: prev, after: status,
  })
  return r
}

// ----- klaim & klinis -----
handlers['GET /claims'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const id = q.get('id')
  if (id) {
    const claim = inH('Claim', hid).find((c) => c.id === id)
    if (!claim) throw new ApiErr(404, 'Klaim tidak ditemukan.')
    return {
      claim,
      anomalies: T('ClaimAnomaly').filter((a) => a.claimId === id),
      evidence: inH('EvidenceItem', hid).filter((e) => e.refId === claim.claimNumber),
    }
  }
  let rows = inH('Claim', hid)
  if (q.get('status')) rows = rows.filter((c) => c.status === q.get('status'))
  if (q.get('careType')) rows = rows.filter((c) => c.careType === q.get('careType'))
  const s = q.get('q')
  if (s) rows = rows.filter((c) => like(c.claimNumber, s) || like(c.sepNumber, s) || like(c.inacbgCode, s))
  rows = sortBy(rows, ['admissionDate', 'desc'])
  const numById = new Map(rows.map((c) => [c.id, c.claimNumber]))
  const anomalies = sortBy(inH('ClaimAnomaly', hid), ['status', 'asc'], ['severity', 'desc']).map((a) => ({
    ...a, claimNumber: numById.get(a.claimId) || '—',
  }))
  return { claims: rows, anomalies, ccs: DB.intel[hid].ccs, total: rows.length }
}

handlers['PATCH /claims'] = ({ user, q, body }) => {
  if (!body) throw new ApiErr(400, 'Payload tidak valid.')
  const hid = resolveHospital(user, q.get('hospitalId'))
  requirePerm(user, 'EDIT')
  if (body.anomalyId) {
    const a = inH('ClaimAnomaly', hid).find((x) => x.id === body.anomalyId)
    if (!a) throw new ApiErr(404, 'Anomali tidak ditemukan.')
    const statusMap: Record<string, string> = { review: 'Ditinjau', confirm: 'Terkonfirmasi Bermasalah', clarify: 'Diklarifikasi' }
    const prev = a.status
    a.status = statusMap[body.action] ?? a.status
    a.reviewedBy = user.name
    const claim = T('Claim').find((c) => c.id === a.claimId)
    logAudit({
      hospitalId: hid, actor: user, action: 'update', entityType: 'Anomali Klaim', entityId: a.id,
      description: `${claim?.claimNumber ?? ''} (${a.category}): ${prev} → ${a.status}`, before: prev, after: a.status,
    })
    return a
  }
  if (body.claimId) {
    const cl = inH('Claim', hid).find((x) => x.id === body.claimId)
    if (!cl) throw new ApiErr(404, 'Klaim tidak ditemukan.')
    const prev = cl.status
    cl.status = body.status
    logAudit({
      hospitalId: hid, actor: user, action: 'update', entityType: 'Klaim', entityId: cl.id,
      description: `${cl.claimNumber}: ${prev} → ${body.status}`, before: prev, after: body.status,
    })
    return cl
  }
  throw new ApiErr(400, 'Payload tidak valid.')
}

handlers['GET /clinical'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const typ = q.get('type') || 'all'
  const intel = DB.intel[hid]
  if (typ === 'pathway') {
    return {
      pathways: sortBy(inH('ClinicalPathway', hid), ['avgAdherence', 'asc']),
      episodes: sortBy(inH('PathwayEpisode', hid), ['adherence', 'asc']),
      cpcs: intel.cpcs,
      coverage: { total: 20, covered: 16 },
    }
  }
  const utils = sortBy(inH('UtilizationStat', hid), ['month', 'asc'])
  const complaints = sortBy(inH('ComplaintStat', hid), ['month', 'asc'])
  if (typ === 'utilization') return { utilization: utils, complaints, ur: intel.ur }
  if (typ === 'safety') {
    return { safety: intel.safety, domain: intel.domains.find((d: Row) => d.key === 'safety') ?? null }
  }
  return {
    pathways: inH('ClinicalPathway', hid), utilization: utils, complaints,
    cpcs: intel.cpcs, ur: intel.ur, safety: intel.safety,
  }
}

// ----- kepatuhan, bukti, audit, pencarian -----
handlers['GET /compliance'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const regs = sortBy(T('Regulation'), ['year', 'desc'])
  const versions = T('RegulationVersion')
  const mappings = sortBy(inH('ComplianceMapping', hid), ['status', 'asc'])
  const regMap = new Map(regs.map((r) => [r.id, r]))
  const count = (s: string) => mappings.filter((m) => m.status === s).length
  return {
    regulations: regs.map((r) => ({ ...r, versions: versions.filter((v) => v.regulationId === r.id) })),
    mappings: mappings.map((m) => {
      const r = regMap.get(m.regulationId)
      return { ...m, regulation: r ? { code: r.code, title: r.title, number: r.number, category: r.category } : null }
    }),
    summary: { patuh: count('Patuh'), sebagian: count('Sebagian'), belum: count('Belum'), ditinjau: count('Sedang Ditinjau') },
    total: mappings.length,
  }
}

handlers['GET /evidence'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  let rows = inH('EvidenceItem', hid)
  if (q.get('refType')) rows = rows.filter((e) => e.refType === q.get('refType'))
  if (q.get('refId')) rows = rows.filter((e) => e.refId === q.get('refId'))
  const s = q.get('q')
  if (s) rows = rows.filter((e) => like(e.title, s) || like(e.metric, s) || like(e.sourceSystem, s))
  rows = sortBy(rows, ['timestamp', 'desc'])
  const byType: Record<string, number> = {}
  for (const e of rows) byType[e.refType] = (byType[e.refType] ?? 0) + 1
  return { items: rows, byType, total: rows.length }
}

handlers['GET /audit'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  let rows = T('AuditLog')
  if (user.role !== 'super_admin') {
    rows = rows.filter((l) => l.hospitalId === hid || (!l.hospitalId && ['super_admin', 'system'].includes(l.actorRole)))
  }
  if (q.get('action')) rows = rows.filter((l) => l.action === q.get('action'))
  if (q.get('entityType')) rows = rows.filter((l) => l.entityType === q.get('entityType'))
  const s = q.get('q')
  if (s) rows = rows.filter((l) => like(l.actorName, s) || like(l.description, s) || like(l.entityId, s))
  const logs = sortBy(rows, ['createdAt', 'desc']).slice(0, 200)
  return { logs, total: logs.length }
}

handlers['GET /search'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const term = q.get('q') ?? ''
  if (term.length < 2) return { results: [], term }
  const results: Row[] = []
  const add = (type: string, icon: string, label: string, sub: string, view: string, id: string) =>
    results.push({ type, icon, label, sub, view, id })

  inH('Practitioner', hid)
    .filter((p) => like(p.fullName, term) || like(p.specialty, term) || like(p.sipNumber, term) || like(p.strNumber, term))
    .slice(0, 5)
    .forEach((p) => add('Praktisi', 'Users', p.fullName, p.profession + (p.specialty ? ` — ${p.specialty}` : ''), 'cred-practitioners', p.id))
  inH('PractitionerCredential', hid)
    .filter((c) => like(c.name, term) || like(c.number, term))
    .slice(0, 5)
    .forEach((c) => add('Kredensial', 'IdCard', c.name, `${c.type} · ${c.status}`, 'cred-credentials', c.id))
  inH('Risk', hid)
    .filter((r) => like(r.title, term) || like(r.code, term))
    .slice(0, 5)
    .forEach((r) => add('Risiko', 'Radar', r.title, `${r.code} · ${r.category} · ${r.level}`, 'risk-radar', r.id))
  inH('Claim', hid)
    .filter((c) => like(c.claimNumber, term) || like(c.sepNumber, term) || like(c.inacbgDesc, term))
    .slice(0, 5)
    .forEach((c) => add('Klaim', 'FileCheck2', c.claimNumber, `${c.inacbgCode} ${c.inacbgDesc} · ${c.status}`, 'claim-anomalies', c.id))
  inH('Action', hid)
    .filter((a) => like(a.title, term) || like(a.code, term) || like(a.owner, term))
    .slice(0, 5)
    .forEach((a) => add('Aksi', 'CheckCircle2', a.title, `${a.code} · ${a.status} · ${a.owner}`, 'actions', a.id))
  T('Regulation')
    .filter((r) => like(r.title, term) || like(r.number, term) || like(r.code, term))
    .slice(0, 4)
    .forEach((r) => add('Regulasi', 'Scale', r.title, r.number, 'compliance', r.id))
  inH('EvidenceItem', hid)
    .filter((e) => like(e.title, term) || like(e.metric, term))
    .slice(0, 4)
    .forEach((e) => add('Bukti', 'Fingerprint', e.title, e.metric != null ? `${e.sourceSystem} · ${e.validationStatus}` : e.sourceSystem, 'evidence', e.id))
  inH('EarlyWarning', hid)
    .filter((w) => like(w.title, term) || like(w.code, term))
    .slice(0, 4)
    .forEach((w) => add('Peringatan', 'AlertTriangle', w.title, `${w.code} · ${w.severity}`, 'early-warnings', w.id))
  return { results, term }
}

// ----- laporan -----
const REPORT_TYPES: Record<string, { title: string; desc: string }> = {
  kredensial: { title: 'Laporan Kredensial RS', desc: 'Skor BDCS, kelompok A–D, prasyarat, kredensial & tenaga medis' },
  'kesehatan-rs': { title: 'Hospital Health Report', desc: '9 domain, tren 6 bulan, rekomendasi perbaikan' },
  risiko: { title: 'Risk Report', desc: 'Risiko aktif, level, owner, dan tenggat' },
  aksi: { title: 'Action & Improvement Report', desc: 'Aksi perbaikan, status, resolusi, dampak skor' },
  kepatuhan: { title: 'Regulatory Compliance Report', desc: 'Matriks kepatuhan, regulasi, versi, temuan' },
  klaim: { title: 'BPJS Claim Integrity Report', desc: 'CCS, sub-skor, anomali aktif, prioritas tinjauan' },
}

const ymd = (s?: string | null) => (s ? String(s).slice(0, 10) : '—')

handlers['GET /reports'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const id = q.get('id')
  if (q.get('export')) throw new ApiErr(400, 'Gunakan unduhan CSV.')
  if (id) {
    const snap = inH('ReportSnapshot', hid).find((r) => r.id === id)
    if (!snap) throw new ApiErr(404, 'Laporan tidak ditemukan.')
    let content: any = null
    try { content = JSON.parse(snap.contentJson) } catch { /* abaikan */ }
    return { report: snap, content }
  }
  return {
    reports: sortBy(inH('ReportSnapshot', hid), ['generatedAt', 'desc']),
    types: Object.entries(REPORT_TYPES).map(([type, m]) => ({ type, title: m.title, desc: m.desc })),
  }
}

handlers['POST /reports'] = ({ user, q, body }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  requirePerm(user, 'CREATE')
  const typ = body?.type || 'kesehatan-rs'
  const meta = REPORT_TYPES[typ]
  if (!meta) throw new ApiErr(400, 'Jenis laporan tidak dikenal.')
  const intel = DB.intel[hid]
  const hospital = T('Hospital').find((h) => h.id === hid) ?? {}
  const period = intel.trend.length ? `${intel.trend[0].period}–${intel.trend[intel.trend.length - 1].period}` : '—'

  const content: Row = {
    ringkasan: {
      namaRS: hospital.name, kelas: hospital.class, kota: hospital.city, periode: period,
      healthScore: intel.healthScore, bcds: intel.bcds.final, stage: intel.bcds.stage,
      trust: intel.trust.score, ccs: intel.ccs.score, cpcs: intel.cpcs.score,
    },
  }
  switch (typ) {
    case 'kredensial': {
      const statusCount: Record<string, number> = {}
      for (const c of inH('PractitionerCredential', hid)) statusCount[c.status] = (statusCount[c.status] ?? 0) + 1
      content.tenagaMedis = { total: inH('Practitioner', hid).length }
      content.kredensial = statusCount
      content.prasyarat = inH('Requirement', hid).map((r) => ({ kode: r.code, nama: r.name, status: r.status, kedaluwarsa: ymd(r.docExpiry) }))
      content.kelompok = intel.kelompok
      content.bdcS = intel.bcds
      content.workforce = intel.workforce
      break
    }
    case 'risiko':
      content.risiko = inH('Risk', hid).map((r) => ({
        kode: r.code, kategori: r.category, level: r.level, judul: r.title,
        owner: r.suggestedOwner, tenggat: ymd(r.deadline), status: r.status,
      }))
      break
    case 'aksi':
      content.aksi = inH('Action', hid).map((a) => ({
        kode: a.code, judul: a.title, prioritas: a.priority, status: a.status, owner: a.owner,
        tenggat: ymd(a.dueDate), resolusi: a.resolutionNote ?? '—',
      }))
      content.catatan = 'Closed-loop: resolusi aksi memicu usulan recalculation skor.'
      break
    case 'kepatuhan': {
      const regs = new Map(T('Regulation').map((r) => [r.id, r]))
      content.matriks = inH('ComplianceMapping', hid).map((m) => ({
        regulasi: regs.get(m.regulationId)?.number ?? '—', pasal: m.article, requirement: m.requirement,
        modul: m.module, status: m.status, owner: m.owner, review: ymd(m.lastReviewedAt),
      }))
      content.disclaimer = 'Designed to support compliance, subject to regulatory and hospital verification.'
      break
    }
    case 'klaim': {
      const claims = inH('Claim', hid)
      const anomalies = inH('ClaimAnomaly', hid)
      content.ccs = intel.ccs
      content.klaim = {
        total: claims.length,
        ditunda: claims.filter((c) => c.status === 'Ditunda').length,
        ditolak: claims.filter((c) => c.status === 'Ditolak').length,
      }
      content.anomali = {
        total: anomalies.length,
        tinggi: anomalies.filter((a) => a.severity === 'Tinggi').length,
        perluTinjauan: anomalies.filter((a) => a.status === 'Ditinjau').length,
      }
      content.disclaimer = 'Sinyal bersifat netral dan memerlukan verifikasi manual — bukan kesimpulan penyalahgunaan.'
      break
    }
    default: {
      content.domain = intel.domains.map((d: Row) => ({ domain: d.label, skor: d.score, bobot: d.weight }))
      content.tren = intel.trend
      content.saran = sortBy<Row>([...intel.domains], ['score', 'asc']).slice(0, 3).map((d) => ({
        domain: d.label, skor: d.score, catatan: d.negatives?.length ? d.negatives[0] : 'Pertahankan kinerja',
      }))
    }
  }

  const now = nowIso()
  const snap: Row = {
    id: newId(), hospitalId: hid, type: typ, title: meta.title, period,
    generatedAt: now, generatedBy: user.name,
    methodology: `Scoring engine ${DB.formulaVersion}; sumber: SIMRS/EMR/BPJS/SATUSEHAT via integration layer; verifikasi manual wajib untuk keputusan final.`,
    version: `v1.${inH('ReportSnapshot', hid).length + 1}`,
    auditRef: `AUD-${Date.now().toString(16).toUpperCase()}`,
    contentJson: JSON.stringify(content),
  }
  T('ReportSnapshot').push(snap)
  logAudit({
    hospitalId: hid, actor: user, action: 'create', entityType: 'Laporan', entityId: snap.id,
    description: `Generate laporan: ${snap.title}`,
  })
  return snap
}

/** Ekspor CSV laporan (dibuat di browser). */
export async function demoCsv(path: string, userId: string): Promise<void> {
  await db()
  const u = new URL(path, 'http://demo.local')
  const user = getUser(userId)
  const hid = resolveHospital(user, u.searchParams.get('hospitalId'))
  requirePerm(user, 'EXPORT')
  const snap = inH('ReportSnapshot', hid).find((r) => r.id === u.searchParams.get('export'))
  if (!snap) throw new ApiErr(404, 'Laporan tidak ditemukan.')
  const content = JSON.parse(snap.contentJson)
  const rows: string[][] = [['Field', 'Nilai']]
  const flat = (obj: Row, prefix: string) => {
    for (const k of Object.keys(obj).sort()) {
      const v = obj[k]
      if (Array.isArray(v)) rows.push([prefix + k, JSON.stringify(v)])
      else if (v && typeof v === 'object') flat(v, `${prefix}${k}.`)
      else rows.push([prefix + k, v === null || v === undefined ? '—' : String(v)])
    }
  }
  flat(content, '')
  const csv = rows.map((r) => r.map((c) => `"${c.replace(/"/g, '""')}"`).join(',')).join('\n') + '\n'
  logAudit({
    hospitalId: hid, actor: user, action: 'export', entityType: 'Laporan', entityId: snap.id,
    description: `Ekspor CSV: ${snap.title}`,
  })
  const blob = new Blob(['﻿' + csv], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `MedCredix-${snap.type}-${snap.period}.csv`
  document.body.appendChild(a)
  a.click()
  a.remove()
  setTimeout(() => URL.revokeObjectURL(url), 1000)
}

// ----- admin, pengaturan, data intelligence -----
handlers['GET /admin'] = ({ user }) => {
  if (user.role !== 'super_admin' && user.role !== 'auditor') {
    throw new ApiErr(403, 'Panel Admin hanya untuk Super Admin & Auditor.')
  }
  const rows = T('Hospital').map((h) => {
    const intel = DB.intel[h.id]
    return {
      id: h.id, name: h.name, code: h.code, class: h.class, city: h.city, province: h.province,
      tier: h.tier, onboarding: h.onboarding,
      healthScore: intel.healthScore, bcds: intel.bcds.final, stage: intel.bcds.stage, lulus: intel.bcds.lulus,
      trust: intel.trust.score, ccs: intel.ccs.score, cpcs: intel.cpcs.score,
      anomalyFlags: inH('DataQualityIssue', h.id).filter((i) => i.anomaly && i.status !== 'Diperbaiki').length,
      underReview: intel.trust.underReview,
      trend: intel.trend.map((t: Row) => ({ period: t.period, bcds: t.bcds, health: t.health })),
    }
  })
  const stageDist = Array.from({ length: 8 }, (_, s) => ({ stage: s, count: rows.filter((r) => r.stage === s).length }))
  const n = rows.length
  return {
    hospitals: rows,
    summary: n
      ? {
          totalRS: n,
          rataBdc: roundHalf(rows.reduce((a, r) => a + r.bcds, 0) / n),
          rataHealth: roundHalf(rows.reduce((a, r) => a + r.healthScore, 0) / n),
          lulus: rows.filter((r) => r.lulus).length,
          onboarding: rows.filter((r) => r.onboarding).length,
          underReview: rows.filter((r) => r.underReview).length,
        }
      : { totalRS: 0, rataBdc: 0, rataHealth: 0, lulus: 0, onboarding: 0, underReview: 0 },
    stageDist,
  }
}

handlers['GET /settings'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const hospital = T('Hospital').find((h) => h.id === hid)
  if (!hospital) throw new ApiErr(404, 'Data tidak ditemukan.')
  const users = sortBy(
    T('User').filter((u) => u.hospitalId === hid || (!u.hospitalId && ['super_admin', 'auditor'].includes(u.role))),
    ['name', 'asc'],
  ).map((u) => ({
    id: u.id, name: u.name, email: u.email, role: u.role, roleLabel: ROLE_LABELS[u.role] ?? u.role,
    title: u.title ?? null, active: u.active, lastLogin: u.lastLogin ?? null,
  }))
  const roles = Object.keys(ROLE_LABELS).map((key) => ({
    key, label: ROLE_LABELS[key], permissions: ROLE_PERMISSIONS[key] ?? ['VIEW'],
  }))
  return { hospital, users, roles, permissions: PERMISSIONS, rules: sortBy(inH('AlertRule', hid), ['name', 'asc']) }
}

handlers['PATCH /settings'] = ({ user, q, body }) => {
  if (!body) throw new ApiErr(400, 'Payload tidak valid.')
  const hid = resolveHospital(user, q.get('hospitalId'))
  requirePerm(user, 'CONFIGURE')
  if (body.target === 'hospital') {
    const h = T('Hospital').find((x) => x.id === hid)
    if (!h) throw new ApiErr(404, 'RS tidak ditemukan.')
    const prevName = h.name
    const map: Record<string, string> = {
      name: 'name', class: 'class', city: 'city', province: 'province', bedLicensed: 'bedLicensed',
      bpjsFacilityCode: 'bpjsFacilityCode', satusehatOrgId: 'satusehatOrgId',
    }
    let changed = false
    for (const [k, f] of Object.entries(map)) {
      if (body[k] !== undefined && body[k] !== null) { h[f] = body[k]; changed = true }
    }
    if (!changed) return h
    logAudit({
      hospitalId: hid, actor: user, action: 'config', entityType: 'Rumah Sakit', entityId: hid,
      description: 'Memperbarui profil RS', before: prevName, after: h.name,
    })
    return h
  }
  if (body.target === 'alert-rule') {
    const rule = inH('AlertRule', hid).find((r) => r.id === body.id)
    if (!rule) throw new ApiErr(404, 'Aturan tidak ditemukan.')
    const prev = rule.active
    rule.active = typeof body.active === 'boolean' ? body.active : !rule.active
    logAudit({
      hospitalId: hid, actor: user, action: 'config', entityType: 'Aturan Peringatan', entityId: rule.id,
      description: `${rule.name}: ${rule.active ? 'aktif' : 'nonaktif'}`, before: String(prev), after: String(rule.active),
    })
    return rule
  }
  throw new ApiErr(400, 'Target tidak dikenal.')
}

handlers['GET /data'] = ({ user, q }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  const typ = q.get('type') || 'all'
  const intel = DB.intel[hid]
  const hospital = T('Hospital').find((h) => h.id === hid)
  if (typ === 'quality') {
    const issues = sortBy(inH('DataQualityIssue', hid), ['status', 'asc'], ['severity', 'desc'])
    return {
      issues,
      summary: {
        terbuka: issues.filter((i) => i.status === 'Terbuka').length,
        total: issues.length,
        anomaly: issues.filter((i) => i.anomaly).length,
      },
    }
  }
  if (typ === 'trust') return { trust: intel.trust, bands: DB.trustBands, bcds: intel.bcds }
  const jobs = sortBy(inH('IntegrationJob', hid), ['startedAt', 'desc']).slice(0, 24)
  if (typ === 'sources') {
    return {
      sources: sortBy(inH('DataSource', hid), ['name', 'asc']),
      jobs,
      facilities: sortBy(inH('Facility', hid), ['category', 'asc']),
      beds: sortBy(inH('BedStat', hid), ['bedClass', 'asc']),
      hospital,
    }
  }
  return {
    issues: inH('DataQualityIssue', hid), sources: inH('DataSource', hid), jobs,
    facilities: inH('Facility', hid), beds: inH('BedStat', hid), hospital, trust: intel.trust,
  }
}

const FHIR_RESOURCES = [
  'Patient', 'Practitioner', 'PractitionerRole', 'Organization', 'Location', 'Encounter',
  'Condition', 'Observation', 'Procedure', 'Medication', 'MedicationRequest', 'MedicationDispense',
  'MedicationAdministration', 'AllergyIntolerance', 'DiagnosticReport', 'ServiceRequest',
  'ImagingStudy', 'CarePlan', 'Composition', 'Coverage', 'Claim', 'DocumentReference',
  'Immunization', 'Appointment',
]

handlers['POST /data'] = ({ user, q, body }) => {
  const hid = resolveHospital(user, q.get('hospitalId'))
  if (body?.action !== 'test-connection') throw new ApiErr(400, 'Aksi tidak dikenal.')
  const simrs = inH('DataSource', hid).find((s) => s.type === 'SIMRS')
  const badFields: Record<string, string> = {
    Composition: 'documentStatus kosong pada 12% record',
    MedicationDispense: 'whenHandedOver kosong pada 9% record',
  }
  const baseFhir = simrs ? simrs.fhirResources : 24
  const coverage = simrs ? simrs.coveragePct : 100
  let ok = 0
  let perlu = 0
  const results = FHIR_RESOURCES.map((r) => {
    const hasIssue = r in badFields && coverage < 100
    if (hasIssue) perlu++
    else ok++
    return {
      resource: r,
      status: hasIssue ? 'Perlu Perbaikan' : 'OK',
      message: hasIssue ? badFields[r] : 'Skema & field wajib valid',
      count: baseFhir * 137 + r.length * 53,
    }
  })
  logAudit({
    hospitalId: hid, actor: user, action: 'config', entityType: 'Integrasi',
    description: 'Test Koneksi 24 resource FHIR dijalankan',
  })
  return {
    testedAt: nowIso(),
    baseUrl: simrs?.configSummary ?? '—',
    results,
    summary: { ok, perlu, gagal: 0 },
  }
}

// ---------- Auth ----------
async function sha256Hex(s: string): Promise<string | null> {
  try {
    if (typeof crypto === 'undefined' || !crypto.subtle) return null
    const buf = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(s))
    return Array.from(new Uint8Array(buf)).map((b) => b.toString(16).padStart(2, '0')).join('')
  } catch {
    return null
  }
}

async function login(body: any) {
  const email = String(body?.email ?? '').toLowerCase().trim()
  const password = String(body?.password ?? '')
  if (!email || !password) throw new ApiErr(400, 'Email dan kata sandi wajib diisi.')
  const u = T('User').find((x) => x.email === email)
  const hash = await sha256Hex('medcredix::' + password)
  // fallback bila Web Crypto tidak tersedia (konteks non-HTTPS): bandingkan sandi demo langsung
  const valid = hash !== null ? hash === DB.passwordHash : password === 'medcredix2026'
  if (!u || !valid) throw new ApiErr(401, 'Email atau kata sandi salah.')
  if (!u.active) throw new ApiErr(403, 'Akun tidak aktif. Hubungi Super Admin.')
  u.lastLogin = nowIso()
  logAudit({
    hospitalId: u.hospitalId ?? null, actor: { name: u.name, role: u.role }, action: 'login',
    entityType: 'Sesi', description: 'Login berhasil',
  })
  return sessionOf(u)
}

// ---------- Entry point ----------
export async function demoApi<T = any>(path: string, method: string, body: unknown, userId: string): Promise<T> {
  await db()
  const u = new URL(path, 'http://demo.local')
  const key = `${method.toUpperCase()} ${u.pathname}`
  if (key === 'POST /auth') return (await login(body)) as T
  const handler = handlers[key]
  if (!handler) throw new ApiErr(404, 'Endpoint tidak ditemukan.')
  const user = getUser(userId)
  // simulasi latensi jaringan ringan agar skeleton/loading tetap terlihat wajar
  await new Promise((r) => setTimeout(r, 40))
  return handler({ user, q: u.searchParams, body }) as T
}
