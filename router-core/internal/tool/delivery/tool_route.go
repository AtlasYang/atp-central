package delivery

import (
	authd "aigendrug.com/router-core/internal/auth/delivery"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupToolRoutes(
	router *gin.Engine,
	db *pgxpool.Pool,
	toolHandler *ToolHandler,
) {
	toolRoutes := router.Group("/v1/tools")
	{
		toolDefaultRoutes := toolRoutes.Group("", authd.DefaultAuthMiddleWare(db))
		{
			toolDefaultRoutes.GET("/:id", toolHandler.GetToolByID)
			toolDefaultRoutes.GET("/uuid/:uuid", toolHandler.GetToolByUUID)
			toolDefaultRoutes.GET("/client", toolHandler.GetAllToolsForClient)
			toolDefaultRoutes.POST("/select", toolHandler.SelectTool)
			toolDefaultRoutes.POST("/:tool_id/execute", toolHandler.ExecuteTool)
		}

		toolAdminRoutes := toolRoutes.Group("", authd.AdminAuthMiddleWare(db))
		{
			toolAdminRoutes.GET("", toolHandler.GetAllTools)
			toolAdminRoutes.GET("/client/:client_id", toolHandler.GetAllToolsByClientID)
			toolAdminRoutes.POST("", toolHandler.CreateTool)
			toolAdminRoutes.PUT("/:id", toolHandler.UpdateTool)
			toolAdminRoutes.DELETE("/:id", toolHandler.DeleteTool)
		}
	}

	// Tool Client Permission routes
	toolPermissionRoutes := router.Group("/v1/tool-permissions")
	{
		toolPermissionDefaultRoutes := toolPermissionRoutes.Group("", authd.DefaultAuthMiddleWare(db))
		{
			toolPermissionDefaultRoutes.GET("/client", toolHandler.GetAllToolClientPermissionsForClient)
		}

		toolPermissionAdminRoutes := toolPermissionRoutes.Group("", authd.AdminAuthMiddleWare(db))
		{
			toolPermissionAdminRoutes.GET("/tool/:tool_id", toolHandler.GetAllToolClientPermissionsByToolID)
			toolPermissionAdminRoutes.GET("/client/:client_id", toolHandler.GetAllToolClientPermissionsByClientID)
			toolPermissionAdminRoutes.POST("", toolHandler.CreateToolClientPermission)
			toolPermissionAdminRoutes.PUT("/:id", toolHandler.UpdateToolClientPermission)
			toolPermissionAdminRoutes.DELETE("/:id", toolHandler.DeleteToolClientPermission)
		}
	}

	// Tool Request routes
	toolRequestRoutes := router.Group("/v1/tool-requests")
	{
		toolRequestDefaultRoutes := toolRequestRoutes.Group("", authd.DefaultAuthMiddleWare(db))
		{
			toolRequestDefaultRoutes.GET("/client", toolHandler.GetAllToolRequestsForClient)
			toolRequestDefaultRoutes.GET("/:id", toolHandler.GetToolRequestByID)
		}

		toolRequestAdminRoutes := toolRequestRoutes.Group("", authd.AdminAuthMiddleWare(db))
		{
			toolRequestAdminRoutes.GET("/tool/:tool_id", toolHandler.GetAllToolRequestsByToolID)
			toolRequestAdminRoutes.GET("/client/:client_id", toolHandler.GetAllToolRequestsByClientID)
			toolRequestAdminRoutes.POST("", toolHandler.CreateToolRequest)
			toolRequestAdminRoutes.PUT("/:id", toolHandler.UpdateToolRequest)
			toolRequestAdminRoutes.DELETE("/:id", toolHandler.DeleteToolRequest)
		}
	}
}
