package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/OVantsevich/Payment-Service/internal/model"
)

// Account postgres entity
type Account struct {
	PgxWithinTransactionRunner
}

// NewAccount creating new account repository
func NewAccount(p PgxWithinTransactionRunner) *Account {
	return &Account{PgxWithinTransactionRunner: p}
}

// CreateAccount create account
func (r *Account) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	account.Created = time.Now()
	account.Updated = time.Now()
	_, err := r.Exec(ctx,
		`insert into accounts (id, "user", created, updated) values ($1, $2, $3, $4);`,
		account.ID, account.User, account.Created, account.Updated)
	if err != nil {
		return nil, fmt.Errorf("account - CreateAccount - Exec: %w", err)
	}

	return account, nil
}

// GetUserAccount get user account
func (r *Account) GetUserAccount(ctx context.Context, userID string) (*model.Account, error) {
	account := model.Account{}
	err := r.QueryRow(ctx, `select id, amount from accounts where "user"=$1 and deleted=false`, userID).Scan(
		&account.ID, &account.Amount)
	if err != nil {
		return nil, fmt.Errorf("account - GetUserAccount - QueryRow: %w", err)
	}

	return &account, nil
}

// GetUserAccountForUpdate get user account for update
func (r *Account) GetUserAccountForUpdate(ctx context.Context, userID string) (*model.Account, error) {
	account := model.Account{}
	err := r.QueryRow(ctx, `select id, amount from accounts where "user"=$1 and deleted=false for update`, userID).Scan(
		&account.ID, &account.Amount)
	if err != nil {
		return nil, fmt.Errorf("account - GetUserAccountForUpdate - QueryRow: %w", err)
	}

	return &account, nil
}

// UpdateAmount update user account amount
func (r *Account) UpdateAmount(ctx context.Context, accountID string, amount float64) error {
	_, err := r.Exec(ctx, `update accounts set amount=amount+$1 where id=$2 and deleted=false`, amount, accountID)
	if err != nil {
		return fmt.Errorf("account - UpdateAmount - Exec: %w", err)
	}

	return nil
}
