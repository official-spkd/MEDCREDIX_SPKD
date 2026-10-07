// Package api — REST API MedCredix (port dari src/app/api/* Next.js routes).
// Envelope respons identik: { ok: true, data } | { ok: false, error }.
package api

import (
        "net/http"
        "os"
        "strings"

        "github.com/gin-gonic/gin"
        "gorm.io/gorm"

        "medcredix/backend/internal/auth"
        "medcredix/backend/internal/db"
        "medcredix/backend/internal/model"
        "medcredix/backend/internal/scoring"
        "medcredix/backend/internal/seed"
)

// GormLoader — implementasi scoring.Loader dengan GORM.
type GormLoader struct{}

func (GormLoader) Load(hospitalID string) (scoring.IntelData, error) {
        var d scoring.IntelData
        g := db.DB
        err := g.Where("hospital_id = ?", hospitalID).Find(&d.Practitioners).Error
        err = firstErr(err,
                g.Where("hospital_id = ?", hospitalID).Find(&d.Credentials).Error,
                g.Where("hospital_id = ?", hospitalID).Find(&d.Stds).Error,
                g.Where("hospital_id = ?", hospitalID).Find(&d.Facilities).Error,
                g.Where("hospital_id = ?", hospitalID).Find(&d.Beds).Error,
                g.Where("hospital_id = ?", hospitalID).Find(&d.DigInd).Error,
                g.Where("hospital_id = ?", hospitalID).Find(&d.Claims).Error,
                g.Where("hospital_id = ?", hospitalID).Find(&d.Pathways).Error,
                g.Where("hospital_id = ?", hospitalID).Find(&d.Episodes).Error,
                g.Order("month asc").Where("hospital_id = ?", hospitalID).Find(&d.Utilizations).Error,
                g.Order("month asc").Where("hospital_id = ?", hospitalID).Find(&d.Complaints).Error,
                g.Where("hospital_id = ?", hospitalID).Find(&d.Trusts).Error,
                g.Where("hospital_id = ?", hospitalID).Find(&d.DqIssues).Error,
                g.Where("hospital_id = ?", hospitalID).Find(&d.Requirements).Error,
        )
        if err != nil {
                return d, err
        }
        var hospital model.Hospital
        if err := g.Where("id = ?", hospitalID).First(&hospital).Error; err == nil {
                d.Hospital = &hospital
        }
        if err := g.Order("period asc").Where("hospital_id = ?", hospitalID).Find(&d.Snapshots).Error; err != nil {
                return d, err
        }
        return d, nil
}

func firstErr(errs ...error) error {
        for _, e := range errs {
                if e != nil {
                        return e
                }
        }
        return nil
}

// ComputeIntel — pintasan hitung intelligence untuk satu rumah sakit.
func ComputeIntel(hospitalID string) (*scoring.Intelligence, error) {
        return scoring.ComputeIntelligence(GormLoader{}, hospitalID)
}

// NewRouter — bangun seluruh rute API.
func NewRouter() *gin.Engine {
        if os.Getenv("GIN_MODE") == "" {
                gin.SetMode(gin.ReleaseMode)
        }
        r := gin.New()
        r.Use(gin.LoggerWithWriter(nil), gin.Recovery())
        r.Use(corsMiddleware())

        r.GET("/api", func(c *gin.Context) {
                c.JSON(http.StatusOK, gin.H{"message": "SPKD MedCredix API — From Hospital Data to Action. (Golang)"})
        })

        api := r.Group("/api")
        {
                api.POST("/auth", handleAuthLogin)

                api.GET("/bootstrap", auth.Wrap(handleBootstrap))
                api.GET("/dashboard", auth.Wrap(handleDashboard))
                api.GET("/scores", auth.Wrap(handleScores))

                api.GET("/risks", auth.Wrap(handleRisksList))
                api.PATCH("/risks", auth.Wrap(handleRiskPatch))

                api.GET("/warnings", auth.Wrap(handleWarningsList))
                api.PATCH("/warnings", auth.Wrap(handleWarningPatch))

                api.GET("/actions", auth.Wrap(handleActionsList))
                api.POST("/actions", auth.Wrap(handleActionCreate))
                api.PATCH("/actions", auth.Wrap(handleActionPatch))

                api.GET("/practitioners", auth.Wrap(handlePractitioners))

                api.GET("/credentials", auth.Wrap(handleCredentialsList))
                api.PATCH("/credentials", auth.Wrap(handleCredentialPatch))

                api.GET("/claims", auth.Wrap(handleClaimsList))
                api.PATCH("/claims", auth.Wrap(handleClaimPatch))

                api.GET("/clinical", auth.Wrap(handleClinical))

                api.GET("/compliance", auth.Wrap(handleCompliance))
                api.GET("/evidence", auth.Wrap(handleEvidence))
                api.GET("/reports", auth.Wrap(handleReportsList))
                api.POST("/reports", auth.Wrap(handleReportGenerate))
                api.GET("/search", auth.Wrap(handleSearch))

                api.GET("/requirements", auth.Wrap(handleRequirementsList))
                api.PATCH("/requirements", auth.Wrap(handleRequirementPatch))

                api.GET("/audit", auth.Wrap(handleAudit))
                api.GET("/admin", auth.Wrap(handleAdmin))
                api.GET("/settings", auth.Wrap(handleSettingsGet))
                api.PATCH("/settings", auth.Wrap(handleSettingsPatch))

                api.GET("/data", auth.Wrap(handleDataGet))
                api.POST("/data", auth.Wrap(handleDataPost))
        }

        // serve frontend build bila tersedia (mode produksi single-binary)
        staticDir := os.Getenv("MEDCREDIX_STATIC")
        if staticDir != "" {
                r.NoRoute(func(c *gin.Context) {
                        path := c.Request.URL.Path
                        if strings.HasPrefix(path, "/api") {
                                c.JSON(404, gin.H{"ok": false, "error": "Endpoint tidak ditemukan."})
                                return
                        }
                        c.File(staticDir + path)
                })
                r.GET("/", func(c *gin.Context) { c.File(staticDir + "/index.html") })
        }
        return r
}

// SeedIfEmpty — seed data demo saat database kosong.
func SeedIfEmpty(g *gorm.DB) error {
        return seed.Run(g)
}

func corsMiddleware() gin.HandlerFunc {
        return func(c *gin.Context) {
                c.Header("Access-Control-Allow-Origin", "*")
                c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
                c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Id")
                if c.Request.Method == http.MethodOptions {
                        c.AbortWithStatus(http.StatusNoContent)
                        return
                }
                c.Next()
        }
}
