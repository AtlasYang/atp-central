package delivery

import (
	_ "embed"
	"net/http"

	"aigendrug.com/router-core/internal/config"
	"github.com/gin-gonic/gin"
)

//go:embed html/index.html
var indexHTML string

type APIClientHandler struct {
	config *config.Config
}

func NewAPIClientHandler(config *config.Config) *APIClientHandler {
	return &APIClientHandler{config: config}
}

func (h *APIClientHandler) ServeIndexHTML(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html", []byte(indexHTML))
}
