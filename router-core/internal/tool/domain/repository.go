package domain

import (
	"context"

	"aigendrug.com/router-core/internal/tool/domain/entity"
	"aigendrug.com/router-core/internal/tool/domain/valueobject"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ToolRepository interface {
	WithTx(ctx context.Context, tx pgx.Tx) ToolRepository

	// Tool
	FindAllTools(ctx context.Context) ([]*entity.Tool, error)
	FindToolByID(ctx context.Context, id int) (*entity.Tool, error)
	FindToolByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Tool, error)
	FindAllToolsByClientID(ctx context.Context, clientID int, permissionLevel valueobject.ToolClientPermissionLevel) ([]*entity.Tool, error)
	CreateTool(ctx context.Context, tool *entity.Tool) (*entity.Tool, error)
	UpdateTool(ctx context.Context, tool *entity.Tool) error
	DeleteTool(ctx context.Context, id int) error

	// ToolClientPermission
	FindAllToolClientPermissionsByToolID(ctx context.Context, toolID int) ([]*entity.ToolClientPermission, error)
	FindAllToolClientPermissionsByClientID(ctx context.Context, clientID int) ([]*entity.ToolClientPermission, error)
	GetToolClientPermissionByToolIDAndClientID(ctx context.Context, toolID int, clientID int) (*entity.ToolClientPermission, error)
	CreateToolClientPermission(ctx context.Context, toolClientPermission *entity.ToolClientPermission) (*entity.ToolClientPermission, error)
	UpdateToolClientPermission(ctx context.Context, toolClientPermission *entity.ToolClientPermission) error
	DeleteToolClientPermission(ctx context.Context, id int) error

	// ToolRequest
	FindToolRequestByID(ctx context.Context, id int) (*entity.ToolRequest, error)
	FindAllToolRequestsByToolID(ctx context.Context, toolID int) ([]*entity.ToolRequest, error)
	FindAllToolRequestsByClientID(ctx context.Context, clientID int) ([]*entity.ToolRequest, error)
	CreateToolRequest(ctx context.Context, toolRequest *entity.ToolRequest) (*entity.ToolRequest, error)
	UpdateToolRequest(ctx context.Context, toolRequest *entity.ToolRequest) error
	DeleteToolRequest(ctx context.Context, id int) error
}
