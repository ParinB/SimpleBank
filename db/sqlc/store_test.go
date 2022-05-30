package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T)  {
	store := NewStore(testDB)
	// Create two random accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	/*
	Send  errors to main go routine that  our test in running on
	We can  use channels it is used to connect concurrent go routines share data with each other with out explicit locking
	We need two channels 1.To receive errors and 2.To receive transactiontx results
	*/
	for i :=0; i < n; i++ {
		go func() {
			result , err := store.TransferTx(context.Background(),TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			errs <- err
			results <- result
		}()
	}
	// Check Results
	for i := 0; i < n ; i ++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID ,transfer.FromAccountID)
		require.Equal(t, account2.ID,transfer.ToAccountID)
		require.Equal(t, amount,transfer.Amount)
		require.NotZero(t, transfer.ID )
		require.NotZero(t, transfer.CreatedAt)
		//check from  entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID,fromEntry.AccountID)
		require.Equal(t, -amount,fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_,err = store.GetEntry(context.Background(),fromEntry.ID)
		require.NoError(t, err)
		//check to entries
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account1.ID,toEntry.AccountID)
		require.Equal(t, -amount,toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_,err = store.GetEntry(context.Background(),toEntry.ID)
		require.NoError(t, err)
		//TOD0 :check accounts' balance


	}

}
