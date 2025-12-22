package services

import (
	"stok-hadiah/models"
	"stok-hadiah/repositories"
)

type PermissionService struct {
	Repo *repositories.PermissionRepository
}

func (s *PermissionService) GetGroupedPermissions() ([]models.PermissionGroup, error) {
	return s.Repo.GetGrouped()
}
