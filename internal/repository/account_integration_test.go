package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/Payment-Service/internal/model"
	"strconv"
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

func TestAccount_GetUserAccountForUpdate(t *testing.T) {
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

	var wg sync.WaitGroup
	for _, p := range testData {
		_, err = testAccountRepository.CreateAccount(ctx, p)
		require.NoError(t, err)

		_, err = testAccountRepository.CreateAccount(ctx, p)
		require.Error(t, err)

		var inTrx int64
		var inGo int64

		wg.Add(1)
		err = testTransactor.WithinTransaction(context.Background(), func(trxCtx context.Context) error {
			var trxErr error
			_, trxErr = testAccountRepository.GetUserAccountForUpdate(trxCtx, p.User)
			require.NoError(t, trxErr)
			go func() {
				_, e := testAccountRepository.GetUserAccountForUpdate(ctx, p.User)
				require.NoError(t, e)
				time.Sleep(time.Second)
				inGo = time.Now().Unix()
				fmt.Println("go: " + strconv.FormatInt(inGo, 10))
				wg.Done()
			}()
			time.Sleep(time.Second * 2)
			inTrx = time.Now().Unix()
			fmt.Println("trx: " + strconv.FormatInt(inTrx, 10))
			trxErr = testAccountRepository.UpdateAmount(trxCtx, p.ID, p.Amount)
			require.NoError(t, trxErr)
			return trxErr
		})
		require.NoError(t, err)
		constr := inTrx < inGo
		require.True(t, constr)
		wg.Wait()
	}
}
