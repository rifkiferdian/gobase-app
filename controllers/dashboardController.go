package controllers

import (
	"github.com/gin-gonic/gin"
)

func DashboardIndex(c *gin.Context) {
	Render(c, "dashboard/index.html", gin.H{
		"Title": "Dashboard User",
		"Page":  "dashboard",
	})

}

func HomeIndex(c *gin.Context) {
	Render(c, "home.html", gin.H{
		"Title": "Home Page",
	})
}
