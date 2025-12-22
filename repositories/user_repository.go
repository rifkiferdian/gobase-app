package repositories

import (
	"database/sql"
	"encoding/json"
	"stok-hadiah/models"
	"strconv"
	"strings"
	"time"
)

type UserRepository struct {
	DB *sql.DB
}

const userModelType = "Models\\User"

// GetAll mengambil seluruh data user beserta parsing store_id JSON dan format tanggal.
func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.DB.Query(`
		SELECT 
			u.id, 
			u.username, 
			u.name, 
			u.email, 
			u.status, 
			u.store_id, 
			u.created_at,
			COALESCE(GROUP_CONCAT(r2.name ORDER BY r2.name SEPARATOR ', '), '') AS role_display
		FROM users u
		LEFT JOIN model_has_roles mhr ON mhr.model_id = u.id AND mhr.model_type = ?
		LEFT JOIN roles r2 ON r2.id = mhr.role_id
		GROUP BY 
			u.id, u.username, u.name, u.email, u.status, u.store_id, u.created_at
		ORDER BY u.created_at DESC
	`, userModelType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var (
			u         models.User
			storeJSON string
			createdAt time.Time
		)

		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Name,
			&u.Email,
			&u.Status,
			&storeJSON,
			&createdAt,
			&u.RoleDisplay,
		); err != nil {
			return nil, err
		}

		u.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		u.CreatedAtDisplay = createdAt.Format("02 Jan 2006 15:04:05")

		if u.Status == "active" {
			u.StatusLabel = "Aktif"
		} else {
			u.StatusLabel = "Non Aktif"
		}

		if storeJSON != "" {
			var storeIDs []int
			if err := json.Unmarshal([]byte(storeJSON), &storeIDs); err == nil {
				u.StoreIDs = storeIDs
				storeNames, err := r.getStoreNames(storeIDs)
				if err == nil && len(storeNames) > 0 {
					u.StoreDisplay = strings.Join(storeNames, ", ")
				} else {
					u.StoreDisplay = joinIntSlice(storeIDs)
				}
			} else {
				u.StoreDisplay = storeJSON
			}
		}

		if u.StoreDisplay == "" {
			u.StoreDisplay = "-"
		}

		if u.RoleDisplay == "" {
			u.RoleDisplay = "-"
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) getStoreNames(ids []int) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "SELECT store_name FROM stores WHERE store_id IN (" + strings.Join(placeholders, ",") + ") ORDER BY store_name"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return names, nil
}

func joinIntSlice(values []int) string {
	if len(values) == 0 {
		return ""
	}

	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = strconv.Itoa(v)
	}

	return strings.Join(parts, ", ")
}
