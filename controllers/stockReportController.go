package controllers

import (
	"github.com/gin-gonic/gin"
)

func StockReportIndex(c *gin.Context) {

	// debug session
	// sess := sessions.Default(c)
	// fmt.Println("DEBUG user_id:", sess.Get("user_id"))
	// fmt.Println("DEBUG user:", sess.Get("user"))

	Render(c, "stock_report.html", gin.H{
		"Title": "Dashboard User",
		"Page":  "dashboard",
	})

}
