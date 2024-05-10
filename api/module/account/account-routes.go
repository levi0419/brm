package account

import (
	"github.com/brm/api/shared"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(server shared.IServer, router *gin.Engine) {

	accountRoutes := router.Group("/api/account")

	accountRoutes.POST("/sign-up", func(ctx *gin.Context) {
		createUser(ctx, server)
	})

	accountRoutes.POST("/login", func(ctx *gin.Context) {
		loginUser(ctx, server)
	})

	
}
