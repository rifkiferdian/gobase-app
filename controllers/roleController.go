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

func RoleFormIndex(c *gin.Context) {
	permissionRepo := &repositories.PermissionRepository{DB: config.DB}
	permissionService := &services.PermissionService{Repo: permissionRepo}

	permissionGroups, err := permissionService.GetGroupedPermissions()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalPermissions := 0
	for _, group := range permissionGroups {
		totalPermissions += len(group.Permissions)
	}

	Render(c, "role_form.html", gin.H{
		"Title":            "Form Role",
		"Page":             "roleForm",
		"PermissionGroups": permissionGroups,
		"TotalPermissions": totalPermissions,
	})

}
