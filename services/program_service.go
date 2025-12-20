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

// SearchPrograms mencari program berdasarkan nama dan/atau rentang tanggal.
func (s *ProgramService) SearchPrograms(name, startDate, endDate string) ([]models.Program, error) {
	return s.Repo.Search(name, startDate, endDate)
}

// CountPrograms mengembalikan jumlah seluruh program.
func (s *ProgramService) CountPrograms() (int, error) {
	return s.Repo.Count()
}

// CountProgramsActiveToday mengembalikan jumlah program yang aktif hari ini.
func (s *ProgramService) CountProgramsActiveToday() (int, error) {
	return s.Repo.CountActiveToday()
}

// CountProgramsInactiveToday mengembalikan jumlah program yang tidak aktif hari ini.
func (s *ProgramService) CountProgramsInactiveToday() (int, error) {
	return s.Repo.CountInactiveToday()
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
