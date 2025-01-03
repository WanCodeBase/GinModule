package api

import (
	"fmt"
	db "github.com/WanCodeBase/GinModule/db/sqlc"
	"github.com/WanCodeBase/GinModule/token"
	"github.com/WanCodeBase/GinModule/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("create token maker failed:%w", err)
	}
	server := &Server{
		tokenMaker: tokenMaker,
		config:     config,
		store:      store,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setRouter()

	return server, nil
}

func (server *Server) setRouter() {
	router := gin.Default()

	// user
	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)

	// add middleware
	authRouters := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouters.GET("/user/:username", server.getUser)

	// account
	authRouters.POST("/account", server.createAccount)
	authRouters.GET("/account/:id", server.getAccount)
	authRouters.GET("/accounts", server.listAccount)

	authRouters.POST("/transfer", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
