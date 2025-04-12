package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Sidsha242/simple_bank/util"
	"github.com/stretchr/testify/require"
)

//unit tests for CRUD operations on account table

//createRandomAccount is a helper function to create a random account for testing ..will not be run during unit test as it does not begin with Test
//t is a testing.T object to log errors and failures
func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t) //creating random user for testing
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	} //struct to pass values to the function

	//random generated data will help to save time and avoid conflicts between unit tests

	account, err := testQueries.CreateAccount(context.Background(), arg)  //testQueries is a global variable in main_test.go to pass queries for testing

	require.NoError(t, err)//require is a package from testify to check if the error is nil
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance) //checking values received with actual values
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}


//dont make unit test relly on each other
//unit test should be independent
func TestGetAccount(t *testing.T) {
	//creating random account for testing
	createAccount := createRandomAccount(t)

	//getting the account from the database using the id of the created account...we will compare both of these account values
	account, err := testQueries.GetAccount(context.Background(), createAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, createAccount.ID, account.ID)
	require.Equal(t, createAccount.Owner, account.Owner)
	require.Equal(t, createAccount.Balance, account.Balance)
	require.Equal(t, createAccount.Currency, account.Currency)
	require.WithinDuration(t, createAccount.CreatedAt, account.CreatedAt, time.Second)//checking if the time is within a second delta is one second
}

func TestUpdateAccount(t *testing.T) {
	createAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      createAccount.ID,
		Balance: util.RandomMoney(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, createAccount.ID, account.ID)
	require.Equal(t, createAccount.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)//only this will change as balanced has been changed
	require.Equal(t, createAccount.Currency, account.Currency)
	require.WithinDuration(t, createAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	createAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), createAccount.ID)

	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), createAccount.ID)
     //check if actually deleted by getting the account

	require.Error(t, err)
	require.EqualError(t, err,sql.ErrNoRows.Error() )//checking if the error is of type sql.ErrNoRows
	require.Empty(t, account)
}

func TestListAccounts(t *testing.T) {
	//create 5 accounts
	for i := 0; i < 5; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams {
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)//checking if the length of the accounts is 5

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}//checking if the accounts are not empty
}