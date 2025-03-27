package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Sidsha242/simple_bank/db/mock" //contains mock implementation of Store interface ; generated using GoMock
	db "github.com/Sidsha242/simple_bank/db/sqlc"
	"github.com/Sidsha242/simple_bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

//unit test for the GetAccountAPI (using GoMock)

func TestGetAccountAPI(t *testing.T) { //test take a testing object as an argument
	account:= randomAccount() //random account generated

	ctrl := gomock.NewController(t) // a controller is created to manage the lifecycle of the mock object
	store:= mockdb.NewMockStore(ctrl) //mock store created

	//build stubs , a stub is a piece of code that replaces the actual code of a method ; GetAccount is replaced by EXCEPT() of mock store. 

	//gomock.Any() ; matches any context passed to the method
	//Times(1) ; specifies the number of times the function will be called
	//Return() ; specifies the return value of the function ; return account object and nil error

	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil) //function will be called 1 time



	server := NewServer(store) //test server created
	recorder:=httptest.NewRecorder() //used to record the http response

	url := fmt.Sprintf("/accounts/%d", account.ID) //url for getaccount api constructed 
	request, _ := http.NewRequest(http.MethodGet, url, nil) //HTTP Get reqyurest is created
	server.router.ServeHTTP(recorder, request) //HTTP request is sent or served to test server using router ; gin router

	//check response
	require.Equal(t, http.StatusOK, recorder.Code)


}

//Helper function - create a new account

func randomAccount() db.Account {
	return db.Account{
		ID: 	 util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}