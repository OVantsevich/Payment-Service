package repository

import (
	"context"
	"github.com/OVantsevich/Payment-Service/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestTransaction_CreateTransaction(t *testing.T) {
	var err error
	ctx := context.Background()
	var testAccount = &model.Account{
		ID:      uuid.NewString(),
		User:    uuid.NewString(),
		Amount:  100,
		Created: time.Now(),
		Updated: time.Now(),
	}
	var testTransaction = &model.Transaction{
		ID:      uuid.NewString(),
		Account: testAccount.ID,
		Amount:  100,
		Created: time.Now(),
		Updated: time.Now(),
	}

	_, err = testTransactionRepository.CreateTransaction(ctx, testTransaction)
	require.Error(t, err)

	_, err = testAccountRepository.CreateAccount(ctx, testAccount)
	require.NoError(t, err)

	_, err = testTransactionRepository.CreateTransaction(ctx, testTransaction)
	require.NoError(t, err)
}

func TestTransaction_GetAccountTransactions(t *testing.T) {
	var err error
	ctx := context.Background()
	var testAccount = &model.Account{
		ID:      uuid.NewString(),
		User:    uuid.NewString(),
		Amount:  100,
		Created: time.Now(),
		Updated: time.Now(),
	}
	var testData = make([]*model.Transaction, 10)
	for i := range testData {
		testData[i] = &model.Transaction{
			ID:      uuid.NewString(),
			Account: testAccount.ID,
			Amount:  rand.Float64(),
			Created: time.Now(),
			Updated: time.Now(),
		}
	}

	_, err = testAccountRepository.CreateAccount(ctx, testAccount)
	require.NoError(t, err)

	for _, trz := range testData {
		_, err = testTransactionRepository.CreateTransaction(ctx, trz)
		require.NoError(t, err)
	}

	trzs, err := testTransactionRepository.GetAccountTransactions(ctx, testAccount.ID)
	require.NoError(t, err)
	require.Equal(t, len(testData), len(trzs))
}
