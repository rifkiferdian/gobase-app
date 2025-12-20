package repositories

import (
	"database/sql"
	"stok-hadiah/models"
)

type ProgramRepository struct {
	DB *sql.DB
}

// GetAll mengambil seluruh data program beserta nama item-nya.
func (r *ProgramRepository) GetAll() ([]models.Program, error) {
	rows, err := r.DB.Query(`
		SELECT p.program_id,
		       p.program_name,
		       p.item_id,
		       i.item_name,
		       DATE_FORMAT(p.start_date, '%Y-%m-%d') AS start_date,
		       DATE_FORMAT(p.end_date, '%Y-%m-%d') AS end_date
		FROM programs p
		JOIN items i ON i.item_id = p.item_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []models.Program

	for rows.Next() {
		var p models.Program
		if err := rows.Scan(
			&p.ProgramID,
			&p.ProgramName,
			&p.ItemID,
			&p.ItemName,
			&p.StartDate,
			&p.EndDate,
		); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}

	return programs, nil
}

// Search mengambil data program berdasarkan kata kunci nama program dan rentang tanggal mulai/selesai.
func (r *ProgramRepository) Search(name, startDate, endDate string) ([]models.Program, error) {
	query := `
		SELECT p.program_id,
		       p.program_name,
		       p.item_id,
		       i.item_name,
		       DATE_FORMAT(p.start_date, '%d-%m-%Y') AS start_date,
		       DATE_FORMAT(p.end_date, '%d-%m-%Y') AS end_date
		FROM programs p
		JOIN items i ON i.item_id = p.item_id
		WHERE 1=1`

	args := []interface{}{}

	if name != "" {
		query += " AND p.program_name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if startDate != "" {
		query += " AND DATE(p.start_date) >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND DATE(p.end_date) <= ?"
		args = append(args, endDate)
	}

	query += " ORDER BY p.program_id DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []models.Program

	for rows.Next() {
		var p models.Program
		if err := rows.Scan(
			&p.ProgramID,
			&p.ProgramName,
			&p.ItemID,
			&p.ItemName,
			&p.StartDate,
			&p.EndDate,
		); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}

	return programs, nil
}

// Count mengembalikan jumlah seluruh data program.
func (r *ProgramRepository) Count() (int, error) {
	row := r.DB.QueryRow(`
		SELECT COUNT(*) FROM programs
	`)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// CountActiveToday menghitung jumlah program yang aktif pada tanggal hari ini (start_date..end_date inclusive).
func (r *ProgramRepository) CountActiveToday() (int, error) {
	row := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM programs
		WHERE CURDATE() BETWEEN DATE(start_date) AND DATE(end_date)
	`)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// CountInactiveToday menghitung jumlah program yang tidak aktif pada tanggal hari ini.
// Ini termasuk yang sudah berakhir maupun yang belum dimulai.
func (r *ProgramRepository) CountInactiveToday() (int, error) {
	row := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM programs
		WHERE CURDATE() NOT BETWEEN DATE(start_date) AND DATE(end_date)
	`)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// GetPaginated mengambil data program dengan pagination menggunakan LIMIT dan OFFSET.
// Data yang diambil sudah termasuk join dengan tabel items untuk mendapatkan nama item.
func (r *ProgramRepository) GetPaginated(limit, offset int) ([]models.Program, error) {
	rows, err := r.DB.Query(`
		SELECT p.program_id,
		       p.program_name,
		       p.item_id,
		       i.item_name,
		       DATE_FORMAT(p.start_date, '%d-%m-%Y') AS start_date,
		       DATE_FORMAT(p.end_date, '%d-%m-%Y') AS end_date
		FROM programs p
		JOIN items i ON i.item_id = p.item_id
		ORDER BY p.program_id DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []models.Program

	for rows.Next() {
		var p models.Program
		if err := rows.Scan(
			&p.ProgramID,
			&p.ProgramName,
			&p.ItemID,
			&p.ItemName,
			&p.StartDate,
			&p.EndDate,
		); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}

	return programs, nil
}

func (r *ProgramRepository) Create(p models.Program) error {
	_, err := r.DB.Exec(`
		INSERT INTO programs (program_name, item_id, start_date, end_date)
		VALUES (?, ?, ?, ?)
	`, p.ProgramName, p.ItemID, p.StartDate, p.EndDate)
	return err
}

func (r *ProgramRepository) Update(p models.Program) error {
	_, err := r.DB.Exec(`
		UPDATE programs
		SET program_name = ?,
		    item_id = ?,
		    start_date = ?,
		    end_date = ?
		WHERE program_id = ?
	`, p.ProgramName, p.ItemID, p.StartDate, p.EndDate, p.ProgramID)
	return err
}

func (r *ProgramRepository) Delete(id int) error {
	_, err := r.DB.Exec(`
		DELETE FROM programs
		WHERE program_id = ?
	`, id)
	return err
}
