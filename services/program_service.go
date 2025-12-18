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

// GetProgramsPaginated mengembalikan data program berdasarkan halaman dan ukuran halaman (pageSize).
// Fungsi ini juga mengembalikan total data program untuk keperluan perhitungan total halaman.
func (s *ProgramService) GetProgramsPaginated(page, pageSize int) ([]models.Program, int, error) {
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

func (s *ProgramService) CreateProgram(p models.Program) error {
	return s.Repo.Create(p)
}

func (s *ProgramService) UpdateProgram(p models.Program) error {
	return s.Repo.Update(p)
}

func (s *ProgramService) DeleteProgram(id int) error {
	return s.Repo.Delete(id)
}
