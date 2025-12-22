package controllers

import (
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/repositories"
	"stok-hadiah/services"

	"github.com/gin-gonic/gin"
)

func UserIndex(c *gin.Context) {
	userRepo := &repositories.UserRepository{DB: config.DB}
	userService := &services.UserService{Repo: userRepo}

	users, err := userService.GetUsers()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "user.html", gin.H{
		"Title": "Daftar User",
		"Page":  "user",
		"users": users,
	})

}
