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
