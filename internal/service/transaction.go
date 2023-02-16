// Package service transaction service
package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"Payment-Service/internal/model"
)

// TransactionRepository repository interface for account service
//
//go:generate mockery --name=TransactionRepository --case=underscore --output=./mocks
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	GetAccountTransactions(ctx context.Context, id string) (map[string]model.Transaction, error)
}

// Transaction transaction service
type Transaction struct {
	rps TransactionRepository
}

// NewTransaction new user service
func NewTransaction(rps TransactionRepository) *Transaction {
	return &Transaction{rps: rps}
}

// CreateTransaction service create transaction
func (t *Transaction) CreateTransaction(ctx context.Context, transaction *model.Transaction) (transactionResult *model.Transaction, err error) {
	transaction.ID = uuid.New().String()
	transactionResult, err = t.rps.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("transaction - CreateTransaction - CreateTransaction: %w", err)
	}

	return
}

// GetAccountTransactions service get account transactions
func (t *Transaction) GetAccountTransactions(ctx context.Context, id string) (transactions map[string]model.Transaction, err error) {
	transactions, err = t.rps.GetAccountTransactions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("transaction - GetAccountTransactions - GetAccountTransactions: %w", err)
	}

	return
}
