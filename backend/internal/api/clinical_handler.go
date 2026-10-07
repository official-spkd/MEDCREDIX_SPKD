package api

// Claim Intelligence & Clinical Intelligence — GET/PATCH.
// Port dari claims/route.ts & clinical/route.ts.
import (
        "sort"

        "github.com/gin-gonic/gin"

        "medcredix/backend/internal/auth"
        "medcredix/backend/internal/db"
        "medcredix/backend/internal/model"
)

type anomalyOut struct {
        model.ClaimAnomaly
        ClaimNumber string `json:"claimNumber"`
}

func handleClaimsList(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        if id := c.Query("id"); id != "" {
                var claim model.Claim
                if err := db.DB.Where("id = ? AND hospital_id = ?", id, hospitalID).First(&claim).Error; err != nil {
                        return 404, "Klaim tidak ditemukan."
                }
                var anomalies []model.ClaimAnomaly
                db.DB.Where("claim_id = ?", id).Find(&anomalies)
                var evidence []model.EvidenceItem
                db.DB.Where("hospital_id = ? AND ref_id = ?", hospitalID, claim.ClaimNumber).Find(&evidence)
                return auth.OK(gin.H{"claim": claim, "anomalies": anomalies, "evidence": evidence})
        }
        q := db.DB.Where("hospital_id = ?", hospitalID)
        if v := c.Query("status"); v != "" {
                q = q.Where("status = ?", v)
        }
        if v := c.Query("careType"); v != "" {
                q = q.Where("care_type = ?", v)
        }
        if s := c.Query("q"); s != "" {
                like := "%" + s + "%"
                q = q.Where("claim_number LIKE ? OR sep_number LIKE ? OR inacbg_code LIKE ?", like, like, like)
        }
        var claims []model.Claim
        if err := q.Order("admission_date desc").Find(&claims).Error; err != nil {
                return auth.Fail(err)
        }
        intel, err := ComputeIntel(hospitalID)
        if err != nil {
                return auth.Fail(err)
        }
        var anomalies []model.ClaimAnomaly
        db.DB.Where("hospital_id = ?", hospitalID).Order("status asc, severity desc").Find(&anomalies)
        claimNumberByID := map[string]string{}
        for _, cl := range claims {
                claimNumberByID[cl.ID] = cl.ClaimNumber
        }
        anomaliesOut := make([]anomalyOut, 0, len(anomalies))
        for _, a := range anomalies {
                num := claimNumberByID[a.ClaimID]
                if num == "" {
                        num = "—"
                }
                anomaliesOut = append(anomaliesOut, anomalyOut{ClaimAnomaly: a, ClaimNumber: num})
        }
        return auth.OK(gin.H{"claims": claims, "anomalies": anomaliesOut, "ccs": intel.Ccs, "total": len(claims)})
}

type claimPatchBody struct {
        AnomalyID string `json:"anomalyId"`
        ClaimID   string `json:"claimId"`
        Action    string `json:"action"`
        Status    string `json:"status"`
}

