package main

import (
	"github.com/me0888/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	svc.Export("900.txt")
	
}
