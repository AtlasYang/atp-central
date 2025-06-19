package delivery

import (
	"net/http"
	"strconv"

	"aigendrug.com/router-core/internal/client/application/dto"
	"aigendrug.com/router-core/internal/client/application/service"
	shared_types "aigendrug.com/router-core/internal/shared/types"
	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	clientService service.ClientService
}

func NewClientHandler(clientService service.ClientService) *ClientHandler {
	return &ClientHandler{clientService: clientService}
}

var _ dto.ReadClientDTO // for swagger

// GetAllClients godoc
// @Summary Get all clients
// @Description Retrieves all clients
// @Tags client
// @Produce json
// @Success 200 {array} dto.ReadClientDTO
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/clients [get]
func (h *ClientHandler) GetAllClients(c *gin.Context) {
	clients, err := h.clientService.GetAllClients(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, clients)
}

// GetClientByID godoc
// @Summary Get client by ID
// @Description Retrieves a client by their ID
// @Tags client
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} dto.ReadClientDTO
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/clients/{id} [get]
func (h *ClientHandler) GetClientByID(c *gin.Context) {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid client ID"})
		return
	}

	client, err := h.clientService.GetClientByID(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	if client == nil {
		c.JSON(http.StatusNotFound, shared_types.HttpErrorResponse{Msg: "Client not found"})
		return
	}
	c.JSON(http.StatusOK, client)
}

// GetCurrentClient godoc
// @Summary Get current client
// @Description Retrieves the current client
// @Tags client
// @Produce json
// @Success 200 {object} dto.ReadClientDTO
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/clients/current [get]
func (h *ClientHandler) GetCurrentClient(c *gin.Context) {
	clientID := c.GetInt("clientID")
	client, err := h.clientService.GetClientByID(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, client)
}

// GetClientByClientIdentifier godoc
// @Summary Get client by identifier
// @Description Retrieves a client by their client identifier
// @Tags client
// @Produce json
// @Param identifier path string true "Client Identifier"
// @Success 200 {object} dto.ReadClientDTO
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/clients/identifier/{identifier} [get]
func (h *ClientHandler) GetClientByClientIdentifier(c *gin.Context) {
	identifier := c.Param("identifier")
	client, err := h.clientService.GetClientByClientIdentifier(c.Request.Context(), identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	if client == nil {
		c.JSON(http.StatusNotFound, shared_types.HttpErrorResponse{Msg: "Client not found"})
		return
	}
	c.JSON(http.StatusOK, client)
}

// CreateClient godoc
// @Summary Create a new client
// @Description Creates a new client
// @Tags client
// @Accept json
// @Produce json
// @Param client body dto.CreateClientDTO true "Client object"
// @Success 201 {object} dto.ReadClientDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/clients [post]
func (h *ClientHandler) CreateClient(c *gin.Context) {
	var client dto.CreateClientDTO
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	createdClient, err := h.clientService.CreateClient(c.Request.Context(), &client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdClient)
}

// UpdateClient godoc
// @Summary Update a client
// @Description Updates an existing client
// @Tags client
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param client body dto.UpdateClientDTO true "Client object"
// @Success 200 {object} dto.ReadClientDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/clients/{id} [put]
func (h *ClientHandler) UpdateClient(c *gin.Context) {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid client ID"})
		return
	}

	var client dto.UpdateClientDTO
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	if err := h.clientService.UpdateClient(c.Request.Context(), clientID, &dto.UpdateClientDTO{
		Name:             client.Name,
		ClientIdentifier: client.ClientIdentifier,
		IsActive:         client.IsActive,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, client)
}

// DeleteClient godoc
// @Summary Delete a client
// @Description Deletes a client by their ID
// @Tags client
// @Produce json
// @Param id path int true "Client ID"
// @Success 204 "No Content"
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/clients/{id} [delete]
func (h *ClientHandler) DeleteClient(c *gin.Context) {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid client ID"})
		return
	}

	if err := h.clientService.DeleteClient(c.Request.Context(), clientID); err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// FindClientApiKeyByClientID godoc
// @Summary Find client API key by client ID
// @Description Finds a client API key by their client ID
// @Tags client
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {array} dto.ReadClientApiKeyDTO
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/clients/{id}/api-keys [get]
func (h *ClientHandler) FindClientApiKeyByClientID(c *gin.Context) {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid client ID"})
		return
	}

	clientApiKeys, err := h.clientService.FindClientApiKeyByClientID(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, clientApiKeys)
}

// CreateClientApiKey godoc
// @Summary Create a new client API key
// @Description Creates a new client API key
// @Tags client
// @Produce json
// @Param id path int true "Client ID"
// @Success 201 {object} dto.ReadClientApiKeyDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/clients/{id}/api-keys [post]
func (h *ClientHandler) CreateClientApiKey(c *gin.Context) {
	id := c.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid client ID"})
		return
	}

	createdClientApiKey, err := h.clientService.CreateClientApiKey(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdClientApiKey)
}

// DeleteClientApiKey godoc
// @Summary Delete a client API key
// @Description Deletes a client API key by their ID
// @Tags client
// @Produce json
// @Param id path int true "Client API Key ID"
// @Success 204 "No Content"
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
func (h *ClientHandler) DeleteClientApiKey(c *gin.Context) {
	id := c.Param("api-key-id")
	clientApiKeyID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid client API key ID"})
		return
	}

	if err := h.clientService.DeleteClientApiKey(c.Request.Context(), clientApiKeyID); err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
