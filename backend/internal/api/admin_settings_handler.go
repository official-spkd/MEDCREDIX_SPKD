package api

// Panel Admin SPKD (multi-RS), Settings, Data Intelligence.
// Port dari admin/route.ts, settings/route.ts, data/route.ts.
import (
	"time"

	"github.com/gin-gonic/gin"

	"medcredix/backend/internal/auth"
	"medcredix/backend/internal/db"
	"medcredix/backend/internal/model"
	"medcredix/backend/internal/scoring"
)

type adminTrendItem struct {
	Period string  `json:"period"`
	Bcds   float64 `json:"bcds"`
	Health float64 `json:"health"`
}

type adminRow struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Code         string          `json:"code"`
	Class        string          `json:"class"`
	City         string          `json:"city"`
	Province     string          `json:"province"`
	Tier         string          `json:"tier"`
	Onboarding   bool            `json:"onboarding"`
	HealthScore  float64         `json:"healthScore"`
	Bcds         float64         `json:"bcds"`
	Stage        int             `json:"stage"`
	Lulus        bool            `json:"lulus"`
	Trust        float64         `json:"trust"`
	Ccs          float64         `json:"ccs"`
	Cpcs         float64         `json:"cpcs"`
	AnomalyFlags int64           `json:"anomalyFlags"`
	UnderReview  bool            `json:"underReview"`
	Trend        []adminTrendItem `json:"trend"`
}

func handleAdmin(c *gin.Context, user *auth.SessionUser) (int, any) {
	if user.Role != "super_admin" && user.Role != "auditor" {
		return 403, "Panel Admin hanya untuk Super Admin & Auditor."
	}
	var hospitals []model.Hospital
	db.DB.Find(&hospitals)
	rows := []adminRow{}
	for _, h := range hospitals {
		intel, err := ComputeIntel(h.ID)
		if err != nil {
			return auth.Fail(err)
		}
		var anomalies int64
		db.DB.Model(&model.DataQualityIssue{}).
			Where("hospital_id = ? AND anomaly = ? AND status != ?", h.ID, true, "Diperbaiki").Count(&anomalies)
		trend := make([]adminTrendItem, 0, len(intel.Trend))
		for _, t := range intel.Trend {
			trend = append(trend, adminTrendItem{Period: t.Period, Bcds: t.Bcds, Health: t.Health})
		}
		rows = append(rows, adminRow{
			ID: h.ID, Name: h.Name, Code: h.Code, Class: h.Class, City: h.City, Province: h.Province,
			Tier: h.Tier, Onboarding: h.Onboarding,
			HealthScore: intel.HealthScore, Bcds: intel.Bcds.Final, Stage: intel.Bcds.Stage,
			Lulus: intel.Bcds.Lulus, Trust: intel.Trust.Score, Ccs: intel.Ccs.Score, Cpcs: intel.Cpcs.Score,
			AnomalyFlags: anomalies, UnderReview: intel.Trust.UnderReview, Trend: trend,
		})
	}
	stageDist := []gin.H{}
	for s := 0; s <= 7; s++ {
		count := 0
		for _, r := range rows {
			if r.Stage == s {
				count++
			}
		}
		stageDist = append(stageDist, gin.H{"stage": s, "count": count})
	}
	if len(rows) == 0 {
		return auth.OK(gin.H{"hospitals": rows, "summary": gin.H{
			"totalRS": 0, "rataBdc": 0, "rataHealth": 0, "lulus": 0, "onboarding": 0, "underReview": 0,
		}, "stageDist": stageDist})
	}
	sumBdc, sumHealth := 0.0, 0.0
	lulus, onboarding, underReview := 0, 0, 0
	for _, r := range rows {
		sumBdc += r.Bcds
		sumHealth += r.HealthScore
		if r.Lulus {
			lulus++
		}
		if r.Onboarding {
			onboarding++
		}
		if r.UnderReview {
			underReview++
		}
	}
	return auth.OK(gin.H{
		"hospitals": rows,
		"summary": gin.H{
			"totalRS": len(rows),
			"rataBdc":  roundHalf(sumBdc / float64(len(rows))),
			"rataHealth": roundHalf(sumHealth / float64(len(rows))),
			"lulus":     lulus, "onboarding": onboarding, "underReview": underReview,
		},
		"stageDist": stageDist,
	})
}

func roundHalf(v float64) float64 {
	return float64(int(v*10+0.5)) / 10
}

