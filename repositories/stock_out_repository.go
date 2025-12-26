package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrItemNotFound     = errors.New("item tidak ditemukan")
	ErrItemNotAllowed   = errors.New("item tidak tersedia untuk store yang diizinkan")
	ErrProgramNotFound  = errors.New("program untuk item belum diatur")
	ErrQuantityNegative = errors.New("quantity tidak boleh negatif")
	ErrQuantityZero     = errors.New("quantity sudah 0")
)

// StockOutRepository menangani penyimpanan quantity keluar serta log per-aksi.
type StockOutRepository struct {
	DB                 *sql.DB
	StoreIDs           []int
	EnforceStoreFilter bool
}

// AdjustQuantity menambah/mengurangi qty untuk item tertentu dan mencatat eventnya.
// Delta harus bernilai +1 atau -1; fungsi akan mengembalikan qty terkini setelah update.
func (r *StockOutRepository) AdjustQuantity(itemID, delta, userID int) (int, error) {
	if delta == 0 {
		return 0, fmt.Errorf("delta tidak boleh 0")
	}
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return 0, ErrItemNotAllowed
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Validasi item dan akses store
	var storeID int
	if err := tx.QueryRow("SELECT store_id FROM items WHERE item_id = ?", itemID).Scan(&storeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrItemNotFound
		}
		return 0, err
	}
	if len(r.StoreIDs) > 0 && !containsInt(r.StoreIDs, storeID) {
		return 0, ErrItemNotAllowed
	}

	// Ambil program terbaru untuk item ini (diasumsikan satu program aktif per item).
	var programID int
	if err := tx.QueryRow(`
		SELECT program_id
		FROM programs
		WHERE item_id = ?
		ORDER BY program_id DESC
		LIMIT 1
	`, itemID).Scan(&programID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrProgramNotFound
		}
		return 0, err
	}

	// Ambil record stock_out terakhir untuk kombinasi user + program.
	var stockOutID int
	var currentQty int
	err = tx.QueryRow(`
		SELECT id, qty
		FROM stock_out
		WHERE program_id = ? AND user_id = ?
		ORDER BY issued_at DESC
		LIMIT 1
	`, programID, userID).Scan(&stockOutID, &currentQty)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	now := time.Now()

	// Jika belum ada record, hanya boleh menambah.
	if errors.Is(err, sql.ErrNoRows) {
		if delta < 0 {
			return 0, ErrQuantityZero
		}
		res, err := tx.Exec(`
			INSERT INTO stock_out (user_id, program_id, issued_at, qty)
			VALUES (?, ?, ?, ?)
		`, userID, programID, now, delta)
		if err != nil {
			return 0, err
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}
		stockOutID = int(lastID)
		currentQty = 0
	}

	newQty := currentQty + delta
	if newQty < 0 {
		return currentQty, ErrQuantityNegative
	}

	if _, err := tx.Exec(`
		UPDATE stock_out
		SET qty = ?, issued_at = ?
		WHERE id = ?
	`, newQty, now, stockOutID); err != nil {
		return 0, err
	}

	if _, err := tx.Exec(`
		INSERT INTO stock_out_events (stock_out_id, user_id, program_id, item_id, event_time, delta_qty)
		VALUES (?, ?, ?, ?, ?, ?)
	`, stockOutID, userID, programID, itemID, now, delta); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newQty, nil
}

func containsInt(list []int, target int) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

// GetTodayQuantities mengambil qty stock_out per item (program terbaru per item) untuk user & tanggal hari ini.
// Jika itemIDs kosong, maka tidak ada data yang diambil.
func (r *StockOutRepository) GetTodayQuantities(itemIDs []int, userID int) (map[int]int, error) {
	result := map[int]int{}
	if userID == 0 || len(itemIDs) == 0 {
		return result, nil
	}
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, nil
	}

	args := []interface{}{}
	itemPlaceholders := []string{}
	for _, id := range itemIDs {
		itemPlaceholders = append(itemPlaceholders, "?")
		args = append(args, id)
	}

	query := `
WITH latest_program AS (
    SELECT p.item_id, MAX(p.program_id) AS program_id
    FROM programs p
    WHERE p.item_id IN (` + strings.Join(itemPlaceholders, ",") + `)
    GROUP BY p.item_id
)
SELECT lp.item_id, COALESCE(so.qty, 0) AS qty
FROM latest_program lp
JOIN items i ON i.item_id = lp.item_id
LEFT JOIN stock_out so ON so.program_id = lp.program_id
    AND so.user_id = ?
    AND DATE(so.issued_at) = CURDATE()
`
	args = append(args, userID)

	conditions := []string{}
	if len(r.StoreIDs) > 0 {
		placeholders := make([]string, len(r.StoreIDs))
		for i, id := range r.StoreIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		conditions = append(conditions, "i.store_id IN ("+strings.Join(placeholders, ",")+")")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var itemID int
		var qty int
		if err := rows.Scan(&itemID, &qty); err != nil {
			return nil, err
		}
		result[itemID] = qty
	}

	return result, nil
}
