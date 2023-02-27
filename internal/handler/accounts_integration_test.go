package handler

import (
	"context"
	"github.com/OVantsevich/Payment-Service/internal/model"
	pr "github.com/OVantsevich/Payment-Service/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAccounts_CreateAccount(t *testing.T) {
	ctx := context.Background()
	var testData = &model.Account{
		User:    uuid.NewString(),
		Amount:  100,
		Created: time.Now(),
		Updated: time.Now(),
	}

	resp, err := testAccountHandler.CreateAccount(ctx, &pr.CreateAccountRequest{UserID: testData.User})
	require.NoError(t, err)
	testData.ID = resp.Account.ID
}

func TestAccounts_GetAccount(t *testing.T) {
	ctx := context.Background()
	var testData = &model.Account{
		User:    uuid.NewString(),
		Amount:  100,
		Created: time.Now(),
		Updated: time.Now(),
	}

	respC, err := testAccountHandler.CreateAccount(ctx, &pr.CreateAccountRequest{UserID: testData.User})
	require.NoError(t, err)
	testData.ID = respC.Account.ID

	respG, err := testAccountHandler.GetAccount(ctx, &pr.GetAccountRequest{UserID: testData.User})
	require.NoError(t, err)
	require.Equal(t, respG.Account.Amount, 0.0)
}

func TestAccounts_IncreaseAmount(t *testing.T) {
	ctx := context.Background()
	var testData = &model.Account{
		User:    uuid.NewString(),
		Amount:  100,
		Created: time.Now(),
		Updated: time.Now(),
	}

	respC, err := testAccountHandler.CreateAccount(ctx, &pr.CreateAccountRequest{UserID: testData.User})
	require.NoError(t, err)
	testData.ID = respC.Account.ID

	_, err = testAccountHandler.IncreaseAmount(ctx, &pr.AmountRequest{AccountID: respC.Account.ID, Amount: testData.Amount})
	require.NoError(t, err)

	respG, err := testAccountHandler.GetAccount(ctx, &pr.GetAccountRequest{UserID: testData.User})
	require.NoError(t, err)
	require.Equal(t, respG.Account.Amount, testData.Amount)
}

func TestAccounts_DecreaseAmount(t *testing.T) {
	ctx := context.Background()
	var testData = &model.Account{
		User:    uuid.NewString(),
		Amount:  100,
		Created: time.Now(),
		Updated: time.Now(),
	}

	respC, err := testAccountHandler.CreateAccount(ctx, &pr.CreateAccountRequest{UserID: testData.User})
	require.NoError(t, err)
	testData.ID = respC.Account.ID

	_, err = testAccountHandler.DecreaseAmount(ctx, &pr.AmountRequest{AccountID: respC.Account.ID, Amount: testData.Amount})
	require.Error(t, err)

	_, err = testAccountHandler.IncreaseAmount(ctx, &pr.AmountRequest{AccountID: respC.Account.ID, Amount: testData.Amount})
	require.NoError(t, err)

	respG, err := testAccountHandler.GetAccount(ctx, &pr.GetAccountRequest{UserID: testData.User})
	require.NoError(t, err)
	require.Equal(t, respG.Account.Amount, testData.Amount)

	_, err = testAccountHandler.DecreaseAmount(ctx, &pr.AmountRequest{AccountID: respC.Account.ID, Amount: testData.Amount})
	require.NoError(t, err)

	_, err = testAccountHandler.DecreaseAmount(ctx, &pr.AmountRequest{AccountID: respC.Account.ID, Amount: testData.Amount})
	require.Error(t, err)
}
