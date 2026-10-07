package api

// GET /api/bootstrap — sesi + ringkasan skor + badge + notifikasi.
// GET /api/dashboard — data lengkap dashboard eksekutif (LEVEL 1-6).
// Port dari src/app/api/bootstrap/route.ts & dashboard/route.ts.
import (
        "fmt"
        "sort"
        "strconv"

        "github.com/gin-gonic/gin"

        "medcredix/backend/internal/auth"
        "medcredix/backend/internal/db"
        "medcredix/backend/internal/model"
        "medcredix/backend/internal/scoring"
)

type bootIntelDomain struct {
        Key    string  `json:"key"`
        Label  string  `json:"label"`
        Score  float64 `json:"score"`
        Weight float64 `json:"weight"`
        Group  string  `json:"group"`
}

type bootIntel struct {
        HealthScore  float64               `json:"healthScore"`
        HealthStatus string                `json:"healthStatus"`
        Bcds         scoring.BcdsInfo      `json:"bcds"`
        Kelompok     scoring.KelompokInfo  `json:"kelompok"`
        Trust        scoring.TrustInfo     `json:"trust"`
        Ccs          float64               `json:"ccs"`
        Cpcs         float64               `json:"cpcs"`
        Ur           scoring.UrInfo        `json:"ur"`
        Domains      []bootIntelDomain    `json:"domains"`
        Trend        []scoring.TrendItem   `json:"trend"`
        Prereq       []scoring.PrereqItem  `json:"prereq"`
}

func handleBootstrap(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        var hospital model.Hospital
        if err := db.DB.Where("id = ?", hospitalID).First(&hospital).Error; err != nil {
                return auth.Fail(err)
        }
        var notifications []model.Notification
        db.DB.Where("hospital_id = ?", hospitalID).Order("created_at desc").Limit(8).Find(&notifications)
        var warningsCount, risksCount, actionsCount int64
        db.DB.Model(&model.EarlyWarning{}).Where("hospital_id = ? AND status = ?", hospitalID, "Baru").Count(&warningsCount)
        db.DB.Model(&model.Risk{}).Where("hospital_id = ? AND status = ?", hospitalID, "Terbuka").Count(&risksCount)
        db.DB.Model(&model.Action{}).Where("hospital_id = ? AND status IN ?", hospitalID, []string{"Baru", "Ditugaskan", "Dalam Pengerjaan", "Menunggu Verifikasi"}).Count(&actionsCount)
        var alertRules []model.AlertRule
        db.DB.Where("hospital_id = ?", hospitalID).Find(&alertRules)
        var users []model.User
        db.DB.Where("hospital_id = ?", hospitalID).Find(&users)
        type userBrief struct {
                ID   string `json:"id"`
                Name string `json:"name"`
                Role string `json:"role"`
        }
        userBriefs := make([]userBrief, 0, len(users))
        for _, u := range users {
                userBriefs = append(userBriefs, userBrief{ID: u.ID, Name: u.Name, Role: u.Role})
        }
        intel, err := ComputeIntel(hospitalID)
        if err != nil {
                return auth.Fail(err)
        }
        bi := bootIntel{
                HealthScore: intel.HealthScore, HealthStatus: intel.HealthStatus,
                Bcds: intel.Bcds, Kelompok: intel.Kelompok, Trust: intel.Trust,
                Ccs: intel.Ccs.Score, Cpcs: intel.Cpcs.Score, Ur: intel.Ur,
                Trend: intel.Trend, Prereq: intel.Prereq,
        }
        for _, d := range intel.Domains {
                bi.Domains = append(bi.Domains, bootIntelDomain{Key: d.Key, Label: d.Label, Score: d.Score, Weight: d.Weight, Group: d.Group})
        }
        var hospitals any
        if user.Role == "super_admin" {
                var rows []model.Hospital
                db.DB.Find(&rows)
                type hBrief struct {
                        ID   string `json:"id"`
                        Name string `json:"name"`
                        City string `json:"city"`
                }
                list := make([]hBrief, 0, len(rows))
                for _, h := range rows {
                        list = append(list, hBrief{ID: h.ID, Name: h.Name, City: h.City})
                }
                hospitals = list
        }
        return auth.OK(gin.H{
                "user": user, "hospital": hospital, "hospitals": hospitals, "intel": bi,
                "badges": gin.H{"warnings": warningsCount, "risks": risksCount, "actions": actionsCount},
                "notifications": notifications, "alertRules": alertRules, "users": userBriefs,
        })
}

