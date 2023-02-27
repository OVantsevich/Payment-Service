package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/Payment-Service/internal/model"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAccount_CreateAccount(t *testing.T) {
	var err error
	ctx := context.Background()
	var testData = []*model.Account{
		{
			ID:      uuid.NewString(),
			User:    uuid.NewString(),
			Amount:  100,
			Created: time.Now(),
			Updated: time.Now(),
		},
	}

	for _, p := range testData {
		_, err = testAccountRepository.CreateAccount(ctx, p)
		require.NoError(t, err)

		_, err = testAccountRepository.CreateAccount(ctx, p)
		require.Error(t, err)
	}
}

func TestAccount_GetUserAccount(t *testing.T) {
	var err error
	ctx := context.Background()
	var testData = []*model.Account{
		{
			ID:      uuid.NewString(),
			User:    uuid.NewString(),
			Amount:  100,
			Created: time.Now(),
			Updated: time.Now(),
		},
	}

	var acc *model.Account
	for _, p := range testData {
		_, err = testAccountRepository.CreateAccount(ctx, p)
		require.NoError(t, err)

		_, err = testAccountRepository.CreateAccount(ctx, p)
		require.Error(t, err)

		acc, err = testAccountRepository.GetUserAccount(ctx, p.User)
		require.NoError(t, err)
		require.Equal(t, acc.ID, p.ID)

		acc, err = testAccountRepository.GetUserAccount(ctx, uuid.NewString())
		require.Error(t, err)
	}
}

func TestAccount_GetAccountForUpdate(t *testing.T) {
	var err error
	ctx := context.Background()
	var testData = &model.Account{
		ID:      uuid.NewString(),
		User:    uuid.NewString(),
		Amount:  100,
		Created: time.Now(),
		Updated: time.Now(),
	}
	_, err = testAccountRepository.CreateAccount(ctx, testData)
	require.NoError(t, err)
	err = testAccountRepository.UpdateAmount(ctx, testData.ID, 100)
	require.NoError(t, err)

	var start sync.Mutex
	var wg sync.WaitGroup
	start.Lock()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go testTransactor.WithinTransaction(context.Background(), func(trxCtx context.Context) error {
			defer wg.Done()
			acc, trxErr := testAccountRepository.GetAccountForUpdate(trxCtx, testData.ID)
			if acc.Amount < 100 {
				return fmt.Errorf("not enough money")
			}
			trxErr = testAccountRepository.UpdateAmount(trxCtx, testData.ID, -100)
			return trxErr
		})
	}

	wg.Wait()
	acc, err := testAccountRepository.GetUserAccount(ctx, testData.User)
	require.NoError(t, err)
	require.Equal(t, 0.0, acc.Amount)
}
