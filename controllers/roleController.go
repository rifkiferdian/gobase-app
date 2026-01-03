package controllers

import (
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/models"
	"stok-hadiah/repositories"
	"stok-hadiah/services"
	"strconv"
	"strings"

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
	renderRoleForm(c, "")
}

// RoleStore menangani penyimpanan role baru dari form.
func RoleStore(c *gin.Context) {
	type roleForm struct {
		Name      string `form:"name" binding:"required"`
		GuardName string `form:"guard_name"`
	}

	var form roleForm
	if err := c.ShouldBind(&form); err != nil {
		renderRoleForm(c, "Form tidak lengkap")
		return
	}

	var permissionIDs []int64
	for _, val := range c.PostFormArray("permissions") {
		if strings.TrimSpace(val) == "" {
			continue
		}
		id, err := strconv.ParseInt(val, 10, 64)
		if err != nil || id <= 0 {
			renderRoleForm(c, "Permission tidak valid")
			return
		}
		permissionIDs = append(permissionIDs, id)
	}

	roleRepo := &repositories.RoleRepository{DB: config.DB}
	roleService := &services.RoleService{Repo: roleRepo}

	input := models.RoleCreateInput{
		Name:          form.Name,
		GuardName:     form.GuardName,
		PermissionIDs: permissionIDs,
	}

	if err := roleService.CreateRole(input); err != nil {
		renderRoleForm(c, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/role")
}

// RoleDelete menghapus role berdasarkan ID.
func RoleDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid role id")
		return
	}

	roleRepo := &repositories.RoleRepository{DB: config.DB}
	roleService := &services.RoleService{Repo: roleRepo}

	if err := roleService.DeleteRole(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/role")
}

func renderRoleForm(c *gin.Context, message string) {
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
		"Error":            message,
	})

}
