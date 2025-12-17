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
