package controllers

import (
	"github.com/gin-gonic/gin"
)

func ItemIndex(c *gin.Context) {
	Render(c, "item/index.html", gin.H{
		"Title": "Item Page",
		"Page":  "item",
	})
}
