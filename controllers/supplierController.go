package controllers

import (
	"github.com/gin-gonic/gin"
)

func SupplierIndex(c *gin.Context) {
	Render(c, "supplier/index.html", gin.H{
		"Title": "Supplier Page",
		"Page":  "supplier",
	})
}
