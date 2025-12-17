package repositories

import (
	"database/sql"
	"stok-hadiah/models"
)

type SupplierRepository struct {
	DB *sql.DB
}

func (r *SupplierRepository) GetAll() ([]models.Supplier, error) {
	rows, err := r.DB.Query(`
		SELECT suppliers_id, supplier_name, active, description
		FROM suppliers
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []models.Supplier

	for rows.Next() {
		var s models.Supplier
		err := rows.Scan(
			&s.SupplierID,
			&s.SupplierName,
			&s.Active,
			&s.Description,
		)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, s)
	}

	return suppliers, nil
}

func (r *SupplierRepository) Create(s models.Supplier) error {
	_, err := r.DB.Exec(`
		INSERT INTO suppliers (supplier_name, active, description)
		VALUES (?, ?, ?)
	`, s.SupplierName, s.Active, s.Description)
	return err
}

func (r *SupplierRepository) Delete(id int) error {
	_, err := r.DB.Exec(`
		DELETE FROM suppliers
		WHERE suppliers_id = ?
	`, id)
	return err
}

func (r *SupplierRepository) Update(s models.Supplier) error {
	_, err := r.DB.Exec(`
		UPDATE suppliers
		SET supplier_name = ?, active = ?, description = ?
		WHERE suppliers_id = ?
	`, s.SupplierName, s.Active, s.Description, s.SupplierID)
	return err
}
