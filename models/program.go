package models

// Program merepresentasikan data program promosi pada tabel `programs`.
// Field ItemName diisi ketika dilakukan join ke tabel items.
// StartDate dan EndDate disimpan dalam format string "YYYY-MM-DD" agar mudah di-bind ke input type="date".
type Program struct {
	ProgramID   int
	ProgramName string
	ItemID      int
	ItemName    string
	StoreID     int
	StoreName   string
	StartDate   string
	EndDate     string
}
