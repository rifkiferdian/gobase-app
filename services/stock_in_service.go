package services

import (
	"stok-hadiah/models"
	"stok-hadiah/repositories"
)

type StockInService struct {
	Repo *repositories.StockInRepository
}

func (s *StockInService) GetStockIns() ([]models.StockIn, error) {
	return s.Repo.GetAll()
}

// GetStockInsPaginated mengembalikan data stok masuk berdasarkan halaman dan ukuran halaman (pageSize).
// Fungsi ini juga mengembalikan total data stok masuk untuk keperluan perhitungan total halaman.
func (s *StockInService) GetStockInsPaginated(page, pageSize int) ([]models.StockIn, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	total, err := s.Repo.Count()
	if err != nil {
		return nil, 0, err
	}

	data, err := s.Repo.GetPaginated(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (s *StockInService) CreateStockIn(in models.StockIn) error {
	return s.Repo.Create(in)
}

func (s *StockInService) UpdateStockIn(in models.StockIn) error {
	return s.Repo.Update(in)
}

func (s *StockInService) DeleteStockIn(id int) error {
	return s.Repo.Delete(id)
}
