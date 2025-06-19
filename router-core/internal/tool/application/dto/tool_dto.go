package dto

import (
	"time"

	"aigendrug.com/router-core/internal/tool/domain/shared_type"
	"aigendrug.com/router-core/internal/tool/domain/valueobject"
	"github.com/google/uuid"
)

type ReadToolDTO struct {
	ID                int                           `json:"id" example:"1"`
	UUID              uuid.UUID                     `json:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name              string                        `json:"name" example:"Tool Name"`
	Version           string                        `json:"version" example:"1.0.0"`
	Description       string                        `json:"description" example:"Tool Description"`
	EngineInterface   shared_type.EngineInterface   `json:"engine_interface"`
	ProviderInterface shared_type.ProviderInterface `json:"provider_interface"`
}

type CreateToolDTO struct {
	Name              string                        `json:"name" example:"Tool Name"`
	Version           string                        `json:"version" example:"1.0.0"`
	Description       string                        `json:"description" example:"Tool Description"`
	EngineInterface   shared_type.EngineInterface   `json:"engine_interface"`
	ProviderInterface shared_type.ProviderInterface `json:"provider_interface"`
}

type UpdateToolDTO struct {
	Description       string                        `json:"description" example:"Tool Description"`
	EngineInterface   shared_type.EngineInterface   `json:"engine_interface"`
	ProviderInterface shared_type.ProviderInterface `json:"provider_interface"`
}

type ReadToolClientPermissionDTO struct {
	ID              int                                   `json:"id" example:"1"`
	ToolID          int                                   `json:"tool_id" example:"1"`
	ClientID        int                                   `json:"client_id" example:"1"`
	PermissionLevel valueobject.ToolClientPermissionLevel `json:"permission_level" example:"1"`
	CreatedAt       time.Time                             `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt       time.Time                             `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type CreateToolClientPermissionDTO struct {
	ToolID          int                                   `json:"tool_id" example:"1"`
	ClientID        int                                   `json:"client_id" example:"1"`
	PermissionLevel valueobject.ToolClientPermissionLevel `json:"permission_level" example:"read"`
}

type UpdateToolClientPermissionDTO struct {
	PermissionLevel valueobject.ToolClientPermissionLevel `json:"permission_level" example:"read"`
}

type ReadToolRequestDTO struct {
	ID           int                                 `json:"id" example:"1"`
	ToolID       int                                 `json:"tool_id" example:"1"`
	ToolName     string                              `json:"tool_name" example:"Tool Name"`
	ClientID     int                                 `json:"client_id" example:"1"`
	RequestData  shared_type.ToolRequestData         `json:"request_data"`
	ResponseData shared_type.ToolRequestResponseData `json:"response_data"`
	Status       valueobject.ToolRequestStatus       `json:"status" example:"pending"`
	CreatedAt    time.Time                           `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt    time.Time                           `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type CreateToolRequestDTO struct {
	ToolID       int                                 `json:"tool_id" example:"1"`
	ClientID     int                                 `json:"client_id" example:"1"`
	RequestData  shared_type.ToolRequestData         `json:"request_data"`
	ResponseData shared_type.ToolRequestResponseData `json:"response_data"`
	Status       valueobject.ToolRequestStatus       `json:"status" example:"pending"`
}

type UpdateToolRequestDTO struct {
	ResponseData shared_type.ToolRequestResponseData `json:"response_data"`
	Status       valueobject.ToolRequestStatus       `json:"status" example:"pending"`
}

type SelectToolRequestDTO struct {
	UserPrompt string `json:"user_prompt" example:"i want to add two numbers"`
}

type SelectToolResponseDTO struct {
	PermissionLevel valueobject.ToolClientPermissionLevel `json:"permission_level"`
	Tool            *ReadToolDTO                          `json:"tool"`
	Message         string                                `json:"message"`
}

type ToolExecutionRequestDTO struct {
	Payload map[string]any `json:"payload"`
}

type ToolExecutionResponseDTO struct {
	Status        valueobject.ToolExecutionStatus `json:"status"`
	Message       string                          `json:"message"`
	ToolRequestID int                             `json:"tool_request_id"`
}
