package seed

// Data risiko, early warning, alert rules, action center, digital, sarana, ranjang, standar spesialis.
// Port setia dari prisma/seed/data-intel.ts.

type RiskSeed struct {
	Code, Category, Level, Title, Description string
	DetectedAt                                int // offset hari
	Domain, AffectedRecords, EvidenceJson, BusinessImpact, RecommendedAction, SuggestedOwner string
	Deadline                                  int
	Status                                    string
}

var RisksNS = []RiskSeed{
	{"RSK-001", "Kredensial", "Tinggi", "13 kredensial mendekati kedaluwarsa dalam 90 hari",
		"Terdeteksi 13 dokumen kredensial (STR/SIP/Sertifikat) yang akan kedaluwarsa dalam 90 hari; 5 di antaranya dalam 30 hari; 2 dokumen telah kedaluwarsa. Praktisi dengan SIP kedaluwarsa tidak dapat menjalankan pelayanan sesuai regulasi.",
		-2, "Credentialing Readiness", `["Kredensial: 14 item"]`, `["EV-CRED-EXPIRY"]`,
		"Potensi penghentian pelayanan praktisi terdampak; risiko temuan saat surveilan akreditasi; gangguan jadwal DPJP.",
		"Ajukan perpanjangan melalui KKMS/ konsil segera; buat aksi perpanjangan dengan tenggat sebelum tanggal kedaluwarsa.",
		"Ratna Puspitasari (Petugas Kredensialing)", 14, "Terbuka"},
	{"RSK-002", "Klaim", "Sedang", "26 klaim memerlukan tinjauan konsistensi",
		"26 dari 42 klaim terakhir memiliki minimal satu sinyal inkonsistensi (warn/bad) pada kategori diagnosis, obat, prosedur, diagnostik, stok, atau timeline. Verifikasi manual wajib sebelum pengajuan.",
		-3, "Claim Integrity", `["Klaim: 26 item"]`, `["EV-CLAIM-CCS"]`,
		"Potensi penundaan/penolakan klaim BPJS; beban kerja verifikator meningkat; potensi kehilangan pendapatan.",
		"Tinjau klaim berstatus Ditunda/Diverifikasi di modul Claim Intelligence; selesaikan inkonsistensi timeline & stok terlebih dahulu.",
		"Silvia Rahmawati (Tim Klaim)", 15, "Terbuka"},
	{"RSK-003", "Klinis", "Sedang", "Kepatuhan Clinical Pathway menurun 7% dalam 3 periode",
		"Rata-rata kepatuhan clinical pathway turun dari 78% (Jun) menjadi 71% (Sep). Penurunan terkonsentrasi pada pathway Stroke Iskemik (64%) dan Pneumonia dengan komorbid.",
		-5, "Clinical Pathway", `["Pathway: 3 dari 8 aktif"]`, `["EV-CP-ADH"]`,
		"Efisiensi lama rawat menurun; biaya per episode naik di atas tarif INA-CBGs; indikasi efisiensi klaim memburuk.",
		"Audit rekam medis episode dengan adherence <70%; sosialisasikan protokol CP ke DPJP terkait; lengkapi dokumentasi alasan deviasi.",
		"dr. Bambang Sutrisno (Komite Medis)", 20, "Ditinjau"},
	{"RSK-004", "Kualitas Data", "Rendah", "42 rekaman menunggu sinkronisasi antar sistem",
		"42 rekaman (Encounter, MedicationDispense) gagal tersinkron pada 2 pull terakhir dari SIMRS; berada dalam antrian retry.",
		-1, "Data Quality", `["Record: 42 item"]`, `["EV-DQ-SYNC"]`,
		"Skor dapat understated karena data belum lengkap; laporan bulanan belum final.",
		"Pantau retry queue di Data Sources; eskalasi ke vendor SIMRS bila antrian > 50 record.",
		"Andika Pratama (IT/SIMRS)", 7, "Terbuka"},
	{"RSK-005", "Medikasi", "Tinggi", "Digitalisasi administrasi obat 68% — di bawah ambang 80%",
		"Siklus administrasi obat (MedicationAdministration) baru terdokumentasi digital pada 68% pemberian obat; di bawah ambang keselamatan 80%.",
		-4, "Medication Management", `["Pemberian obat: ~2.400/bulan"]`, `["EV-MED-ADMIN"]`,
		"Risiko medication error tidak terlacak; closed-loop medication belum terbentuk penuh; skor keselamatan pasien tertahan.",
		"Aktifkan barcode scanning administrasi obat di ruangan tersisa; latih perawat; target 85% dalam 60 hari.",
		"Ns. Dewi Sartika (Ka. Keperawatan)", 45, "Dalam Pengerjaan"},
	{"RSK-006", "Interoperabilitas", "Sedang", "Validasi FHIR 87% — 13% resource ditolak SATUSEHAT",
		"Dari 24 resource FHIR yang dikirim, validasi berhasil 87%. Kegagalan terbesar pada Composition (resume medis) dan MedicationDispense karena field wajib kosong.",
		-6, "Interoperability", `["Resource FHIR: 24 tipe"]`, `["EV-INT-FHIR"]`,
		"Kepatuhan Pelaporan SATUSEHAT menurun; berpotensi memengaruhi penilaian interoperabilitas & sanksi administratif.",
		"Perbaiki mapping field wajib di middleware; jalankan Test Koneksi 24 resource; koordinasi dengan vendor SIMRS.",
		"Andika Pratama (IT/SIMRS)", 12, "Ditinjau"},
	{"RSK-007", "Kredensial", "Sedang", "1 SIP kedaluwarsa & 3 SIP berakhir dalam 30 hari",
		"SIP Ns. Indah Permatasari telah kedaluwarsa 10 hari lalu (praktik ditahan sesuai kebijakan). Tiga SIP lain (dr. Budi Santoso, dr. Yoga Pratama, dr. Galih Saputra) berakhir dalam 25-30 hari.",
		-10, "Credentialing Readiness", `["Praktisi: 4 orang"]`, `["EV-SIP-EXPIRY"]`,
		"Praktisi bersangkutan tidak boleh memberikan pelayanan medis sampai SIP aktif kembali; jadwal layanan hemodialisa berisiko.",
		"Segera urus perpanjangan SIP ke Dinkes; realokasi jadwal sementara; catat sebagai temuan kredensialing.",
		"Tri Handayani (HR/SDM)", 5, "Dalam Pengerjaan"},
	{"RSK-008", "Kepatuhan", "Tinggi", "Akreditasi KARS berakhir dalam 115 hari",
		"Sertifikat akreditasi dasar KARS berlaku hingga 26 Januari 2027. Persiapan surveilan re-akreditasi perlu dimulai segera (standar akreditasi 2022).",
		-7, "Regulatory Compliance", `["Dokumen: P3 Akreditasi"]`, `["EV-REG-KARS"]`,
		"Tidak lolos re-akreditasi = prasyarat kredensial BPJS gugur; berdampak pada kontrak JKN.",
		"Bentuk tim persiapan re-akreditasi; jadwalkan internal assessment; tutup temuan mutu terbuka.",
		"dr. Kartika Dewi (Direktur Medis)", 30, "Ditinjau"},
	{"RSK-009", "Operasional", "Rendah", "BOR 81% mendekati batas atas optimal (85%)",
		"Bed Occupancy Rate September mencapai 81% — masih dalam rentang optimal 60-85% namun meningkat 4 poin dalam 2 bulan.",
		-8, "Utilization", `["Ranjang aktif: 196"]`, `["EV-UR-BOR"]`,
		"Bila melewati 85% risiko refused admission & menurunnya kualitas layanan.",
		"Optimalkan discharge planning; pantau BOR mingguan; siapkan rencana tambah TT kelas 3 bila tren berlanjut.",
		"dr. Aditya Nugraha (Direktur)", 30, "Terkelola"},
	{"RSK-010", "Keselamatan Pasien", "Sedang", "Closed-loop medication 74% (target 80%)",
		"Rantai resep-dispense-administrasi terdokumentasi penuh pada 74% episode; target keselamatan pasien 80%.",
		-4, "Patient Safety", `["Episode medikasi: ~3.100/bulan"]`, `["EV-PS-CLM"]`,
		"Kesempatan deteksi kesalahan obat (high-alert) berkurang.",
		"Terapkan e-MAR dengan verifikasi barcode di semua ruang rawat inap; audit kepatuhan perawat bulanan.",
		"Apt. Dian Safitri (Farmasi Klinik)", 45, "Terbuka"},
	{"RSK-011", "Kualitas Data", "Rendah", "12 rekaman pasien duplikat terdeteksi",
		"Deteksi duplikasi berbasis kemiripan nama+tanggal lahir menemukan 12 potensi rekaman ganda pada master pasien index (MPI).",
		-12, "Data Quality", `["Pasien (MPI): 12 rekaman"]`, `["EV-DQ-DUP"]`,
		"Risiko keselamatan pasien (riwayat medis terpecah); klaim dapat tertolak karena data pasien tidak konsisten.",
		"Verifikasi & merger rekaman melalui prosedur MPI; aktifkan validasi NIK saat registrasi.",
		"Andika Pratama (IT/SIMRS)", 21, "Terbuka"},
	{"RSK-012", "Klinis", "Sedang", "ALOS kasus bedah 5.9 hari melebihi standar pathway",
		"ALOS episode bedah umum 5.9 hari vs standar 4-5 hari; kontributor utama: 6 episode dengan deviasi LOS tanpa alasan terdokumentasi.",
		-9, "Clinical Pathway", `["Episode bedah: 6"]`, `["EV-CP-BEDAH"]`,
		"Biaya per episode 115% dari tarif; margin klaim bedah tertekan.",
		"Tinjau kriteria discharge; evaluasi episode dengan LOS > standar; wajibkan alasan klinis pada deviasi.",
		"dr. Bambang Sutrisno (Komite Medis)", 25, "Terbuka"},
}

