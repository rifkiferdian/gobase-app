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
	// Ambil parameter page dari query string, default 1
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	const pageSize = 10

	repo := &repositories.SupplierRepository{DB: config.DB}
	service := &services.SupplierService{Repo: repo}

	// Ambil parameter filter pencarian
	filterName := c.Query("supplier_name")

	// Jika ada kata kunci pencarian, gunakan search tanpa pagination
	if filterName != "" {
		data, err := service.SearchSuppliersByName(filterName)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		Render(c, "supplier/index.html", gin.H{
			"Title":              "Supplier Page",
			"Page":               "supplier",
			"suppliers":          data,
			"CurrentPage":        1,
			"TotalPages":         1,
			"Pages":              []int{1},
			"PrevPage":           1,
			"NextPage":           1,
			"FilterSupplierName": filterName,
		})
		return
	}

	data, total, err := service.GetSuppliersPaginated(page, pageSize)
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
	pages := []int{}
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

	Render(c, "supplier/index.html", gin.H{
		"Title":              "Supplier Page",
		"Page":               "supplier",
		"suppliers":          data,
		"CurrentPage":        page,
		"TotalPages":         totalPages,
		"Pages":              pages,
		"PrevPage":           prevPage,
		"NextPage":           nextPage,
		"FilterSupplierName": filterName,
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
