package controllers

import (
	"github.com/gin-gonic/gin"
)

func ProgramIndex(c *gin.Context) {
	Render(c, "program/index.html", gin.H{
		"Title": "Program Page",
		"Page":  "program",
	})
}
