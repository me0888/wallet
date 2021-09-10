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
	_, err := svc.RegisterAccount("09")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = svc.FindAccountByID(2)

	if err == nil {
		t.Error(err)
		return
	}

}
