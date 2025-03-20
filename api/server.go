//http api server

package api

import (
	db "github.com/Sidsha242/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store *db.Store //will allow us to interact with database 

	router *gin.Engine //will allow us to define routes and handlers
}

//Will create a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//Define routes and handlers
	router.POST("/accounts", server.createAccount) // path ; middleware ; handler(over here only handler)

	//method of server struct we need access of store object 
	
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	// router.POST("/transfers", server.createTransfer)

	//set router object to server.router
	server.router = router
	return server
}

//runs the http server on the given address
func (server *Server) Run(address string) error {
	return server.router.Run(address)
}
//Note-router field is private and therefore cannot be accessed outside api package

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}//takes error as input and returns a map (Gin.H object: keyvalue pair)