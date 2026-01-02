package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"stok-hadiah/config"
	"stok-hadiah/repositories"

	"github.com/gin-gonic/gin"
)

type stockChartSeries struct {
	Name string `json:"name"`
	Data []int  `json:"data"`
}

type stockChartPayload struct {
	Categories []string           `json:"categories"`
	Series     []stockChartSeries `json:"series"`
}

type stockChartData struct {
	Week  stockChartPayload `json:"week"`
	Month stockChartPayload `json:"month"`
	Year  stockChartPayload `json:"year"`
}

func DashboardIndex(c *gin.Context) {

	// debug session
	// sess := sessions.Default(c)
	// fmt.Println("DEBUG user_id:", sess.Get("user_id"))
	// fmt.Println("DEBUG user:", sess.Get("user"))

	allowedStoreIDs := getAllowedStoreIDs(c)

	stockInRepo := &repositories.StockInRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}
	stockOutRepo := &repositories.StockOutRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}

	totalStockIn, err := stockInRepo.SumQty()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalStockOut, err := stockOutRepo.SumQty()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalStockOutToday, err := stockOutRepo.SumTodayQty()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	chartData, err := buildStockChartData(stockInRepo, stockOutRepo)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	chartJSON, err := json.Marshal(chartData)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "dashboard/index.html", gin.H{
		"Title":              "Dashboard User",
		"Page":               "dashboard",
		"StockChartData":     template.JS(chartJSON),
		"TotalStockIn":       totalStockIn,
		"TotalStockOut":      totalStockOut,
		"TotalStockOutToday": totalStockOutToday,
	})

}

func HomeIndex(c *gin.Context) {
	Render(c, "home.html", gin.H{
		"Title": "Home Page",
	})
}

func buildStockChartData(stockInRepo *repositories.StockInRepository, stockOutRepo *repositories.StockOutRepository) (stockChartData, error) {
	const weekDays = 7
	today := time.Now()

	// Start of current week (Monday)
	startOfWeek := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	weekday := int(startOfWeek.Weekday())
	offset := weekday - int(time.Monday)
	if offset < 0 {
		offset += 7
	}
	startOfWeek = startOfWeek.AddDate(0, 0, -offset)

	weekPayload, err := buildDailyChartPayload(startOfWeek, weekDays, stockInRepo, stockOutRepo)
	if err != nil {
		return stockChartData{}, err
	}

	// Current month range (1st to last day)
	startOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	daysInMonth := time.Date(today.Year(), today.Month()+1, 0, 0, 0, 0, 0, today.Location()).Day()
	monthPayload, err := buildDailyChartPayload(startOfMonth, daysInMonth, stockInRepo, stockOutRepo)
	if err != nil {
		return stockChartData{}, err
	}

	yearPayload, err := buildYearlyChartPayload(stockInRepo, stockOutRepo)
	if err != nil {
		return stockChartData{}, err
	}

	return stockChartData{
		Week:  weekPayload,
		Month: monthPayload,
		Year:  yearPayload,
	}, nil
}

func buildDailyChartPayload(start time.Time, totalDays int, stockInRepo *repositories.StockInRepository, stockOutRepo *repositories.StockOutRepository) (stockChartPayload, error) {
	stockInTotals, err := stockInRepo.DailyTotalsSince(start)
	if err != nil {
		return stockChartPayload{}, err
	}

	stockOutTotals, err := stockOutRepo.DailyTotalsSince(start)
	if err != nil {
		return stockChartPayload{}, err
	}

	categories := make([]string, totalDays)
	stockInSeries := make([]int, totalDays)
	stockOutSeries := make([]int, totalDays)

	for i := 0; i < totalDays; i++ {
		dayTime := start.AddDate(0, 0, i)
		dayKey := dayTime.Format("2006-01-02")
		categories[i] = dayTime.Format("02 Jan")
		stockInSeries[i] = stockInTotals[dayKey]
		stockOutSeries[i] = stockOutTotals[dayKey]
	}

	return stockChartPayload{
		Categories: categories,
		Series: []stockChartSeries{
			{Name: "Stock In", Data: stockInSeries},
			{Name: "Stock Out", Data: stockOutSeries},
		},
	}, nil
}

func buildYearlyChartPayload(stockInRepo *repositories.StockInRepository, stockOutRepo *repositories.StockOutRepository) (stockChartPayload, error) {
	const monthsToShow = 12

	now := time.Now()
	startMonth := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())

	stockInTotals, err := stockInRepo.MonthlyTotalsSince(startMonth)
	if err != nil {
		return stockChartPayload{}, err
	}

	stockOutTotals, err := stockOutRepo.MonthlyTotalsSince(startMonth)
	if err != nil {
		return stockChartPayload{}, err
	}

	categories := make([]string, monthsToShow)
	stockInSeries := make([]int, monthsToShow)
	stockOutSeries := make([]int, monthsToShow)
	for i := 0; i < monthsToShow; i++ {
		monthTime := startMonth.AddDate(0, i, 0)
		monthKey := monthTime.Format("2006-01-02")
		categories[i] = monthTime.Format("Jan 2006")
		stockInSeries[i] = stockInTotals[monthKey]
		stockOutSeries[i] = stockOutTotals[monthKey]
	}

	return stockChartPayload{
		Categories: categories,
		Series: []stockChartSeries{
			{Name: "Stock In", Data: stockInSeries},
			{Name: "Stock Out", Data: stockOutSeries},
		},
	}, nil
}