type timelineItem struct {
        Period string  `json:"period"`
        Type   string  `json:"type"`
        Label  string  `json:"label"`
        Delta  *any    `json:"delta"`
}

func handleDashboard(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        intel, err := ComputeIntel(hospitalID)
        if err != nil {
                return auth.Fail(err)
        }
        var risks []model.Risk
        db.DB.Where("hospital_id = ? AND status != ?", hospitalID, "Ditutup").Order("level desc, detected_at desc").Find(&risks)
        var warnings []model.EarlyWarning
        db.DB.Where("hospital_id = ? AND status IN ?", hospitalID, []string{"Baru", "Diketahui", "Ditugaskan"}).Order("severity desc, detected_at desc").Find(&warnings)
        var actions []model.Action
        db.DB.Where("hospital_id = ? AND status IN ?", hospitalID, []string{"Baru", "Ditugaskan", "Dalam Pengerjaan", "Menunggu Verifikasi"}).Order("priority desc, due_date asc").Find(&actions)
        var evidence []model.EvidenceItem
        db.DB.Where("hospital_id = ?", hospitalID).Order("timestamp desc").Limit(6).Find(&evidence)
        var resolvedActions []model.Action
        db.DB.Where("hospital_id = ? AND status IN ? AND resolved_at IS NOT NULL", hospitalID, []string{"Selesai", "Ditutup"}).Order("resolved_at desc").Limit(6).Find(&resolvedActions)
        var pathwayWorst []model.ClinicalPathway
        db.DB.Where("hospital_id = ?", hospitalID).Order("avg_adherence asc").Limit(3).Find(&pathwayWorst)

        // Improvement timeline: snapshot + aksi terakhir
        timeline := []timelineItem{}
        for _, t := range intel.Trend {
                timeline = append(timeline, timelineItem{
                        Period: t.Period, Type: "Skor",
                        Label: fLabel(t.Health) + " · BDCS " + fLabel(t.Bcds) + " · Trust " + fLabel(t.Trust),
                })
        }
        for _, a := range resolvedActions {
                if a.ResolvedAt != nil {
                        timeline = append(timeline, timelineItem{
                                Period: a.ResolvedAt.Format("2006-01"), Type: "Perbaikan",
                                Label: a.Code + " — " + a.Title,
                        })
                }
        }
        sort.SliceStable(timeline, func(i, j int) bool { return timeline[i].Period < timeline[j].Period })

        topActions := actions
        if len(topActions) > 5 {
                topActions = topActions[:5]
        }
        riskSummary := gin.H{
                "tinggi": countLevel(risks, "Tinggi"),
                "sedang": countLevel(risks, "Sedang"),
                "rendah": countLevel(risks, "Rendah"),
        }
        return auth.OK(gin.H{
                "intel": intel, "risks": risks, "warnings": warnings, "actions": actions,
                "topActions": topActions, "evidence": evidence, "pathwayWorst": pathwayWorst,
                "timeline": timeline, "riskSummary": riskSummary,
        })
}

func countLevel(risks []model.Risk, level string) int {
        n := 0
        for _, r := range risks {
                if r.Level == level {
                        n++
                }
        }
        return n
}

func fLabel(v float64) string {
        return trimFloat(v)
}

// trimFloat — tampilan angka ala template literal TS (tanpa .0).
func trimFloat(v float64) string {
        if v == float64(int(v)) {
                return itos(int(v))
        }
        return fts(v)
}

func itos(n int) string {
        if n == 0 {
                return "0"
        }
        neg := n < 0
        if neg {
                n = -n
        }
        var buf [20]byte
        i := len(buf)
        for n > 0 {
                i--
                buf[i] = byte('0' + n%10)
                n /= 10
        }
        if neg {
                i--
                buf[i] = '-'
        }
        return string(buf[i:])
}

func fts(v float64) string {
        return strconv.FormatFloat(v, 'f', 1, 64)
}

var _ = fmt.Sprintf
