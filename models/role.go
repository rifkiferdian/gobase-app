package models

// Role mewakili data role beserta jumlah permission dan user yang terkait.
type Role struct {
	ID              int
	Name            string
	PermissionCount int
	UserCount       int
	UpdatedAt       string
}

// RoleCreateInput mewakili payload untuk membuat role baru.
type RoleCreateInput struct {
	Name          string
	GuardName     string
	PermissionIDs []int64
}