func handleClaimPatch(c *gin.Context, user *auth.SessionUser) (int, any) {
        var body claimPatchBody
        if err := c.ShouldBindJSON(&body); err != nil {
                return 400, "Payload tidak valid."
        }
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        if err := auth.RequirePermission(user, "EDIT"); err != nil {
                return auth.Fail(err)
        }
        if body.AnomalyID != "" {
                var a model.ClaimAnomaly
                if err := db.DB.Where("id = ? AND hospital_id = ?", body.AnomalyID, hospitalID).First(&a).Error; err != nil {
                        return 404, "Anomali tidak ditemukan."
                }
                statusMap := map[string]string{"review": "Ditinjau", "confirm": "Terkonfirmasi Bermasalah", "clarify": "Diklarifikasi"}
                newStatus := a.Status
                if v, ok := statusMap[body.Action]; ok {
                        newStatus = v
                }
                if err := db.DB.Model(&a).Updates(map[string]any{"status": newStatus, "reviewed_by": user.Name}).Error; err != nil {
                        return auth.Fail(err)
                }
                var claim model.Claim
                db.DB.Where("id = ?", a.ClaimID).First(&claim)
                hid := hospitalID
                auth.LogAudit(auth.AuditOptions{
                        HospitalID: &hid, Actor: user, Action: "update", EntityType: "Anomali Klaim", EntityID: &a.ID,
                        Description: claim.ClaimNumber + " (" + a.Category + "): " + a.Status + " → " + newStatus,
                        Before: sp(a.Status), After: sp(newStatus), C: c,
                })
                var updated model.ClaimAnomaly
                db.DB.Where("id = ?", a.ID).First(&updated)
                return auth.OK(updated)
        }
        if body.ClaimID != "" {
                var cl model.Claim
                if err := db.DB.Where("id = ? AND hospital_id = ?", body.ClaimID, hospitalID).First(&cl).Error; err != nil {
                        return 404, "Klaim tidak ditemukan."
                }
                if err := db.DB.Model(&cl).Update("status", body.Status).Error; err != nil {
                        return auth.Fail(err)
                }
                hid := hospitalID
                auth.LogAudit(auth.AuditOptions{
                        HospitalID: &hid, Actor: user, Action: "update", EntityType: "Klaim", EntityID: &cl.ID,
                        Description: cl.ClaimNumber + ": " + cl.Status + " → " + body.Status,
                        Before: sp(cl.Status), After: sp(body.Status), C: c,
                })
                var updated model.Claim
                db.DB.Where("id = ?", cl.ID).First(&updated)
                return auth.OK(updated)
        }
        return 400, "Payload tidak valid."
}

func handleClinical(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        typ := c.Query("type")
        if typ == "" {
                typ = "all"
        }
        intel, err := ComputeIntel(hospitalID)
        if err != nil {
                return auth.Fail(err)
        }
        if typ == "pathway" {
                var pathways []model.ClinicalPathway
                db.DB.Where("hospital_id = ?", hospitalID).Order("avg_adherence asc").Find(&pathways)
                var episodes []model.PathwayEpisode
                db.DB.Where("hospital_id = ?", hospitalID).Order("adherence asc").Find(&episodes)
                return auth.OK(gin.H{"pathways": pathways, "episodes": episodes, "cpcs": intel.Cpcs,
                        "coverage": gin.H{"total": 20, "covered": 16}})
        }
        if typ == "utilization" {
                var utils []model.UtilizationStat
                db.DB.Where("hospital_id = ?", hospitalID).Order("month asc").Find(&utils)
                var complaints []model.ComplaintStat
                db.DB.Where("hospital_id = ?", hospitalID).Order("month asc").Find(&complaints)
                return auth.OK(gin.H{"utilization": utils, "complaints": complaints, "ur": intel.Ur})
        }
        if typ == "safety" {
                var safetyDomain any
                for i := range intel.Domains {
                        if intel.Domains[i].Key == "safety" {
                                safetyDomain = intel.Domains[i]
                                break
                        }
                }
                return auth.OK(gin.H{"safety": intel.Safety, "domain": safetyDomain})
        }
        var pathways []model.ClinicalPathway
        db.DB.Where("hospital_id = ?", hospitalID).Find(&pathways)
        var utils []model.UtilizationStat
        db.DB.Where("hospital_id = ?", hospitalID).Order("month asc").Find(&utils)
        var complaints []model.ComplaintStat
        db.DB.Where("hospital_id = ?", hospitalID).Order("month asc").Find(&complaints)
        return auth.OK(gin.H{"pathways": pathways, "utilization": utils, "complaints": complaints,
                "cpcs": intel.Cpcs, "ur": intel.Ur, "safety": intel.Safety})
}

var _ = sort.Strings
