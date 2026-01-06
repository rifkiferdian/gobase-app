package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ItemReportIndex(c *gin.Context) {
	itemIDStr := c.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil || itemID <= 0 {
		c.String(http.StatusBadRequest, "ID item tidak valid")
		return
	}

	allowedStoreIDs := getAllowedStoreIDs(c)

	itemRepo := &repositories.ItemRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}

	item, err := itemRepo.GetByID(itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.String(http.StatusNotFound, "Item tidak ditemukan atau tidak diizinkan.")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	stockInRepo := &repositories.StockInRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}
	stockIns, err := stockInRepo.GetByItemID(itemID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	stockOutRepo := &repositories.StockOutRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}
	stockOuts, err := stockOutRepo.GetByItemID(itemID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalIn := 0
	for _, s := range stockIns {
		totalIn += s.Qty
	}
	totalOut := 0
	for _, s := range stockOuts {
		totalOut += s.Qty
	}
	remaining := totalIn - totalOut

	Render(c, "item_report.html", gin.H{
		"Title":       "Item Report",
		"Page":        "itemReport",
		"Item":        item,
		"StockIns":    stockIns,
		"StockOuts":   stockOuts,
		"TotalQtyIn":  totalIn,
		"TotalQtyOut": totalOut,
		"Remaining":   remaining,
	})
}

func ItemOutReportIndex(c *gin.Context) {
	filterName := c.Query("item_name")
	filterDate := c.Query("date")

	storeStr := c.Query("store_id")
	storeID := 0
	if id, err := strconv.Atoi(storeStr); err == nil && id > 0 {
		storeID = id
	}

	allowedStoreIDs := getAllowedStoreIDs(c)

	stockOutRepo := &repositories.StockOutRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}

	stockOuts, err := stockOutRepo.ListReports(filterName, filterDate, storeID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	storeRepo := &repositories.StoreRepository{DB: config.DB}
	stores, _ := storeRepo.GetByIDs(allowedStoreIDs)
	if len(stores) == 0 {
		stores = nil
	}

	Render(c, "item_out_report.html", gin.H{
		"Title":          "Item Out Report",
		"Page":           "itemOutReport",
		"items":          stockOuts,
		"stores":         stores,
		"FilterItemName": filterName,
		"FilterDate":     filterDate,
		"FilterStoreID":  storeID,
	})
}
