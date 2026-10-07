package api

// Compliance, Evidence, Audit, Search — GET.
// Port dari compliance/route.ts, evidence/route.ts, audit/route.ts, search/route.ts.
import (
        "sort"

        "github.com/gin-gonic/gin"

        "medcredix/backend/internal/auth"
        "medcredix/backend/internal/db"
        "medcredix/backend/internal/model"
)

func handleCompliance(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        var regs []model.Regulation
        db.DB.Order("year desc").Find(&regs)
        var versions []model.RegulationVersion
        db.DB.Find(&versions)
        var mappings []model.ComplianceMapping
        db.DB.Where("hospital_id = ?", hospitalID).Order("status asc").Find(&mappings)
        regMap := map[string]model.Regulation{}
        for _, r := range regs {
                regMap[r.ID] = r
        }
        summary := gin.H{"patuh": 0, "sebagian": 0, "belum": 0, "ditinjau": 0}
        for _, m := range mappings {
                switch m.Status {
                case "Patuh":
                        summary["patuh"] = summary["patuh"].(int) + 1
                case "Sebagian":
                        summary["sebagian"] = summary["sebagian"].(int) + 1
                case "Belum":
                        summary["belum"] = summary["belum"].(int) + 1
                case "Sedang Ditinjau":
                        summary["ditinjau"] = summary["ditinjau"].(int) + 1
                }
        }
        type regOut struct {
                model.Regulation
                Versions []model.RegulationVersion `json:"versions"`
        }
        regsOut := make([]regOut, 0, len(regs))
        for _, r := range regs {
                vs := []model.RegulationVersion{}
                for _, v := range versions {
                        if v.RegulationID == r.ID {
                                vs = append(vs, v)
                        }
                }
                regsOut = append(regsOut, regOut{Regulation: r, Versions: vs})
        }
        type mappingOut struct {
                model.ComplianceMapping
                Regulation *struct {
                        Code     string `json:"code"`
                        Title    string `json:"title"`
                        Number   string `json:"number"`
                        Category string `json:"category"`
                } `json:"regulation"`
        }
        mapsOut := make([]mappingOut, 0, len(mappings))
        for _, m := range mappings {
                row := mappingOut{ComplianceMapping: m}
                if r, ok := regMap[m.RegulationID]; ok {
                        row.Regulation = &struct {
                                Code     string `json:"code"`
                                Title    string `json:"title"`
                                Number   string `json:"number"`
                                Category string `json:"category"`
                        }{Code: r.Code, Title: r.Title, Number: r.Number, Category: r.Category}
                }
                mapsOut = append(mapsOut, row)
        }
        return auth.OK(gin.H{
                "regulations": regsOut, "mappings": mapsOut, "summary": summary, "total": len(mappings),
        })
}

func handleEvidence(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        q := db.DB.Where("hospital_id = ?", hospitalID)
        if v := c.Query("refType"); v != "" {
                q = q.Where("ref_type = ?", v)
        }
        if v := c.Query("refId"); v != "" {
                q = q.Where("ref_id = ?", v)
        }
        if s := c.Query("q"); s != "" {
                like := "%" + s + "%"
                q = q.Where("title LIKE ? OR metric LIKE ? OR source_system LIKE ?", like, like, like)
        }
        var items []model.EvidenceItem
        if err := q.Order("timestamp desc").Find(&items).Error; err != nil {
                return auth.Fail(err)
        }
        byType := map[string]int{}
        for _, i := range items {
                byType[i.RefType]++
        }
        return auth.OK(gin.H{"items": items, "byType": byType, "total": len(items)})
}

func handleAudit(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        var q = db.DB.Model(&model.AuditLog{})
        if user.Role == "super_admin" {
                // super admin melihat semua
        } else {
                hid := hospitalID
                q = q.Where("hospital_id = ? OR (hospital_id IS NULL AND actor_role IN ?)", hid, []string{"super_admin", "system"})
        }
        if v := c.Query("action"); v != "" {
                q = q.Where("action = ?", v)
        }
        if v := c.Query("entityType"); v != "" {
                q = q.Where("entity_type = ?", v)
        }
        if s := c.Query("q"); s != "" {
                like := "%" + s + "%"
                q = q.Where("actor_name LIKE ? OR description LIKE ? OR entity_id LIKE ?", like, like, like)
        }
        var logs []model.AuditLog
        if err := q.Order("created_at desc").Limit(200).Find(&logs).Error; err != nil {
                return auth.Fail(err)
        }
        return auth.OK(gin.H{"logs": logs, "total": len(logs)})
}

