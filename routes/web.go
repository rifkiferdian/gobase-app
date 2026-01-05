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
	auth.Use(middleware.AuthRequired(), middleware.PermissionContext())
	{
		auth.GET("/dashboard", controllers.DashboardIndex)

		auth.GET("/profile", controllers.ProfileIndex)
		auth.POST("/change-password", controllers.ChangePassword)

		auth.GET("/users", middleware.RequirePermission("user_management_access"), controllers.UserIndex)
		auth.POST("/users", middleware.RequirePermission("user_create"), controllers.UserStore)
		auth.POST("/users/update", middleware.RequirePermission("user_edit"), controllers.UserUpdate)
		auth.GET("/users/delete/:id", middleware.RequirePermission("user_delete"), controllers.UserDelete)
		auth.GET("/role", controllers.RoleIndex)
		auth.GET("/roleForm", controllers.RoleFormIndex)
		auth.GET("/role/:id/edit", middleware.RequirePermission("role_edit"), controllers.RoleEdit)
		auth.POST("/role", middleware.RequirePermission("role_create"), controllers.RoleStore)
		auth.POST("/role/update", middleware.RequirePermission("role_edit"), controllers.RoleUpdate)
		auth.GET("/role/delete/:id", middleware.RequirePermission("role_delete"), controllers.RoleDelete)

		auth.GET("/suppliers", controllers.SupplierIndex)
		auth.POST("/suppliers", controllers.SupplierStore)
		auth.POST("/suppliers/update", controllers.SupplierUpdate)
		auth.GET("/suppliers/delete/:id", controllers.SupplierDelete)

		auth.GET("/items", controllers.ItemIndex)
		auth.POST("/items", controllers.ItemStore)
		auth.POST("/items/update", controllers.ItemUpdate)
		auth.GET("/items/delete/:id", controllers.ItemDelete)
		auth.GET("/items-report/:id", controllers.ItemReportIndex)

		auth.GET("/item-out", controllers.ItemOutIndex)
		auth.POST("/item-out/update", controllers.ItemOutUpdate)
		auth.POST("/item-out/case", controllers.ItemOutCaseStore)
		auth.DELETE("/item-out/case/:id", controllers.ItemOutCaseDelete)

		auth.GET("/item-in", controllers.ItemInIndex)
		auth.POST("/item-in", controllers.ItemInStore)
		auth.POST("/item-in/update", controllers.ItemInUpdate)
		auth.GET("/item-in/delete/:id", controllers.ItemInDelete)

		auth.GET("/programs", controllers.ProgramIndex)
		auth.POST("/programs", controllers.ProgramStore)
		auth.POST("/programs/update", controllers.ProgramUpdate)
		auth.GET("/programs/delete/:id", controllers.ProgramDelete)

		auth.GET("/stock-report", controllers.StockReportIndex)

		auth.GET("/home", controllers.HomeIndex) // contoh tambahan route
	}
}
