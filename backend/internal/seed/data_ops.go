package seed

// Data operasional: sumber data, integrasi, kualitas data, trust factor, evidence, audit, notifikasi.
// Port setia dari prisma/seed/data-ops.ts.

type DataSourceSeed struct {
	Name, Type, Status, AuthMethod, ConfigSummary string
	SuccessRate float64
	AvgLatencyMs, LastSyncDay, LastSyncHour int
	FhirResources int
	CoveragePct, Uptime30d float64
	Notes *string
}

var DataSourcesNS = []DataSourceSeed{
	{"SIMRS Inti Sehat (Core HIS)", "SIMRS", "Terhubung", "OAuth 2.0 Client Credentials", "Base URL FHIR internal; 24 resource TDS v2.0.0", 98.6, 210, 0, 8, 24, 96, 99.2, nil},
	{"EMR/RME Terpadu", "EMR/RME", "Terhubung", "OAuth 2.0 Client Credentials", "CPPT, resume medis, asuhan keperawatan", 97.2, 260, 0, 8, 19, 92, 98.9, nil},
	{"HRIS Nusantara", "HRIS", "Terhubung", "API Key", "Praktisi, lisensi, sertifikasi (endpoint REST)", 99.1, 140, 0, 6, 3, 98, 99.8, nil},
	{"BPJS VClaim (e-Klaim)", "BPJS VClaim", "Terhubung Sebagian", "OAuth 2.0 + User Key", "Kode faskes 3301R011; INACBGs bundling", 89.3, 480, -1, 14, 2, 84, 97.1, sp("Intermittent timeout pada jam sibuk 10.00-12.00; eskalasi ke vendor middleware.")},
	{"BPJS P-Care / JKN Mobile", "BPJS P-Care", "Terhubung", "User Key + Secret", "Data keluhan & kepuasan JKN", 94.7, 390, -1, 20, 1, 90, 98.2, nil},
	{"SATUSEHAT Platform", "SATUSEHAT", "Terhubung Sebagian", "OAuth 2.0 (Client ID/Secret Kemenkes)", "Organization ID org-034-9f2a-nusantara; cross-reference Encounter/Patient", 88.0, 620, -1, 2, 24, 88, 99.4, sp("2 resource ditolak karena field wajib kosong (Composition, MedicationDispense).")},
	{"LIS (Laboratory Information System)", "LIS", "Terhubung", "API Key", "Hasil lab, critical values, TAT", 96.4, 180, 0, 2, 4, 94, 99.0, nil},
	{"RIS/PACS (Radiologi)", "RIS/PACS", "Terhubung", "DICOM + API Key", "Hasil radiologi & imaging study", 95.8, 320, 0, 2, 3, 92, 98.7, nil},
	{"Sistem Farmasi", "Farmasi", "Terhubung Sebagian", "API Key", "e-Resep, dispensing, administrasi, stok", 91.5, 410, 0, 2, 4, 87, 98.4, sp("42 rekaman dispensasi gagal sinkron (retry queue aktif).")},
	{"Sistem Keuangan & Kasir", "Keuangan", "Terhubung", "API Key", "Piutang BPJS, rekonsiliasi pembayaran", 97.9, 150, -1, 20, 2, 95, 99.5, nil},
}

var DataSourcesMS = []DataSourceSeed{
	{"SIMRS Medika", "SIMRS", "Terhubung Sebagian", "API Key", "FHIR parsial (14 resource)", 88.2, 520, -1, 8, 14, 72, 96.8, sp("10 resource TDS belum tersedia dari vendor SIMRS.")},
	{"SATUSEHAT Platform", "SATUSEHAT", "Terhubung Sebagian", "OAuth 2.0", "Organization ID org-041-7c1e-medika", 85.0, 700, -1, 2, 14, 78, 99.2, nil},
}

