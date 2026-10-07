package seed

// Data tenancy & credentialing: rumah sakit, pengguna, praktisi, prasyarat.
// Port setia dari prisma/seed/data-tenancy.ts.

type HospitalSeed struct {
	ID, Name, Code, Class, City, Province string
	BedLicensed, BedActive, BedNotOperating int
	NotOperatingNote, BpjsFacilityCode, SatusehatOrgID, Tier string
	Onboarding bool
}

var Hospitals = []HospitalSeed{
	{HNusantara, "RS Nusantara Sehat", "RSNS", "B", "Bandung", "Jawa Barat", 240, 196, 24,
		"Renovasi lantai 3 (12 TT), kalibrasi bed monitor (6 TT), penataan ruangan rawat inap (6 TT).",
		"3301R011", "org-034-9f2a-nusantara", "Enterprise", true},
	{HMedika, "RS Medika Sejahtera", "RSMS", "C", "Semarang", "Jawa Tengah", 120, 98, 12,
		"Pemeliharaan AC ruang rawat inap (8 TT), perbaikan plumbing (4 TT).",
		"3374R003", "org-041-7c1e-medika", "Growth", true},
}

type UserSeed struct {
	ID, Email, Name, Role, Title, Phone string
	HospitalID                          *string
}

func sp(s string) *string { return &s }

var Users = []UserSeed{
	{"u_super", "superadmin@spkd.co.id", "Surya Atmaja, M.Kom", "super_admin", "Platform Administrator SPKD", "+62 811-2001-001", nil},
	{"u_admin_ns", "admin@nusantarasehat.co.id", "dr. Aditya Nugraha, MARS", "admin_rs", "Direktur Utama / Admin RS", "+62 811-2001-002", sp(HNusantara)},
	{"u_mgmt_ns", "manajemen@nusantarasehat.co.id", "dr. Kartika Dewi", "hospital_mgmt", "Direktur Medis & Kepatuhan", "+62 811-2001-003", sp(HNusantara)},
	{"u_cred_ns", "kredensial@nusantarasehat.co.id", "Ratna Puspitasari, S.KM", "credential_officer", "Petugas Kredensialing", "+62 811-2001-004", sp(HNusantara)},
	{"u_komed_ns", "komite@nusantarasehat.co.id", "dr. Bambang Sutrisno, Sp.PD-KIC", "medical_committee", "Ketua Komite Kredensial & Medis", "+62 811-2001-005", sp(HNusantara)},
	{"u_sdm_ns", "sdm@nusantarasehat.co.id", "Tri Handayani, S.Psi", "hr_sdm", "Kepala SDM", "+62 811-2001-006", sp(HNusantara)},
	{"u_audit", "auditor@spkd.co.id", "Ir. Gunawan Wibisono, MHSA", "auditor", "Auditor Independen (ditunjuk)", "+62 811-2001-007", sp(HNusantara)},
	{"u_claim_ns", "klaim@nusantarasehat.co.id", "Silvia Rahmawati, AMd.Kep", "claims_team", "Koordinator Klaim JKN", "+62 811-2001-008", sp(HNusantara)},
	{"u_it_ns", "it@nusantarasehat.co.id", "Andika Pratama, S.Kom", "it_simrs", "Kepala IT / SIMRS", "+62 811-2001-009", sp(HNusantara)},
	{"u_view_ns", "viewer@nusantarasehat.co.id", "Dewi Anggraini, S.E", "viewer", "Staf Manajemen Mutu", "+62 811-2001-010", sp(HNusantara)},
	{"u_admin_ms", "admin@medikasejahtera.co.id", "dr. Rina Oktaviani", "admin_rs", "Direktur RS Medika Sejahtera", "+62 812-3001-001", sp(HMedika)},
	{"u_audit_bpjs", "auditorbpjs@spkd.co.id", "Joko Susilo, S.Sos", "auditor", "Auditor BPJS Kesehatan", "+62 812-3001-002", nil},
}

// PractitionerRow — tuple: [nama, profesi, spesialisasi, subspesialis, departemen, tipe, status, strExp(d), sipExp(d)]
type PractitionerRow struct {
	Name, Profession, Specialty, SubSpecialty, Department, EmploymentType, Status string
	StrExp, SipExp int
}