func handleSettingsGet(c *gin.Context, user *auth.SessionUser) (int, any) {
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	var hospital model.Hospital
	if err := db.DB.Where("id = ?", hospitalID).First(&hospital).Error; err != nil {
		return auth.Fail(err)
	}
	var users []model.User
	db.DB.Where("hospital_id = ? OR (hospital_id IS NULL AND role IN ?)", hospitalID, []string{"super_admin", "auditor"}).
		Order("name asc").Find(&users)
	type userOut struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		Email     string     `json:"email"`
		Role      string     `json:"role"`
		RoleLabel string     `json:"roleLabel"`
		Title     *string    `json:"title"`
		Active    bool       `json:"active"`
		LastLogin *time.Time `json:"lastLogin"`
	}
	usersOut := make([]userOut, 0, len(users))
	for _, u := range users {
		label := auth.RoleLabels[u.Role]
		if label == "" {
			label = u.Role
		}
		usersOut = append(usersOut, userOut{ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role,
			RoleLabel: label, Title: u.Title, Active: u.Active, LastLogin: u.LastLogin})
	}
	type roleOut struct {
		Key         string   `json:"key"`
		Label       string   `json:"label"`
		Permissions []string `json:"permissions"`
	}
	roles := []roleOut{}
	for _, key := range []string{"super_admin", "admin_rs", "hospital_mgmt", "credential_officer",
		"medical_committee", "hr_sdm", "auditor", "claims_team", "it_simrs", "viewer"} {
		perms := auth.RolePermissions[key]
		if perms == nil {
			perms = []string{"VIEW"}
		}
		roles = append(roles, roleOut{Key: key, Label: auth.RoleLabels[key], Permissions: perms})
	}
	var rules []model.AlertRule
	db.DB.Where("hospital_id = ?", hospitalID).Order("name asc").Find(&rules)
	return auth.OK(gin.H{"hospital": hospital, "users": usersOut, "roles": roles,
		"permissions": auth.Permissions, "rules": rules})
}

type settingsPatchBody struct {
	Target           string  `json:"target"`
	ID               string  `json:"id"`
	Name             *string `json:"name"`
	Class            *string `json:"class"`
	City             *string `json:"city"`
	Province         *string `json:"province"`
	BedLicensed      *int    `json:"bedLicensed"`
	BpjsFacilityCode *string `json:"bpjsFacilityCode"`
	SatusehatOrgID   *string `json:"satusehatOrgId"`
	Active           *bool   `json:"active"`
}

