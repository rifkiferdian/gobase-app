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

// ProgramIndex menampilkan halaman listing program.
func ProgramIndex(c *gin.Context) {
	programRepo := &repositories.ProgramRepository{DB: config.DB}
	programService := &services.ProgramService{Repo: programRepo}

	itemRepo := &repositories.ItemRepository{DB: config.DB}
	itemService := &services.ItemService{Repo: itemRepo}

	programs, err := programService.GetPrograms()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	items, err := itemService.GetItems()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	Render(c, "program/index.html", gin.H{
		"Title":    "Program Page",
		"Page":     "program",
		"programs": programs,
		"items":    items,
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

		Render(c, "program/index.html", gin.H{
			"Title":    "Program Page",
			"Page":     "program",
			"programs": programs,
			"items":    items,
			"Error":    "Nama program, item, tanggal mulai dan tanggal selesai wajib diisi",
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

		Render(c, "program/index.html", gin.H{
			"Title":    "Program Page",
			"Page":     "program",
			"programs": programs,
			"items":    items,
			"Error":    "Nama program, item, tanggal mulai dan tanggal selesai wajib diisi",
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
