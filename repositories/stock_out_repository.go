package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"stok-hadiah/models"
)

var (
	ErrItemNotFound     = errors.New("item tidak ditemukan")
	ErrItemNotAllowed   = errors.New("item tidak tersedia untuk store yang diizinkan")
	ErrProgramNotFound  = errors.New("program untuk item belum diatur")
	ErrQuantityNegative = errors.New("quantity tidak boleh negatif")
	ErrQuantityZero     = errors.New("quantity sudah 0")
	ErrStockOutNotFound = errors.New("stock out tidak ditemukan")
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

// Count mengembalikan jumlah baris stock_out sesuai filter store (jika ada).
func (r *StockOutRepository) Count() (int, error) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return 0, nil
	}

	args := []interface{}{}
	query := `
SELECT COUNT(*)
FROM stock_out so
JOIN programs p ON p.program_id = so.program_id
JOIN items i ON i.item_id = p.item_id
WHERE 1=1
`

	if len(r.StoreIDs) > 0 {
		placeholders := make([]string, len(r.StoreIDs))
		for i, id := range r.StoreIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		query += " AND i.store_id IN (" + strings.Join(placeholders, ",") + ")"
	}

	row := r.DB.QueryRow(query, args...)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// SumQty menghitung total qty stock out sesuai filter store (jika ada).
func (r *StockOutRepository) SumQty() (int, error) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return 0, nil
	}

	args := []interface{}{}
	query := `
SELECT COALESCE(SUM(so.qty), 0)
FROM stock_out so
JOIN programs p ON p.program_id = so.program_id
JOIN items i ON i.item_id = p.item_id
WHERE 1=1
`

	if len(r.StoreIDs) > 0 {
		placeholders := make([]string, len(r.StoreIDs))
		for i, id := range r.StoreIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		query += " AND i.store_id IN (" + strings.Join(placeholders, ",") + ")"
	}

	row := r.DB.QueryRow(query, args...)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// CountToday menghitung jumlah baris stock out untuk hari ini.
func (r *StockOutRepository) CountToday() (int, error) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return 0, nil
	}

	args := []interface{}{}
	query := `
SELECT COUNT(*)
FROM stock_out so
JOIN programs p ON p.program_id = so.program_id
JOIN items i ON i.item_id = p.item_id
WHERE DATE(so.issued_at) = CURDATE()
`

	if len(r.StoreIDs) > 0 {
		placeholders := make([]string, len(r.StoreIDs))
		for i, id := range r.StoreIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		query += " AND i.store_id IN (" + strings.Join(placeholders, ",") + ")"
	}

	row := r.DB.QueryRow(query, args...)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// SumTodayQty menghitung total qty stock out untuk hari ini.
func (r *StockOutRepository) SumTodayQty() (int, error) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return 0, nil
	}

	args := []interface{}{}
	query := `
SELECT COALESCE(SUM(so.qty), 0)
FROM stock_out so
JOIN programs p ON p.program_id = so.program_id
JOIN items i ON i.item_id = p.item_id
WHERE DATE(so.issued_at) = CURDATE()
`

	if len(r.StoreIDs) > 0 {
		placeholders := make([]string, len(r.StoreIDs))
		for i, id := range r.StoreIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		query += " AND i.store_id IN (" + strings.Join(placeholders, ",") + ")"
	}

	row := r.DB.QueryRow(query, args...)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// DeleteCaseStockOut menghapus entry stock_out dengan alasan (kasus khusus) dan seluruh event terkait.
