package delivery

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"aigendrug.com/router-core/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
)

type APIDocsHandler struct {
	config *config.Config
}

func NewAPIDocsHandler(config *config.Config) *APIDocsHandler {
	return &APIDocsHandler{config: config}
}

func (c *APIDocsHandler) ServeSwaggerJSON(ctx *gin.Context) {
	swaggerJSON, err := os.ReadFile("_docs/swagger.json")
	if err != nil {
		log.Printf("Error reading swagger.json: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not load API specification"})
		return
	}

	switch c.config.Server.Env {
	case "local":
		swaggerJSON, err = sjson.SetBytes(swaggerJSON, "host", fmt.Sprintf("%s:%s", c.config.Server.Host, c.config.Server.ExternalPort))
		if err != nil {
			log.Printf("Error setting host: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not set host"})
			return
		}

		swaggerJSON, err = sjson.SetBytes(swaggerJSON, "schemes", []string{"http"})
		if err != nil {
			log.Printf("Error setting schemes: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not set schemes"})
			return
		}

	case "development":
		swaggerJSON, err = sjson.SetBytes(swaggerJSON, "host", c.config.Server.Host)
		if err != nil {
			log.Printf("Error setting host: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not set host"})
			return
		}

		swaggerJSON, err = sjson.SetBytes(swaggerJSON, "schemes", []string{"https"})
		if err != nil {
			log.Printf("Error setting schemes: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not set schemes"})
			return
		}
	}

	ctx.Data(http.StatusOK, "application/json", swaggerJSON)
}

func (c *APIDocsHandler) ServeDocs(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="referrer" content="same-origin" />
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
    <title>[AIGENDRUG CID 2] atp-central router core API Documentation</title>
    <link href="https://unpkg.com/@stoplight/elements@8.0.0/styles.min.css" rel="stylesheet" />
    <script src="https://unpkg.com/@stoplight/elements@8.0.0/web-components.min.js"
            integrity="sha256-yIhuSFMJJ6mp2XTUAb4SiSYneP3Qav8Uu+7NBhGJW5A="
            crossorigin="anonymous"></script>
  </head>
  <body style="height: 100vh; margin: 0;">
    <elements-api
      apiDescriptionUrl="/swagger.json"
      router="hash"
      hideSchemas="true"
      tryItCredentialsPolicy="same-origin"
    />
  </body>
</html>`))
}