type WarningSeed struct {
	Code, Title, Description, Severity, Domain, ThresholdRule, TriggerValue string
	DetectedAt                                                              int
	Status                                                                  string
	AckBy                                                                   *string
	AckAt                                                                   *int
	Assigned                                                                *string
	Due                                                                     *int
}

func ip(n int) *int { return &n }

var WarningsNS = []WarningSeed{
	{"EW-001", "5 kredensial praktisi berakhir dalam 30 hari",
		"Ambang 30 hari terpicu: 5 dokumen kredensial (STR/SIP/Sertifikat) akan kedaluwarsa. Segera ajukan perpanjangan.",
		"Tinggi", "Credentialing Readiness", "credential_expiry <= 30 hari", "5 dokumen", -1, "Baru", nil, nil, nil, nil},
	{"EW-002", "Sinyal inkonsistensi klaim naik 18% bulan ini",
		"Rasio klaim dengan sinyal warn naik dari 38% ke 56% pada periode September. Perlu ditinjau pola pengkodean & kronologi SEP.",
		"Tinggi", "Claim Integrity", "inkonsistensi_klaim_bulanan >= +15%", "+18%", -2, "Baru", nil, nil, nil, nil},
	{"EW-003", "Kepatuhan Clinical Pathway turun 3 periode berturut-turut",
		"Adherence rata-rata 78% → 75% → 71% (Jul-Sep). Perlu intervensi Komite Medis.",
		"Sedang", "Clinical Pathway", "cp_adherence turun >= 3 periode", "3 periode", -3, "Diketahui",
		sp("dr. Bambang Sutrisno, Sp.PD-KIC"), ip(-2), nil, nil},
	{"EW-004", "Dokumen legal RS mendekati kedaluwarsa (PKS 35 hari; KARS 115 hari)",
		"PKS BPJS berakhir dalam 35 hari dan Sertifikat Akreditasi KARS 115 hari. Siapkan dokumen perpanjangan dan tim re-akreditasi.",
		"Sedang", "Regulatory Compliance", "document_expiry <= 90 hari", "PKS 35h; KARS 115h", -4, "Ditugaskan",
		sp("Ratna Puspitasari, S.KM"), ip(-3), sp("Ratna Puspitasari, S.KM"), ip(10)},
	{"EW-005", "Success rate integrasi SATUSEHAT turun di bawah ambang",
		"Success rate pengiriman resource turun ke 88% (ambang 95%) pada 3 window sinkronisasi terakhir.",
		"Sedang", "Interoperability", "integration_success < 95%", "88%", -2, "Diketahui",
		sp("Andika Pratama, S.Kom"), ip(-1), nil, nil},
	{"EW-006", "Administrasi obat digital di bawah ambang keselamatan",
		"Digitalisasi MedicationAdministration 68% — di bawah ambang 80%. Kombinasi dengan sinyal closed-loop 74%.",
		"Tinggi", "Medication Management", "med_admin_digital < 80%", "68%", -4, "Ditugaskan",
		sp("Ns. Dewi Sartika, S.Kep"), ip(-3), sp("Ns. Dewi Sartika, S.Kep"), ip(30)},
	{"EW-007", "Lonjakan volume data MedicationDispense tidak wajar",
		"Volume pull MedicationDispense naik 42% vs median 30 hari (sinyal anti-manipulasi). Bisa jadi wajar (lonjakan rawat inap) — perlu konfirmasi silang dengan SATUSEHAT.",
		"Sedang", "Data Trust", "volume_spike >= +40% vs median", "+42%", -5, "Diketahui",
		sp("Andika Pratama, S.Kom"), ip(-4), nil, nil},
	{"EW-008", "ALOS rawat inap mendekati target atas",
		"ALOS 5.3 hari (target ≤5). Tren menurun 4 bulan terakhir namun belum mencapai target.",
		"Rendah", "Utilization", "alos >= 5.0 hari", "5.3 hari", -6, "Teratasi",
		sp("dr. Aditya Nugraha, MARS"), ip(-5), nil, nil},
	{"EW-009", "Keluhan kategori antrian turun 3 bulan berturut",
		"Keluhan antrian 6 → 5 → 4 → 3 — tren positif pasung implementasi antrian digital.",
		"Rendah", "Patient Experience", "kpi_positif_tren (informatif)", "turun 3 periode", -7, "Teratasi",
		sp("dr. Kartika Dewi"), ip(-6), nil, nil},
	{"EW-010", "Gap data Encounter terdeteksi pada pull 02.00",
		"Window sinkronisasi 02.00 mencatat gap 4 jam pada resource Encounter (tidak ada data). Data gap detection aktif.",
		"Rendah", "Data Quality", "data_gap > 2 jam", "4 jam", -8, "Teratasi",
		sp("Andika Pratama, S.Kom"), ip(-8), nil, nil},
}

