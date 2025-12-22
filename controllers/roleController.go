package controllers

import (
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/repositories"
	"stok-hadiah/services"

	"github.com/gin-gonic/gin"
)

func RoleIndex(c *gin.Context) {

	roleRepo := &repositories.RoleRepository{DB: config.DB}
	roleService := &services.RoleService{Repo: roleRepo}

	roles, err := roleService.GetRoles()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "role.html", gin.H{
		"Title": "Daftar Role",
		"Page":  "role",
		"roles": roles,
	})

}
