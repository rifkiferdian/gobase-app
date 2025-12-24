package models

// StockIn merepresentasikan data stok barang masuk pada tabel `stock_in`.
// Field ItemName dan UserName diisi ketika dilakukan join ke tabel terkait.
// ReceivedAt disimpan sebagai string yang sudah di-format agar mudah ditampilkan di template.
type StockIn struct {
	ID                int
	UserID            int
	UserName          string
	ItemID            int
	ItemName          string
	StoreID           int
	StoreName         string
	SupplierName      string
	Qty               int
	ReceivedAt        string
	ReceivedAtDisplay string
	Description       string
}
