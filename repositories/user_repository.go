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

type UserCreateParams struct {
	NIP            int
	Username       string
	HashedPassword string
	Name           string
	Email          string
	Status         string
	StoreIDs       []int
}

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

// CreateUserWithRoles menyimpan data user baru beserta assignment rolenya dalam satu transaksi.
func (r *UserRepository) CreateUserWithRoles(params UserCreateParams, roleIDs []int64) (int64, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	storeJSON, err := json.Marshal(params.StoreIDs)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var emailVal interface{}
	if strings.TrimSpace(params.Email) == "" {
		emailVal = nil // simpan NULL jika email kosong
	} else {
		emailVal = params.Email
	}

	res, err := tx.Exec(`
		INSERT INTO users (nip, username, password, name, email, status, store_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, params.NIP, params.Username, params.HashedPassword, params.Name, emailVal, params.Status, string(storeJSON))
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(roleIDs) > 0 {
		stmt, err := tx.Prepare(`
			INSERT INTO model_has_roles (role_id, model_type, model_id)
			VALUES (?, ?, ?)
		`)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		defer stmt.Close()

		for _, roleID := range roleIDs {
			if _, err := stmt.Exec(roleID, userModelType, userID); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}

	return userID, nil
}

// ExistsByUsername mengecek apakah username sudah digunakan.
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(1) FROM users WHERE username = ?`, username).Scan(&count)
	return count > 0, err
}

// ExistsByNIP mengecek apakah NIP sudah digunakan.
func (r *UserRepository) ExistsByNIP(nip int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(1) FROM users WHERE nip = ?`, nip).Scan(&count)
	return count > 0, err
}

// ExistsByEmail mengecek apakah email sudah digunakan (abaikan jika kosong).
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	if strings.TrimSpace(email) == "" {
		return false, nil
	}

	var count int
	err := r.DB.QueryRow(`SELECT COUNT(1) FROM users WHERE email = ?`, email).Scan(&count)
	return count > 0, err
}

// GetRoleIDsByNames mengambil role_id berdasarkan nama role yang diberikan.
func (r *UserRepository) GetRoleIDsByNames(names []string) (map[string]int64, error) {
	result := make(map[string]int64)

	if len(names) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(names))
	args := make([]interface{}, len(names))

	for i, name := range names {
		placeholders[i] = "?"
		args[i] = name
	}

	query := `
		SELECT id, name
		FROM roles
		WHERE name IN (` + strings.Join(placeholders, ",") + `)
	`

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int64
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		result[name] = id
	}

	return result, rows.Err()
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

// DeleteUser removes a user and related role/permission mappings in a single transaction.
func (r *UserRepository) DeleteUser(id int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`DELETE FROM model_has_roles WHERE model_id = ? AND model_type = ?`, id, userModelType); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`DELETE FROM model_has_permissions WHERE model_id = ? AND model_type = ?`, id, userModelType); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`DELETE FROM users WHERE id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
