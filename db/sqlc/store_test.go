package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	randomFromAccount := createRandomAccount(t)
	randomToAccount := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: randomFromAccount.ID,
				ToAccountID:   randomToAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, randomFromAccount.ID, transfer.FromAccountID)
		require.Equal(t, randomToAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, randomFromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, randomToAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, randomFromAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, randomToAccount.ID, toAccount.ID)

		// check accounts' balance
		fromAccountBalanceDiff := randomFromAccount.Balance - fromAccount.Balance
		toAccountBalanceDiff := toAccount.Balance - randomToAccount.Balance
		require.Equal(t, fromAccountBalanceDiff, toAccountBalanceDiff)
		require.True(t, fromAccountBalanceDiff > 0)
		require.True(t, fromAccountBalanceDiff%amount == 0)

		k := int(fromAccountBalanceDiff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedFromAccount, err := store.GetAccount(context.Background(), randomFromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := store.GetAccount(context.Background(), randomToAccount.ID)
	require.NoError(t, err)

	require.Equal(t, randomFromAccount.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, randomToAccount.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	randomFromAccount := createRandomAccount(t)
	randomToAccount := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := randomFromAccount.ID
		toAccountID := randomToAccount.ID

		if i%2 == 1 {
			fromAccountID = randomToAccount.ID
			toAccountID = randomFromAccount.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	// check the final updated balance
	updatedFromAccount, err := store.GetAccount(context.Background(), randomFromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := store.GetAccount(context.Background(), randomToAccount.ID)
	require.NoError(t, err)

	require.Equal(t, randomFromAccount.Balance, updatedFromAccount.Balance)
	require.Equal(t, randomToAccount.Balance, updatedToAccount.Balance)
}
