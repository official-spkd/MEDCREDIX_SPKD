package api

// Reports — generate snapshot, daftar, detail, ekspor CSV.
// Port dari src/app/api/reports/route.ts.
import (
        "encoding/json"
        "fmt"
        "sort"
        "strings"
        "time"

        "github.com/gin-gonic/gin"

        "medcredix/backend/internal/auth"
        "medcredix/backend/internal/db"
        "medcredix/backend/internal/model"
        "medcredix/backend/internal/scoring"
)

type reportTypeMeta struct {
        Type, Title, Desc string
}

var reportTypes = map[string]reportTypeMeta{
        "kredensial":    {"kredensial", "Laporan Kredensial RS", "Skor BDCS, kelompok A–D, prasyarat, kredensial & tenaga medis"},
        "kesehatan-rs":  {"kesehatan-rs", "Hospital Health Report", "9 domain, tren 6 bulan, rekomendasi perbaikan"},
        "risiko":        {"risiko", "Risk Report", "Risiko aktif, level, owner, dan tenggat"},
        "aksi":          {"aksi", "Action & Improvement Report", "Aksi perbaikan, status, resolusi, dampak skor"},
        "kepatuhan":     {"kepatuhan", "Regulatory Compliance Report", "Matriks kepatuhan, regulasi, versi, temuan"},
        "klaim":         {"klaim", "BPJS Claim Integrity Report", "CCS, sub-skor, anomali aktif, prioritas tinjauan"},
}

