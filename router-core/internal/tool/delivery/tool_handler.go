package delivery

import (
	"net/http"
	"strconv"

	shared_types "aigendrug.com/router-core/internal/shared/types"
	"aigendrug.com/router-core/internal/tool/application/dto"
	"aigendrug.com/router-core/internal/tool/application/service"
	"aigendrug.com/router-core/internal/tool/domain/valueobject"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ToolHandler struct {
	toolService service.ToolService
}

func NewToolHandler(toolService service.ToolService) *ToolHandler {
	return &ToolHandler{toolService: toolService}
}

// GetAllTools godoc
// @Summary Get all tools
// @Description Retrieves all tools
// @Tags tool
// @Produce json
// @Success 200 {array} dto.ReadToolDTO
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tools [get]
func (h *ToolHandler) GetAllTools(c *gin.Context) {
	tools, err := h.toolService.GetAllTools(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, tools)
}

// GetToolByID godoc
// @Summary Get tool by ID
// @Description Retrieves a tool by its ID
// @Tags tool
// @Produce json
// @Param id path int true "Tool ID"
// @Success 200 {object} dto.ReadToolDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tools/{id} [get]
func (h *ToolHandler) GetToolByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid tool ID"})
		return
	}

	// if not admin, check if the tool is in the client's permission level
	if !c.GetBool("isAdmin") {
		permissionLevel, err := h.toolService.GetToolClientPermissionByToolIDAndClientID(c.Request.Context(), id, c.GetInt("clientID"))
		if err != nil {
			c.JSON(http.StatusForbidden, shared_types.HttpErrorResponse{Msg: "You are not authorized to access this tool"})
			return
		}
		if permissionLevel.PermissionLevel == valueobject.ToolClientPermissionLevelNone {
			c.JSON(http.StatusForbidden, shared_types.HttpErrorResponse{Msg: "You are not authorized to access this tool"})
			return
		}
	}

	tool, err := h.toolService.GetToolByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, tool)
}

// GetToolByUUID godoc
// @Summary Get tool by UUID
// @Description Retrieves a tool by its UUID
// @Tags tool
// @Produce json
// @Param uuid path string true "Tool UUID"
// @Success 200 {object} dto.ReadToolDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tools/uuid/{uuid} [get]
func (h *ToolHandler) GetToolByUUID(c *gin.Context) {
	uuidStr := c.Param("uuid")
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid UUID format"})
		return
	}

	tool, err := h.toolService.GetToolByUUID(c.Request.Context(), uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, tool)
}

// GetAllToolsByClientID godoc
// @Summary Get all tools by client ID
// @Description Retrieves all tools for a specific client
// @Tags tool
// @Produce json
// @Param client_id path int true "Client ID"
// @Param permission_level query int true "Permission Level"
// @Success 200 {array} dto.ReadToolDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tools/client/{client_id} [get]
func (h *ToolHandler) GetAllToolsByClientID(c *gin.Context) {
	clientID, err := strconv.Atoi(c.Param("client_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid client ID"})
		return
	}

	permissionLevel, err := strconv.Atoi(c.Query("permission_level"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid permission level"})
		return
	}

	tools, err := h.toolService.GetAllToolsByClientID(
		c.Request.Context(), clientID, valueobject.ToolClientPermissionLevel(permissionLevel),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, tools)
}

// GetAllToolsForClient godoc
// @Summary Get all tools for a client
// @Description Retrieves all tools for a specific client
// @Tags tool
// @Produce json
// @Param permission_level query int true "Permission Level"
// @Success 200 {array} dto.ReadToolDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tools/client [get]
func (h *ToolHandler) GetAllToolsForClient(c *gin.Context) {
	clientID := c.GetInt("clientID")

	permissionLevel, err := strconv.Atoi(c.Query("permission_level"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid permission level"})
		return
	}

	tools, err := h.toolService.GetAllToolsByClientID(
		c.Request.Context(), clientID, valueobject.ToolClientPermissionLevel(permissionLevel),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, tools)
}

// CreateTool godoc
// @Summary Create a new tool
// @Description Creates a new tool
// @Tags tool
// @Accept json
// @Produce json
// @Param tool body dto.CreateToolDTO true "Tool to create"
// @Success 201 {object} dto.ReadToolDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tools [post]
func (h *ToolHandler) CreateTool(c *gin.Context) {
	var tool dto.CreateToolDTO
	if err := c.ShouldBindJSON(&tool); err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	createdTool, err := h.toolService.CreateTool(c.Request.Context(), &tool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdTool)
}

// UpdateTool godoc
// @Summary Update a tool
// @Description Updates an existing tool
// @Tags tool
// @Accept json
// @Produce json
// @Param id path int true "Tool ID"
// @Param tool body dto.UpdateToolDTO true "Tool to update"
// @Success 200 {object} shared_types.HttpSuccessResponse
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tools/{id} [put]
func (h *ToolHandler) UpdateTool(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid tool ID"})
		return
	}

	var tool dto.UpdateToolDTO
	if err := c.ShouldBindJSON(&tool); err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	if err := h.toolService.UpdateTool(c.Request.Context(), id, &tool); err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, shared_types.HttpSuccessResponse{Msg: "Tool updated successfully"})
}

