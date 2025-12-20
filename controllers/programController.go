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

// fetchProgramMetrics menghitung total program, total aktif hari ini, dan total nonaktif hari ini.
func fetchProgramMetrics(service *services.ProgramService, knownTotal *int) (int, int, int, error) {
	var total int
	var err error

	if knownTotal != nil {
		total = *knownTotal
	} else {
		total, err = service.CountPrograms()
		if err != nil {
			return 0, 0, 0, err
		}
	}

	activeToday, err := service.CountProgramsActiveToday()
	if err != nil {
		return 0, 0, 0, err
	}

	inactiveToday, err := service.CountProgramsInactiveToday()
	if err != nil {
		return 0, 0, 0, err
	}

	return total, activeToday, inactiveToday, nil
}

// ProgramIndex menampilkan halaman listing program.
func ProgramIndex(c *gin.Context) {
	// Ambil parameter page dari query string, default 1
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	const pageSize = 10

	programRepo := &repositories.ProgramRepository{DB: config.DB}
	programService := &services.ProgramService{Repo: programRepo}

	// Ambil filter pencarian
	filterName := c.Query("program_name")
	filterStartDate := c.Query("start_date")
	filterEndDate := c.Query("end_date")

	// Jika ada filter pencarian, gunakan search tanpa pagination
	if filterName != "" || filterStartDate != "" || filterEndDate != "" {
		programs, err := programService.SearchPrograms(filterName, filterStartDate, filterEndDate)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		totalPrograms, activeToday, inactiveToday, err := fetchProgramMetrics(programService, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Item tetap diambil semua untuk kebutuhan dropdown
		itemRepo := &repositories.ItemRepository{DB: config.DB}
		itemService := &services.ItemService{Repo: itemRepo}

		items, err := itemService.GetItems()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		Render(c, "program/index.html", gin.H{
			"Title":             "Program Page",
			"Page":              "program",
			"programs":          programs,
			"items":             items,
			"TotalPrograms":     totalPrograms,
			"ActiveToday":       activeToday,
			"InactiveToday":     inactiveToday,
			"CurrentPage":       1,
			"TotalPages":        1,
			"Pages":             []int{1},
			"PrevPage":          1,
			"NextPage":          1,
			"FilterProgramName": filterName,
			"FilterStartDate":   filterStartDate,
			"FilterEndDate":     filterEndDate,
		})
		return
	}

	// Ambil data program dengan pagination
	programs, total, err := programService.GetProgramsPaginated(page, pageSize)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalPrograms, activeToday, inactiveToday, err := fetchProgramMetrics(programService, &total)
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

	// Item tetap diambil semua untuk kebutuhan dropdown
	itemRepo := &repositories.ItemRepository{DB: config.DB}
	itemService := &services.ItemService{Repo: itemRepo}

	items, err := itemService.GetItems()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "program/index.html", gin.H{
		"Title":             "Program Page",
		"Page":              "program",
		"programs":          programs,
		"items":             items,
		"TotalPrograms":     totalPrograms,
		"ActiveToday":       activeToday,
		"InactiveToday":     inactiveToday,
		"CurrentPage":       page,
		"TotalPages":        totalPages,
		"Pages":             pages,
		"PrevPage":          prevPage,
		"NextPage":          nextPage,
		"FilterProgramName": filterName,
		"FilterStartDate":   filterStartDate,
		"FilterEndDate":     filterEndDate,
	})
}

// ProgramStore menangani penyimpanan data program baru dari form modal.
func ProgramStore(c *gin.Context) {
	type ProgramForm struct {
		ProgramName string `form:"program_name" binding:"required"`
		ItemID      int    `form:"item_id" binding:"required"`
		StartDate   string `form:"start_date" binding:"required"`
		EndDate     string `form:"end_date" binding:"required"`
	}

	var form ProgramForm
	programRepo := &repositories.ProgramRepository{DB: config.DB}
	programService := &services.ProgramService{Repo: programRepo}

	itemRepo := &repositories.ItemRepository{DB: config.DB}
	itemService := &services.ItemService{Repo: itemRepo}

	if err := c.ShouldBind(&form); err != nil {
		// Jika validasi form gagal, kirim error ke view
		programs, _ := programService.GetPrograms()
		items, _ := itemService.GetItems()
		totalPrograms, activeToday, inactiveToday, _ := fetchProgramMetrics(programService, nil)

		Render(c, "program/index.html", gin.H{
			"Title":         "Program Page",
			"Page":          "program",
			"programs":      programs,
			"items":         items,
			"TotalPrograms": totalPrograms,
			"ActiveToday":   activeToday,
			"InactiveToday": inactiveToday,
			"Error":         "Nama program, item, tanggal mulai dan tanggal selesai wajib diisi",
		})
		return
	}

	program := models.Program{
		ProgramName: form.ProgramName,
		ItemID:      form.ItemID,
		StartDate:   form.StartDate,
		EndDate:     form.EndDate,
	}

	if err := programService.CreateProgram(program); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/programs")
}

// ProgramUpdate menangani update data program dari form edit modal.
func ProgramUpdate(c *gin.Context) {
	type ProgramUpdateForm struct {
		ProgramID   int    `form:"program_id" binding:"required"`
		ProgramName string `form:"program_name" binding:"required"`
		ItemID      int    `form:"item_id" binding:"required"`
		StartDate   string `form:"start_date" binding:"required"`
		EndDate     string `form:"end_date" binding:"required"`
	}

	var form ProgramUpdateForm
	programRepo := &repositories.ProgramRepository{DB: config.DB}
	programService := &services.ProgramService{Repo: programRepo}

	itemRepo := &repositories.ItemRepository{DB: config.DB}
	itemService := &services.ItemService{Repo: itemRepo}

	if err := c.ShouldBind(&form); err != nil {
		programs, _ := programService.GetPrograms()
		items, _ := itemService.GetItems()
		totalPrograms, activeToday, inactiveToday, _ := fetchProgramMetrics(programService, nil)

		Render(c, "program/index.html", gin.H{
			"Title":         "Program Page",
			"Page":          "program",
			"programs":      programs,
			"items":         items,
			"TotalPrograms": totalPrograms,
			"ActiveToday":   activeToday,
			"InactiveToday": inactiveToday,
			"Error":         "Nama program, item, tanggal mulai dan tanggal selesai wajib diisi",
		})
		return
	}

	program := models.Program{
		ProgramID:   form.ProgramID,
		ProgramName: form.ProgramName,
		ItemID:      form.ItemID,
		StartDate:   form.StartDate,
		EndDate:     form.EndDate,
	}

	if err := programService.UpdateProgram(program); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/programs")
}

// ProgramDelete menangani penghapusan data program berdasarkan ID.
func ProgramDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid program id")
		return
	}

	programRepo := &repositories.ProgramRepository{DB: config.DB}
	programService := &services.ProgramService{Repo: programRepo}

	if err := programService.DeleteProgram(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/programs")
}
