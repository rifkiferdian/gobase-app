package repositories

import (
	"database/sql"
	"stok-hadiah/models"
	"strings"
	"time"
)

type StockInRepository struct {
	DB                 *sql.DB
	StoreIDs           []int
	FilterStoreID      *int
	EnforceStoreFilter bool
}

// GetAll mengambil seluruh data stok masuk beserta nama item dan nama petugas.
func (r *StockInRepository) GetAll() ([]models.StockIn, error) {
	args := []interface{}{}
	query := `
		SELECT si.id,
		       si.user_id,
		       u.name,
		       si.item_id,
		       i.item_name,
		       i.store_id,
		       st.store_name,
		       s.supplier_name,
		       si.qty,
		       si.received_at,
		       si.details
		FROM stock_in si
		JOIN users u ON u.id = si.user_id
		JOIN items i ON i.item_id = si.item_id
		LEFT JOIN stores st ON st.store_id = i.store_id
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
	`
	query, skip := r.appendStoreFilter(query, &args, false)
	if skip {
		return []models.StockIn{}, nil
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
			&s.StoreID,
			&s.StoreName,
			&s.SupplierName,
			&s.Qty,
			&s.ReceivedAt,
			&s.Description,
		); err != nil {
			return nil, err
		}

		s.ReceivedAt, s.ReceivedAtDisplay = formatStockInTime(s.ReceivedAt)

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
		       i.store_id,
		       st.store_name,
		       s.supplier_name,
		       si.qty,
		       si.received_at,
		       si.details
		FROM stock_in si
		JOIN users u ON u.id = si.user_id
		JOIN items i ON i.item_id = si.item_id
		LEFT JOIN stores st ON st.store_id = i.store_id
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

	query, skip := r.appendStoreFilter(query, &args, true)
	if skip {
		return []models.StockIn{}, nil
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
			&s.StoreID,
			&s.StoreName,
			&s.SupplierName,
			&s.Qty,
			&s.ReceivedAt,
			&s.Description,
		); err != nil {
			return nil, err
		}

		s.ReceivedAt, s.ReceivedAtDisplay = formatStockInTime(s.ReceivedAt)

		stockIns = append(stockIns, s)
	}

	return stockIns, nil
}

// Count mengembalikan jumlah seluruh data stok masuk.
func (r *StockInRepository) Count() (int, error) {
	args := []interface{}{}
	query, skip := r.appendStoreFilter(`
		SELECT COUNT(*)
		FROM stock_in si
		JOIN items i ON i.item_id = si.item_id
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

// SumQty menghitung total quantity seluruh data stok masuk.
func (r *StockInRepository) SumQty() (int, error) {
	args := []interface{}{}
	query, skip := r.appendStoreFilter(`
		SELECT COALESCE(SUM(si.qty), 0)
		FROM stock_in si
		JOIN items i ON i.item_id = si.item_id
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

// CountToday menghitung total transaksi stok masuk untuk hari ini.
func (r *StockInRepository) CountToday() (int, error) {
	args := []interface{}{}
	query, skip := r.appendStoreFilter(`
		SELECT COUNT(*)
		FROM stock_in si
		JOIN items i ON i.item_id = si.item_id
		WHERE DATE(si.received_at) = CURDATE()
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

// SumTodayQty menghitung total quantity stok masuk untuk hari ini.
func (r *StockInRepository) SumTodayQty() (int, error) {
	args := []interface{}{}
	query, skip := r.appendStoreFilter(`
		SELECT COALESCE(SUM(si.qty), 0)
		FROM stock_in si
		JOIN items i ON i.item_id = si.item_id
		WHERE DATE(si.received_at) = CURDATE()
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

// GetPaginated mengambil data stok masuk dengan pagination menggunakan LIMIT dan OFFSET.
// Data yang diambil sudah termasuk join dengan tabel items, suppliers, dan users.
func (r *StockInRepository) GetPaginated(limit, offset int) ([]models.StockIn, error) {
	args := []interface{}{}
	query := `
		SELECT si.id,
		       si.user_id,
		       u.name,
		       si.item_id,
		       i.item_name,
		       i.store_id,
		       st.store_name,
		       s.supplier_name,
		       si.qty,
		       si.received_at,
		       si.details
		FROM stock_in si
		JOIN users u ON u.id = si.user_id
		JOIN items i ON i.item_id = si.item_id
		LEFT JOIN stores st ON st.store_id = i.store_id
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
	`
	query, skip := r.appendStoreFilter(query, &args, false)
	if skip {
		return []models.StockIn{}, nil
	}

	query += `
		ORDER BY si.received_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, limit, offset)

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
			&s.StoreID,
			&s.StoreName,
			&s.SupplierName,
			&s.Qty,
			&s.ReceivedAt,
			&s.Description,
		); err != nil {
			return nil, err
		}

		s.ReceivedAt, s.ReceivedAtDisplay = formatStockInTime(s.ReceivedAt)

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

	args := []interface{}{s.ItemID, s.Qty, receivedAt, s.Description, s.ID}
	query := `
		UPDATE stock_in si
		JOIN items i ON i.item_id = si.item_id
		SET si.item_id = ?,
		    si.qty = ?,
		    si.received_at = ?,
		    si.details = ?
		WHERE si.id = ?
	`
	query, skip := r.appendStoreFilter(query, &args, true)
	if skip {
		return nil
	}

	_, err := r.DB.Exec(query, args...)
	return err
}

func (r *StockInRepository) Delete(id int) error {
	args := []interface{}{id}
	query := `
		DELETE si
		FROM stock_in si
		JOIN items i ON i.item_id = si.item_id
		WHERE si.id = ?
	`
	query, skip := r.appendStoreFilter(query, &args, true)
	if skip {
		return nil
	}

	_, err := r.DB.Exec(query, args...)
	return err
}

// formatStockInTime mencoba berbagai layout datetime agar tampilan konsisten.
// Mengembalikan pasangan (value untuk input datetime-local, value untuk tampilan tabel).
func formatStockInTime(raw string) (string, string) {
	layouts := []string{
		"2006-01-02 15:04:05", // MySQL DATETIME default
		time.RFC3339,          // e.g. 2025-12-18T00:00:00Z
		"2006-01-02T15:04:05", // ISO tanpa offset
		"2006-01-02",          // tanggal saja
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.Format("2006-01-02T15:04"), t.Format("02-01-2006 15:04:05")
		}
	}
	return raw, raw
}

// appendStoreFilter menambahkan filter store sesuai hak akses StoreIDs dan pilihan filter user.
// Jika EnforceStoreFilter true dan StoreIDs kosong, fungsi akan menandakan skip query (no access).
// hasWhere menentukan apakah query sudah memiliki klausa WHERE sebelumnya.
func (r *StockInRepository) appendStoreFilter(query string, args *[]interface{}, hasWhere bool) (string, bool) {
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
