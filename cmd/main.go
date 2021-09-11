package main

import (
	"github.com/me0888/wallet/pkg/types"
	"github.com/me0888/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	// svc.Import("D:/Users/билол/Desktop")
	pays := []types.Payment{
		// {
		// 	ID:        "1",
		// 	AccountID: 100,
		// 	Amount:    900,
		// 	Category:  "aoto",
		// 	Status:    types.PaymentStatusInProgress,
		// },
		// {
		// 	ID:        "2",
		// 	AccountID: 100,
		// 	Amount:    900,
		// 	Category:  "aoto",
		// 	Status:    types.PaymentStatusInProgress,
		// },
		// {
		// 	ID:        "3",
		// 	AccountID: 100,
		// 	Amount:    900,
		// 	Category:  "aoto",
		// 	Status:    types.PaymentStatusInProgress,
		// },
		// {
		// 	ID:        "4",
		// 	AccountID: 100,
		// 	Amount:    900,
		// 	Category:  "aoto",
		// 	Status:    types.PaymentStatusInProgress,
		// },
		// {
		// 	ID:        "5",
		// 	AccountID: 100,
		// 	Amount:    900,
		// 	Category:  "aoto",
		// 	Status:    types.PaymentStatusInProgress,
		// },
		// {
		// 	ID:        "6",
		// 	AccountID: 100,
		// 	Amount:    900,
		// 	Category:  "aoto",
		// 	Status:    types.PaymentStatusInProgress,
		// },
	}
	svc.HistoryToFiles(pays,"D:/Users/билол/Desktop/1",4)

}
