package models

// SessionUser is a minimal user payload stored in session cookies.
type SessionUser struct {
	UserID          int
	NIP             string
	Name            string
	Username        string
	Role            string
	IsAuthenticated bool
}
