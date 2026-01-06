package models

// UserReportSummary merepresentasikan rekap aktivitas stok per user.
// Nilai TotalIn/TotalOut dan ItemTypes sudah dihitung sesuai filter store yang berlaku.
type UserReportSummary struct {
	UserID     int
	NIP        int
	Name       string
	StoreIDs   []int
	StoreNames string
	ItemTypes  int
	TotalIn    int
	TotalOut   int
}