type AlertRuleSeed struct {
	Name, Domain, Metric, Condition, Severity, Channels string
	TriggeredCount                                       int
	LastTriggered                                        *int
}

var AlertRules = []AlertRuleSeed{
	{"Kredensial mendekati kedaluwarsa (30 hari)", "Credentialing Readiness", "credential_expiry_days", "expiration_days <= 30", "Tinggi", "in-app,email,whatsapp", 5, ip(-1)},
	{"Kenaikan sinyal inkonsistensi klaim", "Claim Integrity", "claim_inconsistency_ratio", "monthly_change >= +15%", "Tinggi", "in-app,email", 3, ip(-2)},
	{"Kepatuhan CP turun beruntun", "Clinical Pathway", "cp_adherence", "decrease >= 3 periode", "Sedang", "in-app,email", 1, ip(-3)},
	{"Dokumen legal RS mendekati kedaluwarsa (90 hari)", "Regulatory Compliance", "document_expiry_days", "expiration_days <= 90", "Sedang", "in-app,email", 2, ip(-4)},
	{"Success rate integrasi turun", "Interoperability", "integration_success_rate", "value < 95%", "Sedang", "in-app", 5, ip(-2)},
	{"Digitalisasi administrasi obat rendah", "Medication Management", "med_admin_digital", "value < 80%", "Tinggi", "in-app,email,whatsapp", 2, ip(-4)},
	{"Lonjakan volume data (anti-manipulasi)", "Data Trust", "resource_volume_spike", "spike >= +40% vs median-30d", "Sedang", "in-app", 1, ip(-5)},
	{"BOR melewati batas optimal", "Utilization", "bor", "value > 85%", "Rendah", "in-app", 0, nil},
}

