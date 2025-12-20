package services

import (
	"stok-hadiah/models"
	"stok-hadiah/repositories"
)

type ItemService struct {
	Repo *repositories.ItemRepository
}

func (s *ItemService) GetItems() ([]models.Item, error) {
	return s.Repo.GetAll()
}

// SearchItems mencari item berdasarkan nama dan/atau kategori.
func (s *ItemService) SearchItems(name, category string) ([]models.Item, error) {
	return s.Repo.Search(name, category)
}

// CountItems mengembalikan jumlah seluruh item.
func (s *ItemService) CountItems() (int, error) {
	return s.Repo.Count()
}

// CountItemsByCategory mengembalikan jumlah item pada kategori tertentu.
func (s *ItemService) CountItemsByCategory(category string) (int, error) {
	return s.Repo.CountByCategory(category)
}

// GetItemsPaginated mengembalikan data item berdasarkan halaman dan ukuran halaman (pageSize).
// Fungsi ini juga mengembalikan total data item untuk keperluan perhitungan total halaman.
func (s *ItemService) GetItemsPaginated(page, pageSize int) ([]models.Item, int, error) {
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

func (s *ItemService) CreateItem(item models.Item) error {
	return s.Repo.Create(item)
}

func (s *ItemService) UpdateItem(item models.Item) error {
	return s.Repo.Update(item)
}

func (s *ItemService) DeleteItem(id int) error {
	return s.Repo.Delete(id)
}
