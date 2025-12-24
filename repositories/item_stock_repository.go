package repositories

import (
	"database/sql"
	"stok-hadiah/models"
	"strings"
)

// ItemStockRepository menyediakan ringkasan stok per item (total masuk, keluar, dan sisa).
type ItemStockRepository struct {
	DB                 *sql.DB
	StoreIDs           []int
	FilterStoreID      *int
	EnforceStoreFilter bool
}

// GetSummaries mengembalikan daftar ringkasan stok item dengan filter nama, kategori, dan supplier opsional.
// Filter store otomatis diterapkan berdasarkan StoreIDs/FilterStoreID.
func (r *ItemStockRepository) GetSummaries(filterName, filterCategory string, supplierID *int) ([]models.ItemStockSummary, error) {
	if r.EnforceStoreFilter && len(r.StoreIDs) == 0 {
		return []models.ItemStockSummary{}, nil
	}

	args := []interface{}{}
	query := `
		SELECT
			i.item_id,
			i.item_name,
			i.category,
			su.supplier_name,
			st.store_name,
			i.description,
			COALESCE(SUM(si.qty), 0) AS total_in,
			COALESCE((
				SELECT SUM(so.qty)
				FROM stock_out so
				JOIN programs p2 ON p2.program_id = so.program_id
				WHERE p2.item_id = i.item_id
			), 0) AS total_out
		FROM items i
		JOIN suppliers su ON su.suppliers_id = i.supplier_id
		LEFT JOIN stores st ON st.store_id = i.store_id
		LEFT JOIN stock_in si ON si.item_id = i.item_id
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
			i.description
		ORDER BY i.item_id DESC
	`

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []models.ItemStockSummary

	for rows.Next() {
		var summary models.ItemStockSummary
		if err := rows.Scan(
			&summary.ItemID,
			&summary.ItemName,
			&summary.Category,
			&summary.SupplierName,
			&summary.StoreName,
			&summary.Description,
			&summary.QtyIn,
			&summary.QtyOut,
		); err != nil {
			return nil, err
		}

		summary.Remaining = summary.QtyIn - summary.QtyOut
		summaries = append(summaries, summary)
	}

	return summaries, nil
}
