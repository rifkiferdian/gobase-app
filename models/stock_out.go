package models

import "time"

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
