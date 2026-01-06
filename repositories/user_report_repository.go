package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"stok-hadiah/models"
)

// UserReportRepository menghitung ringkasan stok per user dengan filter store opsional.
type UserReportRepository struct {
	DB                 *sql.DB
	StoreIDs           []int
	EnforceStoreFilter bool
}

// GetSummaries mengembalikan daftar ringkasan per user (jumlah jenis item, total masuk & keluar).
func (r *UserReportRepository) GetSummaries() ([]models.UserReportSummary, error) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return []models.UserReportSummary{}, nil
	}

	users, err := r.fetchUsers()
	if err != nil {
		return nil, err
	}

	allowed := make(map[int]bool)
	for _, id := range r.StoreIDs {
		allowed[id] = true
	}

	allStoreIDs := make(map[int]bool)
	for i := range users {
		users[i].StoreIDs = filterStoreIDs(users[i].StoreIDs, allowed)
		for _, id := range users[i].StoreIDs {
			allStoreIDs[id] = true
		}
	}

	storeNameMap, err := r.buildStoreNameMap(allStoreIDs)
	if err != nil {
		return nil, err
	}

	stockInTotals, err := r.fetchStockInTotals()
	if err != nil {
		return nil, err
	}

	stockOutTotals, err := r.fetchStockOutTotals()
	if err != nil {
		return nil, err
	}

	itemCounts, err := r.fetchItemCounts()
	if err != nil {
		return nil, err
	}

	var reports []models.UserReportSummary
	for _, u := range users {
		storeNames := joinStoreNames(u.StoreIDs, storeNameMap)
		reports = append(reports, models.UserReportSummary{
			UserID:     u.ID,
			NIP:        u.NIP,
			Name:       u.Name,
			StoreIDs:   u.StoreIDs,
			StoreNames: storeNames,
			ItemTypes:  itemCounts[u.ID],
			TotalIn:    stockInTotals[u.ID],
			TotalOut:   stockOutTotals[u.ID],
		})
	}

	return reports, nil
}

// GetDetail mengambil detail stok masuk dan keluar untuk satu user.
func (r *UserReportRepository) GetDetail(userID int, itemName, date string) (models.UserReportDetail, error) {
	detail := models.UserReportDetail{}
	if userID <= 0 {
		return detail, nil
	}

	user, err := r.fetchUserByID(userID)
	if err != nil {
		return detail, err
	}
	if user.ID == 0 {
		return detail, nil
	}

	allowed := make(map[int]bool)
	for _, id := range r.StoreIDs {
		allowed[id] = true
	}

	user.StoreIDs = filterStoreIDs(user.StoreIDs, allowed)

	storeNameMap, err := r.buildStoreNameMap(storeIDsToMap(user.StoreIDs))
	if err != nil {
		return detail, err
	}

	detail.UserID = user.ID
	detail.NIP = user.NIP
	detail.Name = user.Name
	detail.StoreIDs = user.StoreIDs
	detail.StoreNames = joinStoreNames(user.StoreIDs, storeNameMap)

	// Jika akses store dikunci dan kosong, kembalikan info user saja.
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return detail, nil
	}

	stockIns, totalIn, err := r.fetchStockInsByUser(user.ID, itemName, date)
	if err != nil {
		return detail, err
	}

	stockOuts, totalOut, err := r.fetchStockOutByUser(user.ID, itemName, date)
	if err != nil {
		return detail, err
	}

	detail.StockIns = stockIns
	detail.StockOuts = stockOuts
	detail.TotalIn = totalIn
	detail.TotalOut = totalOut

	return detail, nil
}

type userReportUser struct {
	ID       int
	NIP      int
	Name     string
	StoreIDs []int
}

func (r *UserReportRepository) fetchUserByID(userID int) (userReportUser, error) {
	var (
		user     userReportUser
		storeRaw string
	)

	if userID <= 0 {
		return userReportUser{}, nil
	}

	err := r.DB.QueryRow(`
		SELECT id, nip, name, store_id
		FROM users
		WHERE id = ?
	`, userID).Scan(&user.ID, &user.NIP, &user.Name, &storeRaw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userReportUser{}, nil
		}
		return userReportUser{}, err
	}

	user.StoreIDs = parseStoreIDs(storeRaw)
	return user, nil
}