// JobRow — [tipeDataSource, tipeJob, mulai(d), durasiMs, status, records, offsetMenitSelesai, pesan]
type JobRow struct {
	Source, JobType string
	StartDay        int
	DurationMs      int
	Status          string
	Records         int
	FinishMin       int
	Message         *string
}

var IntegrationJobs = []JobRow{
	{"SIMRS", "pull-data", 0, 42000, "Berhasil", 12480, 12, nil},
	{"EMR/RME", "pull-data", 0, 38000, "Berhasil", 8360, 9, nil},
	{"Farmasi", "pull-data", 0, 51000, "Gagal", 0, 0, sp("Timeout pada batch MedicationDispense #4 — 42 record masuk retry queue")},
	{"Farmasi", "retry", 0, 26000, "Berhasil", 28, 4, sp("28/42 rekaman berhasil disinkron ulang")},
	{"SATUSEHAT", "cross-ref", 0, 92000, "Gagal", 0, 0, sp("2 resource ditolak: field wajib kosong (Composition.documentStatus, MedicationDispense.whenHandedOver)")},
	{"LIS", "pull-data", 0, 21000, "Berhasil", 3120, 6, nil},
	{"RIS/PACS", "pull-data", 0, 29000, "Berhasil", 980, 8, nil},
	{"BPJS VClaim", "pull-data", -1, 84000, "Berhasil", 242, 22, sp("Latensi tinggi 480ms (normal 210ms)")},
	{"BPJS P-Care", "pull-data", -1, 45000, "Berhasil", 156, 14, nil},
	{"HRIS", "pull-data", 0, 9000, "Berhasil", 312, 3, nil},
	{"Keuangan", "pull-data", -1, 12000, "Berhasil", 486, 5, nil},
	{"SIMRS", "validate", 0, 15000, "Berhasil", 12480, 5, sp("Validasi skema FHIR lolos 96%")},
	{"SATUSEHAT", "pull-data", -1, 88000, "Berhasil", 11240, 25, nil},
	{"EMR/RME", "validate", 0, 11000, "Berhasil", 8360, 4, nil},
	{"SIMRS", "scoring", 0, 3000, "Berhasil", 9, 2, sp("Scoring engine inkremental: 9 dimensi diperbarui")},
	{"Farmasi", "pull-data", -1, 48000, "Berhasil", 5240, 15, nil},
	{"SATUSEHAT", "cross-ref", -1, 95000, "Berhasil", 11240, 27, sp("Cross-reference Encounter/Patient cocok 99.1%")},
	{"SIMRS", "pull-data", -1, 44000, "Berhasil", 12110, 13, nil},
}

type DqSeed struct {
	Dimension, Title, Description string
	AffectedRecords int
	SourceSystem, Severity string
	Anomaly bool
	Status string
	DetectedAt int
	ResolvedAt *int
}

