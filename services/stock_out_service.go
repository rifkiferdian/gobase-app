package services

import (
	"errors"
	"strings"

	"stok-hadiah/models"
	"stok-hadiah/repositories"
)

// StockOutService menyediakan layer bisnis untuk penyesuaian stok keluar.
type StockOutService struct {
	Repo *repositories.StockOutRepository
}

// AdjustQuantity memanggil repo untuk menambah/mengurangi qty keluar per item.
// Delta seharusnya bernilai +1 atau -1.
func (s *StockOutService) AdjustQuantity(itemID, delta, userID int) (int, error) {
	if delta == 0 {
		return 0, errors.New("delta harus lebih dari 0")
	}
	return s.Repo.AdjustQuantity(itemID, delta, userID)
}

// CreateCaseStockOut mencatat pengeluaran stok dengan alasan khusus.
func (s *StockOutService) CreateCaseStockOut(itemID, qty, userID int, reason string) (models.StockOutCase, error) {
	reason = strings.TrimSpace(reason)
	if qty <= 0 {
		return models.StockOutCase{}, errors.New("qty harus lebih dari 0")
	}
	if reason == "" {
		return models.StockOutCase{}, errors.New("alasan wajib diisi")
	}
	return s.Repo.CreateCaseStockOut(itemID, qty, userID, reason)
}

// ListCaseStockOuts mengambil daftar pengeluaran stok yang memiliki alasan/keterangan.
func (s *StockOutService) ListCaseStockOuts(userID, limit int) ([]models.StockOutCase, error) {
	return s.Repo.ListCaseStockOuts(userID, limit)
}