func (r *UserReportRepository) fetchUsers() ([]userReportUser, error) {
	rows, err := r.DB.Query(`
		SELECT id, nip, name, store_id
		FROM users
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []userReportUser
	for rows.Next() {
		var (
			u        userReportUser
			storeRaw string
		)
		if err := rows.Scan(&u.ID, &u.NIP, &u.Name, &storeRaw); err != nil {
			return nil, err
		}
		u.StoreIDs = parseStoreIDs(storeRaw)
		users = append(users, u)
	}

	return users, rows.Err()
}

func (r *UserReportRepository) buildStoreNameMap(storeIDs map[int]bool) (map[int]string, error) {
	result := make(map[int]string)
	if len(storeIDs) == 0 {
		return result, nil
	}

	idList := make([]int, 0, len(storeIDs))
	for id := range storeIDs {
		idList = append(idList, id)
	}

	storeRepo := StoreRepository{DB: r.DB}
	stores, err := storeRepo.GetByIDs(idList)
	if err != nil {
		return nil, err
	}

	for _, st := range stores {
		result[st.StoreID] = st.StoreName
	}

	return result, nil
}

func (r *UserReportRepository) fetchStockInsByUser(userID int, itemName, date string) ([]models.StockIn, int, error) {
	result := []models.StockIn{}
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, 0, nil
	}

	args := []interface{}{userID}
	conditions := []string{"si.user_id = ?"}

	if itemName != "" {
		conditions = append(conditions, "i.item_name LIKE ?")
		args = append(args, "%"+itemName+"%")
	}

	if date != "" {
		conditions = append(conditions, "DATE(si.received_at) = ?")
		args = append(args, date)
	}

	if cond, condArgs := r.buildStoreFilter("i.store_id"); cond != "" {
		conditions = append(conditions, cond)
		args = append(args, condArgs...)
	}

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

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY si.received_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		var entry models.StockIn
		if err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.UserName,
			&entry.ItemID,
			&entry.ItemName,
			&entry.StoreID,
			&entry.StoreName,
			&entry.SupplierName,
			&entry.Qty,
			&entry.ReceivedAt,
			&entry.Description,
		); err != nil {
			return nil, total, err
		}

		entry.ReceivedAt, entry.ReceivedAtDisplay = formatStockInTime(entry.ReceivedAt)
		total += entry.Qty
		result = append(result, entry)
	}

	return result, total, rows.Err()
}

func (r *UserReportRepository) fetchStockOutByUser(userID int, itemName, date string) ([]models.StockOutDetail, int, error) {
	result := []models.StockOutDetail{}
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, 0, nil
	}

	args := []interface{}{userID}
	conditions := []string{"so.user_id = ?"}

	if itemName != "" {
		conditions = append(conditions, "i.item_name LIKE ?")
		args = append(args, "%"+itemName+"%")
	}

	if date != "" {
		conditions = append(conditions, "DATE(so.issued_at) = ?")
		args = append(args, date)
	}

	if cond, condArgs := r.buildStoreFilter("i.store_id"); cond != "" {
		conditions = append(conditions, cond)
		args = append(args, condArgs...)
	}

	query := `
		SELECT
			so.id,
			so.user_id,
			u.name,
			so.program_id,
			p.program_name,
			p.item_id,
			i.item_name,
			i.store_id,
			st.store_name,
			s.supplier_name,
			so.qty,
			so.issued_at,
			COALESCE(so.reason, '')
		FROM stock_out so
		JOIN programs p ON p.program_id = so.program_id
		JOIN items i ON i.item_id = p.item_id
		LEFT JOIN stores st ON st.store_id = i.store_id
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
		JOIN users u ON u.id = so.user_id
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY so.issued_at DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		var entry models.StockOutDetail
		if err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.UserName,
			&entry.ProgramID,
			&entry.ProgramName,
			&entry.ItemID,
			&entry.ItemName,
			&entry.StoreID,
			&entry.StoreName,
			&entry.SupplierName,
			&entry.Qty,
			&entry.IssuedAt,
			&entry.Reason,
		); err != nil {
			return nil, total, err
		}

		entry.IssuedAt, entry.IssuedAtDisplay = formatStockOutTime(entry.IssuedAt)
		entry.Reason = strings.TrimSpace(entry.Reason)
		total += entry.Qty
		result = append(result, entry)
	}

	return result, total, rows.Err()
}

