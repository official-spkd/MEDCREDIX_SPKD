package api

// Action Center — closed-loop: create, assign, status, comment, resolve, reopen.
// Port dari src/app/api/actions/route.ts.
import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"medcredix/backend/internal/auth"
	"medcredix/backend/internal/db"
	"medcredix/backend/internal/model"
)

var openActionStatuses = []string{"Baru", "Ditugaskan", "Dalam Pengerjaan", "Menunggu Verifikasi"}

func handleActionsList(c *gin.Context, user *auth.SessionUser) (int, any) {
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	if id := c.Query("id"); id != "" {
		var a model.Action
		if err := db.DB.Where("id = ? AND hospital_id = ?", id, hospitalID).First(&a).Error; err != nil {
			return 404, "Aksi tidak ditemukan."
		}
		var events []model.ActionEvent
		db.DB.Where("action_id = ?", id).Order("created_at asc").Find(&events)
		var users []model.User
		db.DB.Where("hospital_id = ?", hospitalID).Find(&users)
		type userBrief struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Role string `json:"role"`
		}
		ub := make([]userBrief, 0, len(users))
		for _, u := range users {
			ub = append(ub, userBrief{ID: u.ID, Name: u.Name, Role: u.Role})
		}
		return auth.OK(gin.H{"action": a, "events": events, "users": ub})
	}
	q := db.DB.Where("hospital_id = ?", hospitalID)
	if v := c.Query("status"); v != "" {
		q = q.Where("status = ?", v)
	}
	if v := c.Query("priority"); v != "" {
		q = q.Where("priority = ?", v)
	}
	if s := c.Query("q"); s != "" {
		like := "%" + s + "%"
		q = q.Where("title LIKE ? OR code LIKE ? OR owner LIKE ?", like, like, like)
	}
	var actions []model.Action
	if err := q.Order("created_at desc").Find(&actions).Error; err != nil {
		return auth.Fail(err)
	}
	open, high, resolved := 0, 0, 0
	for _, a := range actions {
		if containsStr(openActionStatuses, a.Status) {
			open++
		}
		if a.Priority == "Tinggi" && a.Status != "Selesai" && a.Status != "Ditutup" {
			high++
		}
		if a.Status == "Selesai" || a.Status == "Ditutup" {
			resolved++
		}
	}
	return auth.OK(gin.H{
		"actions": actions,
		"stats": gin.H{"total": len(actions), "open": open, "high": high, "resolved": resolved},
	})
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

type actionCreateBody struct {
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Priority        string  `json:"priority"`
	SourceType      string  `json:"sourceType"`
	SourceRef       *string `json:"sourceRef"`
	Owner           string  `json:"owner"`
	OwnerRole       string  `json:"ownerRole"`
	DueDate         string  `json:"dueDate"`
	ExpectedOutcome *string `json:"expectedOutcome"`
}

func handleActionCreate(c *gin.Context, user *auth.SessionUser) (int, any) {
	var body actionCreateBody
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
	var count int64
	db.DB.Model(&model.Action{}).Where("hospital_id = ?", hospitalID).Count(&count)
	code := "ACT-" + pad3(int(count)+1)
	priority := body.Priority
	if priority == "" {
		priority = "Sedang"
	}
	sourceType := body.SourceType
	if sourceType == "" {
		sourceType = "Manual"
	}
	owner := body.Owner
	if owner == "" {
		owner = "Belum ditugaskan"
	}
	created := model.Action{
		HospitalID: hospitalID, Code: code, Title: body.Title, Description: body.Description,
		Priority: priority, Status: "Baru", SourceType: sourceType, SourceRef: body.SourceRef,
		Owner: owner, OwnerRole: body.OwnerRole, ExpectedOutcome: body.ExpectedOutcome,
	}
	if body.DueDate != "" {
		if t, err := time.Parse(time.RFC3339, body.DueDate); err == nil {
			created.DueDate = &t
		} else if t, err := time.Parse("2006-01-02", body.DueDate); err == nil {
			created.DueDate = &t
		}
	}
	if err := db.DB.Create(&created).Error; err != nil {
		return auth.Fail(err)
	}
	srcRef := ""
	if body.SourceRef != nil {
		srcRef = *body.SourceRef
	}
	event := model.ActionEvent{
		ActionID: created.ID, Type: "created", ToValue: sp("Baru"), ByName: user.Name,
		Note: sp("Aksi dibuat dari " + sourceType + srcRefSuffix(srcRef)),
	}
	db.DB.Create(&event)
	hid := hospitalID
	auth.LogAudit(auth.AuditOptions{
		HospitalID: &hid, Actor: user, Action: "create", EntityType: "Aksi",
		EntityID: &created.ID, Description: "Membuat aksi " + code + ": " + body.Title, C: c,
	})
	return auth.OK(created)
}

