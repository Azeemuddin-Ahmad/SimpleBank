package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
	"time"
)

func createRamdomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	randomAccount := createRandomAccount(t)
	createRamdomEntry(t, randomAccount)
}

func TestGetEntry(t *testing.T) {
	randomAccount := createRandomAccount(t)
	randomEntry := createRamdomEntry(t, randomAccount)

	entry, err := testQueries.GetEntry(context.Background(), randomEntry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, randomEntry.ID, entry.ID)
	require.Equal(t, randomEntry.AccountID, entry.AccountID)
	require.Equal(t, randomEntry.Amount, entry.Amount)

	require.WithinDuration(t, randomEntry.CreatedAt, entry.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRamdomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
