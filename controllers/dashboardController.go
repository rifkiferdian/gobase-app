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
