package db

import (
	"context"
	"database/sql"
	"fmt"
)

// a database transaction is a way to group multiple database operations into a single unit of work that either all succeed or all fail.

type Store struct {
    *Queries
    db *sql.DB
} // Composition of Queries and db
// *Queries for individual query ..only for one table
// for multiple - transaction

func NewStore(db *sql.DB) *Store {
    return &Store{
        db:      db,
        Queries: New(db), // connection
    }
}

// in small letters - private function
// in capital letters - public function
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error { // fn func (*Queries) error : callback function
    tx, err := store.db.BeginTx(ctx, nil) // default isolation level
    if err != nil {
        return err
    }

    q := New(tx) // passing transaction to the queries
    err = fn(q)  // we have queries within transaction call input function (callback function)

    // if error = rollback
    if err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
        }
        return err
    }

    return tx.Commit() // all operations are successful
}

type TransferTxParams struct {
    FromAccountID int64 `json:"from_account_id"`
    ToAccountID   int64 `json:"to_account_id"`
    Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
    Transfer    Transfer `json:"transfer"` // created transfer record
    FromAccount Account  `json:"from_account"`
    // after balance updated
    ToAccount Account `json:"to_account"`
    // after balance updated
    FromEntry Entry `json:"from_entry"`
    // entry recording money moving out
    ToEntry Entry `json:"to_entry"`
    // entry recording money moving in
}

// Lets try money transfer example

// It creates a transfer record, add account entries for the from and to accounts, and update the account balances, all within a single database transaction.

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

    // Step1: Create a new transfer record in the database
    var result TransferTxResult

    err := store.execTx(ctx, func(q *Queries) error {

        var err error

        // q is queries object
        // q.CreateTransfer : calling create transfer function from queries

        result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
            FromAccountID: arg.FromAccountID,
            ToAccountID:   arg.ToAccountID,
            Amount:        arg.Amount,
        })

        if err != nil {
            return err
        }

        // Step2: Add account entries to the accounts of the sender and the receiver

        result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.FromAccountID, // sender
            Amount:    -arg.Amount,       // negative amount
        })
        if err != nil { // if error, rollback
            return err
        }

        result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.ToAccountID, // receiver
            Amount:    arg.Amount,      // positive amount
        })
        if err != nil { // if error, rollback
            return err
        }

        // Step3: Update the account balances of the sender and the receiver
        // [Ensure proper locking mechanism]
        // If you try without locking two transactions can update the same account at the same time

        // By using ForUpdate deadlock will occur if two transactions try to update the same row at the same time.

        // Update accounts in a consistent order to avoid deadlocks
        if arg.FromAccountID < arg.ToAccountID {
            // Update sender account first
            result.FromAccount, result.ToAccount, err = updateAccounts(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
        } else {
            // Update receiver account first
            result.ToAccount, result.FromAccount, err = updateAccounts(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
        }

        if err != nil {
            return err
        }

        return nil

    })

    return result, err
}

func updateAccounts(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (Account, Account, error) {
    account1, err := q.GetAccountForUpdate(ctx, accountID1)
    if err != nil {
        return Account{}, Account{}, err
    }

    account1.Balance += amount1
    _, err = q.UpdateAccount(ctx, UpdateAccountParams{
        ID:      account1.ID,
        Balance: account1.Balance,
    })
    if err != nil {
        return Account{}, Account{}, err
    }

    account2, err := q.GetAccountForUpdate(ctx, accountID2)
    if err != nil {
        return Account{}, Account{}, err
    }

    account2.Balance += amount2
    _, err = q.UpdateAccount(ctx, UpdateAccountParams{
        ID:      account2.ID,
        Balance: account2.Balance,
    })
    if err != nil {
        return Account{}, Account{}, err
    }

    return account1, account2, nil
}

// Closure : a function value that references variables from its surrounding lexical scope. This means that the function "closes over" these variables, capturing them so that they remain accessible even when the function is executed outside of its original context.

// the call back function acts as a closure and has access to the queries object and the transaction object. This allows us to call the queries functions within the transaction.

// https://www.youtube.com/watch?v=gBh__1eFwVI&list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE&index=7

// https://www.youtube.com/watch?v=G2aggv_3Bbg&list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE&index=9