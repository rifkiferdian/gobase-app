package controllers

import (
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/repositories"

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
