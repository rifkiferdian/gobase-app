package models

import "time"

// StockOutDetail merepresentasikan detail keluarnya stok per baris stock_out.
// IssuedAt disimpan sebagai string agar mudah diisi ke input/tabel dan tampilan terformat.
type StockOutDetail struct {
	ID              int
	UserID          int
	UserName        string
	ProgramID       int
	ItemID          int
	ItemName        string
	StoreID         int
	StoreName       string
	Qty             int
	IssuedAt        string
	IssuedAtDisplay string
	Reason          string
}

// StockOutCase merepresentasikan pengeluaran stok dengan alasan/keterangan khusus.
// IssuedAt disimpan sebagai time.Time agar mudah diformat di layer lain.
type StockOutCase struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProgramID int       `json:"program_id"`
	ItemID    int       `json:"item_id"`
	ItemName  string    `json:"item_name"`
	Qty       int       `json:"qty"`
	Reason    string    `json:"reason"`
	IssuedAt  time.Time `json:"issued_at"`
}
