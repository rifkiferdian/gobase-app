package controllers

import (
	"net/http"
	"stok-hadiah/config"
	"stok-hadiah/models"
	"stok-hadiah/repositories"
	"stok-hadiah/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ItemIndex menampilkan halaman listing item.
func ItemIndex(c *gin.Context) {
	itemRepo := &repositories.ItemRepository{DB: config.DB}
	itemService := &services.ItemService{Repo: itemRepo}

	supplierRepo := &repositories.SupplierRepository{DB: config.DB}
	suppliers, err := supplierRepo.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	items, err := itemService.GetItems()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "item/index.html", gin.H{
		"Title":     "Item Page",
		"Page":      "item",
		"items":     items,
		"suppliers": suppliers,
	})
}

// ItemStore menangani penyimpanan data item baru dari form modal.
func ItemStore(c *gin.Context) {
	type ItemForm struct {
		ItemName    string `form:"item_name" binding:"required"`
		Category    string `form:"category" binding:"required"`
		SupplierID  int    `form:"supplier_id" binding:"required"`
		Description string `form:"description"`
	}

	var form ItemForm
	itemRepo := &repositories.ItemRepository{DB: config.DB}
	itemService := &services.ItemService{Repo: itemRepo}

	supplierRepo := &repositories.SupplierRepository{DB: config.DB}

	if err := c.ShouldBind(&form); err != nil {
		// Jika validasi form gagal, kirim error ke view di atas tabel data
		items, _ := itemService.GetItems()
		suppliers, _ := supplierRepo.GetAll()

		Render(c, "item/index.html", gin.H{
			"Title":     "Item Page",
			"Page":      "item",
			"items":     items,
			"suppliers": suppliers,
			"Error":     "Nama item, kategori, dan supplier wajib diisi",
		})
		return
	}

	item := models.Item{
		ItemName:    form.ItemName,
		Category:    form.Category,
		SupplierID:  form.SupplierID,
		Description: form.Description,
	}

	if err := itemService.CreateItem(item); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/items")
}

// ItemUpdate menangani update data item dari form edit modal.
func ItemUpdate(c *gin.Context) {
	type ItemUpdateForm struct {
		ItemID      int    `form:"item_id" binding:"required"`
		ItemName    string `form:"item_name" binding:"required"`
		Category    string `form:"category" binding:"required"`
		SupplierID  int    `form:"supplier_id" binding:"required"`
		Description string `form:"description"`
	}

	var form ItemUpdateForm
	itemRepo := &repositories.ItemRepository{DB: config.DB}
	itemService := &services.ItemService{Repo: itemRepo}

	supplierRepo := &repositories.SupplierRepository{DB: config.DB}

	if err := c.ShouldBind(&form); err != nil {
		items, _ := itemService.GetItems()
		suppliers, _ := supplierRepo.GetAll()

		Render(c, "item/index.html", gin.H{
			"Title":     "Item Page",
			"Page":      "item",
			"items":     items,
			"suppliers": suppliers,
			"Error":     "Nama item, kategori, dan supplier wajib diisi",
		})
		return
	}

	item := models.Item{
		ItemID:      form.ItemID,
		ItemName:    form.ItemName,
		Category:    form.Category,
		SupplierID:  form.SupplierID,
		Description: form.Description,
	}

	if err := itemService.UpdateItem(item); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/items")
}

// ItemDelete menangani penghapusan data item berdasarkan ID.
func ItemDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid item id")
		return
	}

	itemRepo := &repositories.ItemRepository{DB: config.DB}
	itemService := &services.ItemService{Repo: itemRepo}

	if err := itemService.DeleteItem(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/items")
}
