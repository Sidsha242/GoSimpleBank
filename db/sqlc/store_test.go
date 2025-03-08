package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

//tests the TransferTx method of the Store struct, ensuring that multiple concurrent transfer transactions are executed correctly and consistently.
func TestTrasnferTx(t *testing.T){
	testStore := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// We have to maintain concurrency

	//Run n concurrent transfer transactions(go routines)

	n:= 5 //number of concurrent transactions
	amount := int64(10)

	//go channel connects two concurrent GO routines and allow them to safely share data with each other without explicit locking

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i:=0; i<n; i++{
		//go routine
		//txName := fmt.Sprintf("tx %d", i+1) //name of the transaction - for debugging deadlock
		go func(){
			//transfer amount from account1 to account2


			//ctx := context.WithValue(context.Background(), txKey, txName) //Passing context with transaction name 

			ctx := context.Background()

			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			//we cannot use testify require to verify the results as the function is running in a different go routine to the main test function TestTransferTx
			//No guarantee that the test will be stopped if requirement not stopped
			//Correct way to verify is to send it back to the main go routine our test is running on ...using channel

			errs <- err    //channel on left ..data to be sent on right

			results <- result




		}()
	}

	//check the results
	existed := make(map[int]bool) //to check if the key already exists

	for i:=0; i<n; i++{
		err := <- errs     //get error from channel and store it in err 
		require.NoError(t, err)  //We are expecting no error

		result := <- results  //get result from channel and store it in result

		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//check account balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0) //balance will reduce by n * amount



		//ensures that the transfer transactions are executed the expected number of times and that each transfer count is unique
		k := int(diff1 / amount)   //k is the number of times the amount was transferred
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k) //check if the key already exists
		existed[k] = true




	}

	//check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance - int64(n) * amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n) * amount, updatedAccount2.Balance)


}

//will cause deadlock error as two transactions try to update the same row at the same time
//How to deal with this ?
//If you check the db the deadlock is happening as they are trying to acquire a shared lock 
//When we select an account for update it acquires a lock to prevent conflicts and ensure the consistency of the data
//The deadlock is happening exaclty because of foreign key constraint between account and transfer (account id is foreign key in transfer)
//But cannot remove the foreign key constraint as it is necessary to maintain the consistency of the data
//If we look at the UpdateAccount query only the balance is getting changed not the account id , there is no need to acquire the lock If we can somehow tell postgres that the primary key wont be touched ; postgres will not acquire the lock


//But if you change the order example over here account1 is transferred to account2 ..but simulataneously we also have account2 is transferred to account1 ..this will cause deadlock

//to avoid this alwyas update account with the smallest id first