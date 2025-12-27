package controllers

import (
	"errors"
	"net/http"
	"stok-hadiah/config"
	helpers "stok-hadiah/helper"
	"stok-hadiah/repositories"
	"stok-hadiah/services"
	"strconv"
	"strings"

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

	itemQtyMap := map[int]int{}
	stockOutInfo := map[int]repositories.StockOutInfo{}
	priorOutMap := map[int]int{}
	availableTodayMap := map[int]int{}
	userID, _ := getCurrentUserID(c)

	stockOutRepo := &repositories.StockOutRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}

	if userID > 0 && len(items) > 0 {
		itemIDs := make([]int, 0, len(items))
		for _, it := range items {
			itemIDs = append(itemIDs, it.ItemID)
		}
		if infoMap, err := stockOutRepo.GetTodayQuantities(itemIDs, userID); err == nil {
			stockOutInfo = infoMap
			for itemID, info := range infoMap {
				itemQtyMap[itemID] = info.Qty
			}
		}

		if prevMap, err := stockOutRepo.GetQuantityBeforeToday(itemIDs); err == nil {
			priorOutMap = prevMap
		}
		for _, it := range items {
			prior := priorOutMap[it.ItemID]
			available := it.QtyIn - prior
			if available < 0 {
				available = 0
			}
			availableTodayMap[it.ItemID] = available
		}
	}

	type itemOutSummary struct {
		ItemID         int
		ItemName       string
		QtyIn          int
		QtyOut         int
		Remaining      int
		PriorQtyOut    int
		AvailableToday int
	}

	totalOut := 0
	totalIn := 0
	totalPriorOut := 0
	totalAvailable := 0
	totalRemaining := 0
	summaryOut := make([]itemOutSummary, 0)
	for _, it := range items {
		info := stockOutInfo[it.ItemID]
		qty := info.Qty
		if qty <= 0 || strings.TrimSpace(info.Reason) != "" {
			continue
		}
		totalOut += qty
		priorOut := priorOutMap[it.ItemID]
		availableToday := availableTodayMap[it.ItemID]
		summaryOut = append(summaryOut, itemOutSummary{
			ItemID:         it.ItemID,
			ItemName:       it.ItemName,
			QtyIn:          it.QtyIn,
			QtyOut:         qty,
			Remaining:      it.Remaining,
			PriorQtyOut:    priorOut,
			AvailableToday: availableToday,
		})
		totalIn += it.QtyIn
		totalPriorOut += priorOut
		totalAvailable += availableToday
		totalRemaining += it.Remaining
	}

	Render(c, "item_out.html", gin.H{
		"Title":                 "Item Out",
		"Page":                  "item_out",
		"items":                 items,
		"suppliers":             suppliers,
		"FilterItemName":        filterName,
		"FilterCategory":        filterCategory,
		"FilterSupplierID":      selectedSupplier,
		"dateNow":               helpers.DateNowID(),
		"ItemQtyMap":            itemQtyMap,
		"ItemPriorOutMap":       priorOutMap,
		"ItemAvailableTodayMap": availableTodayMap,
		"TotalOut":              totalOut,
		"SummaryOut":            summaryOut,
		"SummaryTotalIn":        totalIn,
		"SummaryTotalPrior":     totalPriorOut,
		"SummaryTotalAvail":     totalAvailable,
		"SummaryTotalRemain":    totalRemaining,
	})
}

// ItemOutUpdate menerima aksi plus/minus via AJAX dan menyimpan ke stock_out & stock_out_events.
func ItemOutUpdate(c *gin.Context) {
	type requestPayload struct {
		ItemID    int    `json:"item_id" binding:"required"`
		Direction string `json:"direction" binding:"required"` // "up" atau "down"
	}

	var payload requestPayload
	if err := c.ShouldBindJSON(&payload); err != nil || payload.ItemID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload tidak valid"})
		return
	}

	var delta int
	switch payload.Direction {
	case "up":
		delta = 1
	case "down":
		delta = -1
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "aksi tidak dikenal"})
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user tidak ditemukan"})
		return
	}

	allowedStoreIDs := getAllowedStoreIDs(c)
	repo := &repositories.StockOutRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}
	service := &services.StockOutService{Repo: repo}

	newQty, err := service.AdjustQuantity(payload.ItemID, delta, userID)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, repositories.ErrItemNotAllowed):
			status = http.StatusForbidden
		case errors.Is(err, repositories.ErrItemNotFound),
			errors.Is(err, repositories.ErrProgramNotFound),
			errors.Is(err, repositories.ErrQuantityNegative),
			errors.Is(err, repositories.ErrQuantityZero):
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"new_qty":   newQty,
		"direction": payload.Direction,
	})
}
