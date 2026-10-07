package api

// POST /api/auth — login (email + password) → sesi demo (X-User-Id).
// Port dari src/app/api/auth/route.ts.
import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"medcredix/backend/internal/auth"
	"medcredix/backend/internal/db"
	"medcredix/backend/internal/model"
)

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func handleAuthLogin(c *gin.Context) {
	var body loginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "Email dan kata sandi wajib diisi."})
		return
	}
	email := strings.ToLower(strings.TrimSpace(body.Email))
	password := body.Password
	if email == "" || password == "" {
		c.JSON(400, gin.H{"ok": false, "error": "Email dan kata sandi wajib diisi."})
		return
	}
	var user model.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil || user.PasswordHash != auth.HashPassword(password) {
		c.JSON(401, gin.H{"ok": false, "error": "Email atau kata sandi salah."})
		return
	}
	if !user.Active {
		c.JSON(403, gin.H{"ok": false, "error": "Akun tidak aktif. Hubungi Super Admin."})
		return
	}
	now := time.Now()
	db.DB.Model(&user).Update("last_login", now)
	su := &auth.SessionUser{ID: user.ID, Name: "", Email: user.Email, Role: user.Role}
	auth.LogAudit(auth.AuditOptions{
		HospitalID: user.HospitalID, Actor: su, Action: "login", EntityType: "Sesi",
		Description: "Login berhasil", C: c,
	})
	perms, ok := auth.RolePermissions[user.Role]
	if !ok {
		perms = []string{"VIEW"}
	}
	label := auth.RoleLabels[user.Role]
	if label == "" {
		label = user.Role
	}
	c.JSON(200, gin.H{"ok": true, "data": gin.H{
		"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role,
		"roleLabel": label, "title": user.Title, "hospitalId": user.HospitalID,
		"permissions": perms,
	}})
}
