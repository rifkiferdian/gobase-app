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

// StockOutInfo merepresentasikan qty dan alasan keluarnya hadiah untuk satu item.
type StockOutInfo struct {
	Qty    int
	Reason string
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
	var createdAt time.Time
	err = tx.QueryRow(`
		SELECT id, qty, created_at
		FROM stock_out
		WHERE program_id = ? AND user_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`, programID, userID).Scan(&stockOutID, &currentQty, &createdAt)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	now := time.Now()
	startNewRecord := false

	switch {
	case errors.Is(err, sql.ErrNoRows):
		startNewRecord = true
	case !sameDay(createdAt, now):
		startNewRecord = true
	}

	// Jika belum ada record atau tanggal terakhir berbeda, mulai record baru dan reset qty.
	if startNewRecord {
		if delta < 0 {
			return 0, ErrQuantityZero
		}
		newQty := delta

		res, err := tx.Exec(`
			INSERT INTO stock_out (user_id, program_id, issued_at, qty)
			VALUES (?, ?, ?, ?)
		`, userID, programID, now, newQty)
		if err != nil {
			return 0, err
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}
		stockOutID = int(lastID)
		currentQty = 0
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

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// GetTodayQuantities mengambil qty stock_out per item (program terbaru per item) untuk user & tanggal hari ini.
// Jika itemIDs kosong, maka tidak ada data yang diambil.
func (r *StockOutRepository) GetTodayQuantities(itemIDs []int, userID int) (map[int]StockOutInfo, error) {
	result := map[int]StockOutInfo{}
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
SELECT lp.item_id, COALESCE(so.qty, 0) AS qty, COALESCE(so.reason, '') AS reason
FROM latest_program lp
JOIN items i ON i.item_id = lp.item_id
LEFT JOIN stock_out so ON so.program_id = lp.program_id
    AND so.user_id = ?
    AND DATE(so.created_at) = CURDATE()
    AND (so.reason IS NULL OR TRIM(so.reason) = '')
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
		var reason string
		if err := rows.Scan(&itemID, &qty, &reason); err != nil {
			return nil, err
		}
		result[itemID] = StockOutInfo{
			Qty:    qty,
			Reason: reason,
		}
	}

	return result, nil
}
