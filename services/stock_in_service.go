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

func (s *StockInService) CreateStockIn(in models.StockIn) error {
	return s.Repo.Create(in)
}

func (s *StockInService) UpdateStockIn(in models.StockIn) error {
	return s.Repo.Update(in)
}

func (s *StockInService) DeleteStockIn(id int) error {
	return s.Repo.Delete(id)
}