type ActionSeed struct {
	Code, Title, Description, Priority, Status, SourceType string
	SourceRef                                              *string
	Owner, OwnerRole                                       string
	Due                                                    int
	Expected                                               string
	Progress                                               *string
	Resolved                                               *string
	ResolvedAt                                             *int
}

var ActionsNS = []ActionSeed{
	{"ACT-001", "Lengkapi 6 kredensial praktisi yang belum lengkap", "6 praktisi memiliki dokumen kredensial belum lengkap (2 STR/SIP, 4 sertifikat kompetensi). Kumpulkan dokumen & ajukan verifikasi.", "Tinggi", "Dalam Pengerjaan", "Risiko", sp("RSK-001"), "Ratna Puspitasari, S.KM", "credential_officer", 12, "Kelengkapan kredensial naik ke ≥95%", sp("4 dari 6 dokumen sudah dikumpulkan; 2 menunggu dari praktisi."), nil, nil},
	{"ACT-002", "Tinjau 26 klaim dengan sinyal inkonsistensi", "Selesaikan tinjauan 26 klaim (prioritas: 5 klaim Ditunda/Ditolak dengan sinyal tinggi pada prosedur, stok, dan diagnosis).", "Tinggi", "Ditugaskan", "Risiko", sp("RSK-002"), "Silvia Rahmawati, AMd.Kep", "claims_team", 15, "CCS naik dari 79 ke ≥85", nil, nil, nil},
	{"ACT-003", "Audit kepatuhan Clinical Pathway Stroke & Pneumonia", "Audit rekam medis episode dengan adherence <70% pada pathway Stroke Iskemik dan Pneumonia; sosialisasikan protokol.", "Sedang", "Ditinjau", "Peringatan", sp("EW-003"), "dr. Bambang Sutrisno, Sp.PD-KIC", "medical_committee", 20, "Adherence kembali ≥78%; dokumentasi deviasi ≥75%", nil, nil, nil},
	{"ACT-004", "Perpanjang SIP Ns. Indah Permatasari", "SIP telah kedaluwarsa 10 hari. Ajukan perpanjangan ke Dinkes Kota Bandung; realokasi jadwal HD sementara.", "Tinggi", "Dalam Pengerjaan", "Risiko", sp("RSK-007"), "Tri Handayani, S.Psi", "hr_sdm", 5, "SIP aktif kembali ≤7 hari", sp("Berkas diserahkan ke Dinkes 2 hari lalu; menunggu verifikasi."), nil, nil},
	{"ACT-005", "Bentuk tim persiapan re-akreditasi KARS", "Akreditasi berakhir 115 hari. Bentuk tim, jadwalkan internal assessment standar 2022, tutup temuan mutu terbuka.", "Tinggi", "Baru", "Risiko", sp("RSK-008"), "dr. Kartika Dewi", "hospital_mgmt", 30, "Internal assessment selesai sebelum surveilan", nil, nil, nil},
	{"ACT-006", "Aktifkan barcode administrasi obat di 4 ruangan tersisa", "Digitalisasi MedicationAdministration 68%. Deploy e-MAR + barcode di ruangan yang belum memakai.", "Tinggi", "Dalam Pengerjaan", "Risiko", sp("RSK-005"), "Ns. Dewi Sartika, S.Kep", "credential_officer", 45, "Administrasi obat digital ≥85%", sp("2 ruangan terdeploy (progress 50%)."), nil, nil},
	{"ACT-007", "Perbaiki mapping field FHIR Composition & MedicationDispense", "Validasi FHIR 87%. Field wajib kosong pada Composition (resume) & MedicationDispense; perbaiki middleware + test koneksi.", "Sedang", "Ditugaskan", "Risiko", sp("RSK-006"), "Andika Pratama, S.Kom", "it_simrs", 12, "Validasi FHIR ≥95%", nil, nil, nil},
	{"ACT-008", "Verifikasi & merger 12 rekaman pasien duplikat", "Gunakan prosedur MPI; verifikasi identitas sebelum merger; aktifkan validasi NIK saat registrasi.", "Sedang", "Baru", "Risiko", sp("RSK-011"), "Andika Pratama, S.Kom", "it_simrs", 21, "Duplikasi < 0,5%", nil, nil, nil},
	{"ACT-009", "Ajukan perpanjangan PKS BPJS Kesehatan", "PKS berakhir dalam 35 hari. Siapkan dokumen perpanjangan & rapat koordinasi dengan KC BPJS setempat.", "Sedang", "Menunggu Verifikasi", "Peringatan", sp("EW-004"), "dr. Aditya Nugraha, MARS", "admin_rs", 10, "PKS baru terbit sebelum berakhirnya masa lama", sp("Draft dokumen terkirim ke manajemen BPJS."), nil, nil},
	{"ACT-010", "Selesaikan 42 rekaman antrean sinkronisasi", "Retry queue sync Encounter/MedicationDispense; eskalasi vendor bila antrian >50.", "Rendah", "Dalam Pengerjaan", "Risiko", sp("RSK-004"), "Andika Pratama, S.Kom", "it_simrs", 7, "Antrian retry <10 record", sp("28/42 rekaman berhasil disinkronkan ulang."), nil, nil},
	{"ACT-011", "Rekrut dokter spesialis Mata", "Spesialis Mata belum tersedia (kebutuhan wajib kelas B). Buka lowongan & koordinasi dengan konsultan HR.", "Sedang", "Ditugaskan", "Skor", sp("DOMAIN:workforce"), "Tri Handayani, S.Psi", "hr_sdm", 60, "Coverage spesialis 13/13", nil, nil, nil},
	{"ACT-012", "Selesaikan 4 klaim Ditunda & 1 Ditolak", "5 klaim bermasalah (4 Ditunda, 1 Ditolak) dengan sinyal timeline, stok, dan diagnosis; lengkapi dokumen lalu ajukan ulang.", "Sedang", "Ditugaskan", "Klaim", sp("CLM-2026-00129"), "Silvia Rahmawati, AMd.Kep", "claims_team", 10, "Klaim Disetujui kembali", nil, nil, nil},
	{"ACT-013", "Sosialisasi kriteria discharge pathway bedah", "ALOS bedah 5.9 hari vs standar 4-5; sosialisasi kriteria discharge & wajibkan alasan klinis deviasi LOS.", "Sedang", "Selesai", "Risiko", sp("RSK-012"), "dr. Budi Santoso, Sp.B", "medical_committee", -5, "ALOS bedah ≤5.2 hari", sp("Sosialisasi selesai; adherence LOS membaik."), sp("Sosialisasi 2 sesi selesai 28 Sep; deviasi LOS terdokumentasi naik 31% → 58%."), ip(-5)},
	{"ACT-014", "Deploy antrian digital poliklinik", "Keluhan antrian tertinggi (Apr-Jun). Deploy sistem antrian digital + SMS notifikasi.", "Sedang", "Ditutup", "Manual", nil, "Andika Pratama, S.Kom", "it_simrs", -30, "Keluhan antrian turun 50%", sp("Teralisasi penuh."), sp("Keluhan antrian turun dari 6 ke 3 per bulan (turun 50%). Skor kepuasan naik 82→87."), ip(-30)},
	{"ACT-015", "Latih perawat input vital sign terstruktur", "Kelengkapan data observasi vital sign sebelumnya 84%; latih perawat input terstruktur via EMR.", "Rendah", "Ditutup", "Skor", sp("DOMAIN:diagnostic"), "Ns. Rina Wulandari, S.Kep", "credential_officer", -45, "Kelengkapan vital sign ≥90%", sp("Selesai."), sp("Kelengkapan vital sign 84% → 93% (di atas ambang 90%)."), ip(-45)},
	{"ACT-016", "Evaluasi kebutuhan tambah TT kelas 3", "BOR 81% mendekati batas atas; kaji kelayakan menambah 12 TT kelas 3 pasca renovasi.", "Rendah", "Ditinjau", "Risiko", sp("RSK-009"), "dr. Aditya Nugraha, MARS", "admin_rs", 30, "Keputusan investasi TT terdokumentasi", nil, nil, nil},
}