func handleSettingsPatch(c *gin.Context, user *auth.SessionUser) (int, any) {
	var body settingsPatchBody
	if err := c.ShouldBindJSON(&body); err != nil {
		return 400, "Payload tidak valid."
	}
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	if err := auth.RequirePermission(user, "CONFIGURE"); err != nil {
		return auth.Fail(err)
	}
	if body.Target == "hospital" {
		var h model.Hospital
		if err := db.DB.Where("id = ?", hospitalID).First(&h).Error; err != nil {
			return 404, "RS tidak ditemukan."
		}
		updates := map[string]any{}
		if body.Name != nil {
			updates["name"] = *body.Name
		}
		if body.Class != nil {
			updates["class"] = *body.Class
		}
		if body.City != nil {
			updates["city"] = *body.City
		}
		if body.Province != nil {
			updates["province"] = *body.Province
		}
		if body.BedLicensed != nil {
			updates["bed_licensed"] = *body.BedLicensed
		}
		if body.BpjsFacilityCode != nil {
			updates["bpjs_facility_code"] = *body.BpjsFacilityCode
		}
		if body.SatusehatOrgID != nil {
			updates["satusehat_org_id"] = *body.SatusehatOrgID
		}
		if len(updates) == 0 {
			return auth.OK(h)
		}
		if err := db.DB.Model(&h).Updates(updates).Error; err != nil {
			return auth.Fail(err)
		}
		var updated model.Hospital
		db.DB.Where("id = ?", hospitalID).First(&updated)
		hid := hospitalID
		auth.LogAudit(auth.AuditOptions{
			HospitalID: &hid, Actor: user, Action: "config", EntityType: "Rumah Sakit", EntityID: &hospitalID,
			Description: "Memperbarui profil RS", Before: sp(h.Name), After: sp(updated.Name), C: c,
		})
		return auth.OK(updated)
	}
	if body.Target == "alert-rule" {
		var rule model.AlertRule
		if err := db.DB.Where("id = ? AND hospital_id = ?", body.ID, hospitalID).First(&rule).Error; err != nil {
			return 404, "Aturan tidak ditemukan."
		}
		active := !rule.Active
		if body.Active != nil {
			active = *body.Active
		}
		if err := db.DB.Model(&rule).Update("active", active).Error; err != nil {
			return auth.Fail(err)
		}
		var updated model.AlertRule
		db.DB.Where("id = ?", rule.ID).First(&updated)
		hid := hospitalID
		auth.LogAudit(auth.AuditOptions{
			HospitalID: &hid, Actor: user, Action: "config", EntityType: "Aturan Peringatan", EntityID: &rule.ID,
			Description: rule.Name + ": " + map[bool]string{true: "aktif", false: "nonaktif"}[active],
			Before: sp(boolStr(rule.Active)), After: sp(boolStr(active)), C: c,
		})
		return auth.OK(updated)
	}
	return 400, "Target tidak dikenal."
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func handleDataGet(c *gin.Context, user *auth.SessionUser) (int, any) {
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
	if typ == "quality" {
		var issues []model.DataQualityIssue
		db.DB.Where("hospital_id = ?", hospitalID).Order("status asc, severity desc").Find(&issues)
		terbuka, anomaly := 0, 0
		for _, i := range issues {
			if i.Status == "Terbuka" {
				terbuka++
			}
			if i.Anomaly {
				anomaly++
			}
		}
		return auth.OK(gin.H{"issues": issues, "summary": gin.H{"terbuka": terbuka, "total": len(issues), "anomaly": anomaly}})
	}
	if typ == "trust" {
		bands := []gin.H{}
		for _, b := range scoring.TrustBands {
			bands = append(bands, gin.H{"minTrust": b.MinTrust, "multiplier": b.Multiplier, "label": b.Label})
		}
		return auth.OK(gin.H{"trust": intel.Trust, "bands": bands, "bcds": intel.Bcds})
	}
	if typ == "sources" {
		var sources []model.DataSource
		db.DB.Where("hospital_id = ?", hospitalID).Order("name asc").Find(&sources)
		var jobs []model.IntegrationJob
		db.DB.Where("hospital_id = ?", hospitalID).Order("started_at desc").Limit(24).Find(&jobs)
		var facilities []model.Facility
		db.DB.Where("hospital_id = ?", hospitalID).Order("category asc").Find(&facilities)
		var beds []model.BedStat
		db.DB.Where("hospital_id = ?", hospitalID).Order("bed_class asc").Find(&beds)
		var hospital model.Hospital
		db.DB.Where("id = ?", hospitalID).First(&hospital)
		return auth.OK(gin.H{"sources": sources, "jobs": jobs, "facilities": facilities, "beds": beds, "hospital": hospital})
	}
	var issues []model.DataQualityIssue
	db.DB.Where("hospital_id = ?", hospitalID).Find(&issues)
	var sources []model.DataSource
	db.DB.Where("hospital_id = ?", hospitalID).Find(&sources)
	var jobs []model.IntegrationJob
	db.DB.Where("hospital_id = ?", hospitalID).Order("started_at desc").Limit(24).Find(&jobs)
	var facilities []model.Facility
	db.DB.Where("hospital_id = ?", hospitalID).Find(&facilities)
	var beds []model.BedStat
	db.DB.Where("hospital_id = ?", hospitalID).Find(&beds)
	var hospital model.Hospital
	db.DB.Where("id = ?", hospitalID).First(&hospital)
	return auth.OK(gin.H{"issues": issues, "sources": sources, "jobs": jobs, "facilities": facilities,
		"beds": beds, "hospital": hospital, "trust": intel.Trust})
}

type dataPostBody struct {
	Action string `json:"action"`
}

var fhirResources = []string{
	"Patient", "Practitioner", "PractitionerRole", "Organization", "Location", "Encounter",
	"Condition", "Observation", "Procedure", "Medication", "MedicationRequest", "MedicationDispense",
	"MedicationAdministration", "AllergyIntolerance", "DiagnosticReport", "ServiceRequest",
	"ImagingStudy", "CarePlan", "Composition", "Coverage", "Claim", "DocumentReference",
	"Immunization", "Appointment",
}

func handleDataPost(c *gin.Context, user *auth.SessionUser) (int, any) {
	var body dataPostBody
	if err := c.ShouldBindJSON(&body); err != nil {
		return 400, "Payload tidak valid."
	}
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	if body.Action != "test-connection" {
		return 400, "Aksi tidak dikenal."
	}
	var sources []model.DataSource
	db.DB.Where("hospital_id = ?", hospitalID).Find(&sources)
	var simrs *model.DataSource
	for i := range sources {
		if sources[i].Type == "SIMRS" {
			simrs = &sources[i]
			break
		}
	}
	badFields := map[string]string{
		"Composition":        "documentStatus kosong pada 12% record",
		"MedicationDispense": "whenHandedOver kosong pada 9% record",
	}
	baseFhir := 24
	coverage := 100.0
	if simrs != nil {
		baseFhir = simrs.FhirResources
		coverage = simrs.CoveragePct
	}
	type resRow struct {
		Resource string `json:"resource"`
		Status   string `json:"status"`
		Message  string `json:"message"`
		Count    int    `json:"count"`
	}
	results := []resRow{}
	okCount, perluCount := 0, 0
	for _, r := range fhirResources {
		hasIssue := false
		if _, bad := badFields[r]; bad && coverage < 100 {
			hasIssue = true
		}
		row := resRow{Resource: r, Status: "OK", Message: "Skema & field wajib valid",
			Count: baseFhir*137 + len(r)*53}
		if hasIssue {
			row.Status = "Perlu Perbaikan"
			row.Message = badFields[r]
			perluCount++
		} else {
			okCount++
		}
		results = append(results, row)
	}
	hid := hospitalID
	auth.LogAudit(auth.AuditOptions{
		HospitalID: &hid, Actor: user, Action: "config", EntityType: "Integrasi",
		Description: "Test Koneksi 24 resource FHIR dijalankan", C: c,
	})
	baseURL := "—"
	if simrs != nil && simrs.ConfigSummary != nil {
		baseURL = *simrs.ConfigSummary
	}
	return auth.OK(gin.H{
		"testedAt": time.Now().Format(time.RFC3339),
		"baseUrl":  baseURL,
		"results":  results,
		"summary":  gin.H{"ok": okCount, "perlu": perluCount, "gagal": 0},
	})
}
