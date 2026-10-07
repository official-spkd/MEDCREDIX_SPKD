// Helper format & konstanta UI (client-safe) — port dari src/lib/format.ts.
const MONTHS_ID = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des']

export function fmtID(n: number, digits = 0): string {
  return (n ?? 0).toLocaleString('id-ID', { minimumFractionDigits: digits, maximumFractionDigits: digits })
}

export function fmtRupiah(n: number): string {
  return 'Rp ' + fmtID(n)
}

export function fmtDate(dt: string | Date | null | undefined): string {
  if (!dt) return '—'
  const d = new Date(dt)
  if (isNaN(d.getTime())) return '—'
  return `${d.getDate()} ${MONTHS_ID[d.getMonth()]} ${d.getFullYear()}`
}

export function fmtDateTime(dt: string | Date | null | undefined): string {
  if (!dt) return '—'
  const d = new Date(dt)
  if (isNaN(d.getTime())) return '—'
  return `${d.getDate()} ${MONTHS_ID[d.getMonth()]} ${d.getFullYear()}, ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

export function daysFromNow(dt: string | Date | null | undefined): number | null {
  if (!dt) return null
  return Math.ceil((new Date(dt).getTime() - Date.now()) / 86400000)
}

export function expiryLabel(
  dt: string | Date | null | undefined,
): { text: string; tone: 'success' | 'warning' | 'danger' | 'muted' } {
  const d = daysFromNow(dt)
  if (d == null) return { text: '—', tone: 'muted' }
  if (d < 0) return { text: `Kedaluwarsa ${Math.abs(d)} hari lalu`, tone: 'danger' }
  if (d <= 30) return { text: `${d} hari lagi`, tone: 'danger' }
  if (d <= 90) return { text: `${d} hari lagi`, tone: 'warning' }
  return { text: `${d} hari lagi`, tone: 'success' }
}

export function statusTone(status: string): 'success' | 'warning' | 'danger' | 'info' | 'muted' {
  if (
    ['Aktif', 'Valid', 'Disetujui', 'Selesai', 'Ditutup', 'Tervalidasi', 'Berhasil', 'Terhubung', 'Patuh', 'Baik', 'Teratasi', 'Diperbaiki', 'Diklarifikasi', 'Terbuka', 'Diverifikasi'].includes(status)
  ) return 'success'
  if (
    ['Segera Berakhir', 'Menunggu Verifikasi', 'Ditinjau', 'Diketahui', 'Dalam Pengerjaan', 'Perlu Perhatian', 'Ditugaskan', 'Sebagian', 'Terhubung Sebagian', 'Diajukan', 'Ditunda', 'Sedang Ditinjau'].includes(status)
  ) return 'warning'
  if (
    ['Kedaluwarsa', 'Ditolak', 'Gagal', 'Kritis', 'Tidak Tervalidasi', 'Ditangguhkan', 'Dicabut', 'Tidak Ada', 'Terputus', 'Belum'].includes(status)
  ) return 'danger'
  return 'info'
}

export const CRED_LIFECYCLE = [
  'Draft', 'Diajukan', 'Tinjauan Dokumen', 'Verifikasi', 'Tinjauan Komite', 'Disetujui',
  'Aktif', 'Segera Berakhir', 'Kedaluwarsa', 'Ditangguhkan', 'Dicabut',
]

export const ACTION_STATUSES = ['Baru', 'Ditugaskan', 'Dalam Pengerjaan', 'Menunggu Verifikasi', 'Selesai', 'Ditolak', 'Dibuka Kembali', 'Ditutup']

export const RISK_CATEGORIES = ['Kredensial', 'Klinis', 'Keselamatan Pasien', 'Klaim', 'Kualitas Data', 'Interoperabilitas', 'Medikasi', 'Operasional', 'Kepatuhan']

export const DOMAIN_LABELS: Record<string, string> = {
  workforce: 'Clinical Workforce',
  credentialing: 'Credentialing Readiness',
  digital: 'Digital Maturity',
  medication: 'Medication Management',
  diagnostic: 'Diagnostic Readiness',
  interop: 'Interoperability',
  claim: 'Claim Integrity',
  pathway: 'Clinical Pathway',
  safety: 'Patient Safety',
  trust: 'Data Trust',
}

export const AI_DISCLAIMER =
  'Insight ini dihasilkan otomatis (assistive). Verifikasikan terhadap data sumber sebelum mengambil keputusan.'
