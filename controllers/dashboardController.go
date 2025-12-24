package controllers

import (
	"github.com/gin-gonic/gin"
)

func DashboardIndex(c *gin.Context) {

	// debug session
	// sess := sessions.Default(c)
	// fmt.Println("DEBUG user_id:", sess.Get("user_id"))
	// fmt.Println("DEBUG user:", sess.Get("user"))

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
