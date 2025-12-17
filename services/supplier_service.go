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

func (s *SupplierService) CreateSupplier(supplier models.Supplier) error {
	return s.Repo.Create(supplier)
}

func (s *SupplierService) DeleteSupplier(id int) error {
	return s.Repo.Delete(id)
}

func (s *SupplierService) UpdateSupplier(supplier models.Supplier) error {
	return s.Repo.Update(supplier)
}
