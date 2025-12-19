package controllers

import (
	"context"
	"errors"
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/models"
	"stok-hadiah/repositories"
	"stok-hadiah/services"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// ItemInIndex menampilkan halaman listing stok barang masuk.
func ItemInIndex(c *gin.Context) {
	// Ambil parameter page dari query string, default 1
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	const pageSize = 10

	stockRepo := &repositories.StockInRepository{DB: config.DB}
	stockService := &services.StockInService{Repo: stockRepo}

	itemRepo := &repositories.ItemRepository{DB: config.DB}

	// =========================
	// Goroutine: items + statistik
	// =========================
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	var (
		items []models.Item // sesuaikan tipe item kamu

		totalQty          int
		totalTransactions int
		todayQty          int
		todayTransactions int

		itemsErr, errTotalQty, errTotalTrans, errTodayQty, errTodayTrans error
	)

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
		}
		items, itemsErr = itemRepo.GetAll()
		if itemsErr != nil {
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
		}
		totalQty, errTotalQty = stockService.TotalQty()
		if errTotalQty != nil {
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
		}
		totalTransactions, errTotalTrans = stockService.TotalTransactions()
		if errTotalTrans != nil {
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
		}
		todayQty, errTodayQty = stockService.TodayQty()
		if errTodayQty != nil {
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
		}
		todayTransactions, errTodayTrans = stockService.TodayTransactions()
		if errTodayTrans != nil {
			cancel()
		}
	}()

	wg.Wait()

	// Cek error setelah semua selesai
	if itemsErr != nil {
		c.String(http.StatusInternalServerError, itemsErr.Error())
		return
	}
	if errTotalQty != nil {
		c.String(http.StatusInternalServerError, errTotalQty.Error())
		return
	}
	if errTotalTrans != nil {
		c.String(http.StatusInternalServerError, errTotalTrans.Error())
		return
	}
	if errTodayQty != nil {
		c.String(http.StatusInternalServerError, errTodayQty.Error())
		return
	}
	if errTodayTrans != nil {
		c.String(http.StatusInternalServerError, errTodayTrans.Error())
		return
	}

	// Ambil parameter filter pencarian
	itemName := c.Query("item_name")
	date := c.Query("date")

	// Jika ada filter, gunakan pencarian tanpa pagination
	if itemName != "" || date != "" {
		stockIns, err := stockService.SearchStockIns(itemName, date)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		Render(c, "item/item_in.html", gin.H{
			"Title":          "Barang Masuk",
			"Page":           "item_in",
			"stockIns":       stockIns,
			"items":          items,
			"TotalQty":       totalQty,
			"TotalTrans":     totalTransactions,
			"TodayQty":       todayQty,
			"TodayTrans":     todayTransactions,
			"CurrentPage":    1,
			"TotalPages":     1,
			"Pages":          []int{1},
			"PrevPage":       1,
			"NextPage":       1,
			"FilterItemName": itemName,
			"FilterDate":     date,
		})
		return
	}

	// Tanpa filter, gunakan pagination standar
	stockIns, total, err := stockService.GetStockInsPaginated(page, pageSize)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Hitung total halaman
	totalPages := 0
	if total > 0 {
		if total%pageSize == 0 {
			totalPages = total / pageSize
		} else {
			totalPages = (total / pageSize) + 1
		}
	}

	// Siapkan slice untuk nomor halaman (1..totalPages)
	pages := make([]int, 0, totalPages)
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	// Hitung halaman sebelumnya dan berikutnya untuk tombol Prev/Next
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	nextPage := page + 1
	if totalPages > 0 && nextPage > totalPages {
		nextPage = totalPages
	}

	Render(c, "item/item_in.html", gin.H{
		"Title":          "Barang Masuk",
		"Page":           "item_in",
		"stockIns":       stockIns,
		"items":          items,
		"TotalQty":       totalQty,
		"TotalTrans":     totalTransactions,
		"TodayQty":       todayQty,
		"TodayTrans":     todayTransactions,
		"CurrentPage":    page,
		"TotalPages":     totalPages,
		"Pages":          pages,
		"PrevPage":       prevPage,
		"NextPage":       nextPage,
		"FilterItemName": itemName,
		"FilterDate":     date,
	})
}