var DqIssuesNS = []DqSeed{
	{"Integrasi Gagal", "42 rekaman dispensasi gagal sinkron", "Batch MedicationDispense #4 mengalami timeout; 42 rekaman menunggu retry (28 sudah berhasil pada sesi retry terakhir).", 42, "Farmasi", "Sedang", false, "Dalam Pengerjaan", -1, nil},
	{"Duplikasi", "12 rekaman pasien duplikat", "Potensi duplikat master pasien terdeteksi via kemiripan nama + tanggal lahir (12 pasangan kandidat).", 12, "SIMRS", "Sedang", false, "Terbuka", -12, nil},
	{"Freshness", "Data klaim tertinggal 1 siklus sinkronisasi", "Data VClaim tertinggal 1 pull (window 02.00) karena timeout; data terakhir 14.00 kemarin.", 242, "BPJS VClaim", "Rendah", false, "Diperbaiki", -1, ip(0)},
	{"Nilai Tidak Valid", "3% NIK tidak valid", "3% rekaman pasien memiliki NIK tidak valid (format/poli checksum) — memengaruhi indikator Validitas NIK 97%.", 386, "SIMRS", "Sedang", false, "Terbuka", -6, nil},
	{"Kelengkapan", "Resume medis belum lengkap 14%", "14% encounter rawat inap selesai belum memiliki Composition (resume medis) pada RME.", 96, "EMR/RME", "Sedang", false, "Terbuka", -5, nil},
	{"Konsistensi", "Selisih tanggal SEP vs Encounter", "9 klaim memiliki selisih tanggal SEP vs tanggal Encounter RME (1-2 hari) — memicu sinyal timeline.", 9, "BPJS VClaim", "Sedang", false, "Terbuka", -3, nil},
	{"Referensi Hilang", "8 order lab tanpa hasil", "8 ServiceRequest lab belum memiliki DiagnosticReport terkait (>72 jam) — perlu konfirmasi apakah dibatalkan.", 8, "LIS", "Rendah", false, "Terbuka", -4, nil},
	{"Data Bas", "Lonjakan volume MedicationDispense +42%", "Volume pull naik 42% di atas median 30 hari (sinyal anti-manipulasi). Cross-reference dengan SATUSEHAT menunjukkan lonjakan rawat inap yang wajar — menunggu konfirmasi final.", 0, "Farmasi", "Sedang", true, "Ditinjau", -5, nil},
	{"Data Bas", "Gap data Encounter 4 jam (window 02.00)", "Data gap detection menemukan tidak ada record Encounter pada 02.00-06.00; terkonfirmasi jadwal maintenance SIMRS.", 0, "SIMRS", "Rendah", true, "Diperbaiki", -8, ip(-7)},
	{"Kelengkapan", "Field wajib FHIR kosong pada 2 resource", "Composition.documentStatus dan MedicationDispense.whenHandedOver kosong pada sebagian record — menyebabkan penolakan SATUSEHAT.", 214, "SATUSEHAT", "Tinggi", false, "Dalam Pengerjaan", -2, nil},
}

type TrustSeed struct {
	Factor string
	Score, Weight float64
	Notes  string
}

var TrustFactorsNS = []TrustSeed{
	{"Kelengkapan Data", 94, 0.15, "94% field wajib terisi lintas resource"},
	{"Konsistensi Data", 91, 0.12, "Konsistensi lintas sistem (SIMRS vs SATUSEHAT) 91%"},
	{"Freshness", 88, 0.12, "Median keterlambatan data 6 jam; target <24 jam"},
	{"Keandalan Sumber", 96, 0.15, "Semua sumber terautentikasi; 0 sumber tanpa audit"},
	{"Tingkat Duplikasi", 87, 0.10, "Duplikasi MPI 0,8% (12/1.480)"},
	{"Kesehatan Sinkronisasi", 93, 0.12, "Antrian retry rendah; 2 job gagal / 18"},
	{"Tingkat Anomali", 90, 0.10, "2 anomaly flag aktif (< 3 → tidak Under Review)"},
	{"Status Validasi", 95, 0.08, "95% data lolos validasi skema FHIR"},
	{"Integritas Audit", 98, 0.06, "Audit log lengkap & immutable; 0 gap"},
}

// EvRow — [refType, refId, title, metric, value, unit, source, fhir, recordCount, validation]
type EvRow struct {
	RefType, RefID, Title, Metric, Value, Unit, Source string
	Fhir *string
	RecordCount int
	Validation string
}