// Hanya mengizinkan penghapusan milik user yang sama dan store yang diizinkan.
func (r *StockOutRepository) DeleteCaseStockOut(id, userID int) (models.StockOutCase, error) {
	result := models.StockOutCase{}
	if id <= 0 || userID <= 0 {
		return result, fmt.Errorf("id atau user tidak valid")
	}
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, ErrItemNotAllowed
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	var storeID int
	var reason sql.NullString
	err = tx.QueryRow(`
		SELECT
			so.program_id,
			p.item_id,
			i.store_id,
			i.item_name,
			so.qty,
			so.reason,
			so.issued_at
		FROM stock_out so
		JOIN programs p ON p.program_id = so.program_id
		JOIN items i ON i.item_id = p.item_id
		WHERE so.id = ? AND so.user_id = ?
	`, id, userID).Scan(
		&result.ProgramID,
		&result.ItemID,
		&storeID,
		&result.ItemName,
		&result.Qty,
		&reason,
		&result.IssuedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, ErrStockOutNotFound
		}
		return result, err
	}

	if len(r.StoreIDs) > 0 && !containsInt(r.StoreIDs, storeID) {
		return result, ErrItemNotAllowed
	}

	result.ID = id
	result.UserID = userID
	result.Reason = strings.TrimSpace(reason.String)

	if result.Reason == "" {
		return result, fmt.Errorf("data bukan kasus khusus sehingga tidak bisa dihapus di sini")
	}

	if _, err := tx.Exec(`DELETE FROM stock_out_events WHERE stock_out_id = ?`, id); err != nil {
		return result, err
	}

	if _, err := tx.Exec(`DELETE FROM stock_out WHERE id = ? AND user_id = ?`, id, userID); err != nil {
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, err
	}

	return result, nil
}

