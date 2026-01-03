package controllers

import (
	"database/sql"
	"net/http"
	"strings"

	"stok-hadiah/config"
	"stok-hadiah/models"
	"stok-hadiah/repositories"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func renderProfilePage(c *gin.Context, data gin.H) {
	allowedStoreIDs := getAllowedStoreIDs(c)

	storeRepo := &repositories.StoreRepository{DB: config.DB}
	storeDisplay := "-"
	if stores, err := storeRepo.GetByIDs(allowedStoreIDs); err == nil {
		storeNames := make([]string, 0, len(stores))
		for _, s := range stores {
			storeNames = append(storeNames, s.StoreName)
		}
		if len(storeNames) > 0 {
			storeDisplay = strings.Join(storeNames, ", ")
		}
	}

	if data == nil {
		data = gin.H{}
	}

	data["StoreDisplay"] = storeDisplay
	data["Title"] = "Profile Page"
	data["Page"] = "profile"

	Render(c, "profile.html", data)
}

func ProfileIndex(c *gin.Context) {
	renderProfilePage(c, gin.H{})
}

// ChangePassword handles password updates for the logged-in user.
func ChangePassword(c *gin.Context) {
	session := sessions.Default(c)

	userID := normalizeSessionID(session.Get("user_id"))
	if userID == 0 {
		if u := session.Get("user"); u != nil {
			switch val := u.(type) {
			case models.SessionUser:
				userID = val.UserID
			case map[string]interface{}:
				if id, ok := val["user_id"]; ok {
					userID = normalizeSessionID(id)
				} else if id, ok := val["UserID"]; ok {
					userID = normalizeSessionID(id)
				}
			case gin.H:
				if id, ok := val["user_id"]; ok {
					userID = normalizeSessionID(id)
				} else if id, ok := val["UserID"]; ok {
					userID = normalizeSessionID(id)
				}
			}
		}
	}

	if userID == 0 {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	currentPassword := strings.TrimSpace(c.PostForm("current_password"))
	newPassword := strings.TrimSpace(c.PostForm("password"))
	confirmPassword := strings.TrimSpace(c.PostForm("password_confirmation"))

	if currentPassword == "" || newPassword == "" || confirmPassword == "" {
		renderProfilePage(c, gin.H{"Error": "Semua kolom wajib diisi"})
		return
	}

	if newPassword != confirmPassword {
		renderProfilePage(c, gin.H{"Error": "Konfirmasi password baru tidak sama"})
		return
	}

	if len(newPassword) < 6 {
		renderProfilePage(c, gin.H{"Error": "Password baru minimal 6 karakter"})
		return
	}

	var storedHash string
	err := config.DB.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&storedHash)
	if err == sql.ErrNoRows {
		renderProfilePage(c, gin.H{"Error": "Pengguna tidak ditemukan"})
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "Terjadi kesalahan saat mengambil data user")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(currentPassword)) != nil {
		renderProfilePage(c, gin.H{"Error": "Password saat ini tidak sesuai"})
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Gagal mengenkripsi password baru")
		return
	}

	if _, err := config.DB.Exec("UPDATE users SET password = ? WHERE id = ?", string(newHash), userID); err != nil {
		c.String(http.StatusInternalServerError, "Gagal memperbarui password")
		return
	}

	renderProfilePage(c, gin.H{"Success": "Password berhasil diperbarui"})
}

func normalizeSessionID(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}
