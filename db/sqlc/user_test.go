package db

import (
	"context"
	"testing"
	"time"

	"github.com/Sidsha242/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err:= util.HashPassword(util.RandomString(6)) 
    require.NoError(t, err)


	arg := CreateUserParams{
		Username: util.RandomOwner(),
		HashedPassword: hashedPassword,  
		FullName:  util.RandomOwner(),
		Email: util.RandomEmail(),
	} //struct to pass values to the function

	//random generated data will help to save time and avoid conflicts between unit tests

	user, err := testQueries.CreateUser(context.Background(), arg)  //testQueries is a global variable in main_test.go to pass queries for testing

	require.NoError(t, err)//require is a package from testify to check if the error is nil
	require.NotEmpty(t, user)

	//checking values received with actual values
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword) 
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}


//dont make unit test relly on each other
//unit test should be independent
func TestGetUser(t *testing.T) {

	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second) //check if the time difference is within 1 second
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}