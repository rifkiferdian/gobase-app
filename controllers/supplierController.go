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

type SupplierController struct {
	Service *services.SupplierService
}

func SupplierIndex(c *gin.Context) {
	repo := &repositories.SupplierRepository{DB: config.DB}
	service := &services.SupplierService{Repo: repo}

	data, err := service.GetSuppliers()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "supplier/index.html", gin.H{
		"Title":     "Supplier Page",
		"Page":      "supplier",
		"suppliers": data,
	})
}

// SupplierStore menangani penyimpanan data supplier baru dari form modal
func SupplierStore(c *gin.Context) {
	type SupplierForm struct {
		SupplierName string `form:"supplier_name" binding:"required"`
		Description  string `form:"description"`
		Active       int    `form:"active"`
	}

	var form SupplierForm
	repo := &repositories.SupplierRepository{DB: config.DB}
	service := &services.SupplierService{Repo: repo}

	if err := c.ShouldBind(&form); err != nil {
		// Jika validasi form gagal, kirim error ke view di atas tabel data
		data, _ := service.GetSuppliers()

		Render(c, "supplier/index.html", gin.H{
			"Title":     "Supplier Page",
			"Page":      "supplier",
			"suppliers": data,
			"Error":     "Nama supplier wajib diisi",
		})
		return
	}

	supplier := models.Supplier{
		SupplierName: form.SupplierName,
		Description:  form.Description,
		Active:       form.Active,
	}

	if err := service.CreateSupplier(supplier); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/suppliers")
}

// SupplierUpdate menangani update data supplier dari form edit modal
func SupplierUpdate(c *gin.Context) {
	type SupplierUpdateForm struct {
		SupplierID   int    `form:"supplier_id" binding:"required"`
		SupplierName string `form:"supplier_name" binding:"required"`
		Description  string `form:"description"`
		Active       int    `form:"active"`
	}

	var form SupplierUpdateForm
	repo := &repositories.SupplierRepository{DB: config.DB}
	service := &services.SupplierService{Repo: repo}

	if err := c.ShouldBind(&form); err != nil {
		// Jika validasi form gagal saat update, kirim error ke view di atas tabel data
		data, _ := service.GetSuppliers()

		Render(c, "supplier/index.html", gin.H{
			"Title":     "Supplier Page",
			"Page":      "supplier",
			"suppliers": data,
			"Error":     "Nama supplier wajib diisi",
		})
		return
	}

	supplier := models.Supplier{
		SupplierID:   form.SupplierID,
		SupplierName: form.SupplierName,
		Description:  form.Description,
		Active:       form.Active,
	}

	if err := service.UpdateSupplier(supplier); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/suppliers")
}

// SupplierDelete menangani penghapusan data supplier berdasarkan ID
func SupplierDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid supplier id")
		return
	}

	repo := &repositories.SupplierRepository{DB: config.DB}
	service := &services.SupplierService{Repo: repo}

	if err := service.DeleteSupplier(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/suppliers")
}
