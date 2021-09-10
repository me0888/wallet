package wallet

import (
	"fmt"

	"github.com/me0888/wallet/pkg/types"
	"reflect"
	"testing"
)

func TestService_RegisterAccount_success(t *testing.T) {

	svc := &Service{}
	result, err := svc.RegisterAccount("09")
	if err != nil {
		fmt.Println(err)
		return
	}
	expected := &types.Account{
		ID:      1,
		Phone:   "09",
		Balance: 0,
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("\ninvalid result, \ngot:  %v, \nwant: %v", expected, result)
	}
}

func TestService_RegisterAccount_alreadyRegistered(t *testing.T) {

	svc := &Service{}
	_, err := svc.RegisterAccount("09")
	_, err = svc.RegisterAccount("09")
	if err == nil {
		t.Error(err)
		return
	}

}

func TestService_FindAccountByID_success(t *testing.T) {

	svc := &Service{}
	_, err := svc.RegisterAccount("09")
	if err != nil {
		t.Error(err)
		return
	}

	result, err := svc.FindAccountByID(1)
	if err != nil {
		t.Error(err)
		return
	}
	expected := &types.Account{
		ID:      1,
		Phone:   "09",
		Balance: 0,
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("\ninvalid result, \ngot:  %v, \nwant: %v", expected, result)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {

	svc := &Service{}
	_, err := svc.FindAccountByID(2000)

	if err != ErrAccountNotFound {
		t.Error(err)
		return
	}

}

func TestService_FindPaymentById_success(t *testing.T) {

	svc := &Service{}
	acc, err := svc.RegisterAccount("09")
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Deposit(acc.ID, 100)
	if err != nil {
		t.Error(err)
		return
	}

	pay, err := svc.Pay(acc.ID, 50, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	pay2, err := svc.FindPaymentByID(pay.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(pay, pay2) {
		t.Errorf("\ninvalid result, \ngot:  %v, \nwant: %v", pay, pay2)
	}
}

func TestService_Reject_success(t *testing.T) {

	svc := &Service{}
	acc, err := svc.RegisterAccount("09")
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Deposit(acc.ID, 100)
	if err != nil {
		t.Error(err)
		return
	}

	pay, err := svc.Pay(acc.ID, 50, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Reject(pay.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if acc.Balance!=100 {
		t.Errorf("\ninvalid result, \ngot:  %v, \nwant: 100", acc.Balance)
	}
}
