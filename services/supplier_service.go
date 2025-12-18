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

// GetSuppliersPaginated mengembalikan data supplier berdasarkan halaman dan ukuran halaman (pageSize).
// Fungsi ini juga mengembalikan total data supplier untuk keperluan perhitungan total halaman.
func (s *SupplierService) GetSuppliersPaginated(page, pageSize int) ([]models.Supplier, int, error) {
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

func (s *SupplierService) CreateSupplier(supplier models.Supplier) error {
	return s.Repo.Create(supplier)
}

func (s *SupplierService) DeleteSupplier(id int) error {
	return s.Repo.Delete(id)
}

func (s *SupplierService) UpdateSupplier(supplier models.Supplier) error {
	return s.Repo.Update(supplier)
}
