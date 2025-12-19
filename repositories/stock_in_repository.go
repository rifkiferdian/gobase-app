package repositories

import (
	"database/sql"
	"stok-hadiah/models"
	"time"
)

type StockInRepository struct {
	DB *sql.DB
}

// GetAll mengambil seluruh data stok masuk beserta nama item dan nama petugas.
func (r *StockInRepository) GetAll() ([]models.StockIn, error) {
	rows, err := r.DB.Query(`
		SELECT si.id,
		       si.user_id,
		       u.name,
		       si.item_id,
		       i.item_name,
		       s.supplier_name,
		       si.qty,
		       si.received_at,
		       si.details
		FROM stock_in si
		JOIN users u ON u.id = si.user_id
		JOIN items i ON i.item_id = si.item_id
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
		ORDER BY si.received_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stockIns []models.StockIn

	for rows.Next() {
		var s models.StockIn
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.UserName,
			&s.ItemID,
			&s.ItemName,
			&s.SupplierName,
			&s.Qty,
			&s.ReceivedAt,
			&s.Description,
		); err != nil {
			return nil, err
		}

		// Parsing string datetime dari database ke format tampilan dan format datetime-local untuk form
		if t, err := time.Parse("2006-01-02 15:04:05", s.ReceivedAt); err == nil {
			// Untuk input <input type="datetime-local"> butuh format YYYY-MM-DDTHH:MM
			s.ReceivedAt = t.Format("2006-01-02T15:04")
			// Untuk tampilan di tabel: dd-mm-YYYY HH:MM:SS
			s.ReceivedAtDisplay = t.Format("02-01-2006 15:04:05")
		} else {
			// Fallback kalau parsing gagal: pakai string apa adanya
			s.ReceivedAtDisplay = s.ReceivedAt
		}

		stockIns = append(stockIns, s)
	}

	return stockIns, nil
}

// Search mengambil data stok masuk yang difilter berdasarkan nama barang dan/atau
// tanggal (bagian tanggal dari kolom received_at). Jika kedua parameter kosong,
// maka fungsi akan mengembalikan seluruh data seperti GetAll.
func (r *StockInRepository) Search(itemName, date string) ([]models.StockIn, error) {
	query := `
		SELECT si.id,
		       si.user_id,
		       u.name,
		       si.item_id,
		       i.item_name,
		       s.supplier_name,
		       si.qty,
		       si.received_at,
		       si.details
		FROM stock_in si
		JOIN users u ON u.id = si.user_id
		JOIN items i ON i.item_id = si.item_id
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
		WHERE 1=1`

	args := []interface{}{}

	if itemName != "" {
		query += " AND i.item_name LIKE ?"
		args = append(args, "%"+itemName+"%")
	}

	if date != "" {
		// Hanya cocokkan bagian tanggal dari kolom received_at
		query += " AND DATE(si.received_at) = ?"
		args = append(args, date)
	}

	query += " ORDER BY si.received_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stockIns []models.StockIn

	for rows.Next() {
		var s models.StockIn
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.UserName,
			&s.ItemID,
			&s.ItemName,
			&s.SupplierName,
			&s.Qty,
			&s.ReceivedAt,
			&s.Description,
		); err != nil {
			return nil, err
		}

		// Parsing string datetime dari database ke format tampilan dan format datetime-local untuk form
		if t, err := time.Parse("2006-01-02 15:04:05", s.ReceivedAt); err == nil {
			// Untuk input <input type="datetime-local"> butuh format YYYY-MM-DDTHH:MM
			s.ReceivedAt = t.Format("2006-01-02T15:04")
			// Untuk tampilan di tabel: dd-mm-YYYY HH:MM:SS
			s.ReceivedAtDisplay = t.Format("02-01-2006 15:04:05")
		} else {
			// Fallback kalau parsing gagal: pakai string apa adanya
			s.ReceivedAtDisplay = s.ReceivedAt
		}

		stockIns = append(stockIns, s)
	}

	return stockIns, nil
}

// Count mengembalikan jumlah seluruh data stok masuk.
func (r *StockInRepository) Count() (int, error) {
	row := r.DB.QueryRow(`
		SELECT COUNT(*) FROM stock_in
	`)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// SumQty menghitung total quantity seluruh data stok masuk.
func (r *StockInRepository) SumQty() (int, error) {
	row := r.DB.QueryRow(`
		SELECT COALESCE(SUM(qty), 0) FROM stock_in
	`)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// CountToday menghitung total transaksi stok masuk untuk hari ini.
func (r *StockInRepository) CountToday() (int, error) {
	row := r.DB.QueryRow(`
		SELECT COUNT(*) FROM stock_in WHERE DATE(received_at) = CURDATE()
	`)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// SumTodayQty menghitung total quantity stok masuk untuk hari ini.
func (r *StockInRepository) SumTodayQty() (int, error) {
	row := r.DB.QueryRow(`
		SELECT COALESCE(SUM(qty), 0) FROM stock_in WHERE DATE(received_at) = CURDATE()
	`)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// GetPaginated mengambil data stok masuk dengan pagination menggunakan LIMIT dan OFFSET.
// Data yang diambil sudah termasuk join dengan tabel items, suppliers, dan users.
func (r *StockInRepository) GetPaginated(limit, offset int) ([]models.StockIn, error) {
	rows, err := r.DB.Query(`
		SELECT si.id,
		       si.user_id,
		       u.name,
		       si.item_id,
		       i.item_name,
		       s.supplier_name,
		       si.qty,
		       si.received_at,
		       si.details
		FROM stock_in si
		JOIN users u ON u.id = si.user_id
		JOIN items i ON i.item_id = si.item_id
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
		ORDER BY si.received_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stockIns []models.StockIn

	for rows.Next() {
		var s models.StockIn
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.UserName,
			&s.ItemID,
			&s.ItemName,
			&s.SupplierName,
			&s.Qty,
			&s.ReceivedAt,
			&s.Description,
		); err != nil {
			return nil, err
		}

		// Parsing string datetime dari database ke format tampilan dan format datetime-local untuk form
		if t, err := time.Parse("2006-01-02 15:04:05", s.ReceivedAt); err == nil {
			// Untuk input <input type="datetime-local"> butuh format YYYY-MM-DDTHH:MM
			s.ReceivedAt = t.Format("2006-01-02T15:04")
			// Untuk tampilan di tabel: dd-mm-YYYY HH:MM:SS
			s.ReceivedAtDisplay = t.Format("02-01-2006 15:04:05")
		} else {
			// Fallback kalau parsing gagal: pakai string apa adanya
			s.ReceivedAtDisplay = s.ReceivedAt
		}

		stockIns = append(stockIns, s)
	}

	return stockIns, nil
}

