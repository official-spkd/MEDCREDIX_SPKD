// Konfigurasi navigasi MedCredix (hash-router SPA) — port dari src/lib/nav.ts.
export type NavItem = { id: string; label: string; icon: string; roles?: string[]; badgeKey?: 'warnings' | 'actions' | 'risks' }
export type NavGroup = { id: string; label: string; items: NavItem[] }

export const NAV: NavGroup[] = [
  {
    id: 'utama',
    label: 'Utama',
    items: [
      { id: 'dashboard', label: 'Dashboard Eksekutif', icon: 'LayoutDashboard' },
      { id: 'risk-radar', label: 'Risk Radar', icon: 'Radar', badgeKey: 'risks' },
      { id: 'early-warnings', label: 'Early Warnings', icon: 'AlertTriangle', badgeKey: 'warnings' },
      { id: 'actions', label: 'Action Center', icon: 'CheckCircle2', badgeKey: 'actions' },
    ],
  },
  {
    id: 'credentialing',
    label: 'Credentialing',
    items: [
      { id: 'cred-practitioners', label: 'Praktisi', icon: 'Users' },
      { id: 'cred-credentials', label: 'Kredensial', icon: 'IdCard' },
      { id: 'cred-documents', label: 'Dokumen', icon: 'FileText' },
      { id: 'cred-verification', label: 'Verifikasi & Persetujuan', icon: 'BadgeCheck', roles: ['super_admin', 'admin_rs', 'credential_officer', 'medical_committee', 'auditor'] },
      { id: 'cred-prerequisites', label: 'Prasyarat Wajib', icon: 'ShieldCheck' },
    ],
  },
  {
    id: 'readiness',
    label: 'Kesiapan Rumah Sakit',
    items: [
      { id: 'ready-workforce', label: 'Clinical Workforce', icon: 'Stethoscope' },
      { id: 'ready-facilities', label: 'Sarana & Prasarana', icon: 'Building2' },
      { id: 'ready-digital', label: 'Digital Maturity', icon: 'MonitorSmartphone' },
      { id: 'ready-interop', label: 'Interoperabilitas', icon: 'Network' },
    ],
  },
  {
    id: 'clinical',
    label: 'Clinical Intelligence',
    items: [
      { id: 'clin-pathway', label: 'Clinical Pathway', icon: 'GitBranch' },
      { id: 'clin-utilization', label: 'Utilization Review', icon: 'BedDouble' },
      { id: 'clin-safety', label: 'Patient Safety', icon: 'HeartPulse' },
    ],
  },
  {
    id: 'claims',
    label: 'Claim Intelligence',
    items: [
      { id: 'claim-integrity', label: 'Claim Integrity', icon: 'FileCheck2' },
      { id: 'claim-anomalies', label: 'Anomali Klaim', icon: 'ScanSearch' },
    ],
  },
  {
    id: 'data',
    label: 'Data Intelligence',
    items: [
      { id: 'data-quality', label: 'Data Quality', icon: 'Database' },
      { id: 'data-trust', label: 'Data Trust', icon: 'ShieldHalf' },
      { id: 'data-sources', label: 'Data Sources & Integrasi', icon: 'Cable' },
    ],
  },
  {
    id: 'assurance',
    label: 'Bukti & Kepatuhan',
    items: [
      { id: 'evidence', label: 'Evidence Center', icon: 'Fingerprint' },
      { id: 'compliance', label: 'Regulatory Compliance', icon: 'Scale' },
      { id: 'audit', label: 'Audit Trail', icon: 'History' },
      { id: 'reports', label: 'Laporan', icon: 'FileBarChart' },
    ],
  },
  {
    id: 'settings',
    label: 'Pengaturan',
    items: [
      { id: 'settings-hospital', label: 'Profil RS & Integrasi', icon: 'Settings' },
      { id: 'settings-users', label: 'Pengguna & Peran', icon: 'UserCog' },
      { id: 'settings-alerts', label: 'Aturan Notifikasi', icon: 'Bell' },
      { id: 'admin', label: 'Panel Admin SPKD', icon: 'Globe2', roles: ['super_admin', 'auditor'] },
    ],
  },
]

export const VIEW_TITLES: Record<string, { title: string; sub: string }> = {
  dashboard: { title: 'Dashboard Eksekutif', sub: 'Ringkasan kesehatan & kesiapan rumah sakit' },
  'risk-radar': { title: 'Hospital Risk Radar', sub: 'Pemantauan risiko lintas 9 kategori' },
  'early-warnings': { title: 'Early Warning Center', sub: 'Deteksi dini berbasis ambang batas' },
  actions: { title: 'Action Center', sub: 'Ubah analitik menjadi tindakan perbaikan' },
  'cred-practitioners': { title: 'Praktisi', sub: 'Master tenaga medis & status lisensi' },
  'cred-credentials': { title: 'Kredensial', sub: 'Siklus hidup kredensial praktisi' },
  'cred-documents': { title: 'Dokumen Kredensial', sub: 'Arsip dokumen & status unggahan' },
  'cred-verification': { title: 'Verifikasi & Persetujuan', sub: 'Alur kerja verifikasi komite' },
  'cred-prerequisites': { title: 'Prasyarat Wajib', sub: 'Requirement Engine — dokumen legal RS' },
  'ready-workforce': { title: 'Clinical Workforce', sub: 'Cakupan spesialis & komposisi tenaga' },
  'ready-facilities': { title: 'Sarana & Prasarana', sub: 'Lifecycle fasilitas & ranjang' },
  'ready-digital': { title: 'Digital Maturity', sub: 'EMR, infrastruktur, medikasi, diagnostik' },
  'ready-interop': { title: 'Interoperabilitas', sub: 'SATUSEHAT, BPJS, validasi FHIR' },
  'clin-pathway': { title: 'Clinical Pathway Compliance', sub: 'Kepatuhan, deviasi & dampak finansial' },
  'clin-utilization': { title: 'Utilization Review', sub: 'ALOS, readmission, BOR vs target BPJS' },
  'clin-safety': { title: 'Patient Safety Intelligence', sub: 'Alergi, high-alert, dokumentasi pulang' },
  'claim-integrity': { title: 'Claim Integrity', sub: 'Claim Consistency Score & 6 sub-skor' },
  'claim-anomalies': { title: 'Anomali Klaim', sub: 'Sinyal netral — verifikasi manual wajib' },
  'data-quality': { title: 'Data Quality Center', sub: 'Kelengkapan, konsistensi, freshness' },
  'data-trust': { title: 'Data Trust Score', sub: 'Transparansi penuh — tanpa multiplier tersembunyi' },
  'data-sources': { title: 'Data Sources & Integration Health', sub: 'Status koneksi & kesehatan sinkronisasi' },
  evidence: { title: 'Evidence Center', sub: 'Telusuri setiap angka ke data sumber' },
  compliance: { title: 'Regulatory Compliance', sub: 'Matriks kepatuhan dengan versioning regulasi' },
  audit: { title: 'Audit Trail', sub: 'Log perubahan immutable' },
  reports: { title: 'MedCredix Reports', sub: 'Laporan resmi siap ekspor' },
  'settings-hospital': { title: 'Profil RS & Integrasi', sub: 'Konfigurasi koneksi API & SATUSEHAT' },
  'settings-users': { title: 'Pengguna & Peran', sub: 'Manajemen akses RBAC' },
  'settings-alerts': { title: 'Aturan Notifikasi', sub: 'Ambang batas & kanal peringatan' },
  admin: { title: 'Panel Admin SPKD', sub: 'Gambaran nasional lintas rumah sakit' },
}
