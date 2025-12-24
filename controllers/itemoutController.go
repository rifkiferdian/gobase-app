package controllers

import (
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/repositories"
	"stok-hadiah/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ItemOutIndex(c *gin.Context) {
	filterName := c.Query("item_name")
	filterCategory := c.Query("category")

	supplierStr := c.Query("supplier_id")
	var supplierID *int
	if id, err := strconv.Atoi(supplierStr); err == nil && id > 0 {
		supplierID = &id
	}

	allowedStoreIDs := getAllowedStoreIDs(c)

	itemStockRepo := &repositories.ItemStockRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}
	itemStockService := &services.ItemStockService{Repo: itemStockRepo}

	supplierRepo := &repositories.SupplierRepository{DB: config.DB}
	suppliers, err := supplierRepo.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	items, err := itemStockService.ListSummaries(filterName, filterCategory, supplierID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	selectedSupplier := 0
	if supplierID != nil {
		selectedSupplier = *supplierID
	}

	Render(c, "item_out.html", gin.H{
		"Title":            "Item Out",
		"Page":             "item_out",
		"items":            items,
		"suppliers":        suppliers,
		"FilterItemName":   filterName,
		"FilterCategory":   filterCategory,
		"FilterSupplierID": selectedSupplier,
	})
}
