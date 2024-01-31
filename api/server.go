package api

import (
	"github.com/gin-gonic/gin"
	db "simplebank/db/sqlc"
)

// Server serves HTTP requests for our service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer create a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := Server{store: store}
	router := gin.Default()

	router.POST("/accounts/create", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.PATCH("/accounts/update", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router
	return &server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func successResponse(message string) gin.H {
	return gin.H{"Success": message}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
