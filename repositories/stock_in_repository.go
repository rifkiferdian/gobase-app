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

		// Parsing string datetime dari database ke format tampilan dan format date untuk form
		if t, err := time.Parse("2006-01-02 15:04:05", s.ReceivedAt); err == nil {
			// Untuk input <input type="date"> butuh format YYYY-MM-DD
			s.ReceivedAt = t.Format("2006-01-02")
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

	_, err := r.DB.Exec(`
		INSERT INTO stock_in (user_id, item_id, received_at, qty, details)
		VALUES (?, ?, ?, ?, ?)
	`, s.UserID, s.ItemID, s.ReceivedAt, s.Qty, s.Description)
	return err
}

func (r *StockInRepository) Update(s models.StockIn) error {
	_, err := r.DB.Exec(`
		UPDATE stock_in
		SET item_id = ?,
		    qty = ?,
		    received_at = ?,
		    details = ?
		WHERE id = ?
	`, s.ItemID, s.Qty, s.ReceivedAt, s.Description, s.ID)
	return err
}

func (r *StockInRepository) Delete(id int) error {
	_, err := r.DB.Exec(`
		DELETE FROM stock_in
		WHERE id = ?
	`, id)
	return err
}