// PractitionersNS — RS Nusantara Sehat (28).
var PractitionersNS = []PractitionerRow{
	{"dr. Agus Prasetyo, Sp.PD", "Dokter Spesialis", "Penyakit Dalam", "", "Internal Medicine", "Tetap", "Aktif", 420, 300},
	{"dr. Ratna Wijaya, Sp.PD", "Dokter Spesialis", "Penyakit Dalam", "", "Internal Medicine", "Tetap", "Aktif", 510, 365},
	{"dr. Budi Santoso, Sp.B", "Dokter Spesialis", "Bedah Umum", "", "Surgery", "Tetap", "Aktif", 380, 25},
	{"dr. Sari Kusuma, Sp.A", "Dokter Spesialis", "Anak", "", "Pediatri", "Tetap", "Aktif", 460, 280},
	{"dr. Andi Wirawan, Sp.A", "Dokter Spesialis", "Anak", "", "Pediatri", "Kontrak", "Aktif", 330, 190},
	{"dr. Maya Laksmi, Sp.OG", "Dokter Spesialis", "Kebidanan & Kandungan", "", "Obstetri & Ginekologi", "Tetap", "Aktif", 550, 410},
	{"dr. Hendra Gunawan, Sp.N", "Dokter Spesialis", "Saraf", "", "Neurologi", "Tetap", "Aktif", 290, 80},
	{"dr. Rudi Hartono, Sp.JP", "Dokter Spesialis", "Jantung & Pembuluh Darah", "", "Kardiologi", "Tetap", "Aktif", 470, 355},
	{"dr. Intan Permata, Sp.P", "Dokter Spesialis", "Paru", "", "Pulmonologi", "Tetap", "Aktif", 350, 120},
	{"dr. Dimas Anggara, Sp.THT-KL", "Dokter Spesialis", "THT-KL", "", "THT-KL", "Tetap", "Aktif", 400, 210},
	{"dr. Yoga Pratama, Sp.OT", "Dokter Spesialis", "Bedah Orthopedi", "", "Orthopedi", "Tetap", "Aktif", 260, 24},
	{"dr. Fitri Handayani, Sp.An", "Dokter Spesialis", "Anestesiologi", "", "Anestesi", "Tetap", "Aktif", 440, 150},
	{"dr. Bagus Nugroho, Sp.Rad", "Dokter Spesialis", "Radiologi", "", "Radiologi", "Tetap", "Aktif", 510, 320},
	{"dr. Luna Kartika, Sp.PK", "Dokter Spesialis", "Patologi Klinik", "", "Laboratorium Klinik", "Kontrak", "Aktif", 300, 88},
	{"dr. Eka Putri, Sp.JP(K)", "Dokter Subspesialis", "Jantung & Pembuluh Darah", "Kardiologi Intervensi", "Kardiologi", "Tetap", "Aktif", 620, 480},
	{"dr. Tono Sugiarto, Sp.PD-KGE", "Dokter Subspesialis", "Penyakit Dalam", "Gastroenterologi-Hepatologi", "Internal Medicine", "Visiting Consultant", "Aktif", 240, 60},
	{"dr. Lina Marlina", "Dokter Umum", "", "", "IGD", "Tetap", "Aktif", 540, 390},
	{"dr. Rizky Ramadhan", "Dokter Umum", "", "", "Rawat Jalan", "Tetap", "Aktif", 480, 260},
	{"dr. Putri Ayu Lestari", "Dokter Umum", "", "", "Medical Check-Up", "Kontrak", "Aktif", 320, 170},
	{"dr. Galih Saputra", "Dokter Umum", "", "", "IGD", "Paruh Waktu", "Aktif", 280, 12},
	{"Ns. Dewi Sartika, S.Kep", "Perawat", "", "", "Rawat Inap", "Tetap", "Aktif", 500, 370},
	{"Ns. Rina Wulandari, S.Kep", "Perawat", "", "", "ICU", "Tetap", "Aktif", 430, 300},
	{"Ns. Fajar Nugroho, S.Kep", "Perawat", "", "", "IGD", "Kontrak", "Aktif", 360, 200},
	{"Ns. Sri Wahyuni, S.Kep", "Perawat", "", "", "Rawat Inap", "Tetap", "Cuti", 410, 150},
	{"Ns. Indah Permatasari, S.Kep", "Perawat", "", "", "Hemodialisa", "Kontrak", "Aktif", 310, -10},
	{"Bd. Nur Aini, A.Md.Keb", "Bidan", "", "", "Kebidanan", "Tetap", "Aktif", 470, 330},
	{"Apt. Dian Safitri, S.Farm", "Tenaga Kefarmasian", "", "", "Farmasi Klinik", "Tetap", "Aktif", 450, 340},
	{"Rudi Hartanto, A.Md.AK", "Tenaga Penunjang", "", "", "Laboratorium", "Tetap", "Aktif", 380, 220},
}

