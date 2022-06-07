package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/parin/simplebank/db/sqlc"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	//Interact with  the database while processing api requests from clients
	store  db.Store
	//Send Api Requests to the correct handler for processing
	router *gin.Engine
}
//  NewServer creates a new Http Server And Set Up Routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v,ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency",validCurrency)
	}
	// add routes to the router
	router.POST("/users",server.createUser)
	router.POST("/accounts",server.createAccount)
	router.GET("/accounts/:id",server.getAccount)
	router.GET("/account",server.listAccounts)
	router.POST("/transfers",server.createTransfer)
	server.router = router
	return server
}
// Start runs the Http server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(err error) gin.H  {
	//gin.H is shortcut for map[string]interface{}
	return gin.H{"error": err.Error()}
}