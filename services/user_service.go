package services

import (
	"stok-hadiah/models"
	"stok-hadiah/repositories"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}
