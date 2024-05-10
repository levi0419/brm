package api

import (
	db "github.com/brm/db/sqlc"
	"github.com/brm/token"
	"github.com/brm/utils"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our AltInvest API.
type Server struct {
	config           utils.Config
	router           *gin.Engine
	tokenMaker       token.ITokenMaker
	store            db.Store
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config *utils.Config, store db.Store) (*Server, error) {                                       
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey) // setup token Maker
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:           *config,
		tokenMaker:       tokenMaker,
		store:            store,
	}
	server.setUpRouter()

	return server, nil
}


func (server *Server) GetConfig() utils.Config {
	return server.config
}

func (server *Server) GetTokenMaker() token.ITokenMaker {
	return server.tokenMaker
}

// GetDbStore returns the Db store
func (server *Server) GetDbStore() db.Store {
	return server.store
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start() error {
	return server.router.Run(server.config.ServerAddress)
}
