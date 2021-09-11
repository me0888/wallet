package wallet

import (
	// "bufio"
	"errors"
	"fmt"
	"io"
	"sync"

	// "io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/me0888/wallet/pkg/types"
)

var ErrPhoneRegistered = errors.New("phone alredy registred")
var ErrAmmountMustBePositive = errors.New("amount must be greater 0")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("Not Enough Balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("Favorite not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}

	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount

	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound

	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound

	}

	return account, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	for _, pay := range s.payments {
		if pay.ID == paymentID {
			payment = pay
		}
	}

	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return payment, nil
}

func (s *Service) Reject(paymentID string) error {
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	pay.Status = types.PaymentStatusFail

	acc, err := s.FindAccountByID(pay.AccountID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	acc.Balance += pay.Amount

	return nil

}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	acc, err := s.FindAccountByID(pay.AccountID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	payID := uuid.New().String()
	payment := &types.Payment{
		ID:        payID,
		AccountID: acc.ID,
		Amount:    pay.Amount,
		Category:  pay.Category,
		Status:    types.PaymentStatusInProgress,
	}

	acc.Balance -= pay.Amount

	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := uuid.New().String()
	newFavorite := &types.Favorite{
		ID:        favorite,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, newFavorite)

	return newFavorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {

	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
			return fav, nil
		}
	}
	return nil, ErrFavoriteNotFound

}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	newPaymntID := uuid.New().String()
	newPayment := &types.Payment{
		ID:        newPaymntID,
		AccountID: favorite.AccountID,
		Amount:    favorite.Amount,
		Category:  favorite.Category,
		Status:    types.PaymentStatusInProgress,
	}

	acc, err := s.FindAccountByID(favorite.AccountID)
	if err != nil {
		return nil, err
	}

	acc.Balance -= favorite.Amount

	s.payments = append(s.payments, newPayment)
	return newPayment, nil

}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()
	text := ""

	for _, data := range s.accounts {
		text += strconv.FormatInt(data.ID, 10) + ";" +
			string(data.Phone) + ";" +
			strconv.FormatInt(int64(data.Balance), 10) + "|"
	}
	text = text[:len(text)-1]
	_, err = file.Write([]byte(text))

	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	contetnt := make([]byte, 0)
	buf := make([]byte, 4)

	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			contetnt = append(contetnt, buf[:read]...)
			break
		}

		if err != nil {
			log.Println(err)
			return err
		}

		contetnt = append(contetnt, buf[:read]...)
	}
	accs := strings.Split(string(contetnt), "|")
	for _, acc := range accs {
		data := strings.Split(acc, ";")
		id, _ := strconv.ParseInt(data[0], 10, 64)
		balance, _ := strconv.ParseInt(data[2], 10, 64)
		s.accounts = append(s.accounts, &types.Account{
			ID:      id,
			Phone:   types.Phone(data[1]),
			Balance: types.Money(balance),
		})
	}

	for _, v := range s.accounts {
		log.Println(v)
	}
	return nil
}

func (s *Service) Export(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModeDir)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	result := ""
	for _, acc := range s.accounts {
		result += strconv.FormatInt(acc.ID, 10) + ";" +
			string(acc.Phone) + ";" +
			strconv.FormatInt(int64(acc.Balance), 10) + "\n"

		if result != "" {
			err = os.WriteFile(dir+"/accounts.dump", []byte(result), 0666)
			if err != nil {
				log.Print(err)
				return err
			}
		}
	}

	result = ""
	for _, payment := range s.payments {

		result += payment.ID + ";" +
			strconv.FormatInt(payment.AccountID, 10) + ";" +
			strconv.FormatInt(int64(payment.Amount), 10) + ";" +
			string(payment.Category) + ";" +
			string(payment.Status) + "\n"

		if result != "" {
			err = os.WriteFile(dir+"/payments.dump", []byte(result), 0666)
			if err != nil {
				log.Print(err)
				return err
			}
		}

	}

	result = ""
	for _, favorite := range s.favorites {
		result += favorite.ID + ";" +
			strconv.FormatInt(favorite.AccountID, 10) + ";" +
			favorite.Name + ";" +
			strconv.FormatInt(int64(favorite.Amount), 10) + ";" +
			string(favorite.Category) + "\n"

		if result != "" {
			err = os.WriteFile(dir+"/favorites.dump", []byte(result), 0666)
			if err != nil {
				log.Print(err)
				return err
			}
		}

	}

	return err
}

func filesExist(path string) bool {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return true
	}
	return false
}

