package usecase

import (
	"errors"
	"log"

	"github.com/jerryaugusto/codepix-go/domain/model"
)

type TransactionUseCase struct {
	TransactionRepository model.TransactionRepositoryInterface
	PixRepository         model.PixKeyRepositoryInterface
}

func (t *TransactionUseCase) Register(accountId string, amount float64, pixKeyTo string, pixKeyKindTo string, description string) (*model.Transaction, error) {
	account, err := t.PixRepository.FindAccount(accountId)
	if err != nil {
		return nil, err
	}

	pixKey, err := t.PixRepository.FindKeyByKind(pixKeyTo, pixKeyKindTo)
	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(account, amount, pixKey, description)
	if err != nil {
		return nil, err
	}

	t.TransactionRepository.TransactionSave(transaction)
	if transaction.ID != "" {
		return transaction, nil
	}

	return nil, errors.New("unable to process this transaction")
}

func (t *TransactionUseCase) Confim(transactionId string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.TransactionFind(transactionId)
	if err != nil {
		log.Println("Transaction not found", transactionId)
		return nil, err
	}

	transaction.Status = model.TransactionConfirmed
	err = t.TransactionRepository.TransactionSave(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
