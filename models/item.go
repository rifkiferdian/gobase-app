package models

// Item merepresentasikan data barang pada tabel `items`.
// Field SupplierName diisi ketika dilakukan join ke tabel suppliers.
type Item struct {
	ItemID       int
	ItemName     string
	Category     string
	SupplierID   int
	SupplierName string
	StoreID      int
	StoreName    string
	Description  string
	CreatedAt    string
}
