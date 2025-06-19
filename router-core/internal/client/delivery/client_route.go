package delivery

import (
	authd "aigendrug.com/router-core/internal/auth/delivery"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupClientRoutes(
	router *gin.Engine,
	db *pgxpool.Pool,
	clientHandler *ClientHandler,
) {
	clientRoutes := router.Group("/v1/clients")
	{
		clientDefaultRoutes := clientRoutes.Group("", authd.DefaultAuthMiddleWare(db))
		{
			clientDefaultRoutes.GET("/identifier/:identifier", clientHandler.GetClientByClientIdentifier)
			clientDefaultRoutes.GET("/current", clientHandler.GetCurrentClient)
		}

		clientAdminRoutes := clientRoutes.Group("", authd.AdminAuthMiddleWare(db))
		{
			clientAdminRoutes.GET("", clientHandler.GetAllClients)
			clientAdminRoutes.GET("/:id", clientHandler.GetClientByID)
			clientAdminRoutes.POST("", clientHandler.CreateClient)
			clientAdminRoutes.PUT("/:id", clientHandler.UpdateClient)
			clientAdminRoutes.DELETE("/:id", clientHandler.DeleteClient)

			clientAdminRoutes.GET("/:id/api-keys", clientHandler.FindClientApiKeyByClientID)
			clientAdminRoutes.POST("/:id/api-keys", clientHandler.CreateClientApiKey)
			clientAdminRoutes.DELETE("/:id/api-keys/:api-key-id", clientHandler.DeleteClientApiKey)
		}
	}
}
