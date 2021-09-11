package wallet

import (
	"fmt"
	"log"

	"reflect"
	"testing"

	"github.com/me0888/wallet/pkg/types"
)

type testService struct {
	*Service
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:   "+992000000007",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{
			amount:   1_000_00,
			category: "auto",
		},
	},
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can`t register account, erro = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can`t deposit account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can`t make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

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

	if acc.Balance != 100 {
		t.Errorf("\ninvalid result, \ngot:  %v, \nwant: 100", acc.Balance)
	}
}

func TestService_Repeat_sucsses(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	_, err = s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}
	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): can not find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusInProgress {
		t.Errorf("Repeat(): status did not changed, error = %v", err)
	}
	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Repeat(): can not find account by id, error = %v", err)
		return
	}

	if savedAccount.Balance == defaultTestAccount.balance {
		t.Errorf("Repeat(): balance did not changed, error = %v", err)
		return
	}
}

func TestService_FavoritePayment_sucsses(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "auto")
	if err != nil {
		t.Errorf("PayFromFavorite(): can not add favorite, error = %v", err)
		return
	}
}

func TestService_PayFromFavorite_sucsses(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	fv, err := s.FavoritePayment(payment.ID, "auto")
	if err != nil {
		t.Errorf("PayFromFavorite(): can not add favorite, error = %v", err)
		return
	}

	_, err = s.PayFromFavorite(fv.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): can not find favorite, error = %v", err)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("PayFromFavorite(): can not find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance == defaultTestAccount.balance {
		t.Errorf("PayFromFavorite(): balance did not changed, old = %v, new = %v", defaultTestAccount.balance, savedAccount.Balance)
		return
	}

}

func Test_Service_FindFavoriteByID_Fail(t *testing.T) {
	s := newTestService()
	_, err := s.FindFavoriteByID("")
	if err == nil {
		t.Error(err)
		return
	}
}


func Test_Service_ExportToFile(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.ExportToFile("acc.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func Test_Service_ImportFromFile(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.ImportFromFile("acc.txt")
	if err != nil {
		t.Error(err)
		return
	}
}

func Test_Service_Export(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "auto")
	if err != nil {
		t.Errorf("PayFromFavorite(): can not add favorite, error = %v", err)
		return
	}

	err = s.Export("cmd")
	if err != nil {
		t.Error(err)
		return
	}
}

func Test_Service_Import(t *testing.T) {
	s := newTestService()

	err := s.Import("cmd")
	if err != nil {
		t.Error(err)
		return
	}
}

func Test_Service_ExportAccountHistory(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = s.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
		return
	}
}

func Test_Service_HistoryToFiles(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	his, _ := s.ExportAccountHistory(1)

	err = s.HistoryToFiles(his, "cmd/3", 7)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.HistoryToFiles(his, "", 7)
	if err == nil {
		t.Error(err)
		return
	}	

}

func Test_Service_SumPayments(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		log.Print(err)
		return
	}

	s.SumPayments(1)
}

func testPays(s *testService) {

	account, _ := s.RegisterAccount("+992000000001")
	s.Deposit(account.ID, 80000)
	_, _ = s.Pay(account.ID, 1, "auto")
	_, _ = s.Pay(account.ID, 2, "auto")
	_, _ = s.Pay(account.ID, 3, "auto")
	_, _ = s.Pay(account.ID, 4, "auto")
	account, _ = s.RegisterAccount("+992000000002")
	s.Deposit(account.ID, 80000)
	_, _ = s.Pay(account.ID, 5, "auto")
	_, _ = s.Pay(account.ID, 6, "auto")
	_, _ = s.Pay(account.ID, 7, "auto")
	_, _ = s.Pay(account.ID, 8, "auto")
}

func BenchmarkSumPayments(b *testing.B) {
	s := newTestService()
	testPays(s)
	want := 36
	result := s.SumPayments(1)
	if result != types.Money(want) {
		b.Fatalf("invalid result, got %v, want %v", result, want)
	}
}