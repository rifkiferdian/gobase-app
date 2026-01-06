package controllers

import (
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserReportListIndex(c *gin.Context) {

	allowedStoreIDs := getAllowedStoreIDs(c)

	reportRepo := &repositories.UserReportRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}

	reports, err := reportRepo.GetSummaries()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "user_report_list.html", gin.H{
		"Title": "User Report",
		"Page":  "user_report",
		"Data":  reports,
	})

}

func UserReportDetailIndex(c *gin.Context) {

	userID, _ := strconv.Atoi(c.Param("id"))
	filterName := c.Query("item_name")
	filterDate := c.Query("date")

	allowedStoreIDs := getAllowedStoreIDs(c)

	reportRepo := &repositories.UserReportRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}

	detail, err := reportRepo.GetDetail(userID, filterName, filterDate)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if detail.UserID == 0 {
		c.String(http.StatusNotFound, "User tidak ditemukan atau tidak memiliki data yang bisa ditampilkan.")
		return
	}

	Render(c, "user_report_detail.html", gin.H{
		"Title":          "User Report Detail",
		"Page":           "user_report_detail",
		"UserDetail":     detail,
		"StockIns":       detail.StockIns,
		"StockOuts":      detail.StockOuts,
		"TotalIn":        detail.TotalIn,
		"TotalOut":       detail.TotalOut,
		"FilterItemName": filterName,
		"FilterDate":     filterDate,
	})

}
