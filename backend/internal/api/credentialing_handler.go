package api

// Praktisi, kredensial, requirements — GET/PATCH.
// Port dari practitioners/route.ts, credentials/route.ts, requirements/route.ts.
import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"medcredix/backend/internal/auth"
	"medcredix/backend/internal/db"
	"medcredix/backend/internal/model"
)

func handlePractitioners(c *gin.Context, user *auth.SessionUser) (int, any) {
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	if id := c.Query("id"); id != "" {
		var p model.Practitioner
		if err := db.DB.Where("id = ? AND hospital_id = ?", id, hospitalID).First(&p).Error; err != nil {
			return 404, "Praktisi tidak ditemukan."
		}
		var creds []model.PractitionerCredential
		db.DB.Where("practitioner_id = ?", id).Order("type asc").Find(&creds)
		var docs []model.CredentialDocument
		db.DB.Where("hospital_id = ?", hospitalID).Find(&docs)
		var standards []model.SpecialtyStandard
		db.DB.Where("hospital_id = ?", hospitalID).Find(&standards)
		credIDs := map[string]bool{}
		for _, cr := range creds {
			credIDs[cr.ID] = true
		}
		filtered := make([]model.CredentialDocument, 0)
		for _, d := range docs {
			if credIDs[d.CredentialID] {
				filtered = append(filtered, d)
			}
		}
		var standard *model.SpecialtyStandard
		if p.Specialty != nil {
			for i := range standards {
				if standards[i].Specialty == *p.Specialty {
					standard = &standards[i]
					break
				}
			}
		}
		return auth.OK(gin.H{"practitioner": p, "credentials": creds, "documents": filtered, "standard": standard})
	}
	q := db.DB.Where("hospital_id = ?", hospitalID)
	if s := c.Query("q"); s != "" {
		like := "%" + s + "%"
		q = q.Where("full_name LIKE ? OR specialty LIKE ? OR sip_number LIKE ? OR str_number LIKE ?", like, like, like, like)
	}
	if v := c.Query("prof"); v != "" {
		q = q.Where("profession = ?", v)
	}
	if v := c.Query("status"); v != "" {
		q = q.Where("status = ?", v)
	}
	if v := c.Query("dept"); v != "" {
		q = q.Where("department LIKE ?", "%"+v+"%")
	}
	var list []model.Practitioner
	if err := q.Order("full_name asc").Find(&list).Error; err != nil {
		return auth.Fail(err)
	}
	deptSet := map[string]bool{}
	depts := []string{}
	for _, p := range list {
		if !deptSet[p.Department] {
			deptSet[p.Department] = true
			depts = append(depts, p.Department)
		}
	}
	return auth.OK(gin.H{"practitioners": list, "depts": depts, "total": len(list)})
}

type credOut struct {
	model.PractitionerCredential
	PractitionerName string  `json:"practitionerName"`
	Profession       string  `json:"profession"`
	Specialty        *string `json:"specialty"`
}

func handleCredentialsList(c *gin.Context, user *auth.SessionUser) (int, any) {
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	q := db.DB.Where("hospital_id = ?", hospitalID)
	if v := c.Query("status"); v != "" {
		q = q.Where("status = ?", v)
	}
	if v := c.Query("type"); v != "" {
		q = q.Where("type = ?", v)
	}
	if s := c.Query("q"); s != "" {
		like := "%" + s + "%"
		q = q.Where("name LIKE ? OR number LIKE ?", like, like)
	}
	var creds []model.PractitionerCredential
	if err := q.Order("expiry_date asc").Find(&creds).Error; err != nil {
		return auth.Fail(err)
	}
	var practitioners []model.Practitioner
	db.DB.Where("hospital_id = ?", hospitalID).Find(&practitioners)
	pmap := map[string]*model.Practitioner{}
	for i := range practitioners {
		pmap[practitioners[i].ID] = &practitioners[i]
	}
	out := make([]credOut, 0, len(creds))
	for _, cr := range creds {
		row := credOut{PractitionerCredential: cr, PractitionerName: "—", Profession: "—"}
		if p, ok := pmap[cr.PractitionerID]; ok {
			row.PractitionerName = p.FullName
			row.Profession = p.Profession
			row.Specialty = p.Specialty
		}
		out = append(out, row)
	}
	return auth.OK(gin.H{"credentials": out, "total": len(out)})
}

type credPatchBody struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

