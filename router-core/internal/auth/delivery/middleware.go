package delivery

import (
	"net/http"

	"aigendrug.com/router-core/internal/auth/application/service"
	"aigendrug.com/router-core/internal/auth/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func WebUserAuthorizationMiddleware(svc *service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, err := svc.AuthorizeAdmin(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		if isAdmin {
			c.Set("isAdmin", true)
			c.Set("clientID", 1)
			c.Next()
			return
		}

		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		clientID, err := svc.Authorize(c, apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("clientID", clientID)
		c.Next()
	}
}

func AdminAuthorizationMiddleware(svc *service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, err := svc.AuthorizeAdmin(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admins only"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func DefaultAuthMiddleWare(db *pgxpool.Pool) gin.HandlerFunc {
	sessionRepo := persistence.NewPgSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepo)

	return WebUserAuthorizationMiddleware(sessionService)
}

func AdminAuthMiddleWare(db *pgxpool.Pool) gin.HandlerFunc {
	sessionRepo := persistence.NewPgSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepo)

	return AdminAuthorizationMiddleware(sessionService)
}
