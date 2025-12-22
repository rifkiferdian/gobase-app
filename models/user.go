package models

// User merepresentasikan data pada tabel users.
// Field StoreIDs berisi daftar id toko dalam bentuk slice setelah parsing JSON.
type User struct {
	ID               int
	Username         string
	Name             string
	Email            string
	Status           string
	StatusLabel      string
	StoreIDs         []int
	StoreDisplay     string
	RoleDisplay      string
	CreatedAt        string
	CreatedAtDisplay string
}
