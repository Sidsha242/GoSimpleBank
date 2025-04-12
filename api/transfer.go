package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Sidsha242/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,min=1"`
	Currency 	string `json:"currency" binding:"required,oneof=USD EUR"`   //can also create a custom validator for this ..but wont do it now
}

func (server *Server) createTransfer(ctx *gin.Context) { //receiving server pointer ; and taking context as input [Gin context]
	var req transferRequest 			//request variable

	if err := ctx.ShouldBindJSON(&req); err != nil { //this function will return error if request body is not in correct format
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	} 

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) { //check if from account is valid
			return 
	}
	if !server.validAccount(ctx, req.ToAccountID, req.Currency) { //check if to account is valid
			return
	}


	//Valid data input into object
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount: 	  req.Amount,

	}

	account, err := server.store.TransferTx(ctx, arg) //calling db Querier function // input context , argument
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} //error in inserting into database

	//no error ..200 status code
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId) //get account from database
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
	
	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	return false
	}

	//check if account currency matches the request currency
	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch: %s vs %s", account.ID, account.Currency, currency)

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true //all good

}

type getTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`  //uri tag is used to bind the id parameter from the uri and minimum value is 1 no negative ; parameter is required
}

func (server *Server) getTransfer(ctx *gin.Context) {
	var req getTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}//binduri will bind the id from the uri to the request object

	transfer, err := server.store.GetTransfer(ctx, req.ID) //get transfer from database
	if err != nil {
		if err == sql.ErrNoRows { //check if transfer not found
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)

}