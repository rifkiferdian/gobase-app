package controllers

import (
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/repositories"
	"stok-hadiah/services"

	"github.com/gin-gonic/gin"
)

type SupplierController struct {
	Service *services.SupplierService
}

func SupplierIndex(c *gin.Context) {
	repo := &repositories.SupplierRepository{DB: config.DB}
	service := &services.SupplierService{Repo: repo}

	data, err := service.GetSuppliers()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "supplier/index.html", gin.H{
		"Title":     "Supplier Page",
		"Page":      "supplier",
		"suppliers": data,
	})
}
