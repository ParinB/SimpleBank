package db

import (
	"context"
	"database/sql"
	"fmt"
)

/*
Embedding queries in store struct (Composition) extending struct functionalities
All individual query function provides by queries will be  available to store
*/

// Store provides all functions to execute sql queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context,arg TransferTxParams)(TransferTxResult,error)
}
// SqlStore provides all functions to execute sql queries and transactions
type SqlStore struct {
	*Queries
	db *sql.DB
}
// NewStore creates a new Store
func NewStore(db *sql.DB)  *SqlStore{
	return  &SqlStore{
		db: db ,
		Queries : New(db),
	}
}
// execTx executes a function within a database transaction
func  (store *SqlStore) execTx(ctx context.Context, fn func( *Queries ) error ) error{
	tx,err := store.db.BeginTx(ctx,nil)
	if  err !=nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err !=nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v",err,rbErr)
		}
		return err
	}
	return tx.Commit()
}
// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID int64  `json:"t\o_account_id"`
	Amount int64 `json:"amount"`
}
// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}
// TransferTxParams contains the input parameters of the transfer transaction
// It creates  a transfer record , add account entries and update account balance within a single database transaction

var txKey = struct {}{}

func (store *SqlStore)  TransferTx(ctx context.Context,arg TransferTxParams)(TransferTxResult,error){
	var result TransferTxResult
	/*
	We are accessing result variable of the outer function   from inside  the callback function  makes callback function become a closure
	Closure is often used when we want to get the result from a  callback function , because callback function does not know the type of the result it should return
	*/
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)

		// Creates a transfer record with the  FromAccountID , ToAccountID , Amount
		fmt.Println(txName,"create transfer")
		result.Transfer,err = q.CreateTransfer(ctx , CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName,"create entry 1")
		// Creates an Entry with the account & amount the money is coming from  (fromAccount)
		result.FromEntry,err =  q.CreateEntry(ctx , CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName,"create entry 2")
		// Creates an entry with account & amount the money is going to (toAccount)
		result.ToEntry,err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID:  arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		if arg.FromAccountID < arg.ToAccountID {
			// Update the balance of the account money is moving from (subtraction)
			result.FromAccount,err =q.AddAccountBalance(ctx,AddAccountBalanceParams{
				ID: arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return  err
			}
			result.ToAccount,err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
				ID : arg.ToAccountID,
				Amount : arg.Amount,
			})
			if err != nil {
				return err
			}
		}else {
			//Update the balance of the account money is moving to
			fmt.Println(txName,"update account 2")
			result.ToAccount,err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
				ID : arg.ToAccountID,
				Amount : arg.Amount,
			})
			if err != nil {
				return err
			}

			// Update the balance of the account money is moving from (subtraction)
			result.FromAccount,err =q.AddAccountBalance(ctx,AddAccountBalanceParams{
				ID: arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return  err
			}
		}
		// Get  the account balance of account money is moving from (FromAccount)
		fmt.Println(txName,"update account 1")


		return  nil
	})
	return result,err
}

