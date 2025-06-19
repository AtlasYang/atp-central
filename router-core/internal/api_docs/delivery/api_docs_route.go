package delivery

import (
	"github.com/gin-gonic/gin"
)

func SetupAPIDocsRoutes(router *gin.Engine, apiDocsHandler *APIDocsHandler) {
	router.GET("/docs", apiDocsHandler.ServeDocs)
	router.GET("/swagger.json", apiDocsHandler.ServeSwaggerJSON)
}
