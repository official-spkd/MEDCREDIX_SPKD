// Package auth — autentikasi & RBAC server-side (port dari lib/auth.ts).
// Demo: header X-User-Id; produksi: JWT/OAuth2 (siap MFA).
package auth

import (
        "crypto/sha256"
        "encoding/hex"
        "errors"
        "fmt"
        "math/rand"
        "net/http"
        "strings"
        "time"

        "github.com/gin-gonic/gin"
        "gorm.io/gorm"

        "medcredix/backend/internal/db"
        "medcredix/backend/internal/model"
)

// Permission — daftar izin platform.
var Permissions = []string{
        "VIEW", "CREATE", "EDIT", "DELETE", "VERIFY", "APPROVE", "SIGN",
        "EXPORT", "PRINT", "ASSIGN", "RESOLVE", "REOPEN", "CONFIGURE",
}

// RoleLabels — label tampilan per peran.
var RoleLabels = map[string]string{
        "super_admin":        "SPKD Super Admin",
        "admin_rs":           "Admin RS",
        "hospital_mgmt":      "Manajemen RS",
        "credential_officer": "Petugas Kredensial",
        "medical_committee":  "Komite Medis",
        "hr_sdm":             "HR / SDM",
        "auditor":            "Auditor",
        "claims_team":        "Tim Klaim",
        "it_simrs":           "IT / SIMRS",
        "viewer":             "Viewer",
}

// RolePermissions — matriks RBAC (server-enforced).
var RolePermissions = map[string][]string{
        "super_admin":        Permissions,
        "admin_rs":           Permissions,
        "hospital_mgmt":      {"VIEW", "EXPORT", "PRINT", "ASSIGN"},
        "credential_officer": {"VIEW", "CREATE", "EDIT", "VERIFY", "ASSIGN", "EXPORT", "PRINT"},
        "medical_committee":  {"VIEW", "VERIFY", "APPROVE", "SIGN", "REOPEN", "RESOLVE", "EXPORT"},
        "hr_sdm":             {"VIEW", "CREATE", "EDIT", "EXPORT"},
        "auditor":            {"VIEW", "EXPORT", "PRINT"},
        "claims_team":        {"VIEW", "CREATE", "EDIT", "ASSIGN", "RESOLVE", "EXPORT"},
        "it_simrs":           {"VIEW", "CREATE", "EDIT", "CONFIGURE", "EXPORT"},
        "viewer":             {"VIEW"},
}

// SessionUser — bentuk sesi yang dikirim ke klien.
type SessionUser struct {
        ID           string   `json:"id"`
        Name         string   `json:"name"`
        Email        string   `json:"email"`
        Role         string   `json:"role"`
        RoleLabel    string   `json:"roleLabel"`
        Title        *string  `json:"title"`
        HospitalID   *string  `json:"hospitalId"`
        Permissions  []string `json:"permissions"`
}

// HashPassword — SHA-256+salt (demo; produksi: bcrypt/argon2).
func HashPassword(pw string) string {
        sum := sha256.Sum256([]byte("medcredix::" + pw))
        return hex.EncodeToString(sum[:])
}

// ApiError — error HTTP bertipe status.
type ApiError struct {
        Status  int
        Message string
}

func (e *ApiError) Error() string { return e.Message }

func Err(status int, msg string) *ApiError { return &ApiError{Status: status, Message: msg} }

// GetUser — resolusi sesi dari header X-User-Id atau query userId (fallback ekspor/unduhan).
func GetUser(c *gin.Context) (*SessionUser, *model.User, error) {
        userID := c.GetHeader("X-User-Id")
        if userID == "" {
                userID = c.Query("userId")
        }
        if userID == "" {
                return nil, nil, errors.New("no session")
        }
        var u model.User
        if err := db.DB.Where("id = ?", userID).First(&u).Error; err != nil {
                return nil, nil, errors.New("no session")
        }
        if !u.Active {
                return nil, nil, errors.New("no session")
        }
        perms, ok := RolePermissions[u.Role]
        if !ok {
                perms = []string{"VIEW"}
        }
        label, ok := RoleLabels[u.Role]
        if !ok {
                label = u.Role
        }
        return &SessionUser{
                ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role,
                RoleLabel: label, Title: u.Title, HospitalID: u.HospitalID, Permissions: perms,
        }, &u, nil
}

