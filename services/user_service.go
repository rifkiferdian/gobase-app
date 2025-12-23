package services

import (
	"database/sql"
	"stok-hadiah/config"
	"stok-hadiah/models"
	"stok-hadiah/repositories"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

const userModelType = "Models\\User"

func UserHasPermission(userID int, perm string) (bool, error) {
	var dummy int
	// Cek permission via role yang dimiliki user
	queryRole := `
		SELECT 1
		FROM model_has_roles mhr
		JOIN role_has_permissions rhp ON rhp.role_id = mhr.role_id
		JOIN permissions p ON p.id = rhp.permission_id
		WHERE mhr.model_id = ? AND mhr.model_type = ? AND p.name = ?
		LIMIT 1
	`
	err := config.DB.QueryRow(queryRole, userID, userModelType, perm).Scan(&dummy)
	if err == nil {
		return true, nil
	}
	if err != sql.ErrNoRows {
		return false, err
	}

	// Fallback: cek permission langsung ke user (model_has_permissions)
	queryDirect := `
		SELECT 1
		FROM model_has_permissions mhp
		JOIN permissions p ON p.id = mhp.permission_id
		WHERE mhp.model_id = ? AND mhp.model_type = ? AND p.name = ?
		LIMIT 1
	`

	err = config.DB.QueryRow(queryDirect, userID, userModelType, perm).Scan(&dummy)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil // tidak punya permission
	}

	return false, err // error lain
}

func GetUserPermissions(userID int) (map[string]bool, error) {
	perms := make(map[string]bool)

	rows, err := config.DB.Query(`
		SELECT DISTINCT p.name
		FROM permissions p
		JOIN role_has_permissions rhp ON rhp.permission_id = p.id
		JOIN model_has_roles mhr ON mhr.role_id = rhp.role_id
		WHERE mhr.model_id = ? AND mhr.model_type = ?

		UNION

		SELECT DISTINCT p2.name
		FROM permissions p2
		JOIN model_has_permissions mhp ON mhp.permission_id = p2.id
		WHERE mhp.model_id = ? AND mhp.model_type = ?
	`, userID, userModelType, userID, userModelType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms[name] = true
	}

	return perms, nil
}
