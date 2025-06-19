package delivery

import "github.com/gin-gonic/gin"

func SetupAPIClientRoutes(router *gin.Engine, apiClientHandler *APIClientHandler) {
	router.GET("/console", apiClientHandler.ServeIndexHTML)
}
