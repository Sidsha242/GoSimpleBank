package api

import (
	"database/sql"
	"net/http"

	db "github.com/Sidsha242/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

//This file defines the HTTP handlers for managing accounts
//Handlers are part of Server struct ;
//Handlers will use Gin to parse HTTP requests and send responses (give gin context as input)
//Handlers will interact with database through Store interface

//make type struct for input parameters then make function which will handle the request
//func will take server pointer and gin context as input

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`  
}
//similar to createAccountParams in db
//no need for balance as initial account will have balance 0
//json body will have owner and currency
//binding tag is used to *validate* the request body (gin will handle validation)


func (server *Server) createAccount(ctx *gin.Context) { //receiving server pointer ; and taking context as input [Gin context]
	var req createAccountRequest //request variable

	if err := ctx.ShouldBindJSON(&req);  err != nil { //this function will return error if request body is not in correct format
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	} //http code , json object(converting error to json)

	//Valid data input into object 
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)//input context , argument
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}//error in inserting into database

	//no error ..200 status code
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`  //uri tag is used to bind the id parameter from the uri and minimum value is 1 no negative ; parameter is required
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}//binduri will bind the id from the uri to the request object

	account, err := server.store.GetAccount(ctx, req.ID) //get account from database
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}//error getting from database


	ctx.JSON(http.StatusOK, account) //account found
}

//Pagination : Dont get all data from database ; divide into pages and client will ge data for only one page


//Query Parameters (get added after ? in url)
//page_id:index of page starting from 1 
//page_size:max number of data in a page 

//URI Parameters (get added in url ; like id in getAccountRequest /:id)

type listAccountRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`  //for query parameter use form for uri parameter use uri
	PageLimit int32 `form:"page_limit" binding:"required,min=5,max=10"`

}


func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}//binding function will tell gin to get the query parameters from the request and bind them to the request object

	arg:= db.ListAccountsParams{
		Limit: req.PageLimit,
		Offset: (req.PageID - 1) * req.PageLimit,
	}//offset is the number of rows to skip before starting to return data

	accounts, err := server.store.ListAccounts(ctx,arg ) //get account from database
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}//error getting from database


	ctx.JSON(http.StatusOK, accounts) 
}

//if you send request with page_id which does not exist reponse is null ..it should be an empty list
//to fix this emit_empty_slices: true in sqlc.yaml