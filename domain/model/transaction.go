package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   string = "pending"
	TransactionCompleted string = "completed"
	TransactionError     string = "error"
	TransactinConfirmed  string = "confirmed"
)

type TransactionRepositoryInterface interface {
	TransactionRegister(transaction *Transaction) error
	TransactionSave(transaction *Transaction) error
	TransactionFind(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []Transaction
}

type Transaction struct {
	Base              `valid:"required" json:"base,omitempty"`
	AccountFrom       *Account `valid:"-"`
	Amount            float64  `json:"amount,omitempty" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	Status            string   `json:"status,omitempty" valid:"notnull"`
	Description       string   `json:"description,omitempty" valid:"notnull"`
	CancelDescription string   `json:"cancel_description,omitempty" valid:"-"`
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)

	if t.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionError {
		return errors.New("invalid status for the transaction")
	}

	if t.PixKeyTo.Account.ID == t.AccountFrom.ID {
		return errors.New("the source and destination account cannot be the same")
	}

	if err != nil {
		return err
	}

	return nil
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKeyTo,
		Status:      TransactionPending,
		Description: description,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.isValid()
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()

	err := t.isValid()

	return err
}

func (t *Transaction) Confirm() error {
	t.Status = TransactinConfirmed
	t.UpdatedAt = time.Now()

	err := t.isValid()

	return err
}

func (t *Transaction) Cancel(description string) error {
	t.Status = TransactionError
	t.UpdatedAt = time.Now()
	t.Description = description

	err := t.isValid()

	return err
}