func (r *UserReportRepository) fetchStockInTotals() (map[int]int, error) {
	result := make(map[int]int)
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, nil
	}

	args := []interface{}{}
	query := `
		SELECT si.user_id, COALESCE(SUM(si.qty), 0) AS total_in
		FROM stock_in si
		JOIN items i ON i.item_id = si.item_id
	`

	if cond, condArgs := r.buildStoreFilter("i.store_id"); cond != "" {
		query += " WHERE " + cond
		args = append(args, condArgs...)
	}

	query += " GROUP BY si.user_id"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID, total int
		if err := rows.Scan(&userID, &total); err != nil {
			return nil, err
		}
		result[userID] = total
	}

	return result, rows.Err()
}

func (r *UserReportRepository) fetchStockOutTotals() (map[int]int, error) {
	result := make(map[int]int)
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, nil
	}

	args := []interface{}{}
	query := `
		SELECT so.user_id, COALESCE(SUM(so.qty), 0) AS total_out
		FROM stock_out so
		JOIN programs p ON p.program_id = so.program_id
		JOIN items i ON i.item_id = p.item_id
	`

	if cond, condArgs := r.buildStoreFilter("i.store_id"); cond != "" {
		query += " WHERE " + cond
		args = append(args, condArgs...)
	}

	query += " GROUP BY so.user_id"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID, total int
		if err := rows.Scan(&userID, &total); err != nil {
			return nil, err
		}
		result[userID] = total
	}

	return result, rows.Err()
}

func (r *UserReportRepository) fetchItemCounts() (map[int]int, error) {
	result := make(map[int]int)
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return result, nil
	}

	args := []interface{}{}
	stockInSelect := `
		SELECT si.user_id, si.item_id
		FROM stock_in si
		JOIN items i ON i.item_id = si.item_id
	`
	if cond, condArgs := r.buildStoreFilter("i.store_id"); cond != "" {
		stockInSelect += " WHERE " + cond
		args = append(args, condArgs...)
	}

	stockOutSelect := `
		SELECT so.user_id, p.item_id
		FROM stock_out so
		JOIN programs p ON p.program_id = so.program_id
		JOIN items i ON i.item_id = p.item_id
	`
	if cond, condArgs := r.buildStoreFilter("i.store_id"); cond != "" {
		stockOutSelect += " WHERE " + cond
		args = append(args, condArgs...)
	}

	query := `
		SELECT user_id, COUNT(DISTINCT item_id) AS item_count
		FROM (
			` + stockInSelect + `
			UNION ALL
			` + stockOutSelect + `
		) AS user_items
		GROUP BY user_id
	`

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID, count int
		if err := rows.Scan(&userID, &count); err != nil {
			return nil, err
		}
		result[userID] = count
	}

	return result, rows.Err()
}

func (r *UserReportRepository) buildStoreFilter(column string) (string, []interface{}) {
	if len(r.StoreIDs) == 0 {
		return "", nil
	}

	placeholders := make([]string, len(r.StoreIDs))
	args := make([]interface{}, len(r.StoreIDs))

	for i, id := range r.StoreIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	return column + " IN (" + strings.Join(placeholders, ",") + ")", args
}

func filterStoreIDs(ids []int, allowed map[int]bool) []int {
	seen := make(map[int]bool)
	var result []int

	for _, id := range ids {
		if len(allowed) > 0 && !allowed[id] {
			continue
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		result = append(result, id)
	}

	return result
}

func joinStoreNames(ids []int, nameMap map[int]string) string {
	if len(ids) == 0 {
		return "-"
	}

	names := make([]string, 0, len(ids))
	for _, id := range ids {
		if name, ok := nameMap[id]; ok {
			names = append(names, name)
		} else {
			names = append(names, strconv.Itoa(id))
		}
	}

	return strings.Join(names, ", ")
}

func parseStoreIDs(raw string) []int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []int{}
	}

	var ids []int
	if err := json.Unmarshal([]byte(raw), &ids); err == nil {
		return ids
	}

	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if id, err := strconv.Atoi(p); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}

func storeIDsToMap(ids []int) map[int]bool {
	result := make(map[int]bool, len(ids))
	for _, id := range ids {
		result[id] = true
	}
	return result
}
