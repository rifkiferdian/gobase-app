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

// fetchItemMetrics mengambil total item dan total per kategori untuk ditampilkan di kartu statistik.
func fetchItemMetrics(itemService *services.ItemService, knownTotal *int) (int, int, int, int, error) {
	var totalItems int
	var err error

	if knownTotal != nil {
		totalItems = *knownTotal
	} else {
		totalItems, err = itemService.CountItems()
		if err != nil {
			return 0, 0, 0, 0, err
		}
	}

	totalFood, err := itemService.CountItemsByCategory("FOOD")
	if err != nil {
		return 0, 0, 0, 0, err
	}

	totalNonFood, err := itemService.CountItemsByCategory("NON FOOD")
	if err != nil {
		return 0, 0, 0, 0, err
	}

	totalDeptStore, err := itemService.CountItemsByCategory("DEPT STORE")
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return totalItems, totalFood, totalNonFood, totalDeptStore, nil
}

// ItemIndex menampilkan halaman listing item.
func ItemIndex(c *gin.Context) {
	// Ambil parameter page dari query string, default 1
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	const pageSize = 10

	// Ambil filter store dari query
	filterStoreStr := c.Query("store_id")
	var filterStoreIDPtr *int
	filterStoreID := 0
	if id, err := strconv.Atoi(filterStoreStr); err == nil && id > 0 {
		filterStoreID = id
		filterStoreIDPtr = &filterStoreID
	}

	allowedStoreIDs := getAllowedStoreIDs(c)
	itemRepo := &repositories.ItemRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		FilterStoreID:      filterStoreIDPtr,
		EnforceStoreFilter: true,
	}
	itemService := &services.ItemService{Repo: itemRepo}
	storeRepo := &repositories.StoreRepository{DB: config.DB}
	stores, _ := storeRepo.GetByIDs(allowedStoreIDs)

	// Ambil filter pencarian
	filterName := c.Query("item_name")
	filterCategory := c.Query("category")

	// Jika ada filter pencarian, gunakan search tanpa pagination
	if filterName != "" || filterCategory != "" {
		items, err := itemService.SearchItems(filterName, filterCategory)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		supplierRepo := &repositories.SupplierRepository{DB: config.DB}
		suppliers, err := supplierRepo.GetAll()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		totalItems, totalFood, totalNonFood, totalDeptStore, err := fetchItemMetrics(itemService, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		Render(c, "item/index.html", gin.H{
			"Title":          "Item Page",
			"Page":           "item",
			"items":          items,
			"suppliers":      suppliers,
			"stores":         stores,
			"TotalItems":     totalItems,
			"TotalFood":      totalFood,
			"TotalNonFood":   totalNonFood,
			"TotalDeptStore": totalDeptStore,
			"CurrentPage":    1,
			"TotalPages":     1,
			"Pages":          []int{1},
			"PrevPage":       1,
			"NextPage":       1,
			"FilterItemName": filterName,
			"FilterCategory": filterCategory,
			"FilterStoreID":  filterStoreID,
		})
		return
	}

	// Ambil data item dengan pagination
	items, total, err := itemService.GetItemsPaginated(page, pageSize)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalItems, totalFood, totalNonFood, totalDeptStore, err := fetchItemMetrics(itemService, &total)
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

	// Supplier tetap diambil semua untuk kebutuhan dropdown
	supplierRepo := &repositories.SupplierRepository{DB: config.DB}
	suppliers, err := supplierRepo.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "item/index.html", gin.H{
		"Title":          "Item Page",
		"Page":           "item",
		"items":          items,
		"suppliers":      suppliers,
		"stores":         stores,
		"TotalItems":     totalItems,
		"TotalFood":      totalFood,
		"TotalNonFood":   totalNonFood,
		"TotalDeptStore": totalDeptStore,
		"CurrentPage":    page,
		"TotalPages":     totalPages,
		"Pages":          pages,
		"PrevPage":       prevPage,
		"NextPage":       nextPage,
		"FilterItemName": filterName,
		"FilterCategory": filterCategory,
		"FilterStoreID":  filterStoreID,
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
	allowedStoreIDs := getAllowedStoreIDs(c)
	itemRepo := &repositories.ItemRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}
	itemService := &services.ItemService{Repo: itemRepo}

	supplierRepo := &repositories.SupplierRepository{DB: config.DB}
	storeRepo := &repositories.StoreRepository{DB: config.DB}
	stores, _ := storeRepo.GetByIDs(allowedStoreIDs)

	if err := c.ShouldBind(&form); err != nil {
		// Jika validasi form gagal, kirim error ke view di atas tabel data
		items, _ := itemService.GetItems()
		suppliers, _ := supplierRepo.GetAll()
		totalItems, totalFood, totalNonFood, totalDeptStore, _ := fetchItemMetrics(itemService, nil)

		Render(c, "item/index.html", gin.H{
			"Title":          "Item Page",
			"Page":           "item",
			"items":          items,
			"suppliers":      suppliers,
			"stores":         stores,
			"TotalItems":     totalItems,
			"TotalFood":      totalFood,
			"TotalNonFood":   totalNonFood,
			"TotalDeptStore": totalDeptStore,
			"Error":          "Nama item, kategori, dan supplier wajib diisi",
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
	allowedStoreIDs := getAllowedStoreIDs(c)
	itemRepo := &repositories.ItemRepository{
		DB:                 config.DB,
		StoreIDs:           allowedStoreIDs,
		EnforceStoreFilter: true,
	}
	itemService := &services.ItemService{Repo: itemRepo}

	supplierRepo := &repositories.SupplierRepository{DB: config.DB}
	storeRepo := &repositories.StoreRepository{DB: config.DB}
	stores, _ := storeRepo.GetByIDs(allowedStoreIDs)

	if err := c.ShouldBind(&form); err != nil {
		items, _ := itemService.GetItems()
		suppliers, _ := supplierRepo.GetAll()
		totalItems, totalFood, totalNonFood, totalDeptStore, _ := fetchItemMetrics(itemService, nil)

		Render(c, "item/index.html", gin.H{
			"Title":          "Item Page",
			"Page":           "item",
			"items":          items,
			"suppliers":      suppliers,
			"stores":         stores,
			"TotalItems":     totalItems,
			"TotalFood":      totalFood,
			"TotalNonFood":   totalNonFood,
			"TotalDeptStore": totalDeptStore,
			"Error":          "Nama item, kategori, dan supplier wajib diisi",
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
