// Package repository transaction
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/OVantsevich/Payment-Service/internal/model"
)

// Transaction postgres entity
type Transaction struct {
	PgxWithinTransactionRunner
}

// NewTransaction creating new transaction repository
func NewTransaction(p PgxWithinTransactionRunner) *Transaction {
	return &Transaction{PgxWithinTransactionRunner: p}
}

// CreateTransaction create transaction
func (r *Transaction) CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	transaction.Created = time.Now()
	transaction.Updated = time.Now()
	_, err := r.Exec(ctx,
		"insert into transactions (id, account, amount, created, updated) values ($1, $2, $3, $4, $5);",
		transaction.ID, transaction.Account, transaction.Amount, transaction.Created, transaction.Updated)
	if err != nil {
		return nil, fmt.Errorf("transaction - CreateTransaction - Exec: %w", err)
	}

	return transaction, nil
}

// GetAccountTransactions get account transaction
func (r *Transaction) GetAccountTransactions(ctx context.Context, accountID string) (map[string]model.Transaction, error) {
	transactions := make(map[string]model.Transaction, 0)
	rows, err := r.Query(ctx, `select id, amount, created from transactions where account=$1 and deleted=false`, accountID)
	if err != nil {
		return nil, fmt.Errorf("transaction - GetAccountTransactions - QueryRow: %w", err)
	}

	for rows.Next() {
		trs := model.Transaction{
			Account: accountID,
		}
		err = rows.Scan(&trs.ID, &trs.Amount, &trs.Created)
		if err != nil {
			return nil, fmt.Errorf("transaction - GetAccountTransactions - Scan: %w", err)
		}
		transactions[trs.ID] = trs
	}

	return transactions, nil
}
