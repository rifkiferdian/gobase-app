package controllers

import (
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/repositories"
	"stok-hadiah/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StockReportIndex(c *gin.Context) {
	filterName := c.Query("item_name")
	filterCategory := c.Query("category")

	supplierStr := c.Query("supplier_id")
	var supplierID *int
	selectedSupplier := 0
	if id, err := strconv.Atoi(supplierStr); err == nil && id > 0 {
		supplierID = &id
		selectedSupplier = id
	}

	storeStr := c.Query("store_id")
	var filterStoreIDPtr *int
	filterStoreID := 0
	if id, err := strconv.Atoi(storeStr); err == nil && id > 0 {
		filterStoreID = id
		filterStoreIDPtr = &filterStoreID
	}

	allowedStoreIDs := getAllowedStoreIDs(c)

	stockRepo := &repositories.ItemStockRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		FilterStoreID:      filterStoreIDPtr,
		EnforceStoreFilter: true,
	}
	stockService := &services.ItemStockService{Repo: stockRepo}

	items, err := stockService.ListSummaries(filterName, filterCategory, supplierID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalIn := 0
	totalOut := 0
	totalRemaining := 0
	for _, it := range items {
		totalIn += it.QtyIn
		totalOut += it.QtyOut
		totalRemaining += it.Remaining
	}

	supplierRepo := &repositories.SupplierRepository{DB: config.DB}
	suppliers, err := supplierRepo.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	storeRepo := &repositories.StoreRepository{DB: config.DB}
	stores, _ := storeRepo.GetByIDs(allowedStoreIDs)

	Render(c, "stock_report.html", gin.H{
		"Title":            "Stock Report",
		"Page":             "stock_report",
		"items":            items,
		"suppliers":        suppliers,
		"stores":           stores,
		"FilterItemName":   filterName,
		"FilterCategory":   filterCategory,
		"FilterSupplierID": selectedSupplier,
		"FilterStoreID":    filterStoreID,
		"TotalQtyIn":       totalIn,
		"TotalQtyOut":      totalOut,
		"TotalRemaining":   totalRemaining,
	})
}