type searchResult struct {
        Type  string `json:"type"`
        Icon  string `json:"icon"`
        Label string `json:"label"`
        Sub   string `json:"sub"`
        View  string `json:"view"`
        ID    string `json:"id"`
}

func handleSearch(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        term := c.Query("q")
        if len(term) < 2 {
                return auth.OK(gin.H{"results": []searchResult{}, "term": term})
        }
        like := "%" + term + "%"
        results := []searchResult{}

        var practitioners []model.Practitioner
        db.DB.Where("hospital_id = ? AND (full_name LIKE ? OR specialty LIKE ? OR sip_number LIKE ? OR str_number LIKE ?)",
                hospitalID, like, like, like, like).Limit(5).Find(&practitioners)
        for _, p := range practitioners {
                sub := p.Profession
                if p.Specialty != nil && *p.Specialty != "" {
                        sub += " — " + *p.Specialty
                }
                results = append(results, searchResult{"Praktisi", "Users", p.FullName, sub, "cred-practitioners", p.ID})
        }

        var credentials []model.PractitionerCredential
        db.DB.Where("hospital_id = ? AND (name LIKE ? OR number LIKE ?)", hospitalID, like, like).Limit(5).Find(&credentials)
        for _, cr := range credentials {
                results = append(results, searchResult{"Kredensial", "IdCard", cr.Name, cr.Type + " · " + cr.Status, "cred-credentials", cr.ID})
        }

        var risks []model.Risk
        db.DB.Where("hospital_id = ? AND (title LIKE ? OR code LIKE ?)", hospitalID, like, like).Limit(5).Find(&risks)
        for _, r := range risks {
                results = append(results, searchResult{"Risiko", "Radar", r.Title, r.Code + " · " + r.Category + " · " + r.Level, "risk-radar", r.ID})
        }

        var claims []model.Claim
        db.DB.Where("hospital_id = ? AND (claim_number LIKE ? OR sep_number LIKE ? OR inacbg_desc LIKE ?)", hospitalID, like, like, like).Limit(5).Find(&claims)
        for _, cl := range claims {
                results = append(results, searchResult{"Klaim", "FileCheck2", cl.ClaimNumber, cl.InacbgCode + " " + cl.InacbgDesc + " · " + cl.Status, "claim-anomalies", cl.ID})
        }

        var actions []model.Action
        db.DB.Where("hospital_id = ? AND (title LIKE ? OR code LIKE ? OR owner LIKE ?)", hospitalID, like, like, like).Limit(5).Find(&actions)
        for _, a := range actions {
                results = append(results, searchResult{"Aksi", "CheckCircle2", a.Title, a.Code + " · " + a.Status + " · " + a.Owner, "actions", a.ID})
        }

        var regulations []model.Regulation
        db.DB.Where("title LIKE ? OR number LIKE ? OR code LIKE ?", like, like, like).Limit(4).Find(&regulations)
        for _, r := range regulations {
                results = append(results, searchResult{"Regulasi", "Scale", r.Title, r.Number, "compliance", r.ID})
        }

        var evidence []model.EvidenceItem
        db.DB.Where("hospital_id = ? AND (title LIKE ? OR metric LIKE ?)", hospitalID, like, like).Limit(4).Find(&evidence)
        for _, e := range evidence {
                sub := e.SourceSystem
                if e.Metric != nil {
                        sub = e.SourceSystem + " · " + e.ValidationStatus
                }
                results = append(results, searchResult{"Bukti", "Fingerprint", e.Title, sub, "evidence", e.ID})
        }

        var warnings []model.EarlyWarning
        db.DB.Where("hospital_id = ? AND (title LIKE ? OR code LIKE ?)", hospitalID, like, like).Limit(4).Find(&warnings)
        for _, w := range warnings {
                results = append(results, searchResult{"Peringatan", "AlertTriangle", w.Title, w.Code + " · " + w.Severity, "early-warnings", w.ID})
        }

        return auth.OK(gin.H{"results": results, "term": term})
}

var _ = sort.Strings
