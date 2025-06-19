package service

import (
	"context"
	"fmt"

	lambda_wrapper "aigendrug.com/router-core/internal/shared/lambda-wrapper"
	"aigendrug.com/router-core/internal/shared/selector"
	"aigendrug.com/router-core/internal/tool/application/dto"
	"aigendrug.com/router-core/internal/tool/domain"
	"aigendrug.com/router-core/internal/tool/domain/entity"
	"aigendrug.com/router-core/internal/tool/domain/shared_type"
	"aigendrug.com/router-core/internal/tool/domain/valueobject"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ToolService interface {
	// Tool
	GetAllTools(ctx context.Context) ([]*dto.ReadToolDTO, error)
	GetToolByID(ctx context.Context, id int) (*dto.ReadToolDTO, error)
	GetToolByUUID(ctx context.Context, uuid uuid.UUID) (*dto.ReadToolDTO, error)
	GetAllToolsByClientID(ctx context.Context, clientID int, permissionLevel valueobject.ToolClientPermissionLevel) ([]*dto.ReadToolDTO, error)
	CreateTool(ctx context.Context, tool *dto.CreateToolDTO) (*dto.ReadToolDTO, error)
	UpdateTool(ctx context.Context, id int, tool *dto.UpdateToolDTO) error
	DeleteTool(ctx context.Context, id int) error

	// ToolClientPermission
	GetAllToolClientPermissionsByToolID(ctx context.Context, toolID int) ([]*dto.ReadToolClientPermissionDTO, error)
	GetAllToolClientPermissionsByClientID(ctx context.Context, clientID int) ([]*dto.ReadToolClientPermissionDTO, error)
	GetToolClientPermissionByToolIDAndClientID(ctx context.Context, toolID int, clientID int) (*dto.ReadToolClientPermissionDTO, error)
	CreateToolClientPermission(ctx context.Context, toolClientPermission *dto.CreateToolClientPermissionDTO) (*dto.ReadToolClientPermissionDTO, error)
	UpdateToolClientPermission(ctx context.Context, id int, toolClientPermission *dto.UpdateToolClientPermissionDTO) error
	DeleteToolClientPermission(ctx context.Context, id int) error

	// ToolRequest
	GetAllToolRequestsByToolID(ctx context.Context, toolID int) ([]*dto.ReadToolRequestDTO, error)
	GetAllToolRequestsByClientID(ctx context.Context, clientID int) ([]*dto.ReadToolRequestDTO, error)
	GetToolRequestByID(ctx context.Context, clientID int, id int) (*dto.ReadToolRequestDTO, error)
	CreateToolRequest(ctx context.Context, toolRequest *dto.CreateToolRequestDTO) (*dto.ReadToolRequestDTO, error)
	UpdateToolRequest(ctx context.Context, id int, toolRequest *dto.UpdateToolRequestDTO) error
	DeleteToolRequest(ctx context.Context, id int) error

	// Selector
	SelectTool(ctx context.Context, clientID int, userPrompt string) (*dto.SelectToolResponseDTO, error)

	// Tool Execution
	ExecuteTool(ctx context.Context, clientID int, toolID int, requestData dto.ToolExecutionRequestDTO) (*dto.ToolExecutionResponseDTO, error)
}

type toolService struct {
	db               *pgxpool.Pool
	toolRepo         domain.ToolRepository
	selectorService  selector.SelectorService
	functionExecutor FunctionExecutor
}

func NewToolService(
	dbPool *pgxpool.Pool,
	toolRepo domain.ToolRepository,
	selectorService selector.SelectorService,
	lambdaClient lambda_wrapper.LambdaWrapperClient,
) ToolService {
	functionExecutor := NewFunctionExecutor(context.Background(), toolRepo, lambdaClient)

	return &toolService{
		db:               dbPool,
		toolRepo:         toolRepo,
		selectorService:  selectorService,
		functionExecutor: functionExecutor,
	}
}

func (s *toolService) GetAllTools(ctx context.Context) ([]*dto.ReadToolDTO, error) {
	tools, err := s.toolRepo.FindAllTools(ctx)
	if err != nil {
		return nil, err
	}

	toolsDTO := make([]*dto.ReadToolDTO, len(tools))
	for i, tool := range tools {
		toolsDTO[i] = tool.ToDTO()
	}
	return toolsDTO, nil
}

func (s *toolService) GetToolByID(ctx context.Context, id int) (*dto.ReadToolDTO, error) {
	tool, err := s.toolRepo.FindToolByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return tool.ToDTO(), nil
}

func (s *toolService) GetToolByUUID(ctx context.Context, uuid uuid.UUID) (*dto.ReadToolDTO, error) {
	tool, err := s.toolRepo.FindToolByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return tool.ToDTO(), nil
}

func (s *toolService) GetAllToolsByClientID(
	ctx context.Context, clientID int,
	permissionLevel valueobject.ToolClientPermissionLevel,
) ([]*dto.ReadToolDTO, error) {
	tools, err := s.toolRepo.FindAllToolsByClientID(ctx, clientID, permissionLevel)
	if err != nil {
		return nil, err
	}

	toolsDTO := make([]*dto.ReadToolDTO, len(tools))
	for i, tool := range tools {
		toolsDTO[i] = tool.ToDTO()
	}
	return toolsDTO, nil
}

func (s *toolService) CreateTool(
	ctx context.Context, tool *dto.CreateToolDTO,
) (*dto.ReadToolDTO, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	toolEntity := &entity.Tool{
		UUID:              newUUID,
		Name:              tool.Name,
		Version:           tool.Version,
		Description:       tool.Description,
		EngineInterface:   tool.EngineInterface,
		ProviderInterface: tool.ProviderInterface,
	}

	createdTool, err := s.toolRepo.CreateTool(ctx, toolEntity)
	if err != nil {
		return nil, err
	}
	return createdTool.ToDTO(), nil
}

func (s *toolService) UpdateTool(ctx context.Context, id int, tool *dto.UpdateToolDTO) error {
	toolEntity := &entity.Tool{
		ID:                id,
		Description:       tool.Description,
		EngineInterface:   tool.EngineInterface,
		ProviderInterface: tool.ProviderInterface,
	}

	return s.toolRepo.UpdateTool(ctx, toolEntity)
}

func (s *toolService) DeleteTool(ctx context.Context, id int) error {
	return s.toolRepo.DeleteTool(ctx, id)
}

func (s *toolService) GetAllToolClientPermissionsByToolID(
	ctx context.Context, toolID int,
) ([]*dto.ReadToolClientPermissionDTO, error) {
	toolClientPermissions, err := s.toolRepo.FindAllToolClientPermissionsByToolID(ctx, toolID)
	if err != nil {
		return nil, err
	}

	toolClientPermissionsDTO := make([]*dto.ReadToolClientPermissionDTO, len(toolClientPermissions))
	for i, toolClientPermission := range toolClientPermissions {
		toolClientPermissionsDTO[i] = toolClientPermission.ToDTO()
	}
	return toolClientPermissionsDTO, nil
}

func (s *toolService) GetAllToolClientPermissionsByClientID(
	ctx context.Context, clientID int,
) ([]*dto.ReadToolClientPermissionDTO, error) {
	toolClientPermissions, err := s.toolRepo.FindAllToolClientPermissionsByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	toolClientPermissionsDTO := make([]*dto.ReadToolClientPermissionDTO, len(toolClientPermissions))
	for i, toolClientPermission := range toolClientPermissions {
		toolClientPermissionsDTO[i] = toolClientPermission.ToDTO()
	}
	return toolClientPermissionsDTO, nil
}

func (s *toolService) GetToolClientPermissionByToolIDAndClientID(
	ctx context.Context, toolID int, clientID int,
) (*dto.ReadToolClientPermissionDTO, error) {
	toolClientPermission, err := s.toolRepo.GetToolClientPermissionByToolIDAndClientID(ctx, toolID, clientID)
	if err != nil {
		return nil, err
	}
	return toolClientPermission.ToDTO(), nil
}

func (s *toolService) CreateToolClientPermission(
	ctx context.Context, toolClientPermission *dto.CreateToolClientPermissionDTO,
) (*dto.ReadToolClientPermissionDTO, error) {
	toolClientPermissionEntity := &entity.ToolClientPermission{
		ToolID:          toolClientPermission.ToolID,
		ClientID:        toolClientPermission.ClientID,
		PermissionLevel: toolClientPermission.PermissionLevel,
	}

	createdToolClientPermission, err := s.toolRepo.CreateToolClientPermission(ctx, toolClientPermissionEntity)
	if err != nil {
		return nil, err
	}
	return createdToolClientPermission.ToDTO(), nil
}

func (s *toolService) UpdateToolClientPermission(
	ctx context.Context, id int, toolClientPermission *dto.UpdateToolClientPermissionDTO,
) error {
	toolClientPermissionEntity := &entity.ToolClientPermission{
		ID:              id,
		PermissionLevel: toolClientPermission.PermissionLevel,
	}

	return s.toolRepo.UpdateToolClientPermission(ctx, toolClientPermissionEntity)
}

func (s *toolService) DeleteToolClientPermission(ctx context.Context, id int) error {
	return s.toolRepo.DeleteToolClientPermission(ctx, id)
}

func (s *toolService) GetAllToolRequestsByToolID(
	ctx context.Context, toolID int,
) ([]*dto.ReadToolRequestDTO, error) {
	toolRequests, err := s.toolRepo.FindAllToolRequestsByToolID(ctx, toolID)
	if err != nil {
		return nil, err
	}

	toolRequestsDTO := make([]*dto.ReadToolRequestDTO, len(toolRequests))
	for i, toolRequest := range toolRequests {
		toolRequestsDTO[i] = toolRequest.ToDTO()
	}
	return toolRequestsDTO, nil
}

func (s *toolService) GetAllToolRequestsByClientID(
	ctx context.Context, clientID int,
) ([]*dto.ReadToolRequestDTO, error) {
	toolRequests, err := s.toolRepo.FindAllToolRequestsByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	toolRequestsDTO := make([]*dto.ReadToolRequestDTO, len(toolRequests))
	for i, toolRequest := range toolRequests {
		toolRequestsDTO[i] = toolRequest.ToDTO()
	}
	return toolRequestsDTO, nil
}

func (s *toolService) GetToolRequestByID(
	ctx context.Context, clientID int, id int,
) (*dto.ReadToolRequestDTO, error) {
	toolRequest, err := s.toolRepo.FindToolRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}

	toolClientPermission, err := s.toolRepo.GetToolClientPermissionByToolIDAndClientID(
		ctx, toolRequest.ToolID, clientID,
	)
	if err != nil {
		return nil, err
	}

	if toolClientPermission.PermissionLevel == valueobject.ToolClientPermissionLevelNone {
		return nil, fmt.Errorf("you don't have permission to access this tool request")
	}

	return toolRequest.ToDTO(), nil
}

