package routes

import (
	"stok-hadiah/controllers"
	"stok-hadiah/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWebRoutes(r *gin.Engine) {
	r.Use(middleware.UserMiddleware())

	r.GET("/", controllers.LoginPage)
	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.LoginPost)
	r.POST("/register", controllers.CreateUser)
	r.GET("/logout", controllers.Logout)

	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/dashboard", controllers.DashboardIndex)
	}
}
