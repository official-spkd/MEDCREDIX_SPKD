package seed

// Data regulasi, versi regulasi, compliance mapping, snapshot skor.
// Port setia dari prisma/seed/data-compliance.ts.

type RegulationSeed struct {
        Code, Title, Number, Authority string
        Year                           int
        Category, Status, Summary, Source string
}

var Regulations = []RegulationSeed{
        {"UU-17-2023", "Undang-Undang Kesehatan", "UU No. 17 Tahun 2023", "DPR RI", 2023, "Undang-Undang", "Berlaku",
                "Fondasi hukum transformasi kesehatan: satu data kesehatan, SATUSEHAT, sanksi administratif, kewajiban fasilitas kesehatan.", "https://jdih.kemkes.go.id"},
        {"PP-28-2024", "Peraturan Pelaksanaan UU Kesehatan", "PP No. 28 Tahun 2024", "Pemerintah RI", 2024, "Peraturan Pemerintah", "Berlaku",
                "Rincian teknis penyelenggaraan RS: perizinan, rekam medis elektronik, kerahasiaan data, kewajiban interoperabilitas.", "https://jdih.setneg.go.id"},
        {"UU-27-2022", "Perlindungan Data Pribadi (PDP)", "UU No. 27 Tahun 2022", "DPR RI", 2022, "Undang-Undang", "Berlaku",
                "Dasar hukum pemrosesan data pribadi pasien: dasar pemrosesan, hak subjek data, keamanan, penyalahgunaan.", "https://jdih.kominfo.go.id"},
        {"PERPRES-82-2018", "Jaminan Kesehatan", "Perpres No. 82 Tahun 2018", "Presiden RI", 2018, "Peraturan Presiden", "Perubahan Terbaru",
                "Kerangka JKN: kredensial faskes, klaim, sanksi. Telah diubah dengan Perpres 59/2024.", "https://jdih.setneg.go.id"},
        {"PERMENKES-56-2014", "Klasifikasi dan Perizinan Rumah Sakit", "Permenkes No. 56 Tahun 2014", "Menteri Kesehatan", 2014, "Permenkes", "Berlaku",
                "Standar SDM minimal per kelas RS (kelas B: 13 spesialis dasar), sarana, tempat tidur, perizinan.", "https://jdih.kemkes.go.id"},
        {"PERMENKES-13-2022", "Rekam Medis Elektronik", "Permenkes No. 13 Tahun 2022", "Menteri Kesehatan", 2022, "Permenkes", "Berlaku",
                "Kewajiban RME: keamanan, kerahasiaan, integritas, retensi 5 tahun (10 tahun ringkasan), pertukaran data.", "https://jdih.kemkes.go.id"},
        {"PERMENKES-24-2022", "Penyelenggaraan Satu Data Bidang Kesehatan (SATUSEHAT)", "Permenkes No. 24 Tahun 2022", "Menteri Kesehatan", 2022, "Permenkes", "Berlaku",
                "Integrasi data kesehatan via platform SATUSEHAT; standar interoperabilitas & FHIR; sanksi non-integrasi.", "https://jdih.kemkes.go.id"},
        {"PERMENKES-6-2026", "Pelayanan Minimal Rumah Sakit", "Permenkes No. 6 Tahun 2026", "Menteri Kesehatan", 2026, "Permenkes", "Berlaku",
                "RS sekurang-kurangnya menyelenggarakan: pelayanan medik, intensif, bedah, keperawatan/kebidanan, kefarmasian, laboratorium, radiologi, darah, gizi, pemulasaraan jenazah, sterilisasi sentral, serta pemeliharaan sarana/prasarana dan alat kesehatan.", "https://jdih.kemkes.go.id/documents/peraturan-menteri-kesehatan-nomor-6-tahun-2026"},
        {"KEPMENKES-84-2026", "Petunjuk Teknis Cek Kesehatan Gratis (CKG)", "Kepmenkes HK.01.07/MENKES/84/2026", "Menteri Kesehatan", 2026, "Kepmenkes", "Berlaku",
                "Penyelenggaraan CKG tahunan sepanjang siklus hidup; pencatatan & pelaporan terintegrasi via SATUSEHAT; peran fasilitas pelayanan kesehatan & penyelenggara sistem elektronik.", "https://jdih.kemkes.go.id"},
        {"KEPMENKES-1128-2022", "Standar Akreditasi Rumah Sakit", "Kepmenkes HK.01.07/MENKES/1128/2022", "Menteri Kesehatan", 2022, "Kepmenkes", "Berlaku",
                "Versi standar akreditasi 2022 (KARS): 12 bab versi nasional — termasuk manajemen data & informasi (MRP).", "https://jdih.kemkes.go.id"},
        {"KARS", "Akreditasi RS (KARS) — Program Nasional", "Sertifikat KARS", "Komisi Akreditasi Rumah Sakit", 2024, "Akreditasi", "Berlaku",
                "Sertifikasi eksternal mutu RS: akreditasi dasar, madya, utama. Prasyarat kredensial BPJS.", "https://kars.or.id"},
}

