//http api server

package api

import (
	db "github.com/Sidsha242/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store db.Store //will allow us to interact with database 

	router *gin.Engine //will allow us to define routes and handlers
}


/* Why struct
- structs are used to group related data and methods together
- Server struct groups together fields that are essential for running the HTTP server: store and router
*/

//Will create a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store} //Intialize a Server instance with provided store
	router := gin.Default() //Create new Gin router

	//Define routes and handlers
	router.POST("/accounts", server.createAccount) // path ; middleware ; handler(over here only handler)

	//method of server struct we need access of store object 
	
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)
	router.GET("/transfers/:id", server.getTransfer)

	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)



	//assign router to the server instance we made earlier (server.router)
	server.router = router
	return server //return initialized server instance
}

//runs the http server on the given address
func (server *Server) Run(address string) error {
	return server.router.Run(address)
}

/*Why is the func structed this way?

Receiver: (server *Server)  ...since defined with this receiver Run method is now associated with the Server struct.
Input Parameter: (address string)

//In Go, methods are functions that are associated with a specific type (like a struct) through a receiver.

The receiver allows us to access the fields of the Server struct within the method. 

Why Pointer Receiver?
allows the method to modify the Server instance if needed (not required here). avoids copying the Server struct.

The (gin).Engine.Run method starts the HTTP server and listens for incoming requests on the specified address.


*/
//Note: router field is private and therefore cannot be accessed outside api package

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}//Converts an error into a JSON response. (Gin.H object: keyvalue pair)