// ItemInStore menangani penyimpanan data stok masuk baru dari form modal.
func ItemInStore(c *gin.Context) {
	type StockInForm struct {
		ItemID      int    `form:"item_id" binding:"required"`
		Qty         int    `form:"qty" binding:"required"`
		ReceivedAt  string `form:"received_at" binding:"required"`
		Description string `form:"description"`
	}

	var form StockInForm
	stockRepo := &repositories.StockInRepository{DB: config.DB}
	stockService := &services.StockInService{Repo: stockRepo}

	itemRepo := &repositories.ItemRepository{DB: config.DB}

	if err := c.ShouldBind(&form); err != nil || form.Qty <= 0 {
		// Jika validasi form gagal, kirim error ke view di atas tabel data
		stockIns, _ := stockService.GetStockIns()
		items, _ := itemRepo.GetAll()
		totalQty, _ := stockService.TotalQty()
		totalTransactions, _ := stockService.TotalTransactions()
		todayQty, _ := stockService.TodayQty()
		todayTransactions, _ := stockService.TodayTransactions()

		Render(c, "item/item_in.html", gin.H{
			"Title":      "Barang Masuk",
			"Page":       "item_in",
			"stockIns":   stockIns,
			"items":      items,
			"TotalQty":   totalQty,
			"TotalTrans": totalTransactions,
			"TodayQty":   todayQty,
			"TodayTrans": todayTransactions,
			"Error":      "Item, tanggal dan quantity wajib diisi dan quantity harus lebih dari 0",
		})
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	entry := models.StockIn{
		UserID:      userID,
		ItemID:      form.ItemID,
		Qty:         form.Qty,
		ReceivedAt:  form.ReceivedAt,
		Description: form.Description,
	}

	if err := stockService.CreateStockIn(entry); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/item-in")
}

// ItemInUpdate menangani update data stok masuk dari form edit modal.
func ItemInUpdate(c *gin.Context) {
	type StockInUpdateForm struct {
		ID          int    `form:"id" binding:"required"`
		ItemID      int    `form:"item_id" binding:"required"`
		Qty         int    `form:"qty" binding:"required"`
		ReceivedAt  string `form:"received_at" binding:"required"`
		Description string `form:"description"`
	}

	var form StockInUpdateForm
	stockRepo := &repositories.StockInRepository{DB: config.DB}
	stockService := &services.StockInService{Repo: stockRepo}

	itemRepo := &repositories.ItemRepository{DB: config.DB}

	if err := c.ShouldBind(&form); err != nil || form.Qty <= 0 {
		stockIns, _ := stockService.GetStockIns()
		items, _ := itemRepo.GetAll()
		totalQty, _ := stockService.TotalQty()
		totalTransactions, _ := stockService.TotalTransactions()
		todayQty, _ := stockService.TodayQty()
		todayTransactions, _ := stockService.TodayTransactions()

		Render(c, "item/item_in.html", gin.H{
			"Title":      "Barang Masuk",
			"Page":       "item_in",
			"stockIns":   stockIns,
			"items":      items,
			"TotalQty":   totalQty,
			"TotalTrans": totalTransactions,
			"TodayQty":   todayQty,
			"TodayTrans": todayTransactions,
			"Error":      "Item, tanggal dan quantity wajib diisi dan quantity harus lebih dari 0",
		})
		return
	}

	entry := models.StockIn{
		ID:          form.ID,
		ItemID:      form.ItemID,
		Qty:         form.Qty,
		ReceivedAt:  form.ReceivedAt,
		Description: form.Description,
	}

	if err := stockService.UpdateStockIn(entry); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/item-in")
}

// ItemInDelete menangani penghapusan data stok masuk berdasarkan ID.
func ItemInDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid stock in id")
		return
	}

	stockRepo := &repositories.StockInRepository{DB: config.DB}
	stockService := &services.StockInService{Repo: stockRepo}

	if err := stockService.DeleteStockIn(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/item-in")
}

// getCurrentUserID mengambil ID user yang sedang login berdasarkan username di session/context.
func getCurrentUserID(c *gin.Context) (int, error) {
	userVal, exists := c.Get("User")
	if !exists || userVal == nil {
		return 0, errors.New("user not found in context")
	}

	username, ok := userVal.(string)
	if !ok {
		return 0, errors.New("invalid user value in context")
	}

	var id int
	err := config.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
