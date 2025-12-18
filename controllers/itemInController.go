package controllers

import (
	"errors"
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/models"
	"stok-hadiah/repositories"
	"stok-hadiah/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ItemInIndex menampilkan halaman listing stok barang masuk.
func ItemInIndex(c *gin.Context) {
	stockRepo := &repositories.StockInRepository{DB: config.DB}
	stockService := &services.StockInService{Repo: stockRepo}

	itemRepo := &repositories.ItemRepository{DB: config.DB}

	stockIns, err := stockService.GetStockIns()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	items, err := itemRepo.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "item/item_in.html", gin.H{
		"Title":    "Barang Masuk",
		"Page":     "item_in",
		"stockIns": stockIns,
		"items":    items,
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

		Render(c, "item/item_in.html", gin.H{
			"Title":    "Barang Masuk",
			"Page":     "item_in",
			"stockIns": stockIns,
			"items":    items,
			"Error":    "Item, tanggal dan quantity wajib diisi dan quantity harus lebih dari 0",
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

		Render(c, "item/item_in.html", gin.H{
			"Title":    "Barang Masuk",
			"Page":     "item_in",
			"stockIns": stockIns,
			"items":    items,
			"Error":    "Item, tanggal dan quantity wajib diisi dan quantity harus lebih dari 0",
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
