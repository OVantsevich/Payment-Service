// Package service account service
package service

import (
	"context"
	"fmt"

	"Payment-Service/internal/model"

	"github.com/google/uuid"
)

// AccountRepository repository interface for account service
//
//go:generate mockery --name=AccountRepository --case=underscore --output=./mocks
type AccountRepository interface {
	CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	GetUserAccount(ctx context.Context, id string) (*model.Account, error)
	GetUserAccountForUpdate(ctx context.Context, id string) (*model.Account, error)
	UpdateAmount(ctx context.Context, id string, amount float64) error
}

// Account account service
type Account struct {
	rps AccountRepository
}

// NewAccount new user service
func NewAccount(rps AccountRepository) *Account {
	return &Account{rps: rps}
}

// CreateAccount service create account
func (a *Account) CreateAccount(ctx context.Context, account *model.Account) (accountResult *model.Account, err error) {
	account.ID = uuid.New().String()
	accountResult, err = a.rps.CreateAccount(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("account - CreateAccount - CreateAccount: %w", err)
	}

	return
}

// GetUserAccount service get user account
func (a *Account) GetUserAccount(ctx context.Context, id string) (account *model.Account, err error) {
	account, err = a.rps.GetUserAccount(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("account - GetUserAccount - GetUserAccount: %w", err)
	}

	return
}

// GetUserAccountForUpdate service get user account for update
func (a *Account) GetUserAccountForUpdate(ctx context.Context, id string) (account *model.Account, err error) {
	account, err = a.rps.GetUserAccountForUpdate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("account - GetUserAccountForUpdate - GetUserAccountForUpdate: %w", err)
	}

	return
}

// UpdateAmount update account amount
func (a *Account) UpdateAmount(ctx context.Context, id string, amount float64) error {
	err := a.rps.UpdateAmount(ctx, id, amount)
	if err != nil {
		return fmt.Errorf("account - GetUserAccount - GetUserAccount: %w", err)
	}

	return nil
}
