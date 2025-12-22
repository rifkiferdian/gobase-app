package services

import (
	"stok-hadiah/models"
	"stok-hadiah/repositories"
)

type RoleService struct {
	Repo *repositories.RoleRepository
}

func (s *RoleService) GetRoles() ([]models.Role, error) {
	return s.Repo.GetAll()
}