func (s *toolService) CreateToolRequest(
	ctx context.Context, toolRequest *dto.CreateToolRequestDTO,
) (*dto.ReadToolRequestDTO, error) {
	toolRequestEntity := &entity.ToolRequest{
		ToolID:       toolRequest.ToolID,
		ClientID:     toolRequest.ClientID,
		RequestData:  toolRequest.RequestData,
		ResponseData: toolRequest.ResponseData,
		Status:       toolRequest.Status,
	}

	createdToolRequest, err := s.toolRepo.CreateToolRequest(ctx, toolRequestEntity)
	if err != nil {
		return nil, err
	}
	return createdToolRequest.ToDTO(), nil
}

func (s *toolService) UpdateToolRequest(
	ctx context.Context, id int, toolRequest *dto.UpdateToolRequestDTO,
) error {
	toolRequestEntity := &entity.ToolRequest{
		ID:           id,
		ResponseData: toolRequest.ResponseData,
		Status:       toolRequest.Status,
	}

	return s.toolRepo.UpdateToolRequest(ctx, toolRequestEntity)
}

func (s *toolService) DeleteToolRequest(ctx context.Context, id int) error {
	return s.toolRepo.DeleteToolRequest(ctx, id)
}

func (s *toolService) SelectTool(
	ctx context.Context, clientID int, userPrompt string,
) (*dto.SelectToolResponseDTO, error) {
	selectorResponse, err := s.selectorService.Select(ctx, selector.SelectorRequest{
		UserPrompt: userPrompt,
	})
	if err != nil {
		return nil, err
	}

	tool, err := s.toolRepo.FindToolByID(ctx, selectorResponse.ToolID)
	if err != nil {
		return nil, fmt.Errorf("tool not found")
	}

	toolClientPermission, err := s.toolRepo.GetToolClientPermissionByToolIDAndClientID(
		ctx, tool.ID, clientID,
	)
	if err != nil {
		return nil, err
	}

	if toolClientPermission.PermissionLevel != valueobject.ToolClientPermissionLevelWrite {
		selectorResponse.Message += "\n\nYou don't have permission to use this tool. Please contact the administrator."
	}

	return &dto.SelectToolResponseDTO{
		PermissionLevel: toolClientPermission.PermissionLevel,
		Tool:            tool.ToDTO(),
		Message:         selectorResponse.Message,
	}, nil
}