type RegVersionSeed struct {
        Reg, Version string
        IssuedAt     string // tanggal ISO
        Status, ChangeSummary string
        ImpactLevel  *string
        SupersededBy *string
}

var RegulationVersions = []RegVersionSeed{
        {"PERPRES-82-2018", "asli-2018", "2018-10-15", "Superseded", "Versi asli JKN", nil, sp("Perpres 59/2024")},
        {"PERPRES-82-2018", "perubahan-3 (Perpres 59/2024)", "2024-05-20", "Berlaku", "Perubahan ketiga: iuran JKN, klasifikasi faskes, tata cara kredensial ulang", sp("Sedang"), nil},
        {"KEPMENKES-84-2026", "Kepmenkes 33/2025 (versi lama)", "2025-03-01", "Superseded", "Petunjuk teknis CKG Hari Ulang Tahun (versi awal)", nil, sp("Kepmenkes 84/2026")},
        {"KEPMENKES-84-2026", "Kepmenkes 84/2026", "2026-09-22", "Berlaku", "Penyesuaian cakupan & mutu CKG: sasaran siklus hidup, integrasi SATUSEHAT Mobile, pelaporan", sp("Tinggi"), nil},
        {"PERMENKES-6-2026", "Permenkes 3/2020 (versi lama)", "2020-02-04", "Superseded", "Standar pelayanan RS lama (11 jenis pelayanan)", nil, sp("Permenkes 6/2026")},
        {"PERMENKES-6-2026", "Permenkes 6/2026", "2026-02-10", "Berlaku", "Rumusan ulang pelayanan minimal RS termasuk pemeliharaan sarana/prasarana & alat kesehatan; penyelarasan UU 17/2023", sp("Tinggi"), nil},
        {"PERMENKES-56-2014", "Permenkes 56/2014", "2014-12-02", "Berlaku", "Standar klasifikasi & perizinan RS (masih berlaku untuk standar SDM)", nil, nil},
}

// ComplianceRow — [regCode, pasal, requirement, module, feature, dataField, validation, evidenceRef, audit, owner, status, lastReview(d), reviewer]
type ComplianceRow struct {
        Reg, Article, Requirement, Module, Feature, DataField, Validation, EvidenceRef string
        AuditRequired bool
        Owner, Status string
        LastReview int
        Reviewer string
}