// ResolveHospitalID — resolusi tenant dengan isolasi lintas RS.
func ResolveHospitalID(user *SessionUser, requested string) (string, error) {
        if user.Role == "super_admin" {
                if requested != "" {
                        return requested, nil
                }
                return "h_nusantara", nil
        }
        if user.HospitalID == nil || *user.HospitalID == "" {
                return "", Err(403, "Akun Anda tidak terikat ke rumah sakit mana pun.")
        }
        if requested != "" && requested != *user.HospitalID {
                return "", Err(403, "Akses lintas rumah sakit tidak diizinkan (isolasi tenant).")
        }
        return *user.HospitalID, nil
}

// HasPermission — cek izin peran.
func HasPermission(user *SessionUser, perm string) bool {
        for _, p := range user.Permissions {
                if p == perm {
                        return true
                }
        }
        return false
}

// RequirePermission — panic-free: kembalikan error bila izin tidak dimiliki.
func RequirePermission(user *SessionUser, perm string) error {
        if !HasPermission(user, perm) {
                return Err(403, fmt.Sprintf("Peran %s tidak memiliki izin %s.", user.RoleLabel, perm))
        }
        return nil
}

// AuditOptions — masukan LogAudit.
type AuditOptions struct {
        HospitalID  *string
        Actor       *SessionUser // nil → system
        Action      string
        EntityType  string
        EntityID    *string
        Description string
        Before      *string
        After       *string
        C           *gin.Context
}

// LogAudit — catat aktivitas ke AuditLog (immutable).
func LogAudit(opts AuditOptions) {
        actorName, actorRole := "MedCredix Engine", "system"
        if opts.Actor != nil {
                actorName, actorRole = opts.Actor.Name, opts.Actor.Role
        }
        ip := "10.20.0.1"
        if opts.C != nil {
                if fwd := opts.C.GetHeader("X-Forwarded-For"); fwd != "" {
                        ip = strings.TrimSpace(fwd)
                }
        }
        corr := fmt.Sprintf("corr-%x-%04x", time.Now().UnixMilli(), rand.Intn(0xffff))
        var entityID *string
        if opts.EntityID != nil {
                v := *opts.EntityID
                entityID = &v
        }
        entry := model.AuditLog{
                HospitalID: opts.HospitalID, ActorName: actorName, ActorRole: actorRole,
                Action: opts.Action, EntityType: opts.EntityType, EntityID: entityID,
                Description: opts.Description, Before: opts.Before, After: opts.After,
                IP: &ip, CorrelationID: &corr,
        }
        if err := db.DB.Create(&entry).Error; err != nil {
                // audit tidak boleh memblokir operasi utama; cukup dicatat di log
                fmt.Println("[audit] gagal menyimpan log:", err)
        }
}

// WithUser — middleware Gin: resolusi sesi + penanganan error seragam.
type UserHandler func(user *SessionUser) (int, any)

// Wrap mengubah UserHandler menjadi gin.HandlerFunc dengan envelope {ok,data}/{ok,error}.
func Wrap(fn func(c *gin.Context, user *SessionUser) (int, any)) gin.HandlerFunc {
        return func(c *gin.Context) {
                user, _, err := GetUser(c)
                if err != nil {
                        c.JSON(401, gin.H{"ok": false, "error": "Sesi tidak valid. Silakan masuk kembali."})
                        return
                }
                status, payload := fn(c, user)
                if c.Writer.Written() {
                        return // handler sudah menulis respons langsung (mis. unduhan CSV)
                }
                if status >= 400 {
                        c.AbortWithStatusJSON(status, gin.H{"ok": false, "error": payload})
                        return
                }
                c.JSON(status, gin.H{"ok": true, "data": payload})
        }
}

// OK — helper payload sukses.
func OK(data any) (int, any) { return http.StatusOK, data }

// Fail — helper payload error bertipe *ApiError.
func Fail(err error) (int, any) {
        var apiErr *ApiError
        if errors.As(err, &apiErr) {
                return apiErr.Status, apiErr.Message
        }
        if errors.Is(err, gorm.ErrRecordNotFound) {
                return 404, "Data tidak ditemukan."
        }
        fmt.Println("[api] error:", err)
        return 500, "Terjadi kesalahan internal."
}
