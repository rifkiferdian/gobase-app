package controllers

import (
	"github.com/gin-gonic/gin"
)

func DashboardIndex(c *gin.Context) {
	Render(c, "base.html", gin.H{
		"Title": "Dashboard User",
	})
}