var ComplianceMappings = []ComplianceRow{
        {"UU-17-2023", "Pasal 171", "RS wajib menyelenggarakan satu data kesehatan (SATUSEHAT)", "Interoperabilitas", "Integrasi SATUSEHAT", "DataSource.successRate", "successRate >= 95%", "EV-INT-FHIR", true, "IT/SIMRS", "Sebagian", -30, "Andika Pratama"},
        {"UU-27-2022", "Pasal 20-22", "Pemrosesan data pribadi pasien berdasarkan persetujuan & tujuan terbatas", "Keamanan Data", "Masking identitas & kontrol ekspor", "Claim.patientInit", "Inisial saja pada laporan; NIK tidak diekspor", "EV-DQ-DUP", true, "IT/SIMRS", "Patuh", -20, "Andika Pratama"},
        {"PP-28-2024", "Pasal 252", "RME wajib memenuhi standar keamanan & pertukaran data", "Digital Maturity", "Adopsi RME", "DigitalIndicator.EMR", "adopsi >= 85%", "EV-EMR-ADOPT", true, "IT/SIMRS", "Sebagian", -15, "dr. Kartika Dewi"},
        {"PP-28-2024", "Pasal 266", "Kerahasiaan rekam medis & kontrol akses berbasis peran", "Audit & Keamanan", "RBAC + Audit Trail", "AuditLog", "Semua akses tercatat; ekspor terkontrol", "EV-AUDIT-RBAC", true, "Admin RS", "Patuh", -10, "dr. Aditya Nugraha"},
        {"PERPRES-82-2018", "Perubahan ketiga", "Fasilitas kesehatan mitra JKN wajib menjaga mutu & akuntabilitas klaim", "Claim Intelligence", "Claim Consistency Score", "Claim.flagsJson", "CCS >= 75 (GOOD)", "EV-CLAIM-CCS", true, "Tim Klaim", "Sebagian", -8, "Silvia Rahmawati"},
        {"PERMENKES-56-2014", "Lampiran 1.3", "Kelas B: minimal 13 spesialis dasar", "Clinical Workforce", "Specialist Coverage", "Practitioner.specialty", "coverage >= 13/13", "EV-WF-COV", true, "HR/SDM", "Sebagian", -12, "Tri Handayani"},
        {"PERMENKES-56-2014", "Lampiran 1.5", "Sarana wajib RS kelas B lengkap & laik pakai", "Sarana & Prasarana", "Facility Lifecycle", "Facility.feasible", "Semua fasilitas wajib Laik Pakai", "EV-FAC-LIFE", true, "Admin RS", "Patuh", -12, "dr. Aditya Nugraha"},
        {"PERMENKES-13-2022", "Pasal 6-8", "Keamanan & integritas RME (audit perubahan)", "Digital Maturity", "Audit Trail RME", "AuditLog.entityType", "Perubahan RME tercatat immutable", "EV-AUDIT-RBAC", true, "IT/SIMRS", "Patuh", -10, "Andika Pratama"},
        {"PERMENKES-13-2022", "Pasal 14", "Retensi rekam medis 5 tahun (ringkasan 10 tahun)", "Digital Maturity", "Retensi arsip", "-", "Kebijakan retensi terdokumentasi", "EV-EMR-ADOPT", true, "Rekam Medis", "Patuh", -60, "dr. Kartika Dewi"},
        {"PERMENKES-24-2022", "Pasal 7", "Kewajiban integrasi data kesehatan ke SATUSEHAT", "Interoperabilitas", "Validasi FHIR 24 resource", "DigitalIndicator.Interop", "validasi >= 95%", "EV-INT-FHIR", true, "IT/SIMRS", "Sebagian", -30, "Andika Pratama"},
        {"PERMENKES-24-2022", "Pasal 11", "Kualitas data yang dikirimkan (lengkap, sah, mutakhir)", "Data Intelligence", "Data Quality Center", "DataQualityIssue", "Freshness <= 24 jam; duplikasi < 1%", "EV-DQ-SYNC", true, "IT/SIMRS", "Sebagian", -15, "Andika Pratama"},
        {"PERMENKES-6-2026", "Pasal 3", "Pelayanan minimal: medik, intensif, bedah, keperawatan/kebidanan, kefarmasian, lab, radiologi, darah, gizi, pemulasaraan jenazah, sterilisasi sentral", "Sarana & Prasarana", "Checklist 12 layanan wajib", "Facility.category", "12 kategori terpenuhi", "EV-FAC-LIFE", true, "Admin RS", "Patuh", -5, "dr. Aditya Nugraha"},
        {"PERMENKES-6-2026", "Pasal 3 (huruf l)", "Pemeliharaan sarana/prasarana & alat kesehatan", "Sarana & Prasarana", "Kalibrasi & pemeliharaan", "Facility.calibration", "Kalibrasi valid untuk alat wajib", "EV-FAC-LIFE", true, "Teknik/UMPT", "Sebagian", -5, "dr. Aditya Nugraha"},
        {"KEPMENKES-84-2026", "Diktum KETIGA", "CKG terintegrasi dengan pelayanan kesehatan sepanjang siklus hidup", "Clinical Intelligence", "Kesiapan data CKG", "-", "Pemetaan alur data CKG ke SATUSEHAT", "EV-CKG-READY", true, "IT/SIMRS", "Sedang Ditinjau", -3, "Andika Pratama"},
        {"KEPMENKES-84-2026", "BAB V", "Pencatatan & pelaporan hasil CKG via SATUSEHAT", "Interoperabilitas", "Endpoint CKG", "DataSource.SATUSEHAT", "Resource CKG tersedia di API RS", "EV-CKG-READY", true, "IT/SIMRS", "Sedang Ditinjau", -3, "Andika Pratama"},
        {"KEPMENKES-1128-2022", "Bab VII (MRP)", "Manajemen rekam medis & informasi terstandar", "Digital Maturity", "Kelengkapan resume medis", "DigitalIndicator.EMR", "resume >= 90%", "EV-EMR-ADOPT", true, "Rekam Medis", "Sebagian", -15, "dr. Kartika Dewi"},
        {"KARS", "Standar 1 (APIK)", "Hak pasien & keluarga terpenuhi termasuk keluhan", "Patient Experience", "Keluhan JKN", "ComplaintStat", "kepuasan >= 85%", "EV-UR-BOR", true, "Manajemen Mutu", "Patuh", -20, "Dewi Anggraini"},
        {"KARS", "Standar 12 (MPTK)", "Keselamatan pasien: manajemen risiko & pelaporan insiden", "Patient Safety", "Patient Safety Intelligence", "DigitalIndicator.Keselamatan", "alergi >= 95%; closed-loop >= 80%", "EV-PS-CLM", true, "Komite Mutu", "Sebagian", -10, "dr. Kartika Dewi"},
        {"KARS", "Surveilan", "Re-akreditasi sebelum sertifikat berakhir", "Regulatory Compliance", "Prasyarat P3", "Requirement.P3", "Sertifikat aktif", "EV-REG-KARS", true, "Direktur Medis", "Sebagian", -7, "dr. Kartika Dewi"},
        {"PERMENKES-56-2014", "Pasal 8", "Perizinan tenaga medis: STR & SIP aktif", "Credentialing", "Expiry Monitoring", "Practitioner.sipExpiry", "SIP aktif 100%", "EV-SIP-EXPIRY", true, "HR/SDM", "Sebagian", -10, "Tri Handayani"},
        {"UU-17-2023", "Pasal 193", "Sanksi administratif bagi fasyankes yang tidak memenuhi standar", "Regulatory Compliance", "Compliance Matrix", "-", "Monitoring status kepatuhan", "EV-REG-KARS", true, "Direktur", "Sedang Ditinjau", -25, "dr. Kartika Dewi"},
        {"PP-28-2024", "Pasal 260", "Tenaga kesehatan wajib memiliki izin praktik", "Credentialing", "Verifikasi & Approval", "PractitionerCredential", "Semua praktisi aktif ber-SIP", "EV-SIP-EXPIRY", true, "Komite Medis", "Sebagian", -10, "dr. Bambang Sutrisno"},
}

// SnapRow — [periode, health, bcds, trust, ccs, cpcs, ur, safety]
type SnapRow struct {
        Period string
        Health, Bcds, Trust, Ccs, Cpcs, Ur, Safety float64
}

var SnapshotsNS = []SnapRow{
        {"2026-04", 72, 71, 88, 72, 66, 84, 82},
        {"2026-05", 74, 73, 89, 74, 69, 87, 83},
        {"2026-06", 76, 76, 91, 77, 72, 89, 84},
        {"2026-07", 75, 74, 90, 76, 70, 91, 85},
        {"2026-08", 79, 78, 92, 78, 71, 93, 86},
        {"2026-09", 81, 80, 92, 79, 71, 95, 87},
}

var SnapshotsMS = []SnapRow{
        {"2026-04", 61, 58, 82, 63, 58, 78, 76},
        {"2026-05", 63, 60, 84, 65, 60, 80, 78},
        {"2026-06", 64, 61, 85, 66, 61, 82, 79},
        {"2026-07", 66, 63, 86, 68, 63, 84, 80},
        {"2026-08", 67, 64, 87, 69, 64, 86, 81},
        {"2026-09", 69, 66, 88, 70, 65, 88, 82},
}
