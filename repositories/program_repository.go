package repositories

import (
	"database/sql"
	"stok-hadiah/models"
	"strings"
)

type ProgramRepository struct {
	DB                 *sql.DB
	StoreIDs           []int
	FilterStoreID      *int
	EnforceStoreFilter bool
}

// GetAll mengambil seluruh data program beserta nama item-nya.
func (r *ProgramRepository) GetAll() ([]models.Program, error) {
	args := []interface{}{}
	query := `
		SELECT p.program_id,
		       p.program_name,
		       p.item_id,
		       i.item_name,
		       st.store_name,
		       DATE_FORMAT(p.start_date, '%Y-%m-%d') AS start_date,
		       DATE_FORMAT(p.end_date, '%Y-%m-%d') AS end_date
		FROM programs p
		JOIN items i ON i.item_id = p.item_id
		LEFT JOIN stores st ON st.store_id = i.store_id
	`
	query, skip := r.appendStoreFilter(query, &args, false)
	if skip {
		return []models.Program{}, nil
	}

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
			&p.StoreName,
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
		       st.store_name,
		       DATE_FORMAT(p.start_date, '%d-%m-%Y') AS start_date,
		       DATE_FORMAT(p.end_date, '%d-%m-%Y') AS end_date
		FROM programs p
		JOIN items i ON i.item_id = p.item_id
		LEFT JOIN stores st ON st.store_id = i.store_id
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

	query, skip := r.appendStoreFilter(query, &args, true)
	if skip {
		return []models.Program{}, nil
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
			&p.StoreName,
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
	args := []interface{}{}
	query, skip := r.appendStoreFilter(`
		SELECT COUNT(*) FROM programs p
		JOIN items i ON i.item_id = p.item_id
	`, &args, false)
	if skip {
		return 0, nil
	}

	row := r.DB.QueryRow(query, args...)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// CountActiveToday menghitung jumlah program yang aktif pada tanggal hari ini (start_date..end_date inclusive).
func (r *ProgramRepository) CountActiveToday() (int, error) {
	args := []interface{}{}
	query, skip := r.appendStoreFilter(`
		SELECT COUNT(*)
		FROM programs p
		JOIN items i ON i.item_id = p.item_id
		WHERE CURDATE() BETWEEN DATE(p.start_date) AND DATE(p.end_date)
	`, &args, true)
	if skip {
		return 0, nil
	}

	row := r.DB.QueryRow(query, args...)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// CountInactiveToday menghitung jumlah program yang tidak aktif pada tanggal hari ini.
// Ini termasuk yang sudah berakhir maupun yang belum dimulai.
func (r *ProgramRepository) CountInactiveToday() (int, error) {
	args := []interface{}{}
	query, skip := r.appendStoreFilter(`
		SELECT COUNT(*)
		FROM programs p
		JOIN items i ON i.item_id = p.item_id
		WHERE CURDATE() NOT BETWEEN DATE(p.start_date) AND DATE(p.end_date)
	`, &args, true)
	if skip {
		return 0, nil
	}

	row := r.DB.QueryRow(query, args...)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// GetPaginated mengambil data program dengan pagination menggunakan LIMIT dan OFFSET.
// Data yang diambil sudah termasuk join dengan tabel items untuk mendapatkan nama item.
func (r *ProgramRepository) GetPaginated(limit, offset int) ([]models.Program, error) {
	args := []interface{}{}
	query := `
		SELECT p.program_id,
		       p.program_name,
		       p.item_id,
		       i.item_name,
		       st.store_name,
		       DATE_FORMAT(p.start_date, '%d-%m-%Y') AS start_date,
		       DATE_FORMAT(p.end_date, '%d-%m-%Y') AS end_date
		FROM programs p
		JOIN items i ON i.item_id = p.item_id
		LEFT JOIN stores st ON st.store_id = i.store_id
	`
	query, skip := r.appendStoreFilter(query, &args, false)
	if skip {
		return []models.Program{}, nil
	}
	query += `
		ORDER BY p.program_id DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, limit, offset)

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
			&p.StoreName,
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

// appendStoreFilter menambahkan filter store sesuai hak akses StoreIDs dan pilihan filter user.
// Jika EnforceStoreFilter true dan StoreIDs kosong, fungsi akan menandakan skip query (no access).
// hasWhere menentukan apakah query sudah mengandung WHERE.
func (r *ProgramRepository) appendStoreFilter(query string, args *[]interface{}, hasWhere bool) (string, bool) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return "", true
	}

	if len(r.StoreIDs) > 0 {
		placeholders := make([]string, len(r.StoreIDs))
		for i, id := range r.StoreIDs {
			placeholders[i] = "?"
			*args = append(*args, id)
		}
		keyword := " WHERE "
		if hasWhere {
			keyword = " AND "
		}
		query += keyword + "i.store_id IN (" + strings.Join(placeholders, ",") + ")"
		hasWhere = true
	}

	if r.FilterStoreID != nil {
		keyword := " WHERE "
		if hasWhere {
			keyword = " AND "
		}
		query += keyword + "i.store_id = ?"
		*args = append(*args, *r.FilterStoreID)
	}

	return query, false
}
