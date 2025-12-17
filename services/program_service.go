package services

import (
	"stok-hadiah/models"
	"stok-hadiah/repositories"
)

type ProgramService struct {
	Repo *repositories.ProgramRepository
}

func (s *ProgramService) GetPrograms() ([]models.Program, error) {
	return s.Repo.GetAll()
}

func (s *ProgramService) CreateProgram(p models.Program) error {
	return s.Repo.Create(p)
}

func (s *ProgramService) UpdateProgram(p models.Program) error {
	return s.Repo.Update(p)
}

func (s *ProgramService) DeleteProgram(id int) error {
	return s.Repo.Delete(id)
}
