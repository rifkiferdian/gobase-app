package services

import (
	"stok-hadiah/models"
	"stok-hadiah/repositories"
)

type ItemStockService struct {
	Repo *repositories.ItemStockRepository
}

// ListSummaries mengambil ringkasan stok item dengan filter nama, kategori, dan supplier opsional.
func (s *ItemStockService) ListSummaries(name, category string, supplierID *int) ([]models.ItemStockSummary, error) {
	return s.Repo.GetSummaries(name, category, supplierID)
}

// ListSummariesPaginated mengambil ringkasan stok item dengan pagination.
func (s *ItemStockService) ListSummariesPaginated(name, category string, supplierID *int, limit, offset int) ([]models.ItemStockSummary, error) {
	return s.Repo.GetSummariesPaginated(name, category, supplierID, limit, offset)
}

// CountSummaries menghitung total item sesuai filter.
func (s *ItemStockService) CountSummaries(name, category string, supplierID *int) (int, error) {
	return s.Repo.CountSummaries(name, category, supplierID)
}

// GetSummaryTotals mengembalikan total qty in/out/remaining dari seluruh hasil filter.
func (s *ItemStockService) GetSummaryTotals(name, category string, supplierID *int) (int, int, int, error) {
	return s.Repo.GetSummaryTotals(name, category, supplierID)
}