// DigRow — [dimensi, indikator, nilai, unit, ambang, bobot, status, fhir, sumber, arah]
type DigRow struct {
	Dimension, Indicator string
	Value                float64
	Unit                 string
	Threshold, Weight    float64
	Status               string
	Fhir                 *string
	Source               string
	Direction            string
}

var DigitalNS = []DigRow{
	{"EMR", "Adopsi RME Rawat Inap", 88, "%", 85, 0.4, "Baik", sp("Encounter"), "EMR/RME", "up"},
	{"EMR", "Adopsi RME Rawat Jalan", 79, "%", 85, 0.3, "Perlu Perhatian", sp("Encounter"), "EMR/RME", "up"},
	{"EMR", "Kelengkapan Coding ICD-10", 91, "%", 90, 0.1, "Baik", sp("Condition"), "SIMRS", "up"},
	{"EMR", "Validitas NIK Pasien", 97, "%", 98, 0.1, "Perlu Perhatian", sp("Patient"), "SIMRS", "up"},
	{"EMR", "Kelengkapan Resume Medis", 86, "%", 90, 0.1, "Perlu Perhatian", sp("Composition"), "EMR/RME", "up"},
	{"Infrastruktur", "Uptime API RS (30 hari)", 99.2, "%", 99, 0.3, "Baik", sp("-"), "SIMRS", "up"},
	{"Infrastruktur", "Success Rate Submission", 96.4, "%", 98, 0.3, "Perlu Perhatian", sp("-"), "SATUSEHAT", "up"},
	{"Infrastruktur", "Latensi API Rata-rata", 340, "ms", 300, 0.2, "Perlu Perhatian", sp("-"), "SIMRS", "down"},
	{"Infrastruktur", "Cakupan Multi-Departemen", 84, "%", 90, 0.2, "Perlu Perhatian", sp("Encounter"), "SIMRS", "up"},
	{"Medikasi", "CPOE (Computerized Provider Order Entry)", 82, "%", 90, 0.25, "Perlu Perhatian", sp("MedicationRequest"), "Farmasi", "up"},
	{"Medikasi", "Dispensing Digital", 76, "%", 85, 0.2, "Perlu Perhatian", sp("MedicationDispense"), "Farmasi", "up"},
	{"Medikasi", "Administrasi Obat Digital", 68, "%", 80, 0.25, "Kritis", sp("MedicationAdministration"), "Farmasi", "up"},
	{"Medikasi", "Closed-Loop Medication", 74, "%", 85, 0.2, "Perlu Perhatian", sp("MedicationAdministration"), "Farmasi", "up"},
	{"Medikasi", "Kelengkapan Data Alergi", 93, "%", 95, 0.1, "Perlu Perhatian", sp("AllergyIntolerance"), "EMR/RME", "up"},
	{"Diagnostik", "Digitalisasi Laboratorium", 92, "%", 90, 0.25, "Baik", sp("DiagnosticReport"), "LIS", "up"},
	{"Diagnostik", "Digitalisasi Radiologi", 89, "%", 90, 0.25, "Perlu Perhatian", sp("ImagingStudy"), "RIS/PACS", "up"},
	{"Diagnostik", "Kelengkapan DiagnosticReport", 88, "%", 90, 0.2, "Perlu Perhatian", sp("DiagnosticReport"), "LIS", "up"},
	{"Diagnostik", "Lab Turnaround Time", 3.2, "jam", 4, 0.15, "Baik", sp("DiagnosticReport"), "LIS", "down"},
	{"Diagnostik", "Critical Value Flagging", 85, "%", 95, 0.15, "Perlu Perhatian", sp("Observation"), "LIS", "up"},
	{"Interoperabilitas", "Submission SATUSEHAT", 88, "%", 95, 0.35, "Perlu Perhatian", sp("All (TDS v2)"), "SATUSEHAT", "up"},
	{"Interoperabilitas", "Validasi FHIR", 84, "%", 95, 0.35, "Perlu Perhatian", sp("All (TDS v2)"), "SATUSEHAT", "up"},
	{"Interoperabilitas", "Sukses Integrasi BPJS", 89, "%", 95, 0.3, "Perlu Perhatian", sp("Claim"), "BPJS VClaim", "up"},
	{"Keselamatan", "Keamanan Alergi Obat", 96, "%", 95, 0.4, "Baik", sp("AllergyIntolerance"), "EMR/RME", "up"},
	{"Keselamatan", "High-Alert & Closed-Loop", 78, "%", 80, 0.3, "Perlu Perhatian", sp("MedicationAdministration"), "Farmasi", "up"},
	{"Keselamatan", "Dokumentasi Pemulangan Pasien", 85, "%", 80, 0.3, "Baik", sp("Composition"), "EMR/RME", "up"},
}