// DeleteTool godoc
// @Summary Delete a tool
// @Description Deletes an existing tool
// @Tags tool
// @Produce json
// @Param id path int true "Tool ID"
// @Success 200 {object} shared_types.HttpSuccessResponse
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tools/{id} [delete]
func (h *ToolHandler) DeleteTool(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid tool ID"})
		return
	}

	if err := h.toolService.DeleteTool(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, shared_types.HttpSuccessResponse{Msg: "Tool deleted successfully"})
}

// GetAllToolClientPermissionsByToolID godoc
// @Summary Get all tool client permissions by tool ID
// @Description Retrieves all tool client permissions for a specific tool
// @Tags tool-client-permission
// @Produce json
// @Param tool_id path int true "Tool ID"
// @Success 200 {array} dto.ReadToolClientPermissionDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-permissions/tool/{tool_id} [get]
func (h *ToolHandler) GetAllToolClientPermissionsByToolID(c *gin.Context) {
	toolID, err := strconv.Atoi(c.Param("tool_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid tool ID"})
		return
	}

	permissions, err := h.toolService.GetAllToolClientPermissionsByToolID(c.Request.Context(), toolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

// GetAllToolClientPermissionsByClientID godoc
// @Summary Get all tool client permissions by client ID
// @Description Retrieves all tool client permissions for a specific client
// @Tags tool-client-permission
// @Produce json
// @Param client_id path int true "Client ID"
// @Success 200 {array} dto.ReadToolClientPermissionDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-permissions/client/{client_id} [get]
func (h *ToolHandler) GetAllToolClientPermissionsByClientID(c *gin.Context) {
	clientID, err := strconv.Atoi(c.Param("client_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid client ID"})
		return
	}

	permissions, err := h.toolService.GetAllToolClientPermissionsByClientID(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

// GetAllToolClientPermissionsForClient godoc
// @Summary Get all tool client permissions for a client
// @Description Retrieves all tool client permissions for a specific client
// @Tags tool-client-permission
// @Produce json
// @Success 200 {array} dto.ReadToolClientPermissionDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-permissions/client [get]
func (h *ToolHandler) GetAllToolClientPermissionsForClient(c *gin.Context) {
	clientID := c.GetInt("clientID")

	permissions, err := h.toolService.GetAllToolClientPermissionsByClientID(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

// CreateToolClientPermission godoc
// @Summary Create a new tool client permission
// @Description Creates a new tool client permission
// @Tags tool-client-permission
// @Accept json
// @Produce json
// @Param permission body dto.CreateToolClientPermissionDTO true "Permission to create"
// @Success 201 {object} dto.ReadToolClientPermissionDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-permissions [post]
func (h *ToolHandler) CreateToolClientPermission(c *gin.Context) {
	var permission dto.CreateToolClientPermissionDTO
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	createdPermission, err := h.toolService.CreateToolClientPermission(c.Request.Context(), &permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdPermission)
}

// UpdateToolClientPermission godoc
// @Summary Update a tool client permission
// @Description Updates an existing tool client permission
// @Tags tool-client-permission
// @Accept json
// @Produce json
// @Param id path int true "Permission ID"
// @Param permission body dto.UpdateToolClientPermissionDTO true "Permission to update"
// @Success 200 {object} shared_types.HttpSuccessResponse
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-permissions/{id} [put]
func (h *ToolHandler) UpdateToolClientPermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid permission ID"})
		return
	}

	var permission dto.UpdateToolClientPermissionDTO
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	if err := h.toolService.UpdateToolClientPermission(c.Request.Context(), id, &permission); err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, shared_types.HttpSuccessResponse{Msg: "Permission updated successfully"})
}

// DeleteToolClientPermission godoc
// @Summary Delete a tool client permission
// @Description Deletes an existing tool client permission
// @Tags tool-client-permission
// @Produce json
// @Param id path int true "Permission ID"
// @Success 200 {object} shared_types.HttpSuccessResponse
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-permissions/{id} [delete]
func (h *ToolHandler) DeleteToolClientPermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid permission ID"})
		return
	}

	if err := h.toolService.DeleteToolClientPermission(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, shared_types.HttpSuccessResponse{Msg: "Permission deleted successfully"})
}

