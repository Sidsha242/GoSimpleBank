package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/Sidsha242/simple_bank/db/sqlc"
	"github.com/Sidsha242/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"` //alphanumeric username
	Password string `json:"password" binding:"required,min=6"` //password should be at least 6 characters long
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"` //email tag provided by validator
}

//we dont want to expose the password in the response so we will not include it in the response struct
type createUserResponse struct {
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt     time.Time `json:"created_at"`
}




func (server *Server) createUser(ctx *gin.Context) { 
	var req createUserRequest //request variable

	if err := ctx.ShouldBindJSON(&req);  err != nil { 
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return 
	}

	hashedPassword, err := util.HashPassword(req.Password) //hashing password using bcrypt
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//Valid data input into object 
	arg := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
		
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		pqErr, _ := err.(*pq.Error) //check if error is of type pqError (Postgres error)
		{
			switch pqErr.Code.Name() {
			case "unique_violation":   //username and email unique
				ctx.JSON(http.StatusForbidden, errorResponse(err)) 
				return
			}
				
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getUser(ctx *gin.Context) {
	username := ctx.Param("username") //get username from url parameter

	user, err := server.store.GetUser(ctx, username) //get user from database
	if err != nil {
		if err == sql.ErrNoRows { //check if user not found
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
