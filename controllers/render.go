package controllers

import (
	"net/http"
	"stok-hadiah/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, name string, data gin.H) {

	session := sessions.Default(c)
	if u := session.Get("user"); u != nil {
		switch val := u.(type) {
		case models.SessionUser:
			data["User"] = gin.H{
				"user_id":          val.UserID,
				"nip":              val.NIP,
				"name":             val.Name,
				"username":         val.Username,
				"role":             val.Role,
				"is_authenticated": val.IsAuthenticated,
			}
		case map[string]interface{}:
			data["User"] = val
		case gin.H:
			data["User"] = map[string]interface{}(val)
		default:
			data["User"] = gin.H{
				"name":     u,
				"username": u,
			}
		}
	}

	// ambil permissions dari context
	permsAny, _ := c.Get("Permissions")
	perms, _ := permsAny.(map[string]bool)

	if data == nil {
		data = gin.H{}
	}

	// inject global data (biar semua halaman dapat)
	data["Permissions"] = perms

	c.HTML(http.StatusOK, name, data)
}