// CreateCaseStockOut menyimpan pengeluaran stok dengan alasan khusus ke tabel stock_out dan mencatat event.
func (r *StockOutRepository) CreateCaseStockOut(itemID, qty, userID int, reason string) (models.StockOutCase, error) {
	var result models.StockOutCase
	if qty <= 0 {
		return result, fmt.Errorf("qty harus lebih dari 0")
	}
	if strings.TrimSpace(reason) == "" {
		return result, fmt.Errorf("alasan wajib diisi")
	}
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, ErrItemNotAllowed
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	// Validasi item dan akses store
	var storeID int
	var itemName string
	if err := tx.QueryRow("SELECT store_id, item_name FROM items WHERE item_id = ?", itemID).Scan(&storeID, &itemName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, ErrItemNotFound
		}
		return result, err
	}
	if len(r.StoreIDs) > 0 && !containsInt(r.StoreIDs, storeID) {
		return result, ErrItemNotAllowed
	}

	// Ambil program terbaru untuk item ini.
	var programID int
	if err := tx.QueryRow(`
		SELECT program_id
		FROM programs
		WHERE item_id = ?
		ORDER BY program_id DESC
		LIMIT 1
	`, itemID).Scan(&programID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, ErrProgramNotFound
		}
		return result, err
	}

	now := time.Now()
	res, err := tx.Exec(`
		INSERT INTO stock_out (user_id, program_id, issued_at, qty, reason, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, userID, programID, now, qty, reason, now, now)
	if err != nil {
		return result, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return result, err
	}

	if _, err := tx.Exec(`
		INSERT INTO stock_out_events (stock_out_id, user_id, program_id, item_id, event_time, delta_qty)
		VALUES (?, ?, ?, ?, ?, ?)
	`, lastID, userID, programID, itemID, now, qty); err != nil {
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, err
	}

	result = models.StockOutCase{
		ID:        int(lastID),
		UserID:    userID,
		ProgramID: programID,
		ItemID:    itemID,
		ItemName:  itemName,
		Qty:       qty,
		Reason:    reason,
		IssuedAt:  now,
	}
	return result, nil
}

// ListCaseStockOuts mengembalikan daftar stock out yang memiliki alasan/keterangan (kasus khusus).
// Jika limit <= 0 maka default 20.
func (r *StockOutRepository) ListCaseStockOuts(userID, limit int) ([]models.StockOutCase, error) {
	result := []models.StockOutCase{}
	if userID == 0 {
		return result, nil
	}
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, nil
	}
	if limit <= 0 {
		limit = 20
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	args := []interface{}{userID, startOfDay, endOfDay}
	query := `
SELECT
	so.id,
	so.user_id,
	so.program_id,
	p.item_id,
	i.item_name,
	so.qty,
	so.reason,
	so.issued_at
FROM stock_out so
JOIN programs p ON p.program_id = so.program_id
JOIN items i ON i.item_id = p.item_id
WHERE so.user_id = ?
  AND so.created_at >= ?
  AND so.created_at < ?
  AND TRIM(COALESCE(so.reason, '')) <> ''
`

	if len(r.StoreIDs) > 0 {
		placeholders := make([]string, len(r.StoreIDs))
		for i, id := range r.StoreIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		query += " AND i.store_id IN (" + strings.Join(placeholders, ",") + ")"
	}

	query += " ORDER BY so.issued_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.StockOutCase
		if err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.ProgramID,
			&entry.ItemID,
			&entry.ItemName,
			&entry.Qty,
			&entry.Reason,
			&entry.IssuedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, entry)
	}

	return result, nil
}

// DailyTotalsSince mengembalikan total qty stock out per hari sejak tanggal start (inklusif).
// Data dijumlahkan berdasarkan issued_at agar mencerminkan waktu barang keluar.
func (r *StockOutRepository) DailyTotalsSince(start time.Time) (map[string]int, error) {
	result := map[string]int{}
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, nil
	}

	args := []interface{}{start}
	query := `
SELECT DATE(so.issued_at) AS day,
       COALESCE(SUM(so.qty), 0) AS total_qty
FROM stock_out so
JOIN programs p ON p.program_id = so.program_id
JOIN items i ON i.item_id = p.item_id
WHERE so.issued_at >= ?
`

	if len(r.StoreIDs) > 0 {
		placeholders := make([]string, len(r.StoreIDs))
		for i, id := range r.StoreIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		query += " AND i.store_id IN (" + strings.Join(placeholders, ",") + ")"
	}

	query += `
GROUP BY day
ORDER BY day
`

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var day string
		var total int
		if err := rows.Scan(&day, &total); err != nil {
			return nil, err
		}
		result[day] = total
	}

	return result, nil
}

// MonthlyTotalsSince mengembalikan total qty stock out per bulan sejak tanggal start (inklusif).
// Data dijumlahkan berdasarkan issued_at agar mencerminkan waktu barang keluar.
func (r *StockOutRepository) MonthlyTotalsSince(start time.Time) (map[string]int, error) {
	result := map[string]int{}
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, nil
	}

	args := []interface{}{start}
	query := `
SELECT DATE_FORMAT(so.issued_at, '%Y-%m-01') AS month_start,
       COALESCE(SUM(so.qty), 0) AS total_qty
FROM stock_out so
JOIN programs p ON p.program_id = so.program_id
JOIN items i ON i.item_id = p.item_id
WHERE so.issued_at >= ?
`

	if len(r.StoreIDs) > 0 {
		placeholders := make([]string, len(r.StoreIDs))
		for i, id := range r.StoreIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		query += " AND i.store_id IN (" + strings.Join(placeholders, ",") + ")"
	}

	query += `
GROUP BY month_start
ORDER BY month_start
`

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var month string
		var total int
		if err := rows.Scan(&month, &total); err != nil {
			return nil, err
		}
		result[month] = total
	}

	return result, nil
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
),
today_stock_out AS (
    SELECT so.program_id, SUM(so.qty) AS qty
    FROM stock_out so
    WHERE so.user_id = ?
      AND DATE(so.created_at) = CURDATE()
      AND (so.reason IS NULL OR TRIM(so.reason) = '')
    GROUP BY so.program_id
)
SELECT lp.item_id, COALESCE(ts.qty, 0) AS qty, '' AS reason
FROM latest_program lp
JOIN items i ON i.item_id = lp.item_id
LEFT JOIN today_stock_out ts ON ts.program_id = lp.program_id
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

// GetQuantityBeforeToday mengembalikan total qty keluar per item sebelum hari ini.
func (r *StockOutRepository) GetQuantityBeforeToday(itemIDs []int) (map[int]int, error) {
	result := map[int]int{}
	if len(itemIDs) == 0 {
		return result, nil
	}
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, nil
	}

	args := []interface{}{}
	itemPlaceholders := make([]string, 0, len(itemIDs))
	for _, id := range itemIDs {
		itemPlaceholders = append(itemPlaceholders, "?")
		args = append(args, id)
	}

	query := `
SELECT p.item_id, COALESCE(SUM(so.qty), 0) AS qty_out_before_today
FROM stock_out so
JOIN programs p ON p.program_id = so.program_id
JOIN items i ON i.item_id = p.item_id
WHERE p.item_id IN (` + strings.Join(itemPlaceholders, ",") + `)
  AND DATE(so.created_at) < CURDATE()
`
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
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " GROUP BY p.item_id"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var itemID int
		var qtyBefore int
		if err := rows.Scan(&itemID, &qtyBefore); err != nil {
			return nil, err
		}
		result[itemID] = qtyBefore
	}

	return result, nil
}