// FacRow — [nama, kategori, availability, functional, feasible, calibration, suitable, standar, catatan]
type FacRow struct {
	Name, Category, Availability, Functional, Feasible, Calibration, Suitable, Standard string
	Notes *string
}

var FacilitiesNS = []FacRow{
	{"Instalasi Gawat Darurat (IGD)", "IGD", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: pelayanan medik darurat", nil},
	{"Kamar Operasi (6 meja)", "Kamar Operasi", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: pelayanan bedah", sp("2 OK sedang upgrade laparoskopi (tetap beroperasi)")},
	{"Intensive Care Unit (10 bed)", "ICU", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: pelayanan intensif", nil},
	{"Laboratorium Klinik", "Laboratorium", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: laboratorium", sp("Kalibrasi alat hematology analyzer berikutnya Nov 2026")},
	{"Instalasi Radiologi", "Radiologi", "Ada", "Berfungsi", "Perlu Perbaikan", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: radiologi", sp("CT-Scan 128 slice: pending part replacement — layanan dialihkan sementara")},
	{"Instalasi Farmasi", "Farmasi", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: kefarmasian", nil},
	{"Poliklinik Rawat Jalan", "Poliklinik", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: pelayanan medik", nil},
	{"Unit Hemodialisa (8 mesin)", "Hemodialisa", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes: yankes subspesial", sp("2 mesin maintenance terjadwal Okt")},
	{"CSSD (Sterilisasi Sentral)", "CSSD", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: sterilisasi sentral", nil},
	{"Unit Transfusi Darah", "Bank Darah", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: pelayanan darah", sp("Kerja sama dengan UTD PMI")},
	{"Instalasi Gizi", "Gizi", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: gizi", nil},
	{"Pemulasaraan Jenazah", "Pemulasaraan Jenazah", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Sesuai Kebutuhan", "Permenkes 6/2026: pemulasaraan jenazah", nil},
	{"Cathlab", "Cathlab", "Tidak Ada", "Tidak Berfungsi", "Tidak Laik", "Belum Kalibrasi", "Kurang", "Fasilitas opsional (bukan wajib kelas B)", sp("Rencana pengadaan 2027 (opsional — tidak mengurangi skor wajib)")},
	{"MRI 1.5 Tesla", "Radiologi", "Ada", "Berfungsi", "Laik Pakai", "Valid", "Melebihi Kebutuhan", "Opsional — di atas standar minimal", sp("Nilai tambah di atas standar; meningkatkan grade sarana")},
}

// BedRow — [ward, kelas, total, aktif, layakPakai, alasan]
type BedRow struct {
	Ward, BedClass  string
	Total, Active, Usable int
	Reason          *string
}

var BedsNS = []BedRow{
	{"Anggrek", "VIP", 8, 8, 8, nil},
	{"Melati", "Kelas 1", 32, 30, 30, nil},
	{"Mawar", "Kelas 2", 64, 58, 58, nil},
	{"Tulip Lt.3", "Kelas 3", 96, 84, 72, sp("Renovasi lantai 3: 12 TT belum aktif")},
	{"ICU", "ICU", 12, 10, 8, sp("Kalibrasi bed monitor: 2 TT tidak aktif")},
	{"Perinatologi", "Perinatologi", 8, 6, 6, nil},
}

// SpecStd — [spesialis, wajib, catatan] — standar kelas B (Permenkes 56/2014).
type SpecStd struct {
	Specialty string
	Required  int
	Note      *string
}

var SpecialtyStandards = []SpecStd{
	{"Penyakit Dalam", 2, sp("Dua SDM karena beban rawat inap tinggi (kebijakan internal)")},
	{"Bedah Umum", 1, nil},
	{"Anak", 2, sp("Dua SDM karena beban pelayanan anak tinggi")},
	{"Kebidanan & Kandungan", 1, nil},
	{"Saraf", 1, nil},
	{"Jantung & Pembuluh Darah", 1, nil},
	{"Paru", 1, nil},
	{"THT-KL", 1, nil},
	{"Mata", 1, sp("Kosong — rekrutmen berjalan (ACT-011)")},
	{"Bedah Orthopedi", 1, nil},
	{"Anestesiologi", 1, nil},
	{"Radiologi", 1, nil},
	{"Patologi Klinik", 1, nil},
}
