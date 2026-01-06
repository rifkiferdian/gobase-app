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
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	const pageSize = 10

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

	totalItems, err := stockService.CountSummaries(filterName, filterCategory, supplierID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * pageSize

	items, err := stockService.ListSummariesPaginated(filterName, filterCategory, supplierID, pageSize, offset)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalIn, totalOut, totalRemaining, err := stockService.GetSummaryTotals(filterName, filterCategory, supplierID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	supplierRepo := &repositories.SupplierRepository{DB: config.DB}
	suppliers, err := supplierRepo.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	storeRepo := &repositories.StoreRepository{DB: config.DB}
	stores, _ := storeRepo.GetByIDs(allowedStoreIDs)

	pages := make([]int, 0, totalPages)
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := page + 1
	if totalPages > 0 && nextPage > totalPages {
		nextPage = totalPages
	}

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
		"CurrentPage":      page,
		"TotalPages":       totalPages,
		"Pages":            pages,
		"PrevPage":         prevPage,
		"NextPage":         nextPage,
		"TotalItems":       totalItems,
		"RowOffset":        offset + 1,
	})
}