// PractitionersMS — RS Medika Sejahtera (10).
var PractitionersMS = []PractitionerRow{
	{"dr. Arif Budiman, Sp.PD", "Dokter Spesialis", "Penyakit Dalam", "", "Internal Medicine", "Tetap", "Aktif", 300, 150},
	{"dr. Sinta Nuriyah, Sp.A", "Dokter Spesialis", "Anak", "", "Pediatri", "Tetap", "Aktif", 280, 100},
	{"dr. Wahyu Kurniawan, Sp.B", "Dokter Spesialis", "Bedah Umum", "", "Surgery", "Tetap", "Aktif", 340, 40},
	{"dr. Lastri Anjani, Sp.OG", "Dokter Spesialis", "Kebidanan & Kandungan", "", "Obstetri & Ginekologi", "Kontrak", "Aktif", 260, 75},
	{"dr. Doni Prasetyo", "Dokter Umum", "", "", "Rawat Jalan", "Tetap", "Aktif", 400, 210},
	{"dr. Maya Sari", "Dokter Umum", "", "", "IGD", "Kontrak", "Aktif", 350, 120},
	{"Ns. Rani Oktaviani, S.Kep", "Perawat", "", "", "Rawat Inap", "Tetap", "Aktif", 420, 280},
	{"Ns. Bimo Prakoso, S.Kep", "Perawat", "", "", "IGD", "Tetap", "Aktif", 390, 160},
	{"Bd. Siti Aminah, A.Md.Keb", "Bidan", "", "", "Kebidanan", "Tetap", "Aktif", 370, 190},
	{"Apt. Fajar Ramadhan, S.Farm", "Tenaga Kefarmasian", "", "", "Farmasi", "Tetap", "Aktif", 310, 85},
}

// ExtraCredential — tuple: [idxPraktisi, tipe, nama, terbit(d), exp(d), status]
type ExtraCredential struct {
	PractitionerIdx int
	Type, Name      string
	Issue, Exp      int
	Status          string
}

var ExtraCredentialsNS = []ExtraCredential{
	{0, "Sertifikat Kompetensi", "Sertifikat Kompetensi Spesialis Penyakit Dalam", -700, 130, "Aktif"},
	{0, "Pelatihan", "ACLS Certification 2024", -640, 460, "Aktif"},
	{2, "Sertifikat Kompetensi", "Sertifikat Kompetensi Spesialis Bedah", -800, -20, "Kedaluwarsa"},
	{3, "Pelatihan", "Pelatihan PALS (Pediatric Advanced Life Support)", -540, 390, "Aktif"},
	{5, "Sertifikat Kompetensi", "Sertifikat Kompetensi Spesialis Obstetri & Ginekologi", -900, 210, "Aktif"},
	{6, "Pelatihan", "Pelatihan USG Transkraniel", -400, 90, "Segera Berakhir"},
	{10, "Sertifikat Kompetensi", "Sertifikat Kompetensi Spesialis Orthopedi", -750, 85, "Segera Berakhir"},
	{13, "Sertifikat Kompetensi", "Sertifikat Kompetensi Patologi Klinik", -820, 45, "Segera Berakhir"},
	{14, "Pelatihan", "Pelatihan Kateterisasi Jantung (PCI) Level 2", -300, 500, "Aktif"},
	{16, "Pelatihan", "Pelatihan Manajemen Bencana RS", -350, 60, "Segera Berakhir"},
	{18, "Pelatihan", "Pelatihan Interpretasi EKG Dasar", -500, 25, "Segera Berakhir"},
	{20, "Sertifikat Kompetensi", "Sertifikat Profesi Keperawatan (PPNI)", -880, 300, "Aktif"},
	{24, "Sertifikat Kompetensi", "Sertifikat Profesi Keperawatan (PPNI)", -870, 75, "Segera Berakhir"},
	{26, "Sertifikat Kompetensi", "Struktur Sertifikasi Apoteker", -720, 150, "Aktif"},
	{27, "Pelatihan", "Pelatihan Analis Kesehatan Bank Darah", -430, 20, "Segera Berakhir"},
}