func handleCredentialPatch(c *gin.Context, user *auth.SessionUser) (int, any) {
	var body credPatchBody
	if err := c.ShouldBindJSON(&body); err != nil {
		return 400, "Payload tidak valid."
	}
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	var cred model.PractitionerCredential
	if err := db.DB.Where("id = ? AND hospital_id = ?", body.ID, hospitalID).First(&cred).Error; err != nil {
		return 404, "Kredensial tidak ditemukan."
	}
	type transition struct {
		to   string
		perm string
	}
	transitions := map[string]transition{
		"submit":     {"Diajukan", "EDIT"},
		"verify":     {"Verifikasi", "VERIFY"},
		"review":     {"Tinjauan Komite", "VERIFY"},
		"approve":    {"Disetujui", "APPROVE"},
		"activate":   {"Aktif", "APPROVE"},
		"suspend":    {"Ditangguhkan", "APPROVE"},
		"reactivate": {"Aktif", "APPROVE"},
	}
	t, ok := transitions[body.Action]
	if !ok {
		return 400, "Aksi tidak dikenal."
	}
	if err := auth.RequirePermission(user, t.perm); err != nil {
		return auth.Fail(err)
	}
	// Larang modifikasi diam-diam pada record yang disetujui → versioning
	bump := false
	if cred.Status == "Disetujui" || cred.Status == "Aktif" || cred.Status == "Segera Berakhir" {
		if body.Action != "reactivate" {
			bump = true
		}
	}
	updates := map[string]any{"status": t.to}
	if bump {
		updates["version"] = cred.Version + 1
	}
	if body.Action == "verify" || body.Action == "review" {
		updates["verified_by_id"] = user.ID
		updates["verified_at"] = time.Now()
	}
	if body.Action == "approve" || body.Action == "activate" || body.Action == "reactivate" || body.Action == "suspend" {
		updates["approved_by_id"] = user.ID
		updates["approved_at"] = time.Now()
	}
	if err := db.DB.Model(&cred).Updates(updates).Error; err != nil {
		return auth.Fail(err)
	}
	var updated model.PractitionerCredential
	db.DB.Where("id = ?", cred.ID).First(&updated)
	auditAction := "update"
	if body.Action == "approve" || body.Action == "activate" {
		auditAction = "approve"
	} else if body.Action == "verify" || body.Action == "review" {
		auditAction = "verify"
	}
	hid := hospitalID
	auth.LogAudit(auth.AuditOptions{
		HospitalID: &hid, Actor: user, Action: auditAction, EntityType: "Kredensial",
		EntityID: &cred.ID, Description: body.Action + " → " + t.to + " (" + cred.Name + ")",
		Before: sp(cred.Status), After: sp(t.to), C: c,
	})
	return auth.OK(updated)
}

func handleRequirementsList(c *gin.Context, user *auth.SessionUser) (int, any) {
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	var reqs []model.Requirement
	if err := db.DB.Where("hospital_id = ?", hospitalID).Order("code asc").Find(&reqs).Error; err != nil {
		return auth.Fail(err)
	}
	return auth.OK(gin.H{"requirements": reqs})
}

type reqPatchBody struct {
	ID        string  `json:"id"`
	Action    string  `json:"action"`
	DocNumber *string `json:"docNumber"`
	IssuedAt  *string `json:"issuedAt"`
	DocExpiry *string `json:"docExpiry"`
	FileName  *string `json:"fileName"`
	Notes     *string `json:"notes"`
}

func handleRequirementPatch(c *gin.Context, user *auth.SessionUser) (int, any) {
	var body reqPatchBody
	if err := c.ShouldBindJSON(&body); err != nil {
		return 400, "Payload tidak valid."
	}
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	perm := "EDIT"
	if body.Action == "verify" {
		perm = "VERIFY"
	}
	if err := auth.RequirePermission(user, perm); err != nil {
		return auth.Fail(err)
	}
	var r model.Requirement
	if err := db.DB.Where("id = ? AND hospital_id = ?", body.ID, hospitalID).First(&r).Error; err != nil {
		return 404, "Prasyarat tidak ditemukan."
	}
	status := r.Status
	switch body.Action {
	case "submit":
		status = "Menunggu Verifikasi"
	case "verify":
		status = "Valid"
	case "invalidate":
		status = "Belum Dilengkapi"
	}
	updates := map[string]any{"status": status, "updated_at": time.Now()}
	if body.DocNumber != nil {
		updates["doc_number"] = *body.DocNumber
	}
	if body.IssuedAt != nil {
		if t, err := time.Parse("2006-01-02", *body.IssuedAt); err == nil {
			updates["issued_at"] = t
		}
	}
	if body.DocExpiry != nil {
		if t, err := time.Parse("2006-01-02", *body.DocExpiry); err == nil {
			updates["doc_expiry"] = t
		}
	}
	if body.FileName != nil {
		updates["file_name"] = *body.FileName
	}
	if body.Notes != nil {
		updates["notes"] = *body.Notes
	}
	if err := db.DB.Model(&r).Updates(updates).Error; err != nil {
		return auth.Fail(err)
	}
	var updated model.Requirement
	db.DB.Where("id = ?", r.ID).First(&updated)
	auditAction := "update"
	if body.Action == "verify" {
		auditAction = "verify"
	}
	hid := hospitalID
	auth.LogAudit(auth.AuditOptions{
		HospitalID: &hid, Actor: user, Action: auditAction, EntityType: "Prasyarat",
		EntityID: &r.ID, Description: r.Code + " → " + status,
		Before: sp(r.Status), After: sp(status), C: c,
	})
	return auth.OK(updated)
}

var _ = strings.TrimSpace