var EvidenceNS = []EvRow{
	{"Skor", "DOMAIN:credentialing", "Kelengkapan kredensial praktisi", "Kelengkapan", "85.7", "%", "HRIS", sp("Practitioner"), 28, "Tervalidasi"},
	{"Skor", "DOMAIN:credentialing", "STR/SIP aktif", "Lisensi aktif", "96.4", "%", "HRIS", sp("Practitioner"), 28, "Tervalidasi"},
	{"Skor", "DOMAIN:credentialing", "Kredensial segera berakhir (≤90 hari)", "Expiry ≤90d", "14", "dokumen", "HRIS", sp("Practitioner + Credential"), 14, "Tervalidasi"},
	{"Skor", "DOMAIN:workforce", "Cakupan spesialis wajib (kelas B)", "Coverage", "92.3", "%", "HRIS", sp("Practitioner"), 13, "Tervalidasi"},
	{"Skor", "DOMAIN:workforce", "Komposisi tenaga medis", "Komposisi", "28", "praktisi", "HRIS", sp("Practitioner"), 28, "Tervalidasi"},
	{"Skor", "DOMAIN:claim", "Konsistensi diagnosis klaim", "Diagnosis", "82.1", "%", "BPJS VClaim", sp("Claim"), 42, "Tervalidasi"},
	{"Skor", "DOMAIN:claim", "Konsistensi obat klaim", "Medication", "76.2", "%", "BPJS VClaim", sp("Claim + MedicationRequest"), 42, "Tervalidasi"},
	{"Skor", "DOMAIN:claim", "Konsistensi stok-klaim (triangle)", "Stock", "75.0", "%", "Farmasi", sp("MedicationDispense"), 42, "Tervalidasi"},
	{"Skor", "DOMAIN:pathway", "Kepatuhan clinical pathway rata-rata", "Adherence", "71.4", "%", "EMR/RME", sp("CarePlan"), 26, "Perlu Validasi"},
	{"Skor", "DOMAIN:pathway", "Dokumentasi deviasi CP", "Deviasi terdokumentasi", "61.5", "%", "EMR/RME", sp("CarePlan"), 26, "Perlu Validasi"},
	{"Skor", "DOMAIN:safety", "Keamanan alergi obat", "Alergi", "96", "%", "EMR/RME", sp("AllergyIntolerance"), 3120, "Tervalidasi"},
	{"Skor", "DOMAIN:safety", "High-alert & closed-loop", "Closed-loop", "78", "%", "Farmasi", sp("MedicationAdministration"), 3100, "Tervalidasi"},
	{"Skor", "DOMAIN:safety", "Dokumentasi pemulangan", "Discharge", "85", "%", "EMR/RME", sp("Composition"), 684, "Tervalidasi"},
	{"Skor", "DOMAIN:interop", "Submission SATUSEHAT", "Submission", "88", "%", "SATUSEHAT", sp("All (TDS v2)"), 11240, "Tervalidasi"},
	{"Skor", "DOMAIN:interop", "Validasi FHIR", "Validasi", "84", "%", "SATUSEHAT", sp("All (TDS v2)"), 11240, "Tervalidasi"},
	{"Skor", "DOMAIN:medication", "Administrasi obat digital", "e-MAR", "68", "%", "Farmasi", sp("MedicationAdministration"), 2400, "Tervalidasi"},
	{"Skor", "DOMAIN:diagnostic", "Digitalisasi laboratorium", "Lab digital", "92", "%", "LIS", sp("DiagnosticReport"), 3120, "Tervalidasi"},
	{"Skor", "DOMAIN:digital", "Adopsi RME rawat inap", "RME RI", "88", "%", "EMR/RME", sp("Encounter"), 1140, "Tervalidasi"},
	{"Skor", "DOMAIN:digital", "Adopsi RME rawat jalan", "RME RJ", "79", "%", "EMR/RME", sp("Encounter"), 8360, "Tervalidasi"},
	{"Skor", "DOMAIN:trust", "Trust Score agregat", "Trust", "92.3", "/100", "MedCredix Engine", sp("-"), 0, "Tervalidasi"},
	{"Risiko", "RSK-001", "Daftar kredensial kedaluwarsa ≤90 hari", "Dokumen", "14", "item", "HRIS", sp("PractitionerCredential"), 14, "Tervalidasi"},
	{"Risiko", "RSK-002", "Klaim dengan sinyal inkonsistensi", "Klaim", "23", "item", "BPJS VClaim", sp("Claim"), 42, "Tervalidasi"},
	{"Risiko", "RSK-005", "Administrasi obat digital bulan berjalan", "e-MAR", "68", "%", "Farmasi", sp("MedicationAdministration"), 2400, "Tervalidasi"},
	{"Risiko", "RSK-008", "Sertifikat akreditasi KARS", "Masa berlaku", "115", "hari", "Manual Upload", sp("DocumentReference"), 1, "Tervalidasi"},
	{"Klaim", "CLM-2026-00128", "MedicationDispense vs item klaim", "Selisih", "2", "item", "Farmasi", sp("MedicationDispense"), 1, "Perlu Validasi"},
	{"Klaim", "CLM-2026-00128", "Prosedur vs laporan operasi", "Cocok", "0", "dokumen", "EMR/RME", sp("Procedure"), 1, "Perlu Validasi"},
	{"Klaim", "CLM-2026-00142", "Selisih stok vs klaim (triangle)", "Selisih", "3", "item", "Farmasi", sp("MedicationDispense"), 1, "Perlu Validasi"},
	{"Peringatan", "EW-007", "Volume pull MedicationDispense 30 hari", "Spike", "+42", "%", "Farmasi", sp("MedicationDispense"), 5240, "Perlu Validasi"},
	{"Regulasi", "KEPMENKES-84-2026", "Kesiapan resource CKG pada API RS", "Resource", "0", "dari 6", "SIMRS", sp("QuestionnaireResponse"), 6, "Tidak Tervalidasi"},
	{"Indikator", "DIG-EMR-RI", "Tren adopsi RME rawat inap 6 bulan", "Tren", "+11", "poin", "EMR/RME", sp("Encounter"), 1140, "Tervalidasi"},
}