func (r *StockInRepository) Create(s models.StockIn) error {
	if s.ReceivedAt == "" {
		_, err := r.DB.Exec(`
			INSERT INTO stock_in (user_id, item_id, received_at, qty, details)
			VALUES (?, ?, NOW(), ?, ?)
		`, s.UserID, s.ItemID, s.Qty, s.Description)
		return err
	}

	// Normalisasi format datetime dari form menjadi format MySQL (YYYY-MM-DD HH:MM:SS)
	receivedAt := s.ReceivedAt
	if t, err := time.Parse("2006-01-02T15:04", s.ReceivedAt); err == nil {
		receivedAt = t.Format("2006-01-02 15:04:05")
	} else if t, err := time.Parse("2006-01-02", s.ReceivedAt); err == nil {
		// Jika hanya tanggal yang dikirim, set jam ke 00:00:00
		receivedAt = t.Format("2006-01-02 15:04:05")
	}

	_, err := r.DB.Exec(`
		INSERT INTO stock_in (user_id, item_id, received_at, qty, details)
		VALUES (?, ?, ?, ?, ?)
	`, s.UserID, s.ItemID, receivedAt, s.Qty, s.Description)
	return err
}

func (r *StockInRepository) Update(s models.StockIn) error {
	// Normalisasi format datetime dari form menjadi format MySQL (YYYY-MM-DD HH:MM:SS)
	receivedAt := s.ReceivedAt
	if t, err := time.Parse("2006-01-02T15:04", s.ReceivedAt); err == nil {
		receivedAt = t.Format("2006-01-02 15:04:05")
	} else if t, err := time.Parse("2006-01-02", s.ReceivedAt); err == nil {
		// Jika hanya tanggal yang dikirim, set jam ke 00:00:00
		receivedAt = t.Format("2006-01-02 15:04:05")
	}

	_, err := r.DB.Exec(`
		UPDATE stock_in
		SET item_id = ?,
		    qty = ?,
		    received_at = ?,
		    details = ?
		WHERE id = ?
	`, s.ItemID, s.Qty, receivedAt, s.Description, s.ID)
	return err
}

func (r *StockInRepository) Delete(id int) error {
	_, err := r.DB.Exec(`
		DELETE FROM stock_in
		WHERE id = ?
	`, id)
	return err
}