type actionPatchBody struct {
	ID             string  `json:"id"`
	Action         string  `json:"action"`
	Owner          string  `json:"owner"`
	OwnerRole      string  `json:"ownerRole"`
	DueDate        string  `json:"dueDate"`
	Status         string  `json:"status"`
	Note           *string `json:"note"`
	ResolutionNote *string `json:"resolutionNote"`
}

func handleActionPatch(c *gin.Context, user *auth.SessionUser) (int, any) {
	var body actionPatchBody
	if err := c.ShouldBindJSON(&body); err != nil {
		return 400, "Payload tidak valid."
	}
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	var a model.Action
	if err := db.DB.Where("id = ? AND hospital_id = ?", body.ID, hospitalID).First(&a).Error; err != nil {
		return 404, "Aksi tidak ditemukan."
	}

	switch body.Action {
	case "assign":
		if err := auth.RequirePermission(user, "ASSIGN"); err != nil {
			return auth.Fail(err)
		}
		owner := orDefault(body.Owner, a.Owner)
		ownerRole := orDefault(body.OwnerRole, a.OwnerRole)
		status := a.Status
		if status == "Baru" {
			status = "Ditugaskan"
		}
		updates := map[string]any{"owner": owner, "owner_role": ownerRole, "status": status}
		if body.DueDate != "" {
			if t, err := time.Parse(time.RFC3339, body.DueDate); err == nil {
				updates["due_date"] = t
			} else if t, err := time.Parse("2006-01-02", body.DueDate); err == nil {
				updates["due_date"] = t
			}
		}
		if err := db.DB.Model(&a).Updates(updates).Error; err != nil {
			return auth.Fail(err)
		}
		db.DB.Create(&model.ActionEvent{ActionID: a.ID, Type: "assigned", FromValue: sp(a.Owner),
			ToValue: sp(owner), ByName: user.Name, Note: body.Note})
		hid := hospitalID
		auth.LogAudit(auth.AuditOptions{
			HospitalID: &hid, Actor: user, Action: "update", EntityType: "Aksi", EntityID: &a.ID,
			Description: a.Code + " ditugaskan ke " + owner, Before: sp(a.Owner), After: sp(owner), C: c,
		})
		var updated model.Action
		db.DB.Where("id = ?", a.ID).First(&updated)
		return auth.OK(updated)

	case "status":
		if err := auth.RequirePermission(user, "EDIT"); err != nil {
			return auth.Fail(err)
		}
		updates := map[string]any{"status": body.Status}
		if body.Note != nil {
			updates["progress_note"] = *body.Note
		}
		if body.Status == "Selesai" && body.ResolutionNote != nil {
			updates["resolution_note"] = *body.ResolutionNote
		}
		if body.Status == "Selesai" || body.Status == "Ditutup" {
			now := time.Now()
			updates["resolved_at"] = now
		}
		if err := db.DB.Model(&a).Updates(updates).Error; err != nil {
			return auth.Fail(err)
		}
		db.DB.Create(&model.ActionEvent{ActionID: a.ID, Type: "status", FromValue: sp(a.Status),
			ToValue: sp(body.Status), ByName: user.Name, Note: body.Note})
		hid := hospitalID
		auth.LogAudit(auth.AuditOptions{
			HospitalID: &hid, Actor: user, Action: "update", EntityType: "Aksi", EntityID: &a.ID,
			Description: a.Code + ": " + a.Status + " → " + body.Status,
			Before: sp(a.Status), After: sp(body.Status), C: c,
		})
		var updated model.Action
		db.DB.Where("id = ?", a.ID).First(&updated)
		return auth.OK(updated)

	case "comment":
		note := ""
		if body.Note != nil {
			note = *body.Note
		}
		ev := model.ActionEvent{ActionID: a.ID, Type: "comment", ByName: user.Name, Note: sp(note)}
		if err := db.DB.Create(&ev).Error; err != nil {
			return auth.Fail(err)
		}
		return auth.OK(ev)
	}
	return 400, "Aksi tidak dikenal."
}

func orDefault(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}
