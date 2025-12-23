package repositories

import (
	"database/sql"
	"stok-hadiah/models"
	"strings"
	"time"
)

type ItemRepository struct {
	DB                 *sql.DB
	StoreIDs           []int
	FilterStoreID      *int
	EnforceStoreFilter bool
}

// GetAll mengambil seluruh data item beserta nama supplier-nya.
func (r *ItemRepository) GetAll() ([]models.Item, error) {
	args := []interface{}{}
	query := `
		SELECT i.item_id,
		       i.item_name,
		       i.category,
		       i.supplier_id,
		       s.supplier_name,
		       i.store_id,
		       st.store_name,
		       i.description,
		       i.created_at
		FROM items i
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
		LEFT JOIN stores st ON st.store_id = i.store_id
	`
	query, skip := r.appendStoreFilter(query, &args, false)
	if skip {
		return []models.Item{}, nil
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var (
			it        models.Item
			createdAt time.Time
		)
		if err := rows.Scan(
			&it.ItemID,
			&it.ItemName,
			&it.Category,
			&it.SupplierID,
			&it.SupplierName,
			&it.StoreID,
			&it.StoreName,
			&it.Description,
			&createdAt,
		); err != nil {
			return nil, err
		}
		it.CreatedAt = createdAt.Format("02-01-2006 15:04:05")
		items = append(items, it)
	}

	return items, nil
}

// Search mengambil data item berdasarkan kata kunci nama dan/atau kategori.
// Jika kedua parameter kosong, akan mengembalikan seluruh data seperti GetAll.
func (r *ItemRepository) Search(name, category string) ([]models.Item, error) {
	query := `
		SELECT i.item_id,
		       i.item_name,
		       i.category,
		       i.supplier_id,
		       s.supplier_name,
		       i.store_id,
		       st.store_name,
		       i.description,
		       i.created_at
		FROM items i
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
		LEFT JOIN stores st ON st.store_id = i.store_id
		WHERE 1=1`

	args := []interface{}{}

	if name != "" {
		query += " AND i.item_name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if category != "" {
		query += " AND i.category = ?"
		args = append(args, category)
	}

	query, skip := r.appendStoreFilter(query, &args, true)
	if skip {
		return []models.Item{}, nil
	}

	query += " ORDER BY i.item_id DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var (
			it        models.Item
			createdAt time.Time
		)
		if err := rows.Scan(
			&it.ItemID,
			&it.ItemName,
			&it.Category,
			&it.SupplierID,
			&it.SupplierName,
			&it.StoreID,
			&it.StoreName,
			&it.Description,
			&createdAt,
		); err != nil {
			return nil, err
		}
		it.CreatedAt = createdAt.Format("02-01-2006 15:04:05")
		items = append(items, it)
	}

	return items, nil
}

// Count mengembalikan jumlah seluruh data item.
func (r *ItemRepository) Count() (int, error) {
	args := []interface{}{}
	query, skip := r.appendStoreFilter(`
		SELECT COUNT(*) FROM items i
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

// CountByCategory menghitung jumlah item berdasarkan kategori (case-insensitive).
func (r *ItemRepository) CountByCategory(category string) (int, error) {
	args := []interface{}{category}
	query, skip := r.appendStoreFilter(`
		SELECT COUNT(*) FROM items i
		WHERE UPPER(i.category) = UPPER(?)
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

// GetPaginated mengambil data item dengan pagination menggunakan LIMIT dan OFFSET.
// Data yang diambil sudah termasuk join dengan tabel suppliers untuk mendapatkan nama supplier.
func (r *ItemRepository) GetPaginated(limit, offset int) ([]models.Item, error) {
	args := []interface{}{}
	query := `
		SELECT i.item_id,
		       i.item_name,
		       i.category,
		       i.supplier_id,
		       s.supplier_name,
		       i.store_id,
		       st.store_name,
		       i.description,
		       i.created_at
		FROM items i
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
		LEFT JOIN stores st ON st.store_id = i.store_id
	`
	query, skip := r.appendStoreFilter(query, &args, false)
	if skip {
		return []models.Item{}, nil
	}

	query += `
		ORDER BY i.item_id DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var (
			it        models.Item
			createdAt time.Time
		)
		if err := rows.Scan(
			&it.ItemID,
			&it.ItemName,
			&it.Category,
			&it.SupplierID,
			&it.SupplierName,
			&it.StoreID,
			&it.StoreName,
			&it.Description,
			&createdAt,
		); err != nil {
			return nil, err
		}
		it.CreatedAt = createdAt.Format("02-01-2006 15:04:05")
		items = append(items, it)
	}

	return items, nil
}

func (r *ItemRepository) Create(i models.Item) error {
	_, err := r.DB.Exec(`
		INSERT INTO items (item_name, category, supplier_id, store_id, description)
		VALUES (?, ?, ?, ?, ?)
	`, i.ItemName, i.Category, i.SupplierID, i.StoreID, i.Description)
	return err
}

func (r *ItemRepository) Update(i models.Item) error {
	_, err := r.DB.Exec(`
		UPDATE items
		SET item_name = ?,
		    category = ?,
		    supplier_id = ?,
		    store_id = ?,
		    description = ?
		WHERE item_id = ?
	`, i.ItemName, i.Category, i.SupplierID, i.StoreID, i.Description, i.ItemID)
	return err
}

func (r *ItemRepository) Delete(id int) error {
	_, err := r.DB.Exec(`
		DELETE FROM items
		WHERE item_id = ?
	`, id)
	return err
}

// appendStoreFilter menambahkan filter store_id sesuai StoreIDs dan FilterStoreID.
// Jika EnforceStoreFilter true dan StoreIDs kosong, fungsi akan menandakan untuk skip query (no access).
// hasWhere menentukan apakah query sudah memiliki klausa WHERE sebelumnya.
func (r *ItemRepository) appendStoreFilter(query string, args *[]interface{}, hasWhere bool) (string, bool) {
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