// CredentialOverride — tuple: [idxPraktisi, jenisDoc(STR/SIP), status]
type CredentialOverride struct {
	PractitionerIdx int
	DocType, Status string
}

var CredentialOverridesNS = []CredentialOverride{
	{7, "SIP", "Verifikasi"},
	{13, "STR", "Tinjauan Komite"},
	{19, "SIP", "Tinjauan Komite"},
	{23, "SIP", "Verifikasi"},
	{15, "STR", "Verifikasi"},
	{11, "SIP", "Ditangguhkan"},
}

type RequirementSeed struct {
	Code, Name, Category, RegulatorySource, Version, DocType, ValidationRule, ResponsibleRole, Status string
	Required, EvidenceRequired, ApprovalRequired bool
	DocNumber, FileName, Notes                     string
	IssuedAt, DocExpiry                            int // offset hari; DocExpiry 1<<30 → null
	HasExpiry                                      bool
}

var Requirements = []RequirementSeed{
	{"P1", "Izin Operasional Rumah Sakit", "Legal", "UU 17/2023 & PP 28/2024", "v2026.1", "Izin Operasional",
		"Masa berlaku aktif dari Kemenkes/Dinkes", "admin_rs", "Valid", true, true, true,
		"445/RS-BDG/OPS/2024", "Izin-Operasional-RSNS-2024.pdf", "Diterbitkan Dinkes Jabar; berlaku s.d. Nov 2030.", -500, 1500, true},
	{"P2", "Perjanjian Kerja Sama (PKS) dengan BPJS Kesehatan", "Legal", "Perpres 82/2018 jo. 59/2024", "v2026.1", "PKS",
		"PKS aktif & tidak melewati tanggal berakhir", "admin_rs", "Valid", true, true, true,
		"PKS/0310/1026/2024", "PKS-BPJS-RSNS-2024.pdf", "Perpanjangan PKS tahunan — diajukan 45 hari sebelum berakhir.", -330, 35, true},
	{"P3", "Sertifikat Akreditasi RS (KARS)", "Akreditasi", "Kepmenkes 1128/2022 (standar akreditasi)", "v2026.1", "Sertifikat Akreditasi",
		"Status Dasar/Lanjutan/Madya — masa berlaku 3 tahun", "medical_committee", "Valid", true, true, true,
		"KARS-2024-11-0087", "Akreditasi-KARS-RSNS.pdf", "Akreditasi Dasar; re-akreditasi berikutnya perlu surveilan. Peringatan 90 hari aktif.", -980, 115, true},
	{"P4", "Izin Operasional Radiologi (PKR/PKRT)", "Legal", "Permenkes 3/2025 (keselamatan radiasi)", "v2026.1", "Izin Radiologi",
		"SK PKR aktif, operator bersertifikat", "it_simrs", "Menunggu Verifikasi", true, true, true,
		"PKR/RSNS/2026/014", "Izin-PKR-2026.pdf", "Dokumen diunggah, menunggu verifikasi Admin RS.", -20, 1430, true},
	{"P5", "Izin Pengelolaan Limbah B3 Medis", "Lingkungan", "PP 22/2021 & PermenLHK", "v2026.1", "Izin LB3",
		"Kontrak pengelolaan limbah aktif", "admin_rs", "Valid", true, true, true,
		"LB3/032/2025", "Izin-LB3-2025.pdf", "Kontrak vendor limbah s.d. Mar 2027.", -200, 530, true},
	{"P6", "Pakta Integritas Data", "Integritas", "Kebijakan SPKD MedCredix", "v2026.1", "Pernyataan",
		"Ditandatangani direktur; data operasional nyata", "admin_rs", "Valid", true, true, true,
		"PID-RSNS-2026", "Pakta-Integritas-Data.pdf", "Pernyataan bahwa data yang diakses MedCredix adalah data operasional nyata.", -120, 0, false},
}