func handleReportsList(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        if exportID := c.Query("export"); exportID != "" {
                if err := auth.RequirePermission(user, "EXPORT"); err != nil {
                        return auth.Fail(err)
                }
                var snap model.ReportSnapshot
                if err := db.DB.Where("id = ? AND hospital_id = ?", exportID, hospitalID).First(&snap).Error; err != nil {
                        return 404, "Laporan tidak ditemukan."
                }
                var content map[string]any
                if err := json.Unmarshal([]byte(snap.ContentJson), &content); err != nil {
                        return 500, "Konten laporan tidak valid."
                }
                rows := [][]string{{"Field", "Nilai"}}
                var flat func(obj map[string]any, prefix string)
                flat = func(obj map[string]any, prefix string) {
                        keys := make([]string, 0, len(obj))
                        for k := range obj {
                                keys = append(keys, k)
                        }
                        sort.Strings(keys)
                        for _, k := range keys {
                                v := obj[k]
                                switch val := v.(type) {
                                case map[string]any:
                                        flat(val, prefix+k+".")
                                case []any:
                                        jb, _ := json.Marshal(val)
                                        rows = append(rows, []string{prefix + k, string(jb)})
                                default:
                                        s := "—"
                                        if val != nil {
                                                s = fmt.Sprintf("%v", val)
                                        }
                                        rows = append(rows, []string{prefix + k, s})
                                }
                        }
                }
                flat(content, "")
                var b strings.Builder
                for _, r := range rows {
                        cells := make([]string, len(r))
                        for i, cell := range r {
                                cells[i] = `"` + strings.ReplaceAll(cell, `"`, `""`) + `"`
                        }
                        b.WriteString(strings.Join(cells, ","))
                        b.WriteString("\n")
                }
                hid := hospitalID
                auth.LogAudit(auth.AuditOptions{
                        HospitalID: &hid, Actor: user, Action: "export", EntityType: "Laporan", EntityID: &snap.ID,
                        Description: "Ekspor CSV: " + snap.Title, C: c,
                })
                c.Header("Content-Type", "text/csv; charset=utf-8")
                c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="MedCredix-%s-%s.csv"`, snap.Type, snap.Period))
                c.String(200, "\uFEFF"+b.String())
                return -1, nil // respons sudah ditulis langsung
        }
        if id := c.Query("id"); id != "" {
                var snap model.ReportSnapshot
                if err := db.DB.Where("id = ? AND hospital_id = ?", id, hospitalID).First(&snap).Error; err != nil {
                        return 404, "Laporan tidak ditemukan."
                }
                var content any
                json.Unmarshal([]byte(snap.ContentJson), &content)
                return auth.OK(gin.H{"report": snap, "content": content})
        }
        var snaps []model.ReportSnapshot
        db.DB.Where("hospital_id = ?", hospitalID).Order("generated_at desc").Find(&snaps)
        types := []gin.H{}
        for _, k := range []string{"kredensial", "kesehatan-rs", "risiko", "aksi", "kepatuhan", "klaim"} {
                t := reportTypes[k]
                types = append(types, gin.H{"type": t.Type, "title": t.Title, "desc": t.Desc})
        }
        return auth.OK(gin.H{"reports": snaps, "types": types})
}

type reportGenBody struct {
        Type string `json:"type"`
}

func handleReportGenerate(c *gin.Context, user *auth.SessionUser) (int, any) {
        var body reportGenBody
        if err := c.ShouldBindJSON(&body); err != nil {
                return 400, "Payload tidak valid."
        }
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        if err := auth.RequirePermission(user, "CREATE"); err != nil {
                return auth.Fail(err)
        }
        typ := body.Type
        if typ == "" {
                typ = "kesehatan-rs"
        }
        meta, ok := reportTypes[typ]
        if !ok {
                return 400, "Jenis laporan tidak dikenal."
        }
        intel, err := ComputeIntel(hospitalID)
        if err != nil {
                return auth.Fail(err)
        }
        var hospital model.Hospital
        db.DB.Where("id = ?", hospitalID).First(&hospital)
        period := "—"
        if len(intel.Trend) > 0 {
                period = intel.Trend[0].Period + "–" + intel.Trend[len(intel.Trend)-1].Period
        }

        content := map[string]any{
                "ringkasan": gin.H{
                        "namaRS": hospital.Name, "kelas": hospital.Class, "kota": hospital.City, "periode": period,
                        "healthScore": intel.HealthScore, "bcds": intel.Bcds.Final, "stage": intel.Bcds.Stage,
                        "trust": intel.Trust.Score, "ccs": intel.Ccs.Score, "cpcs": intel.Cpcs.Score,
                },
        }
        switch typ {
        case "kredensial":
                var practitioners []model.Practitioner
                db.DB.Where("hospital_id = ?", hospitalID).Find(&practitioners)
                var creds []model.PractitionerCredential
                db.DB.Where("hospital_id = ?", hospitalID).Find(&creds)
                var reqs []model.Requirement
                db.DB.Where("hospital_id = ?", hospitalID).Find(&reqs)
                statusCount := gin.H{}
                for _, cr := range creds {
                        if v, ok := statusCount[cr.Status].(int); ok {
                                statusCount[cr.Status] = v + 1
                        } else {
                                statusCount[cr.Status] = 1
                        }
                }
                prereqOut := []gin.H{}
                for _, r := range reqs {
                        exp := "—"
                        if r.DocExpiry != nil {
                                exp = r.DocExpiry.Format("2006-01-02")
                        }
                        prereqOut = append(prereqOut, gin.H{"kode": r.Code, "nama": r.Name, "status": r.Status, "kedaluwarsa": exp})
                }
                content["tenagaMedis"] = gin.H{"total": len(practitioners)}
                content["kredensial"] = statusCount
                content["prasyarat"] = prereqOut
                content["kelompok"] = intel.Kelompok
                content["bdcS"] = intel.Bcds
                content["workforce"] = intel.Workforce
        case "risiko":
                var risks []model.Risk
                db.DB.Where("hospital_id = ?", hospitalID).Find(&risks)
                out := []gin.H{}
                for _, r := range risks {
                        tenggat := "—"
                        if r.Deadline != nil {
                                tenggat = r.Deadline.Format("2006-01-02")
                        }
                        out = append(out, gin.H{"kode": r.Code, "kategori": r.Category, "level": r.Level,
                                "judul": r.Title, "owner": r.SuggestedOwner, "tenggat": tenggat, "status": r.Status})
                }
                content["risiko"] = out
        case "aksi":
                var actions []model.Action
                db.DB.Where("hospital_id = ?", hospitalID).Find(&actions)
                out := []gin.H{}
                for _, a := range actions {
                        tenggat := "—"
                        if a.DueDate != nil {
                                tenggat = a.DueDate.Format("2006-01-02")
                        }
                        resolusi := "—"
                        if a.ResolutionNote != nil {
                                resolusi = *a.ResolutionNote
                        }
                        out = append(out, gin.H{"kode": a.Code, "judul": a.Title, "prioritas": a.Priority,
                                "status": a.Status, "owner": a.Owner, "tenggat": tenggat, "resolusi": resolusi})
                }
                content["aksi"] = out
                content["catatan"] = "Closed-loop: resolusi aksi memicu usulan recalculation skor."
        case "kepatuhan":
                var mappings []model.ComplianceMapping
                db.DB.Where("hospital_id = ?", hospitalID).Find(&mappings)
                regs := map[string]model.Regulation{}
                var regRows []model.Regulation
                db.DB.Find(&regRows)
                for _, r := range regRows {
                        regs[r.ID] = r
                }
                out := []gin.H{}
                for _, m := range mappings {
                        number := "—"
                        if r, ok := regs[m.RegulationID]; ok {
                                number = r.Number
                        }
                        review := "—"
                        if m.LastReviewedAt != nil {
                                review = m.LastReviewedAt.Format("2006-01-02")
                        }
                        out = append(out, gin.H{"regulasi": number, "pasal": m.Article, "requirement": m.Requirement,
                                "modul": m.Module, "status": m.Status, "owner": m.Owner, "review": review})
                }
                content["matriks"] = out
                content["disclaimer"] = "Designed to support compliance, subject to regulatory and hospital verification."
        case "klaim":
                var claims []model.Claim
                db.DB.Where("hospital_id = ?", hospitalID).Find(&claims)
                var anomalies []model.ClaimAnomaly
                db.DB.Where("hospital_id = ?", hospitalID).Find(&anomalies)
                ditunda, ditolak := 0, 0
                for _, cl := range claims {
                        if cl.Status == "Ditunda" {
                                ditunda++
                        }
                        if cl.Status == "Ditolak" {
                                ditolak++
                        }
                }
                anomaliTinggi, perluTinjauan := 0, 0
                for _, a := range anomalies {
                        if a.Severity == "Tinggi" {
                                anomaliTinggi++
                        }
                        if a.Status == "Ditinjau" {
                                perluTinjauan++
                        }
                }
                content["ccs"] = intel.Ccs
                content["klaim"] = gin.H{"total": len(claims), "ditunda": ditunda, "ditolak": ditolak}
                content["anomali"] = gin.H{"total": len(anomalies), "tinggi": anomaliTinggi, "perluTinjauan": perluTinjauan}
                content["disclaimer"] = "Sinyal bersifat netral dan memerlukan verifikasi manual — bukan kesimpulan penyalahgunaan."
        default:
                domains := []gin.H{}
                for _, d := range intel.Domains {
                        domains = append(domains, gin.H{"domain": d.Label, "skor": d.Score, "bobot": d.Weight})
                }
                sorted := append([]scoring.Domain{}, intel.Domains...)
                sort.Slice(sorted, func(i, j int) bool { return sorted[i].Score < sorted[j].Score })
                saran := []gin.H{}
                for i, d := range sorted {
                        if i >= 3 {
                                break
                        }
                        catatan := "Pertahankan kinerja"
                        if len(d.Negatives) > 0 {
                                catatan = d.Negatives[0]
                        }
                        saran = append(saran, gin.H{"domain": d.Label, "skor": d.Score, "catatan": catatan})
                }
                content["domain"] = domains
                content["tren"] = intel.Trend
                content["saran"] = saran
        }

        contentJson, _ := json.Marshal(content)
        var count int64
        db.DB.Model(&model.ReportSnapshot{}).Where("hospital_id = ?", hospitalID).Count(&count)
        snap := model.ReportSnapshot{
                HospitalID: hospitalID, Type: typ, Title: meta.Title, Period: period,
                GeneratedAt: time.Now(), GeneratedBy: user.Name,
                Methodology: "Scoring engine " + scoring.FormulaVersion + "; sumber: SIMRS/EMR/BPJS/SATUSEHAT via integration layer; verifikasi manual wajib untuk keputusan final.",
                Version:     fmt.Sprintf("v1.%d", count+1),
                AuditRef:    fmt.Sprintf("AUD-%X", time.Now().UnixMilli()),
                ContentJson: string(contentJson),
        }
        if err := db.DB.Create(&snap).Error; err != nil {
                return auth.Fail(err)
        }
        hid := hospitalID
        auth.LogAudit(auth.AuditOptions{
                HospitalID: &hid, Actor: user, Action: "create", EntityType: "Laporan", EntityID: &snap.ID,
                Description: "Generate laporan: " + snap.Title, C: c,
        })
        return auth.OK(snap)
}
