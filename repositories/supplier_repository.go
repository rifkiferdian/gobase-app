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
