package seed

// Run — orkestrasi seed data demo (port dari prisma/seed.ts main()).
// Menjalankan: bersihkan tabel → RS & pengguna → praktisi & kredensial →
// prasyarat & regulasi → klaim & anomali → pathway → risiko/warning/aksi →
// sarana & digital → snapshot → data ops → tenant kedua.

import (
        "encoding/json"
        "log"
        "math"
        "strconv"
        "strings"
        "time"

        "gorm.io/gorm"

        "medcredix/backend/internal/model"
)

var CATEGORIES = []string{"Diagnosis", "Medication", "Procedure", "Diagnostic", "Stock", "Timeline"}

// Run menjalankan seed bila database masih kosong (idempoten: skip bila sudah ada data).
func Run(db *gorm.DB) error {
        var count int64
        if err := db.Model(&model.Hospital{}).Count(&count).Error; err != nil {
                return err
        }
        if count > 0 {
                log.Println("[seed] database sudah terisi — lewati seeding")
                return nil
        }
        return Force(db)
}

// Force — seed tanpa memeriksa isi database (menghapus data lama terlebih dahulu).
func Force(db *gorm.DB) error {
        log.Println("🧹 Membersihkan database...")
        tables := []string{"notifications", "report_snapshots", "audit_logs", "evidence_items", "trust_factors",
                "data_quality_issues", "integration_jobs", "data_sources", "complaint_stats", "utilization_stats",
                "pathway_episodes", "clinical_pathways", "claim_anomalies", "claims", "action_events", "actions",
                "alert_rules", "early_warnings", "risks", "score_snapshots", "compliance_mappings",
                "regulation_versions", "regulations", "requirements", "credential_documents",
                "practitioner_credentials", "practitioners", "users", "hospitals", "bed_stats", "facilities",
                "specialty_standards", "digital_indicators"}
        for _, t := range tables {
                if err := db.Exec("DELETE FROM " + t).Error; err != nil {
                        // nama tabel mungkin berbeda — coba bentuk lain nanti; log saja
                        log.Printf("[seed] skip %s: %v", t, err)
                }
        }

        pwHash := hashPassword(PW)

        // ============ RUMAH SAKIT & PENGGUNA ============
        log.Println("🏥 Rumah sakit & pengguna...")
        for _, h := range Hospitals {
                hosp := model.Hospital{
                        ID: h.ID, Name: h.Name, Code: h.Code, Class: h.Class, City: h.City, Province: h.Province,
                        BedLicensed: h.BedLicensed, BedActive: h.BedActive, BedNotOperating: h.BedNotOperating,
                        NotOperatingNote: sp(h.NotOperatingNote), BpjsFacilityCode: sp(h.BpjsFacilityCode),
                        SatusehatOrgID: sp(h.SatusehatOrgID), Tier: h.Tier, Onboarding: h.Onboarding,
                }
                if err := db.Create(&hosp).Error; err != nil {
                        return err
                }
        }
        for _, u := range Users {
                lastLogin := D(-1, 8)
                user := model.User{
                        ID: u.ID, HospitalID: u.HospitalID, Email: u.Email, PasswordHash: pwHash,
                        Name: u.Name, Role: u.Role, Title: sp(u.Title), Phone: sp(u.Phone),
                        Active: true, LastLogin: &lastLogin,
                }
                if err := db.Create(&user).Error; err != nil {
                        return err
                }
        }

        // ============ PRAKTISI & KREDENSIAL ============
        log.Println("🩺 Praktisi & kredensial...")
        type PCred struct {
                row model.PractitionerCredential
        }
        var creds []model.PractitionerCredential

        buildPractitioner := func(rows []PractitionerRow, idPrefix, hospitalID string, withOverrides bool) []model.Practitioner {
                out := make([]model.Practitioner, 0, len(rows))
                for i, p := range rows {
                        pid := idPrefix + "_" + pad2(i)
                        strStatus := ""
                        sipStatus := ""
                        if withOverrides {
                                for _, o := range CredentialOverridesNS {
                                        if o.PractitionerIdx == i && o.DocType == "STR" {
                                                strStatus = o.Status
                                        }
                                        if o.PractitionerIdx == i && o.DocType == "SIP" {
                                                sipStatus = o.Status
                                        }
                                }
                        }
                        if strStatus == "" {
                                if p.StrExp <= 90 {
                                        strStatus = "Segera Berakhir"
                                } else {
                                        strStatus = "Aktif"
                                }
                        }
                        if sipStatus == "" {
                                if p.SipExp < 0 {
                                        sipStatus = "Kedaluwarsa"
                                } else if p.SipExp <= 90 {
                                        sipStatus = "Segera Berakhir"
                                } else {
                                        sipStatus = "Aktif"
                                }
                        }
                        var spec, subSpec *string
                        if p.Specialty != "" {
                                spec = sp(p.Specialty)
                        }
                        if p.SubSpecialty != "" {
                                subSpec = sp(p.SubSpecialty)
                        }
                        strNum := "STR-D-20" + itoa(19+(i%6)) + "-" + itoa(70000+i*31)
                        sipNum := "SIP/" + map[bool]string{true: "RSNS", false: "RSMS"}[hospitalID == HNusantara] + "/" + itoa(100+i) + "/2024"
                        if hospitalID == HMedika {
                                strNum = "STR-D-2020-" + itoa(80000+i*17)
                                sipNum = "SIP/RSMS/" + itoa(200+i) + "/2025"
                        }
                        nik := "32********" + pad4(4000+i)
                        if hospitalID == HMedika {
                                nik = "33********" + pad4(2000+i)
                        }
                        phone := "+62 8" + itoa(11+i) + "-5" + itoa(2000+i*7) + "-" + pad4clip(1000+i*13)
                        if hospitalID == HMedika {
                                phone = "+62 8" + itoa(12+i) + "-3" + itoa(3000+i*5) + "-" + pad4clip(2000+i*11)
                        }
                        strExp := D(p.StrExp)
                        sipExp := D(p.SipExp)
                        join := D(-(400 + i*23))
                        if hospitalID == HMedika {
                                join = D(-(300 + i*19))
                        }
                        out = append(out, model.Practitioner{
                                ID: pid, HospitalID: hospitalID, FullName: p.Name, NIKMasked: sp(nik),
                                Profession: p.Profession, Specialty: spec, SubSpecialty: subSpec, Department: p.Department,
                                EmploymentType: p.EmploymentType, Status: p.Status,
                                StrNumber: sp(strNum), StrExpiry: &strExp, SipNumber: sp(sipNum), SipExpiry: &sipExp,
                                Phone: sp(phone), JoinDate: &join, FhirID: sp("practitioner-" + pid),
                        })
                        shortName := strings.Split(p.Name, ",")[0]
                        strIssued := D(p.StrExp - 1825)
                        sipIssued := D(p.SipExp - 730)
                        creds = append(creds,
                                model.PractitionerCredential{
                                        ID: "c_" + idPrefix + pad2(i) + "_str", HospitalID: hospitalID, PractitionerID: pid,
                                        Type: "STR", Name: "Surat Tanda Registrasi — " + shortName, Number: sp(strNum),
                                        IssueDate: &strIssued, ExpiryDate: &strExp, Status: strStatus,
                                },
                                model.PractitionerCredential{
                                        ID: "c_" + idPrefix + pad2(i) + "_sip", HospitalID: hospitalID, PractitionerID: pid,
                                        Type: "SIP", Name: "Surat Izin Praktik — " + shortName, Number: sp(sipNum),
                                        IssueDate: &sipIssued, ExpiryDate: &sipExp, Status: sipStatus,
                                })
                }
                return out
        }

        nsBase := buildPractitioner(PractitionersNS, "pns", HNusantara, true)
        msBase := buildPractitioner(PractitionersMS, "pms", HMedika, false)
        for i := range nsBase {
                if err := db.Create(&nsBase[i]).Error; err != nil {
                        return err
                }
        }
        for i := range msBase {
                if err := db.Create(&msBase[i]).Error; err != nil {
                        return err
                }
        }

        // kredensial tambahan (sertifikat/pelatihan)
        for idx, e := range ExtraCredentialsNS {
                pid := "pns_" + pad2(e.PractitionerIdx)
                issue := D(e.Issue)
                exp := D(e.Exp)
                creds = append(creds, model.PractitionerCredential{
                        ID: "c_pns_extra_" + pad2(idx), HospitalID: HNusantara, PractitionerID: pid,
                        Type: e.Type, Name: e.Name, Number: sp("CERT-" + itoa(2024+(idx%3)) + "-" + itoa(500+idx)),
                        IssueDate: &issue, ExpiryDate: &exp, Status: e.Status,
                })
        }

        // verifikator/persetujuan + dokumen
        var docs []model.CredentialDocument
        for idx := range creds {
                c := &creds[idx]
                isVerified := c.Status == "Aktif" || c.Status == "Segera Berakhir" || c.Status == "Kedaluwarsa"
                if isVerified {
                        verAt := D(-30-(idx%60))
                        appAt := D(-28-(idx%60))
                        c.VerifiedByID = sp("u_komed_ns")
                        c.VerifiedAt = &verAt
                        c.ApprovedByID = sp("u_komed_ns")
                        c.ApprovedAt = &appAt
                }
                ext := "pdf"
                if idx%5 != 0 && idx%3 == 0 {
                        ext = "jpg"
                }
                uploader := "u_admin_ms"
                if c.HospitalID == HNusantara {
                        uploader = "u_cred_ns"
                }
                docStatus := "Terunggah"
                if isVerified {
                        docStatus = "Diverifikasi"
                }
                docs = append(docs, model.CredentialDocument{
                        HospitalID: c.HospitalID, CredentialID: c.ID,
                        FileName: strings.ReplaceAll(c.Type, " ", "-") + "-" + c.PractitionerID + "." + ext,
                        FileType: ext, SizeKb: 180 + (idx*37)%900,
                        UploadedByID: sp(uploader), Status: docStatus, UploadedAt: D(-10-(idx%20), 10),
                })
        }
        for i := range creds {
                if err := db.Create(&creds[i]).Error; err != nil {
                        return err
                }
        }
        for i := range docs {
                if err := db.Create(&docs[i]).Error; err != nil {
                        return err
                }
        }

        // ============ PRASYARAT & REGULASI ============
        log.Println("📋 Prasyarat & regulasi...")
        for _, r := range Requirements {
                issued := D(r.IssuedAt)
                row := model.Requirement{
                        HospitalID: HNusantara, Code: r.Code, Name: r.Name, Category: r.Category,
                        RegulatorySource: r.RegulatorySource, Version: r.Version, Required: r.Required,
                        DocType: r.DocType, ValidationRule: sp(r.ValidationRule), EvidenceRequired: r.EvidenceRequired,
                        ApprovalRequired: r.ApprovalRequired, ResponsibleRole: r.ResponsibleRole, Status: r.Status,
                        DocNumber: sp(r.DocNumber), IssuedAt: &issued,
                        FileName: sp(r.FileName), Notes: sp(r.Notes),
                }
                if r.HasExpiry {
                        exp := D(r.DocExpiry)
                        row.DocExpiry = &exp
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        regID := map[string]string{}
        for _, r := range Regulations {
                row := model.Regulation{
                        Code: r.Code, Title: r.Title, Number: r.Number, Authority: r.Authority,
                        Year: r.Year, Category: r.Category, Status: r.Status, Summary: r.Summary, Source: r.Source,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
                regID[r.Code] = row.ID
        }
        for _, v := range RegulationVersions {
                issued, _ := time.ParseInLocation("2006-01-02", v.IssuedAt, WIB)
                row := model.RegulationVersion{
                        RegulationID: regID[v.Reg], Version: v.Version, IssuedAt: issued, Status: v.Status,
                        ChangeSummary: sp(v.ChangeSummary), ImpactLevel: v.ImpactLevel, SupersededBy: v.SupersededBy,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for _, m := range ComplianceMappings {
                review := D(m.LastReview)
                row := model.ComplianceMapping{
                        RegulationID: regID[m.Reg], HospitalID: HNusantara, Article: m.Article, Requirement: m.Requirement,
                        Module: m.Module, Feature: m.Feature, DataField: sp(m.DataField), Validation: sp(m.Validation),
                        EvidenceRef: sp(m.EvidenceRef), AuditRequired: m.AuditRequired, Owner: m.Owner, Status: m.Status,
                        LastReviewedAt: &review, Reviewer: sp(m.Reviewer),
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }

        // ============ KLAIM & ANOMALI ============
        log.Println("💰 Klaim & anomali...")
        claimIDByNumber := map[string]string{}
        seedClaims := func(rows []ClaimRow, hospitalID string, anomalyStatus func(c ClaimRow) string, reviewedBy func(c ClaimRow) *string) error {
                for _, c := range rows {
                        flags := map[string]string{
                                "Diagnosis": c.Dx, "Medication": c.Med, "Procedure": c.Proc,
                                "Diagnostic": c.Diag, "Stock": c.Stock, "Timeline": c.Tl,
                        }
                        flagsJson, _ := json.Marshal(flags)
                        fj := string(flagsJson)
                        adm := D(c.In)
                        dis := D(c.Out)
                        claim := model.Claim{
                                HospitalID: hospitalID, ClaimNumber: c.No, SepNumber: c.Sep, PatientInit: c.Init,
                                InacbgCode: c.InacbgCode, InacbgDesc: c.Desc, CareType: c.CareType,
                                AdmissionDate: adm, DischargeDate: &dis, BilledAmount: c.Amount, Status: c.Status, FlagsJson: &fj,
                        }
                        if err := db.Create(&claim).Error; err != nil {
                                return err
                        }
                        claimIDByNumber[c.No] = claim.ID
                        for _, cat := range CATEGORIES {
                                flag := flags[cat]
                                if flag == "ok" {
                                        continue
                                }
                                t := AnomalyTemplates[cat]
                                evJson, _ := json.Marshal(map[string][]string{
                                        "source": {"BPJS VClaim", "EMR/RME", "Farmasi"},
                                        "fhir":   {"Claim", "MedicationRequest", "Procedure", "DiagnosticReport", "Encounter"},
                                })
                                evj := string(evJson)
                                detected := D(c.Out + 1)
                                anomaly := model.ClaimAnomaly{
                                        ClaimID: claim.ID, HospitalID: hospitalID, Category: cat, Severity: t.Sev[flag],
                                        Title: t.Title, Description: t.Desc(c), EvidenceJson: &evj, Recommendation: t.Rec,
                                        Status: anomalyStatus(c), DetectedAt: detected, ReviewedBy: reviewedBy(c),
                                }
                                if err := db.Create(&anomaly).Error; err != nil {
                                        return err
                                }
                        }
                }
                return nil
        }
        nsAnomalyStatus := func(c ClaimRow) string {
                if c.Out < -25 {
                        return "Baru"
                }
                if strings.HasSuffix(c.No, "102") || strings.HasSuffix(c.No, "103") || strings.HasSuffix(c.No, "123") {
                        return "Diklarifikasi"
                }
                return "Ditinjau"
        }
        nsReviewedBy := func(c ClaimRow) *string {
                if strings.HasSuffix(c.No, "102") {
                        return sp("Silvia Rahmawati, AMd.Kep")
                }
                return nil
        }
        if err := seedClaims(ClaimsNS, HNusantara, nsAnomalyStatus, nsReviewedBy); err != nil {
                return err
        }

        // ============ CLINICAL PATHWAY ============
        log.Println("🧬 Clinical pathway & episode...")
        for _, pw := range PathwaysNS {
                n := len(pw.Episodes)
                sum := 0.0
                costSum := 0.0
                for _, e := range pw.Episodes {
                        sum += e.Adherence
                        costSum += e.CostRatio
                }
                avg := sum / float64(n)
                devCount := 0
                docCount := 0
                for _, e := range pw.Episodes {
                        if e.DeviationType != "Tidak Ada" {
                                devCount++
                                if e.Documented {
                                        docCount++
                                }
                        }
                }
                documented := 100.0
                if devCount > 0 {
                        documented = math.Round(float64(docCount) / float64(devCount) * 100)
                }
                row := model.ClinicalPathway{
                        HospitalID: HNusantara, Diagnosis: pw.Diagnosis, Icd10: pw.Icd10,
                        StandardLos: pw.StandardLos, TargetAdherence: pw.TargetAdherence,
                        EpisodeCount: n, AvgAdherence: round1(avg),
                        MedicationAdherence: round1(avg + 4), LabAdherence: round1(avg + 6),
                        ProcedureAdherence: round1(avg + 2), LosAdherence: round1(avg - 2),
                        DeviationCount: devCount, DocumentedDeviation: documented,
                        CostRatio: math.Round(costSum/float64(n)*100) / 100, Status: "Aktif",
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
                for _, e := range pw.Episodes {
                        ep := model.PathwayEpisode{
                                PathwayID: row.ID, HospitalID: HNusantara, EpisodeNumber: e.No, PatientInit: e.Init,
                                Dpjp: e.Dpjp, Adherence: e.Adherence, LosActual: e.Los, LosStandard: pw.StandardLos,
                                DeviationType: e.DeviationType, DeviationReason: e.Reason, Documented: e.Documented,
                                CostRatio: e.CostRatio, Period: "2026-09",
                        }
                        if err := db.Create(&ep).Error; err != nil {
                                return err
                        }
                }
        }
        for _, u := range UtilizationNS {
                row := model.UtilizationStat{
                        HospitalID: HNusantara, Month: u.Month, MonthLabel: u.Label,
                        Alos: u.Alos, Readmission: u.Readm, Bor: u.Bor, Score: 0,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for _, c := range ComplaintsNS {
                row := model.ComplaintStat{
                        HospitalID: HNusantara, Month: c.Month, Total: c.Total, Antrian: c.Antrian,
                        Pelayanan: c.Pelayanan, Obat: c.Obat, Administrasi: c.Administrasi, Fasilitas: c.Fasilitas,
                        Lainnya: c.Lainnya, Satisfaction: c.Satisfaction, Per1000: c.Per1000, PeerAvg: c.PeerAvg,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }

        // ============ RISIKO, WARNING, AKSI ============
        log.Println("🚨 Risiko & peringatan dini...")
        for _, r := range RisksNS {
                detected := D(r.DetectedAt)
                deadline := D(r.Deadline)
                row := model.Risk{
                        HospitalID: HNusantara, Code: r.Code, Category: r.Category, Level: r.Level,
                        Title: r.Title, Description: r.Description, DetectedAt: detected, Domain: r.Domain,
                        AffectedRecords: sp(r.AffectedRecords), EvidenceJson: sp(r.EvidenceJson),
                        BusinessImpact: r.BusinessImpact, RecommendedAction: r.RecommendedAction,
                        SuggestedOwner: r.SuggestedOwner, Deadline: &deadline, Status: r.Status,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for _, w := range WarningsNS {
                detected := D(w.DetectedAt)
                row := model.EarlyWarning{
                        HospitalID: HNusantara, Code: w.Code, Title: w.Title, Description: w.Description,
                        Severity: w.Severity, Domain: w.Domain, ThresholdRule: w.ThresholdRule,
                        TriggerValue: w.TriggerValue, DetectedAt: detected, Status: w.Status,
                        AcknowledgedBy: w.AckBy, AssignedTo: w.Assigned,
                }
                if w.AckAt != nil {
                        t := D(*w.AckAt)
                        row.AcknowledgedAt = &t
                }
                if w.Due != nil {
                        t := D(*w.Due)
                        row.DueDate = &t
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for _, a := range AlertRules {
                row := model.AlertRule{
                        HospitalID: HNusantara, Name: a.Name, Domain: a.Domain, Metric: a.Metric,
                        Condition: a.Condition, Severity: a.Severity, Channels: a.Channels,
                        Active: true, TriggeredCount: a.TriggeredCount,
                }
                if a.LastTriggered != nil {
                        t := D(*a.LastTriggered)
                        row.LastTriggeredAt = &t
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }

        log.Println("🎯 Action center...")
        for _, a := range ActionsNS {
                due := D(a.Due)
                row := model.Action{
                        HospitalID: HNusantara, Code: a.Code, Title: a.Title, Description: a.Description,
                        Priority: a.Priority, Status: a.Status, SourceType: a.SourceType, SourceRef: a.SourceRef,
                        Owner: a.Owner, OwnerRole: a.OwnerRole, DueDate: &due,
                        ExpectedOutcome: sp(a.Expected), ProgressNote: a.Progress,
                        ResolutionNote: a.Resolved,
                }
                if a.ResolvedAt != nil {
                        t := D(*a.ResolvedAt)
                        row.ResolvedAt = &t
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
                srcRef := ""
                if a.SourceRef != nil {
                        srcRef = *a.SourceRef
                }
                events := []model.ActionEvent{
                        {ActionID: row.ID, Type: "created", ToValue: sp("Baru"), ByName: "MedCredix Engine",
                                Note: sp("Aksi dibuat dari " + a.SourceType + srcRefSuffix(srcRef)), CreatedAt: D(-12)},
                        {ActionID: row.ID, Type: "assigned", ToValue: sp(a.Owner), ByName: "MedCredix Engine",
                                Note: sp("Ditugaskan ke " + a.Owner + " (" + a.OwnerRole + ")"), CreatedAt: D(-11)},
                }
                if a.Status == "Ditugaskan" || a.Status == "Ditinjau" {
                        events = append(events, model.ActionEvent{ActionID: row.ID, Type: "status", FromValue: sp("Baru"),
                                ToValue: sp(a.Status), ByName: a.Owner, CreatedAt: D(-6)})
                }
                if a.Status == "Dalam Pengerjaan" {
                        events = append(events, model.ActionEvent{ActionID: row.ID, Type: "status", FromValue: sp("Ditugaskan"),
                                ToValue: sp("Dalam Pengerjaan"), ByName: a.Owner, CreatedAt: D(-4)})
                }
                if a.Progress != nil {
                        events = append(events, model.ActionEvent{ActionID: row.ID, Type: "comment",
                                ByName: a.Owner, Note: a.Progress, CreatedAt: D(-2)})
                }
                if a.Status == "Selesai" || a.Status == "Ditutup" {
                        resolvedAt := D(-5)
                        if a.ResolvedAt != nil {
                                resolvedAt = D(*a.ResolvedAt)
                        }
                        events = append(events,
                                model.ActionEvent{ActionID: row.ID, Type: "status", FromValue: sp("Dalam Pengerjaan"),
                                        ToValue: sp("Menunggu Verifikasi"), ByName: a.Owner, CreatedAt: D(-6)},
                                model.ActionEvent{ActionID: row.ID, Type: "resolved", FromValue: sp("Menunggu Verifikasi"),
                                        ToValue: sp(a.Status), ByName: "dr. Kartika Dewi", Note: a.Resolved, CreatedAt: resolvedAt})
                }
                for i := range events {
                        if err := db.Create(&events[i]).Error; err != nil {
                                return err
                        }
                }
        }

        // ============ SARANA, RANJANG, DIGITAL ============
        log.Println("🏥 Sarana, ranjang, indikator digital...")
        for _, f := range FacilitiesNS {
                row := model.Facility{
                        HospitalID: HNusantara, Name: f.Name, Category: f.Category, Availability: f.Availability,
                        Functional: f.Functional, Feasible: f.Feasible, Calibration: f.Calibration,
                        Suitable: f.Suitable, StandardRef: f.Standard, Notes: f.Notes,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for _, b := range BedsNS {
                row := model.BedStat{
                        HospitalID: HNusantara, Ward: b.Ward, BedClass: b.BedClass,
                        Total: b.Total, Active: b.Active, Usable: b.Usable, NonOperatingReason: b.Reason,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for _, s := range SpecialtyStandards {
                row := model.SpecialtyStandard{
                        HospitalID: HNusantara, RsClass: "B", Specialty: s.Specialty, Required: s.Required, Note: s.Note,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for _, i := range DigitalNS {
                sync := D(0, 8)
                row := model.DigitalIndicator{
                        HospitalID: HNusantara, Dimension: i.Dimension, Indicator: i.Indicator, Value: i.Value,
                        Unit: i.Unit, Threshold: i.Threshold, Weight: i.Weight, Status: i.Status,
                        FhirResource: i.Fhir, SourceSystem: i.Source, LastSyncAt: sync, Direction: i.Direction,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }

        // ============ SNAPSHOT SKOR ============
        log.Println("📊 Snapshot skor historis...")
        snapSeed := func(rows []SnapRow, hospitalID string) error {
                for _, s := range rows {
                        dataAt := time.Date(mustYear(s.Period), time.Month(mustMonth(s.Period)), 28, 4, 0, 0, 0, WIB)
                        calcAt := dataAt.Add(30 * time.Minute)
                        row := model.ScoreSnapshot{
                                HospitalID: hospitalID, Period: s.Period, HealthScore: s.Health, Bcds: s.Bcds,
                                Trust: s.Trust, Ccs: s.Ccs, Cpcs: s.Cpcs, UrScore: s.Ur, SafetyScore: s.Safety,
                                FormulaVersion: "medcredix-scorer/1.0", CalcAt: calcAt, DataAt: dataAt,
                        }
                        if err := db.Create(&row).Error; err != nil {
                                return err
                        }
                }
                return nil
        }
        if err := snapSeed(SnapshotsNS, HNusantara); err != nil {
                return err
        }
        if err := snapSeed(SnapshotsMS, HMedika); err != nil {
                return err
        }

        // ============ DATA OPS ============
        log.Println("🔌 Sumber data, integrasi, kualitas data...")
        dsID := map[string]string{}
        seedSources := func(rows []DataSourceSeed, hospitalID string) error {
                for _, s := range rows {
                        sync := D(s.LastSyncDay, s.LastSyncHour)
                        row := model.DataSource{
                                HospitalID: hospitalID, Name: s.Name, Type: s.Type, Status: s.Status,
                                AuthMethod: s.AuthMethod, ConfigSummary: sp(s.ConfigSummary),
                                SuccessRate: s.SuccessRate, AvgLatencyMs: s.AvgLatencyMs, LastSyncAt: &sync,
                                FhirResources: s.FhirResources, CoveragePct: s.CoveragePct, Uptime30d: s.Uptime30d, Notes: s.Notes,
                        }
                        if err := db.Create(&row).Error; err != nil {
                                return err
                        }
                        dsID[s.Name] = row.ID
                }
                return nil
        }
        if err := seedSources(DataSourcesNS, HNusantara); err != nil {
                return err
        }
        if err := seedSources(DataSourcesMS, HMedika); err != nil {
                return err
        }
        jobTypeMap := map[string]string{
                "SIMRS": "SIMRS Inti Sehat (Core HIS)", "EMR/RME": "EMR/RME Terpadu",
                "Farmasi": "Sistem Farmasi", "SATUSEHAT": "SATUSEHAT Platform",
                "LIS": "LIS (Laboratory Information System)", "RIS/PACS": "RIS/PACS (Radiologi)",
                "BPJS VClaim": "BPJS VClaim (e-Klaim)", "BPJS P-Care": "BPJS P-Care / JKN Mobile",
                "HRIS": "HRIS Nusantara", "Keuangan": "Sistem Keuangan & Kasir",
        }
        for i, j := range IntegrationJobs {
                hour := []int{2, 8, 14, 20}[i%4]
                started := D(j.StartDay, hour)
                finished := started.Add(time.Duration(j.DurationMs) * time.Millisecond)
                dur := j.DurationMs
                row := model.IntegrationJob{
                        HospitalID: HNusantara, DataSourceID: dsID[jobTypeMap[j.Source]], JobType: j.JobType,
                        StartedAt: started, FinishedAt: &finished, Status: j.Status,
                        RecordsPulled: j.Records, DurationMs: &dur, Message: j.Message,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for _, q := range DqIssuesNS {
                detected := D(q.DetectedAt)
                row := model.DataQualityIssue{
                        HospitalID: HNusantara, Dimension: q.Dimension, Title: q.Title, Description: q.Description,
                        AffectedRecords: q.AffectedRecords, SourceSystem: q.SourceSystem, Severity: q.Severity,
                        Anomaly: q.Anomaly, Status: q.Status, DetectedAt: detected,
                }
                if q.ResolvedAt != nil {
                        t := D(*q.ResolvedAt)
                        row.ResolvedAt = &t
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for _, t := range TrustFactorsNS {
                row := model.TrustFactor{
                        HospitalID: HNusantara, Factor: t.Factor, Score: t.Score, Weight: t.Weight, Notes: sp(t.Notes),
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }

        log.Println("🧾 Evidence & audit...")
        for _, e := range EvidenceNS {
                sync := D(0, 8)
                row := model.EvidenceItem{
                        HospitalID: HNusantara, RefType: e.RefType, RefID: e.RefID, Title: e.Title,
                        Metric: sp(e.Metric), Value: sp(e.Value), Unit: sp(e.Unit), SourceSystem: e.Source,
                        FhirResource: e.Fhir, RecordCount: ip(e.RecordCount), LastSyncAt: &sync,
                        ValidationStatus: e.Validation, Timestamp: sync, DemoData: true,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for i, a := range AuditLogs {
                hospID := sp(HNusantara)
                if i%9 == 0 {
                        hospID = nil
                }
                ipAddr := "10.20." + itoa(i%5) + "." + itoa(10+i)
                row := model.AuditLog{
                        HospitalID: hospID, ActorName: a.Actor, ActorRole: a.Role, Action: a.Action,
                        EntityType: a.Entity, EntityID: sp(a.EntityID), Description: a.Description,
                        Before: a.Before, After: a.After, IP: sp(ipAddr),
                        CorrelationID: sp("corr-" + pad3(i)),
                        CreatedAt:     D(-i/2, 7+(i%10)),
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }
        for _, n := range Notifications {
                created := D(-1, 8)
                row := model.Notification{
                        HospitalID: HNusantara, Title: n.Title, Body: n.Body, Severity: n.Severity,
                        Channel: n.Channel, Link: sp(n.Link), CreatedAt: created,
                }
                if err := db.Create(&row).Error; err != nil {
                        return err
                }
        }

        // ============ RS MEDIKA (tenant kedua — ringkas) ============
        log.Println("🏥 Data ringkas RS Medika Sejahtera...")
        msAnomalyStatus := func(c ClaimRow) string { return "Ditinjau" }
        msReviewedBy := func(c ClaimRow) *string { return nil }
        if err := seedClaims(MedikaClaims, HMedika, msAnomalyStatus, msReviewedBy); err != nil {
                return err
        }

        // ringkasan jumlah baris
        var counts struct {
                Practitioners, Credentials, Documents, Claims, Anomalies int64
                Pathways, Episodes, Risks, Warnings, Actions int64
        }
        db.Model(&model.Practitioner{}).Count(&counts.Practitioners)
        db.Model(&model.PractitionerCredential{}).Count(&counts.Credentials)
        db.Model(&model.CredentialDocument{}).Count(&counts.Documents)
        db.Model(&model.Claim{}).Count(&counts.Claims)
        db.Model(&model.ClaimAnomaly{}).Count(&counts.Anomalies)
        db.Model(&model.ClinicalPathway{}).Count(&counts.Pathways)
        db.Model(&model.PathwayEpisode{}).Count(&counts.Episodes)
        db.Model(&model.Risk{}).Count(&counts.Risks)
        db.Model(&model.EarlyWarning{}).Count(&counts.Warnings)
        db.Model(&model.Action{}).Count(&counts.Actions)
        log.Printf("✅ Seed selesai: praktisi=%d kredensial=%d dokumen=%d klaim=%d anomali=%d pathway=%d episode=%d risiko=%d warning=%d aksi=%d",
                counts.Practitioners, counts.Credentials, counts.Documents, counts.Claims, counts.Anomalies,
                counts.Pathways, counts.Episodes, counts.Risks, counts.Warnings, counts.Actions)
        return nil
}

func hashPassword(pw string) string {
        sum := sha256Hex("medcredix::" + pw)
        return sum
}

func srcRefSuffix(ref string) string {
        if ref == "" {
                return ""
        }
        return " (" + ref + ")"
}

func pad4(n int) string {
        s := itoa(n)
        for len(s) < 4 {
                s = "0" + s
        }
        // slice(-4): ambil 4 karakter terakhir
        if len(s) > 4 {
                s = s[len(s)-4:]
        }
        return s
}

func pad4clip(n int) string {
        s := itoa(n)
        if len(s) > 4 {
                s = s[:4]
        }
        return s
}

func pad3(n int) string {
        s := itoa(n)
        for len(s) < 3 {
                s = "0" + s
        }
        return s
}

func round1(v float64) float64 { return math.Round(v*10) / 10 }

func mustYear(period string) int  { y, _ := strconv.Atoi(period[:4]); return y }
func mustMonth(period string) int { m, _ := strconv.Atoi(period[5:7]); return m }