// AuditRow — [aktor, peran, aksi, entitas, id, deskripsi, before, after]
type AuditRow struct {
	Actor, Role, Action, Entity, EntityID, Description string
	Before, After *string
}

var AuditLogs = []AuditRow{
	{"Andika Pratama, S.Kom", "it_simrs", "login", "Sesi", "-", "Login berhasil dari 10.20.4.18", nil, nil},
	{"Ratna Puspitasari, S.KM", "credential_officer", "login", "Sesi", "-", "Login berhasil dari 10.20.3.52", nil, nil},
	{"Silvia Rahmawati, AMd.Kep", "claims_team", "update", "Klaim", "CLM-2026-00128", "Menandai anomali prosedul 'Ditinjau'", sp("Baru"), sp("Ditinjau")},
	{"dr. Bambang Sutrisno, Sp.PD-KIC", "medical_committee", "approve", "Kredensial", "PC-014", "Menyetujui STR dr. Luna Kartika (komite medis)", sp("Tinjauan Komite"), sp("Disetujui")},
	{"Tri Handayani, S.Psi", "hr_sdm", "create", "Aksi", "ACT-004", "Membuat aksi perpanjangan SIP Ns. Indah Permatasari", nil, sp("Baru")},
	{"Andika Pratama, S.Kom", "it_simrs", "config", "Integrasi", "SATUSEHAT", "Mengubah konfigurasi retry timeout 30s → 60s", sp("timeout=30s"), sp("timeout=60s")},
	{"MedCredix Engine", "system", "scoring", "Skor", "SNAPSHOT-SEP", "Kalkulasi skor periodik: Health 81, BDCS 80, Trust 92", nil, sp("health=81")},
	{"MedCredix Engine", "system", "sync", "Integrasi", "Farmasi", "Pull-data farmasi gagal (timeout) — retry dijadwalkan", nil, sp("retry")},
	{"dr. Aditya Nugraha, MARS", "admin_rs", "verify", "Prasyarat", "P4", "Memverifikasi Izin Radiologi (PKR)", sp("Menunggu Verifikasi"), sp("Valid")},
	{"dr. Kartika Dewi", "hospital_mgmt", "export", "Laporan", "RPT-KRED-SEP", "Mengekspor Laporan Kredensial (PDF)", nil, sp("exported")},
	{"Silvia Rahmawati, AMd.Kep", "claims_team", "login", "Sesi", "-", "Login berhasil dari 10.20.5.31", nil, nil},
	{"Ratna Puspitasari, S.KM", "credential_officer", "update", "Kredensial", "PC-008", "Mengunggah dokumen sertifikat ACLS dr. Agus", nil, sp("Terunggah")},
	{"MedCredix Engine", "system", "sync", "Integrasi", "SATUSEHAT", "Cross-reference 02.45 selesai — cocok 99.1%", nil, nil},
	{"Ir. Gunawan Wibisono, MHSA", "auditor", "login", "Sesi", "-", "Login auditor (akses read-only) dari 10.30.1.7", nil, nil},
	{"Ir. Gunawan Wibisono, MHSA", "auditor", "view", "Evidence", "EV-SIP-EXPIRY", "Melihat bukti SIP kedaluwarsa", nil, nil},
	{"Andika Pratama, S.Kom", "it_simrs", "update", "Aksi", "ACT-010", "Update progres: 28/42 rekaman tersinkron", sp("0/42"), sp("28/42")},
	{"Ns. Dewi Sartika, S.Kep", "credential_officer", "update", "Aksi", "ACT-006", "Deploy e-MAR di ruangan Mawar & Tulip", sp("2 ruangan"), sp("4 ruangan")},
	{"MedCredix Engine", "system", "scoring", "Skor", "DOMAIN:claim", "Recalculation Claim Integrity: 79 (cccs subscore diperbarui)", sp("78.2"), sp("79.0")},
	{"dr. Budi Santoso, Sp.B", "medical_committee", "resolve", "Aksi", "ACT-013", "Menyelesaikan aksi sosialisasi discharge pathway bedah", sp("Ditinjau"), sp("Selesai")},
	{"Surya Atmaja, M.Kom", "super_admin", "config", "Platform", "ALERT-RULES", "Mengaktifkan aturan peringatan volume spike", sp("nonaktif"), sp("aktif")},
	{"Silvia Rahmawati, AMd.Kep", "claims_team", "update", "Anomali Klaim", "CA-003", "Mengklarifikasi anomali stok CLM-2026-00103 (stock opname)", sp("Ditinjau"), sp("Diklarifikasi")},
	{"Tri Handayani, S.Psi", "hr_sdm", "create", "Praktisi", "PR-029", "Menambah praktisi: dr. Spesialis Mata (kandidat rekrutmen)", nil, sp("Draft")},
	{"Ratna Puspitasari, S.KM", "credential_officer", "ack", "Peringatan", "EW-004", "Acknowledge peringatan dokumen RS mendekati kedaluwarsa", sp("Baru"), sp("Diketahui")},
	{"MedCredix Engine", "system", "sync", "Integrasi", "SIMRS", "Pull-data 02.00 selesai: 12.480 record", nil, nil},
	{"dr. Aditya Nugraha, MARS", "admin_rs", "verify", "Prasyarat", "P2", "Verifikasi PKS BPJS tercatat (masa berlaku 35 hari)", nil, sp("Valid")},
	{"Andika Pratama, S.Kom", "it_simrs", "create", "Aksi", "ACT-007", "Membuat aksi perbaikan mapping FHIR", nil, sp("Baru")},
	{"MedCredix Engine", "system", "scoring", "Skor", "DOMAIN:pathway", "Recalculation Clinical Pathway: 72.8 (episode baru masuk)", sp("71.4"), sp("72.8")},
	{"Silvia Rahmawati, AMd.Kep", "claims_team", "export", "Laporan", "RPT-CLAIM-SEP", "Ekspor Claim Integrity Report (CSV)", nil, sp("exported")},
	{"dr. Kartika Dewi", "hospital_mgmt", "approve", "Aksi", "ACT-013", "Verifikasi hasil resolusi aksi ACT-013", sp("Menunggu Verifikasi"), sp("Selesai")},
	{"MedCredix Engine", "system", "detect", "Risiko", "RSK-002", "Deteksi risiko klaim: 23 klaim perlu tinjauan", nil, sp("Terbuka")},
	{"MedCredix Engine", "system", "detect", "Peringatan", "EW-002", "Peringatan dini: sinyal inkonsistensi +18%", nil, sp("Baru")},
	{"Andika Pratama, S.Kom", "it_simrs", "config", "Integrasi", "BPJS VClaim", "Rotasi API key (scheduled rotation)", sp("key-2026-07"), sp("key-2026-10")},
	{"Ratna Puspitasari, S.KM", "credential_officer", "export", "Laporan", "RPT-KRED-OKT", "Generate Laporan Kredensial Oktober", nil, sp("generated")},
	{"dr. Bambang Sutrisno, Sp.PD-KIC", "medical_committee", "verify", "Kredensial", "PC-021", "Verifikasi dokumen kompetensi dr. Intan Permata", sp("Verifikasi"), sp("Disetujui")},
	{"Tri Handayani, S.Psi", "hr_sdm", "update", "Praktisi", "PR-025", "Perbarui status Ns. Indah: jadwal realokasi (SIP kedaluwarsa)", sp("Aktif"), sp("Aktif (realokasi)")},
	{"MedCredix Engine", "system", "sync", "Integrasi", "LIS", "Pull-data 02.00 LIS: 3.120 hasil lab", nil, nil},
	{"dr. Aditya Nugraha, MARS", "admin_rs", "login", "Sesi", "-", "Login berhasil dari 10.20.1.5", nil, nil},
	{"Silvia Rahmawati, AMd.Kep", "claims_team", "update", "Klaim", "CLM-2026-00138", "Melengkapi dokumen timeline klaim ditunda", sp("Ditunda"), sp("Diverifikasi")},
	{"MedCredix Engine", "system", "scoring", "Skor", "DOMAIN:safety", "Recalculation Patient Safety: 87.3", sp("86.8"), sp("87.3")},
	{"Surya Atmaja, M.Kom", "super_admin", "create", "Pengguna", "u_audit_bpjs", "Membuat akun Auditor BPJS (akses lintas RS)", nil, sp("aktif")},
}