// Core function to execute a tool
// 1. Check if the client has permission to use the tool
// 2. Check if the tool exists
// 3. Call FunctionExecutor to execute the tool
// 4. Create a tool request
// 5. Return the tool request ID
func (s *toolService) ExecuteTool(
	ctx context.Context, clientID int, toolID int, requestData dto.ToolExecutionRequestDTO,
) (*dto.ToolExecutionResponseDTO, error) {
	toolClientPermission, err := s.toolRepo.GetToolClientPermissionByToolIDAndClientID(
		ctx, toolID, clientID,
	)
	if err != nil {
		return &dto.ToolExecutionResponseDTO{
			Status:  valueobject.ToolExecutionStatusUnauthorized,
			Message: "You don't have permission to use this tool. Please contact the administrator.",
		}, nil
	}

	if toolClientPermission.PermissionLevel != valueobject.ToolClientPermissionLevelWrite {
		return &dto.ToolExecutionResponseDTO{
			Status: valueobject.ToolExecutionStatusUnauthorized,
			Message: fmt.Sprintf("You don't have permission to use this tool. Please contact the administrator. (permission level: %d)",
				toolClientPermission.PermissionLevel.Int()),
		}, nil
	}

	tool, err := s.toolRepo.FindToolByID(ctx, toolID)
	if err != nil {
		return &dto.ToolExecutionResponseDTO{
			Status:  valueobject.ToolExecutionStatusFailed,
			Message: "Tool not found",
		}, nil
	}

	toolRequestEntity := &entity.ToolRequest{
		ToolID:   toolID,
		ClientID: clientID,
		RequestData: shared_type.ToolRequestData{
			Payload: requestData.Payload,
		},
		ResponseData: shared_type.ToolRequestResponseData{},
		Status:       valueobject.ToolRequestStatusPending,
	}

	createdToolRequest, err := s.toolRepo.CreateToolRequest(ctx, toolRequestEntity)
	if err != nil {
		return nil, err
	}

	go s.functionExecutor.Sync(ctx, tool, createdToolRequest.ID, requestData)

	return &dto.ToolExecutionResponseDTO{
		Status:        valueobject.ToolExecutionStatusSuccess,
		Message:       "Tool execution started",
		ToolRequestID: createdToolRequest.ID,
	}, nil
}
