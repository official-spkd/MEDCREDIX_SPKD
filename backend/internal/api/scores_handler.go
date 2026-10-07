package api

// GET /api/scores?domain=key — WHY panel: kalkulasi + kontributor + bukti + risiko
// + rekomendasi + aksi + tren. Port dari src/app/api/scores/route.ts.
import (
	"sort"
	"strings"

	"github.com/gin-gonic/gin"

	"medcredix/backend/internal/auth"
	"medcredix/backend/internal/db"
	"medcredix/backend/internal/model"
	"medcredix/backend/internal/scoring"
)

// domainLabels — peta domain → label (port dari lib/format.ts DOMAIN_LABELS).
var domainLabels = map[string]string{
	"workforce": "Clinical Workforce", "credentialing": "Credentialing Readiness",
	"digital": "Digital Maturity", "medication": "Medication Management",
	"diagnostic": "Diagnostic Readiness", "interop": "Interoperability",
	"claim": "Claim Integrity", "pathway": "Clinical Pathway",
	"safety": "Patient Safety", "trust": "Data Trust",
}

type stageTableItem struct {
	Min   float64 `json:"min"`
	Stage int     `json:"stage"`
}

type whyDomainRow struct {
	Key          string  `json:"key"`
	Label        string  `json:"label"`
	Score        float64 `json:"score"`
	Weight       float64 `json:"weight"`
	Contribution float64 `json:"contribution"`
}

func handleScores(c *gin.Context, user *auth.SessionUser) (int, any) {
	hospitalID, err := auth.ResolveHospitalID(user, c.Query("hospitalId"))
	if err != nil {
		return auth.Fail(err)
	}
	domainKey := c.Query("domain")
	if domainKey == "" {
		domainKey = "health"
	}
	intel, err := ComputeIntel(hospitalID)
	if err != nil {
		return auth.Fail(err)
	}

	if domainKey == "health" {
		var risks []model.Risk
		db.DB.Where("hospital_id = ? AND status != ?", hospitalID, "Ditutup").Limit(6).Find(&risks)
		var evidence []model.EvidenceItem
		db.DB.Where("hospital_id = ? AND ref_type = ?", hospitalID, "Skor").Limit(10).Find(&evidence)
		var actions []model.Action
		db.DB.Where("hospital_id = ? AND status IN ?", hospitalID, []string{"Baru", "Ditugaskan", "Dalam Pengerjaan"}).Limit(6).Find(&actions)

		domains := make([]whyDomainRow, 0, len(intel.Domains))
		for _, d := range intel.Domains {
			domains = append(domains, whyDomainRow{
				Key: d.Key, Label: d.Label, Score: d.Score, Weight: d.Weight,
				Contribution: round1Of(d.Score * d.Weight),
			})
		}
		stages := make([]stageTableItem, 0, 8)
		for _, min := range []float64{90, 80, 66, 60, 50, 35, 20, 0} {
			stages = append(stages, stageTableItem{Min: min, Stage: scoring.StageOf(min)})
		}
		// rekomendasi dari domain terlemah
		sorted := append([]scoring.Domain{}, intel.Domains...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].Score < sorted[j].Score })
		recs := []scoring.Recommendation{}
		if len(sorted) > 0 {
			recs = scoring.DomainRecommendations(&sorted[0], intel)
			if len(recs) > 3 {
				recs = recs[:3]
			}
		}
		return auth.OK(gin.H{
			"kind": "health",
			"formula": "Hospital Health Score = Σ(skor domain × bobot domain). Bobot: Claim Integrity 16%, Credentialing 14%, Patient Safety 14%, Clinical Workforce 12%, Clinical Pathway 10%, Digital Maturity 10%, Medication 8%, Diagnostik 8%, Interop 8%.",
			"formulaVersion": scoring.FormulaVersion,
			"domains":        domains,
			"healthScore":    intel.HealthScore, "bcds": intel.Bcds, "trust": intel.Trust, "kelompok": intel.Kelompok,
			"stageTable": stages,
			"risks": risks, "evidence": evidence, "actions": actions, "trend": intel.Trend,
			"recommendations": recs,
		})
	}

	if domainKey == "bcds" {
		bands := make([]gin.H, 0, len(scoring.TrustBands))
		for _, b := range scoring.TrustBands {
			bands = append(bands, gin.H{"minTrust": b.MinTrust, "multiplier": b.Multiplier, "label": b.Label})
		}
		return auth.OK(gin.H{
			"kind":    "bcds",
			"formula": "BPJS Credential Readiness (BDCS) = (Kelompok A × 25%) + (B × 30%) + (C × 35%) + (D × 10%), lalu dikalikan Trust Multiplier — ditampilkan transparan.",
			"kelompok": intel.Kelompok, "bcds": intel.Bcds, "trust": intel.Trust,
			"trustBands": bands,
			"stage": intel.Bcds.Stage, "lulus": intel.Bcds.Lulus, "prereqOk": intel.Bcds.PrereqOk,
			"trend": intel.Trend, "formulaVersion": scoring.FormulaVersion,
		})
	}

	var domain *scoring.Domain
	for i := range intel.Domains {
		if intel.Domains[i].Key == domainKey {
			domain = &intel.Domains[i]
			break
		}
	}
	if domain == nil {
		return 404, "Domain tidak dikenal."
	}

	var risks []model.Risk
	db.DB.Where("hospital_id = ? AND status != ?", hospitalID, "Ditutup").Find(&risks)
	var evidence []model.EvidenceItem
	db.DB.Where("hospital_id = ? AND ref_id = ?", hospitalID, "DOMAIN:"+domainKey).Find(&evidence)
	var actions []model.Action
	db.DB.Where("hospital_id = ?", hospitalID).Find(&actions)

	var relatedRisks []model.Risk
	dl := domainLabels[domainKey]
	for _, r := range risks {
		if r.Domain == dl || strings.Contains(strings.ToLower(r.Domain), strings.ToLower(domainKey)) {
			relatedRisks = append(relatedRisks, r)
		}
	}
	var relatedActions []model.Action
	for _, a := range actions {
		if a.SourceRef != nil && *a.SourceRef == "DOMAIN:"+domainKey {
			relatedActions = append(relatedActions, a)
			continue
		}
		for _, r := range relatedRisks {
			if a.SourceRef != nil && *a.SourceRef == r.Code {
				relatedActions = append(relatedActions, a)
				break
			}
		}
	}
	if evidence == nil {
		evidence = []model.EvidenceItem{}
	}

	return auth.OK(gin.H{
		"kind": "domain", "domain": domain, "formulaVersion": scoring.FormulaVersion,
		"evidence": evidence, "risks": relatedRisks, "actions": relatedActions,
		"recommendations": scoring.DomainRecommendations(domain, intel),
		"trend":           intel.Trend, "healthScore": intel.HealthScore,
	})
}

func round1Of(v float64) float64 {
	return float64(int(v*10+0.5*sign(v))) / 10
}

func sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}
