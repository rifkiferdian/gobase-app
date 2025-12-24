package controllers

import (
	"github.com/gin-gonic/gin"
)

func ItemOutIndex(c *gin.Context) {

	Render(c, "item_out.html", gin.H{
		"Title": "Item Out",
		"Page":  "itemout",
	})

}
