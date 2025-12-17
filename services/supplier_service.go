package services

import (
	"stok-hadiah/models"
	"stok-hadiah/repositories"
)

type SupplierService struct {
	Repo *repositories.SupplierRepository
}

func (s *SupplierService) GetSuppliers() ([]models.Supplier, error) {
	return s.Repo.GetAll()
}