// GetAllToolRequestsByToolID godoc
// @Summary Get all tool requests by tool ID
// @Description Retrieves all tool requests for a specific tool
// @Tags tool-request
// @Produce json
// @Param tool_id path int true "Tool ID"
// @Success 200 {array} dto.ReadToolRequestDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-requests/tool/{tool_id} [get]
func (h *ToolHandler) GetAllToolRequestsByToolID(c *gin.Context) {
	toolID, err := strconv.Atoi(c.Param("tool_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid tool ID"})
		return
	}

	requests, err := h.toolService.GetAllToolRequestsByToolID(c.Request.Context(), toolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// GetAllToolRequestsByClientID godoc
// @Summary Get all tool requests by client ID
// @Description Retrieves all tool requests for a specific client
// @Tags tool-request
// @Produce json
// @Param client_id path int true "Client ID"
// @Success 200 {array} dto.ReadToolRequestDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-requests/client/{client_id} [get]
func (h *ToolHandler) GetAllToolRequestsByClientID(c *gin.Context) {
	clientID, err := strconv.Atoi(c.Param("client_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid client ID"})
		return
	}

	requests, err := h.toolService.GetAllToolRequestsByClientID(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// GetAllToolRequestsForClient godoc
// @Summary Get all tool requests for a client
// @Description Retrieves all tool requests for a specific client
// @Tags tool-request
// @Produce json
// @Success 200 {array} dto.ReadToolRequestDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-requests/client [get]
func (h *ToolHandler) GetAllToolRequestsForClient(c *gin.Context) {
	clientID := c.GetInt("clientID")

	requests, err := h.toolService.GetAllToolRequestsByClientID(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// GetToolRequestByID godoc
// @Summary Get a tool request by ID
// @Description Retrieves a tool request by its ID
// @Tags tool-request
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} dto.ReadToolRequestDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-requests/{id} [get]
func (h *ToolHandler) GetToolRequestByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid request ID"})
		return
	}

	request, err := h.toolService.GetToolRequestByID(c.Request.Context(), c.GetInt("clientID"), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, request)
}

// CreateToolRequest godoc
// @Summary Create a new tool request
// @Description Creates a new tool request
// @Tags tool-request
// @Accept json
// @Produce json
// @Param request body dto.CreateToolRequestDTO true "Request to create"
// @Success 201 {object} dto.ReadToolRequestDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-requests [post]
func (h *ToolHandler) CreateToolRequest(c *gin.Context) {
	var request dto.CreateToolRequestDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	createdRequest, err := h.toolService.CreateToolRequest(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdRequest)
}

// UpdateToolRequest godoc
// @Summary Update a tool request
// @Description Updates an existing tool request
// @Tags tool-request
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param request body dto.UpdateToolRequestDTO true "Request to update"
// @Success 200 {object} shared_types.HttpSuccessResponse
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-requests/{id} [put]
func (h *ToolHandler) UpdateToolRequest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid request ID"})
		return
	}

	var request dto.UpdateToolRequestDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	if err := h.toolService.UpdateToolRequest(c.Request.Context(), id, &request); err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, shared_types.HttpSuccessResponse{Msg: "Request updated successfully"})
}

// DeleteToolRequest godoc
// @Summary Delete a tool request
// @Description Deletes an existing tool request
// @Tags tool-request
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} shared_types.HttpSuccessResponse
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 404 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tool-requests/{id} [delete]
func (h *ToolHandler) DeleteToolRequest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid request ID"})
		return
	}

	if err := h.toolService.DeleteToolRequest(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, shared_types.HttpSuccessResponse{Msg: "Request deleted successfully"})
}

// SelectTool godoc
// @Summary Select a tool
// @Description Selects a tool based on user prompt
// @Tags tool
// @Accept json
// @Produce json
// @Param prompt body dto.SelectToolRequestDTO true "User prompt"
// @Success 200 {object} dto.SelectToolResponseDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tools/select [post]
func (h *ToolHandler) SelectTool(c *gin.Context) {
	var request dto.SelectToolRequestDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	response, err := h.toolService.SelectTool(c.Request.Context(), c.GetInt("clientID"), request.UserPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// ExecuteTool godoc
// @Summary Execute a tool
// @Description Executes a tool based on user prompt
// @Tags tool
// @Accept json
// @Produce json
// @Param tool_id path int true "Tool ID"
// @Param request body dto.ToolExecutionRequestDTO true "Request to execute"
// @Success 200 {object} dto.ToolExecutionResponseDTO
// @Failure 400 {object} shared_types.HttpErrorResponse
// @Failure 500 {object} shared_types.HttpErrorResponse
// @Router /v1/tools/{tool_id}/execute [post]
func (h *ToolHandler) ExecuteTool(c *gin.Context) {
	toolID, err := strconv.Atoi(c.Param("tool_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: "Invalid tool ID"})
		return
	}

	var request dto.ToolExecutionRequestDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}

	response, err := h.toolService.ExecuteTool(c.Request.Context(), c.GetInt("clientID"), toolID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared_types.HttpErrorResponse{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
