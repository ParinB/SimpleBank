package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/parin/simplebank/db/sqlc"
	"github.com/parin/simplebank/db/util"
	"github.com/parin/simplebank/token"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config util.Config
	//Interact with  the database while processing api requests from clients
	store  db.Store
	tokenMaker  token.Maker
	//Send Api Requests to the correct handler for processing
	router *gin.Engine
}
//  NewServer creates a new Http Server And Set Up Routing
func NewServer(config util.Config,store db.Store) (*Server,error) {
	tokenMaker , err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err !=nil {
		return nil, fmt.Errorf("cannot create token maker: %w",err)
	}

	server := &Server{
		store: store,
		tokenMaker : tokenMaker,
	}
	if v,ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency",validCurrency)
	}
	server.setupRouter()
	return server,nil
}
func (server *Server) setupRouter()  {
	router := gin.Default()
	// add routes to the router
	router.POST("/users",server.createUser)
	router.POST("/users/login",server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts",server.createAccount)
	authRoutes.GET("/accounts/:id",server.getAccount)
	authRoutes.GET("/account",server.listAccounts)
	authRoutes.POST("/transfers",server.createTransfer)
	server.router = router
}
// Start runs the Http server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(err error) gin.H  {
	//gin.H is shortcut for map[string]interface{}
	return gin.H{"error": err.Error()}
}