package api

import (
		"github.com/gin-gonic/gin"
		"github.com/brm/api/module/account"

)

func (server *Server) setUpRouter() {
	router := gin.Default()

	// // Call module routes
	// saf.SetUpRoutes(server, router)
	// test.SetUpRoutes(server, router)
	account.SetUpRoutes(server, router)

	server.router = router
}