func (s *Service) Import(dir string) error {

	if filesExist(dir + "/accounts.dump") {
		src, err := os.ReadFile(dir + "/accounts.dump")
		if err != nil {
			log.Println(err)
			return err
		}

		line := strings.TrimSuffix(string(src), "\n")
		line = strings.TrimSuffix(line, "\r")
		accs := strings.Split(line, "\n")

		for _, acc := range accs {
			data := strings.Split(acc, ";")
			id, _ := strconv.ParseInt(data[0], 10, 64)
			balance, _ := strconv.ParseInt(data[2], 10, 64)
			acc, err := s.FindAccountByID(id)
			if err != nil {
				acc2, err := s.RegisterAccount(types.Phone(data[1]))
				if err != nil {
					log.Print(err)
					return err
				}
				acc2.Balance = types.Money(balance)
			} else {
				acc.Phone = types.Phone(data[1])
				acc.Balance = types.Money(balance)
			}
		}
	}

	if filesExist(dir + "/payments.dump") {
		src, err := os.ReadFile(dir + "/payments.dump")
		if err != nil {
			log.Println(err)
			return err
		}
		line := strings.TrimSuffix(string(src), "\n")
		line = strings.TrimSuffix(line, "\r")
		pays := strings.Split(line, "\n")

		for _, acc := range pays {
			data := strings.Split(acc, ";")
			accountID, _ := strconv.ParseInt(data[1], 10, 64)
			amount, _ := strconv.ParseInt(data[2], 10, 64)
			category := types.PaymentCategory(data[3])
			status := types.PaymentStatus(data[4])

			pay, err := s.FindPaymentByID(data[0])
			if err != nil {
				payment := &types.Payment{
					ID:        data[0],
					AccountID: accountID,
					Amount:    types.Money(amount),
					Category:  category,
					Status:    status,
				}
				s.payments = append(s.payments, payment)
			} else {
				pay.AccountID = accountID
				pay.Amount = types.Money(amount)
				pay.Category = category
				pay.Status = status
			}

		}
	}

	if filesExist(dir + "/favorites.dump") {
		src, err := os.ReadFile(dir + "/favorites.dump")
		if err != nil {
			log.Println(err)
			return err
		}

		line := strings.TrimSuffix(string(src), "\n")
		line = strings.TrimSuffix(line, "\r")
		favs := strings.Split(line, "\n")
		for _, acc := range favs {
			data := strings.Split(acc, ";")

			accountID, _ := strconv.ParseInt(data[1], 10, 64)
			amount, _ := strconv.ParseInt(data[3], 10, 64)
			id := data[0]
			name := data[2]
			category := data[4]

			favorite, err := s.FindFavoriteByID(id)
			if err != nil {
				nfavorite := &types.Favorite{
					ID:        id,
					AccountID: accountID,
					Name:      name,
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(category),
				}
				s.favorites = append(s.favorites, nfavorite)

			} else {
				favorite.AccountID = accountID
				favorite.Name = name
				favorite.Amount = types.Money(amount)
				favorite.Category = types.PaymentCategory(category)
			}

		}
	}

	return nil
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		log.Println(err)
		return nil, ErrAccountNotFound
	}

	result := []types.Payment{}
	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			result = append(result, *payment)
		}
	}
	return result, nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	if len(payments) == 0 {
		return nil
	}

	if !filesExist(dir) {
		err := os.Mkdir(dir, os.ModeDir)
		if err != nil {
			log.Print(err)
			return err
		}
	}
	var content []string
	for _, payment := range payments {
		content = append(content,
			payment.ID+";"+
				strconv.FormatInt(payment.AccountID, 10)+";"+
				strconv.FormatInt(int64(payment.Amount), 10)+";"+
				string(payment.Category)+";"+
				string(payment.Status)+"\n")

	}
	log.Println(len(payments))
	if len(payments) <= records {
		err := os.WriteFile(dir+"/payments.dump", []byte(strings.Join(content, "")), 0666)
		if err != nil {
			log.Print("err from HistoryToFiles", err)
			return err
		}
		return nil
	}

	length := len(content)
	ln := length / records
	if length%records != 0 {
		ln = ln + 1
	}
	start := 0
	end := records

	for i := 0; i < ln; i++ {
		err := os.WriteFile(dir+"/payments"+strconv.Itoa(i+1)+".dump", []byte(strings.Join(content[start:end], "")), 0666)
		if err != nil {
			log.Print("err from HistoryToFiles", err)
			return err
		}
		start += records
		end += records
		if end > length {
			end = length
		}
	}
	return nil

}



func (s *Service) SumPayments(goroutines int) types.Money {
	if goroutines == 0 {
		goroutines = 1
	}
	wg := sync.WaitGroup{}
	wg.Add(goroutines)
	mu := sync.Mutex{}
	sum := 0
	lengthOfPayments := len(s.payments)
	count := lengthOfPayments / goroutines
	if lengthOfPayments%goroutines != 0 {
		count = count + 1
	}
	for i := 0; i < goroutines; i++ {
		start := i * count
		end := (i + 1) * count
		if end >= lengthOfPayments {
			end = lengthOfPayments
		}
		if start >= lengthOfPayments {
			break
		}
		go func() {
			defer wg.Done()
			val := int64(0)
			for i := start; i < end; i++ {
				val += int64(s.payments[i].Amount)
			}
			mu.Lock()
			sum += int(val)
			mu.Unlock()
		}()

	}

	wg.Wait()
	return types.Money(sum)
}

func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}
	if goroutines == 0 {
		goroutines = 1
	}
	wg := sync.WaitGroup{}
	wg.Add(goroutines)
	mu := sync.Mutex{}
	filter := []types.Payment{}
	lengthOfPayments := len(s.payments)
	count := lengthOfPayments / goroutines
	if lengthOfPayments%goroutines != 0 {
		count = count + 1
	}
	for i := 0; i < goroutines; i++ {
		start := i * count
		end := (i + 1) * count
		if end >= lengthOfPayments {
			end = lengthOfPayments
		}
		if start >= lengthOfPayments {
			break
		}
		go func() {
			defer wg.Done()
			for i := start; i < end; i++ {
				if s.payments[i].AccountID == accountID {
					mu.Lock()
					filter = append(filter, types.Payment{
						ID:        s.payments[i].ID,
						AccountID: s.payments[i].AccountID,
						Amount:    s.payments[i].Amount,
						Category:  s.payments[i].Category,
						Status:    s.payments[i].Status,
					})
					mu.Unlock()
				}
			}
		}()

	}

	wg.Wait()
	return filter, nil
}