type NotifSeed struct {
	Title, Body, Severity, Channel, Link string
}

var Notifications = []NotifSeed{
	{"12 kredensial berakhir dalam 30 hari", "Segera ajukan perpanjangan STR/SIP/Sertifikat. Lihat Risk Radar untuk daftar lengkap.", "Tinggi", "in-app", "#/warnings"},
	{"Sinyal inkonsistensi klaim naik 18%", "Rasio klaim dengan sinyal naik dari 38% ke 56% pada September.", "Tinggi", "in-app", "#/claims/anomalies"},
	{"Sertifikat Akreditasi KARS: 115 hari tersisa", "Mulai persiapan re-akreditasi. Aksi ACT-005 telah dibuat.", "Sedang", "in-app", "#/compliance/matrix"},
	{"Integrasi SATUSEHAT menolak 2 resource", "Field wajib kosong pada Composition & MedicationDispense.", "Sedang", "in-app", "#/data/sources"},
	{"Skor bulanan tersedia: Health Score 81", "Snapshot September tersedia. Lihat tren perbaikan 6 bulan.", "Rendah", "in-app", "#/dashboard"},
	{"ACT-013 diselesaikan", "Sosialisasi kriteria discharge pathway bedah selesai; deviasi terdokumentasi naik 31%.", "Rendah", "in-app", "#/actions"},
}
