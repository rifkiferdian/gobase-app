package repositories

import (
	"database/sql"
	"strings"
	"time"

	helpers "stok-hadiah/helper"
	"stok-hadiah/models"
)

func formatDateListID(dateList string) string {
	trimmed := strings.TrimSpace(dateList)
	if trimmed == "" {
		return ""
	}

	parts := strings.Split(trimmed, ",")
	formatted := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		t, err := time.Parse("2006-01-02", p)
		if err != nil {
			formatted = append(formatted, p)
			continue
		}
		formatted = append(formatted, helpers.FormatDateID(t))
	}

	return strings.Join(formatted, ", ")
}

// ItemStockRepository menyediakan ringkasan stok per item (total masuk, keluar, dan sisa).
type ItemStockRepository struct {
	DB                 *sql.DB
	StoreIDs           []int
	FilterStoreID      *int
	EnforceStoreFilter bool
}

func (r *ItemStockRepository) buildSummaryQuery(filterName, filterCategory string, supplierID *int) (string, []interface{}) {
	args := []interface{}{}
	query := `
		SELECT
			i.item_id,
			i.item_name,
			i.category,
			su.supplier_name,
			COALESCE(p.program_names, '') AS program_names,
			COALESCE(p.program_start_dates, '') AS program_start_dates,
			COALESCE(p.program_end_dates, '') AS program_end_dates,
			st.store_name,
			i.description,
			COALESCE(SUM(si.qty), 0) AS total_in,
			COALESCE(sot.total_out, 0) AS total_out,
			COALESCE(SUM(si.qty), 0) - COALESCE(sot.total_out, 0) AS remaining
		FROM items i
		JOIN suppliers su ON su.suppliers_id = i.supplier_id
		LEFT JOIN (
			SELECT
				item_id,
				GROUP_CONCAT(DISTINCT program_name ORDER BY program_id DESC SEPARATOR ', ') AS program_names,
				GROUP_CONCAT(DISTINCT DATE_FORMAT(start_date, '%Y-%m-%d') ORDER BY program_id DESC SEPARATOR ', ') AS program_start_dates,
				GROUP_CONCAT(DISTINCT DATE_FORMAT(end_date, '%Y-%m-%d') ORDER BY program_id DESC SEPARATOR ', ') AS program_end_dates
			FROM programs
			GROUP BY item_id
		) p ON p.item_id = i.item_id
		LEFT JOIN stores st ON st.store_id = i.store_id
		LEFT JOIN stock_in si ON si.item_id = i.item_id
		LEFT JOIN (
			SELECT p2.item_id, SUM(so.qty) AS total_out
			FROM stock_out so
			JOIN programs p2 ON p2.program_id = so.program_id
			GROUP BY p2.item_id
		) sot ON sot.item_id = i.item_id
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

	if r.FilterStoreID != nil {
		conditions = append(conditions, "i.store_id = ?")
		args = append(args, *r.FilterStoreID)
	}

	if filterName != "" {
		conditions = append(conditions, "i.item_name LIKE ?")
		args = append(args, "%"+filterName+"%")
	}

	if filterCategory != "" {
		conditions = append(conditions, "i.category = ?")
		args = append(args, filterCategory)
	}

	if supplierID != nil {
		conditions = append(conditions, "i.supplier_id = ?")
		args = append(args, *supplierID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY
			i.item_id,
			i.item_name,
			i.category,
			su.supplier_name,
			st.store_name,
			i.description,
			p.program_names,
			p.program_start_dates,
			p.program_end_dates,
			sot.total_out
		HAVING (COALESCE(SUM(si.qty), 0) - COALESCE(sot.total_out, 0)) > 0
	`

	return query, args
}

func scanItemStockSummaries(rows *sql.Rows) ([]models.ItemStockSummary, error) {
	var summaries []models.ItemStockSummary
	for rows.Next() {
		var summary models.ItemStockSummary
		if err := rows.Scan(
			&summary.ItemID,
			&summary.ItemName,
			&summary.Category,
			&summary.SupplierName,
			&summary.ProgramNames,
			&summary.ProgramStartDates,
			&summary.ProgramEndDates,
			&summary.StoreName,
			&summary.Description,
			&summary.QtyIn,
			&summary.QtyOut,
			&summary.Remaining,
		); err != nil {
			return nil, err
		}

		summary.ProgramStartDisplay = formatDateListID(summary.ProgramStartDates)
		summary.ProgramEndDisplay = formatDateListID(summary.ProgramEndDates)
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

// GetSummaries mengembalikan daftar ringkasan stok item dengan filter nama, kategori, dan supplier opsional.
// Filter store otomatis diterapkan berdasarkan StoreIDs/FilterStoreID.
func (r *ItemStockRepository) GetSummaries(filterName, filterCategory string, supplierID *int) ([]models.ItemStockSummary, error) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return []models.ItemStockSummary{}, nil
	}

	baseQuery, args := r.buildSummaryQuery(filterName, filterCategory, supplierID)
	query := baseQuery + "\nORDER BY i.item_id DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanItemStockSummaries(rows)
}

// GetSummariesPaginated mengembalikan daftar ringkasan stok item dengan pagination.
func (r *ItemStockRepository) GetSummariesPaginated(filterName, filterCategory string, supplierID *int, limit, offset int) ([]models.ItemStockSummary, error) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return []models.ItemStockSummary{}, nil
	}

	baseQuery, baseArgs := r.buildSummaryQuery(filterName, filterCategory, supplierID)
	args := append([]interface{}{}, baseArgs...)
	query := baseQuery + "\nORDER BY i.item_id DESC"
	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanItemStockSummaries(rows)
}

// CountSummaries menghitung total item yang memiliki stok tersisa sesuai filter.
func (r *ItemStockRepository) CountSummaries(filterName, filterCategory string, supplierID *int) (int, error) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return 0, nil
	}

	baseQuery, args := r.buildSummaryQuery(filterName, filterCategory, supplierID)
	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") summary"

	var total int
	if err := r.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// GetSummaryTotals mengembalikan total qty in/out dan remaining untuk seluruh hasil filter.
func (r *ItemStockRepository) GetSummaryTotals(filterName, filterCategory string, supplierID *int) (int, int, int, error) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return 0, 0, 0, nil
	}

	baseQuery, args := r.buildSummaryQuery(filterName, filterCategory, supplierID)
	totalQuery := `
		SELECT
			COALESCE(SUM(total_in), 0) AS total_qty_in,
			COALESCE(SUM(total_out), 0) AS total_qty_out,
			COALESCE(SUM(remaining), 0) AS total_remaining
		FROM (` + baseQuery + `) summary
	`

	var totalIn, totalOut, totalRemaining int
	if err := r.DB.QueryRow(totalQuery, args...).Scan(&totalIn, &totalOut, &totalRemaining); err != nil {
		return 0, 0, 0, err
	}

	return totalIn, totalOut, totalRemaining, nil
}
