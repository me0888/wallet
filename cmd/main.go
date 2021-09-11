package main

import (
	"fmt"

	"github.com/me0888/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	acc, err := svc.RegisterAccount("09")
	if err != nil {
		fmt.Println(err)
		return
	}

	acc, err = svc.RegisterAccount("08")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(acc.ID, 900)
	if err != nil {
		switch err {
		case wallet.ErrAccountNotFound:
			fmt.Println("acc not found")

		case wallet.ErrAmmountMustBePositive:
			fmt.Println("amount mast be positive")

		}
		return
	}
	svc.ExportToFile("900.txt")
	
}
