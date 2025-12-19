package repositories

import (
	"database/sql"
	"stok-hadiah/models"
)

type ItemRepository struct {
	DB *sql.DB
}

// GetAll mengambil seluruh data item beserta nama supplier-nya.
func (r *ItemRepository) GetAll() ([]models.Item, error) {
	rows, err := r.DB.Query(`
		SELECT i.item_id,
		       i.item_name,
		       i.category,
		       i.supplier_id,
		       s.supplier_name,
		       i.description
		FROM items i
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var it models.Item
		if err := rows.Scan(
			&it.ItemID,
			&it.ItemName,
			&it.Category,
			&it.SupplierID,
			&it.SupplierName,
			&it.Description,
		); err != nil {
			return nil, err
		}
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
		       i.description
		FROM items i
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
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

	query += " ORDER BY i.item_id DESC"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var it models.Item
		if err := rows.Scan(
			&it.ItemID,
			&it.ItemName,
			&it.Category,
			&it.SupplierID,
			&it.SupplierName,
			&it.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	return items, nil
}

// Count mengembalikan jumlah seluruh data item.
func (r *ItemRepository) Count() (int, error) {
	row := r.DB.QueryRow(`
		SELECT COUNT(*) FROM items
	`)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

// GetPaginated mengambil data item dengan pagination menggunakan LIMIT dan OFFSET.
// Data yang diambil sudah termasuk join dengan tabel suppliers untuk mendapatkan nama supplier.
func (r *ItemRepository) GetPaginated(limit, offset int) ([]models.Item, error) {
	rows, err := r.DB.Query(`
		SELECT i.item_id,
		       i.item_name,
		       i.category,
		       i.supplier_id,
		       s.supplier_name,
		       i.description
		FROM items i
		JOIN suppliers s ON s.suppliers_id = i.supplier_id
		ORDER BY i.item_id DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var it models.Item
		if err := rows.Scan(
			&it.ItemID,
			&it.ItemName,
			&it.Category,
			&it.SupplierID,
			&it.SupplierName,
			&it.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	return items, nil
}

func (r *ItemRepository) Create(i models.Item) error {
	_, err := r.DB.Exec(`
		INSERT INTO items (item_name, category, supplier_id, description)
		VALUES (?, ?, ?, ?)
	`, i.ItemName, i.Category, i.SupplierID, i.Description)
	return err
}

func (r *ItemRepository) Update(i models.Item) error {
	_, err := r.DB.Exec(`
		UPDATE items
		SET item_name = ?,
		    category = ?,
		    supplier_id = ?,
		    description = ?
		WHERE item_id = ?
	`, i.ItemName, i.Category, i.SupplierID, i.Description, i.ItemID)
	return err
}

func (r *ItemRepository) Delete(id int) error {
	_, err := r.DB.Exec(`
		DELETE FROM items
		WHERE item_id = ?
	`, id)
	return err
}
