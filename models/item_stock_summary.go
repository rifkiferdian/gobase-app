package models

// ItemStockSummary merepresentasikan ringkasan stok per item,
// termasuk total masuk, total keluar, dan sisa yang tersedia.
type ItemStockSummary struct {
	ItemID       int
	ItemName     string
	Category     string
	SupplierName string
	StoreName    string
	Description  string
	QtyIn        int
	QtyOut       int
	Remaining    int
}
