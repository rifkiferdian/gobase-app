package services

import (
	"errors"
	"stok-hadiah/repositories"
)

// StockOutService menyediakan layer bisnis untuk penyesuaian stok keluar.
type StockOutService struct {
	Repo *repositories.StockOutRepository
}

// AdjustQuantity memanggil repo untuk menambah/mengurangi qty keluar per item.
// Delta seharusnya bernilai +1 atau -1.
func (s *StockOutService) AdjustQuantity(itemID, delta, userID int) (int, error) {
	if delta == 0 {
		return 0, errors.New("delta harus lebih dari 0")
	}
	return s.Repo.AdjustQuantity(itemID, delta, userID)
}
