package api

// Risk Radar & Early Warning Center — GET/PATCH.
// Port dari src/app/api/risks/route.ts & warnings/route.ts.
import (
        "time"

        "github.com/gin-gonic/gin"

        "medcredix/backend/internal/auth"
        "medcredix/backend/internal/db"
        "medcredix/backend/internal/model"
)

func handleRisksList(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        q := db.DB.Where("hospital_id = ?", hospitalID)
        if v := c.Query("level"); v != "" {
                q = q.Where("level = ?", v)
        }
        if v := c.Query("category"); v != "" {
                q = q.Where("category = ?", v)
        }
        if v := c.Query("status"); v != "" {
                q = q.Where("status = ?", v)
        }
        var risks []model.Risk
        if err := q.Order("detected_at desc").Find(&risks).Error; err != nil {
                return auth.Fail(err)
        }
        var evidence []model.EvidenceItem
        db.DB.Where("hospital_id = ? AND ref_type = ?", hospitalID, "Risiko").Find(&evidence)
        return auth.OK(gin.H{"risks": risks, "evidence": evidence, "total": len(risks)})
}

type riskPatchBody struct {
        ID     string `json:"id"`
        Action string `json:"action"`
        Status string `json:"status"`
        Owner  string `json:"owner"`
        Deadline string `json:"deadline"`
}

func handleRiskPatch(c *gin.Context, user *auth.SessionUser) (int, any) {
        var body riskPatchBody
        if err := c.ShouldBindJSON(&body); err != nil {
                return 400, "Payload tidak valid."
        }
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        var risk model.Risk
        if err := db.DB.Where("id = ? AND hospital_id = ?", body.ID, hospitalID).First(&risk).Error; err != nil {
                return 404, "Risiko tidak ditemukan."
        }
        if body.Action == "ack" || body.Action == "assign" {
                if err := auth.RequirePermission(user, "ASSIGN"); err != nil {
                        return auth.Fail(err)
                }
        } else if err := auth.RequirePermission(user, "EDIT"); err != nil {
                return auth.Fail(err)
        }
        updates := map[string]any{}
        if body.Status != "" {
                updates["status"] = body.Status
        } else if body.Action == "ack" {
                updates["status"] = "Ditinjau"
        }
        if body.Owner != "" {
                updates["suggested_owner"] = body.Owner
        }
        if body.Deadline != "" {
                if t, err := time.Parse(time.RFC3339, body.Deadline); err == nil {
                        updates["deadline"] = t
                } else if t, err := time.Parse("2006-01-02", body.Deadline); err == nil {
                        updates["deadline"] = t
                }
        }
        if len(updates) == 0 {
                return auth.OK(risk)
        }
        if err := db.DB.Model(&risk).Updates(updates).Error; err != nil {
                return auth.Fail(err)
        }
        var updated model.Risk
        db.DB.Where("id = ?", risk.ID).First(&updated)
        hid := hospitalID
        auth.LogAudit(auth.AuditOptions{
                HospitalID: &hid, Actor: user, Action: "update", EntityType: "Risiko",
                EntityID: &risk.ID, Description: risk.Code + ": " + body.Action + ownerSuffix(body.Owner),
                Before: sp(risk.Status), After: sp(updated.Status), C: c,
        })
        return auth.OK(updated)
}

func ownerSuffix(owner string) string {
        if owner == "" {
                return ""
        }
        return " → " + owner
}

func sp(s string) *string { return &s }

func handleWarningsList(c *gin.Context, user *auth.SessionUser) (int, any) {
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        q := db.DB.Where("hospital_id = ?", hospitalID)
        if v := c.Query("severity"); v != "" {
                q = q.Where("severity = ?", v)
        }
        if v := c.Query("status"); v != "" {
                q = q.Where("status = ?", v)
        }
        var warnings []model.EarlyWarning
        if err := q.Order("status asc, detected_at desc").Find(&warnings).Error; err != nil {
                return auth.Fail(err)
        }
        var rules []model.AlertRule
        db.DB.Where("hospital_id = ?", hospitalID).Order("name asc").Find(&rules)
        return auth.OK(gin.H{"warnings": warnings, "rules": rules, "total": len(warnings)})
}

type warningPatchBody struct {
        ID       string `json:"id"`
        Action   string `json:"action"`
        Assignee string `json:"assignee"`
        DueDate  string `json:"dueDate"`
}

func handleWarningPatch(c *gin.Context, user *auth.SessionUser) (int, any) {
        var body warningPatchBody
        if err := c.ShouldBindJSON(&body); err != nil {
                return 400, "Payload tidak valid."
        }
        hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
        if err != nil {
                return auth.Fail(err)
        }
        var w model.EarlyWarning
        if err := db.DB.Where("id = ? AND hospital_id = ?", body.ID, hospitalID).First(&w).Error; err != nil {
                return 404, "Peringatan tidak ditemukan."
        }
        if body.Action == "resolve" {
                if err := auth.RequirePermission(user, "RESOLVE"); err != nil {
                        return auth.Fail(err)
                }
        } else if err := auth.RequirePermission(user, "ASSIGN"); err != nil {
                return auth.Fail(err)
        }
        statusMap := map[string]string{"ack": "Diketahui", "assign": "Ditugaskan", "resolve": "Teratasi"}
        newStatus := w.Status
        if v, ok := statusMap[body.Action]; ok {
                newStatus = v
        }
        updates := map[string]any{"status": newStatus}
        if body.Action == "ack" || body.Action == "assign" {
                updates["acknowledged_by"] = user.Name
                updates["acknowledged_at"] = time.Now()
        }
        if body.Action == "assign" {
                assignee := body.Assignee
                if assignee == "" && w.AcknowledgedBy != nil {
                        assignee = *w.AcknowledgedBy
                }
                if assignee == "" && w.AssignedTo != nil {
                        assignee = *w.AssignedTo
                }
                if assignee != "" {
                        updates["assigned_to"] = assignee
                }
        }
        if body.DueDate != "" {
                if t, err := time.Parse(time.RFC3339, body.DueDate); err == nil {
                        updates["due_date"] = t
                } else if t, err := time.Parse("2006-01-02", body.DueDate); err == nil {
                        updates["due_date"] = t
                }
        }
        if err := db.DB.Model(&w).Updates(updates).Error; err != nil {
                return auth.Fail(err)
        }
        var updated model.EarlyWarning
        db.DB.Where("id = ?", w.ID).First(&updated)
        hid := hospitalID
        auth.LogAudit(auth.AuditOptions{
                HospitalID: &hid, Actor: user, Action: "update", EntityType: "Peringatan",
                EntityID: &w.ID, Description: w.Code + ": " + body.Action,
                Before: sp(w.Status), After: sp(updated.Status), C: c,
        })
        return auth.OK(updated)
}
