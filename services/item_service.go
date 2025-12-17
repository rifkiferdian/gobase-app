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

func (s *ItemService) CreateItem(item models.Item) error {
	return s.Repo.Create(item)
}

func (s *ItemService) UpdateItem(item models.Item) error {
	return s.Repo.Update(item)
}

func (s *ItemService) DeleteItem(id int) error {
	return s.Repo.Delete(id)
}
