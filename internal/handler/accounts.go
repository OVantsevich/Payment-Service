// Package handler handler
package handler

import (
	"context"
	"fmt"

	"github.com/OVantsevich/Payment-Service/internal/model"
	"github.com/OVantsevich/Payment-Service/internal/repository"
	pr "github.com/OVantsevich/Payment-Service/proto"

	"github.com/sirupsen/logrus"
)

// AccountService account service interface
//
//go:generate mockery --name=AccountService --case=underscore --output=./mocks
type AccountService interface {
	CreateAccount(ctx context.Context, account *model.Account) (accountResult *model.Account, err error)
	GetUserAccount(ctx context.Context, id string) (*model.Account, error)
	GetUserAccountForUpdate(ctx context.Context, id string) (*model.Account, error)
	UpdateAmount(ctx context.Context, id string, amount float64) error
}

// TransactionsService transaction service interface
//
//go:generate mockery --name=TransactionsService --case=underscore --output=./mocks
type TransactionsService interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	GetAccountTransactions(ctx context.Context, id string) (map[string]model.Transaction, error)
}

// Accounts handler
type Accounts struct {
	pr.UnimplementedPaymentServiceServer
	trService TransactionsService
	acService AccountService

	transactor repository.PgxTransactor
}

// NewAccountsHandler new user handler
func NewAccountsHandler(tr TransactionsService, ac AccountService, trx repository.PgxTransactor) *Accounts {
	return &Accounts{trService: tr, acService: ac, transactor: trx}
}

// CreateAccount handler create account
func (h *Accounts) CreateAccount(ctx context.Context, request *pr.CreateAccountRequest) (response *pr.CreateAccountResponse, err error) {
	user := &model.Account{
		User: request.UserID,
	}

	var accountResponse *model.Account
	response = &pr.CreateAccountResponse{}
	accountResponse, err = h.acService.CreateAccount(ctx, user)
	if err != nil {
		err = fmt.Errorf("user - CreateAccount - CreateAccount: %w", err)
		logrus.Error(err)
		return
	}

	response.Account = &pr.Account{
		ID:     accountResponse.ID,
		UserID: accountResponse.User,
		Amount: accountResponse.Amount,
	}

	return
}

// GetAccount handler get user account
func (h *Accounts) GetAccount(ctx context.Context, request *pr.GetAccountRequest) (response *pr.GetAccountResponse, err error) {
	response = &pr.GetAccountResponse{}
	var accountResponse *model.Account
	accountResponse, err = h.acService.GetUserAccount(ctx, request.UserID)
	if err != nil {
		err = fmt.Errorf("user - GetAccount - GetAccount: %w", err)
		logrus.Error(err)
		return
	}
	response.Account = &pr.Account{
		ID:     accountResponse.ID,
		UserID: accountResponse.User,
		Amount: accountResponse.Amount,
	}

	return
}

// IncreaseAmount handler increase amount
func (h *Accounts) IncreaseAmount(ctx context.Context, request *pr.AmountRequest) (response *pr.AmountResponse, err error) {
	err = h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		trxErr := h.acService.UpdateAmount(ctx, request.AccountID, request.Amount)
		if trxErr != nil {
			trxErr = fmt.Errorf("user - IncreaseAmount - UpdateAmount: %w", trxErr)
			logrus.Error(trxErr)
			return trxErr
		}

		_, trxErr = h.trService.CreateTransaction(ctx, &model.Transaction{
			Account: request.AccountID,
			Amount:  request.Amount,
		})
		if trxErr != nil {
			trxErr = fmt.Errorf("user - IncreaseAmount - CreateTransaction: %w", trxErr)
			logrus.Error(trxErr)
			return trxErr
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("user - IncreaseAmount - WithinTransaction: %w", err)
	}

	return
}

// DecreaseAmount handler decrease amount
func (h *Accounts) DecreaseAmount(ctx context.Context, request *pr.AmountRequest) (response *pr.AmountResponse, err error) {
	err = h.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		account, trxErr := h.acService.GetUserAccountForUpdate(ctx, request.AccountID)
		if trxErr != nil {
			trxErr = fmt.Errorf("user - DecreaseAmount - GetUserAccountForUpdate: %w", trxErr)
			logrus.Error(trxErr)
			return trxErr
		}
		if account.Amount < request.Amount {
			trxErr = fmt.Errorf("user - DecreaseAmount: not enough money")
			logrus.Error(trxErr)
			return trxErr
		}
		trxErr = h.acService.UpdateAmount(ctx, request.AccountID, -request.Amount)
		if trxErr != nil {
			trxErr = fmt.Errorf("user - DecreaseAmount - UpdateAmount: %w", trxErr)
			logrus.Error(trxErr)
			return trxErr
		}

		_, trxErr = h.trService.CreateTransaction(ctx, &model.Transaction{
			Account: request.AccountID,
			Amount:  -request.Amount,
		})
		if trxErr != nil {
			trxErr = fmt.Errorf("user - DecreaseAmount - CreateTransaction: %w", trxErr)
			logrus.Error(trxErr)
			return trxErr
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("user - DecreaseAmount - WithinTransaction: %w", err)
	}

	return response, nil
}
