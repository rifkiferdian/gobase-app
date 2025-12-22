package repositories

import (
	"database/sql"
	"stok-hadiah/models"
)

type RoleRepository struct {
	DB *sql.DB
}

// GetAll mengambil seluruh data role beserta jumlah permission dan user yang terkait.
func (r *RoleRepository) GetAll() ([]models.Role, error) {
	rows, err := r.DB.Query(`
		SELECT 
			r.id,
			r.name,
			COUNT(DISTINCT rhp.permission_id) AS permission_count,
			COUNT(DISTINCT mhr.model_id) AS user_count,
			r.updated_at
		FROM roles r
		LEFT JOIN role_has_permissions rhp ON rhp.role_id = r.id
		LEFT JOIN model_has_roles mhr ON mhr.role_id = r.id
		GROUP BY r.id, r.name, r.updated_at
		ORDER BY r.updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role

	for rows.Next() {
		var (
			role      models.Role
			updatedAt sql.NullTime
		)

		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.PermissionCount,
			&role.UserCount,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		if updatedAt.Valid {
			role.UpdatedAt = updatedAt.Time.Format("01-02-2006 15:04:05")
		} else {
			role.UpdatedAt = "-"
		}

		roles = append(roles, role)
	}

	return roles, rows.Err()
}
