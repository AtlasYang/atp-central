package persistence

import (
	"context"

	"aigendrug.com/router-core/internal/shared/database/postgres"
	"aigendrug.com/router-core/internal/tool/domain"
	"aigendrug.com/router-core/internal/tool/domain/entity"
	"aigendrug.com/router-core/internal/tool/domain/valueobject"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgToolRepository struct {
	db postgres.DbExecutor
}

func NewPgToolRepository(dbPool *pgxpool.Pool) domain.ToolRepository {
	return &pgToolRepository{db: dbPool}
}

func (r *pgToolRepository) WithTx(ctx context.Context, tx pgx.Tx) domain.ToolRepository {
	return &pgToolRepository{db: tx}
}

func (r *pgToolRepository) FindAllTools(ctx context.Context) ([]*entity.Tool, error) {
	query := `
		SELECT 
			id, uuid, name, 
			version, description, engine_interface, 
			provider_interface, created_at, updated_at
		FROM tools
	`

	var tools []*entity.ToolRow
	if err := pgxscan.Select(ctx, r.db, &tools, query); err != nil {
		return nil, err
	}

	if len(tools) == 0 {
		return []*entity.Tool{}, nil
	}

	toolsEntity := make([]*entity.Tool, len(tools))
	for i, tool := range tools {
		toolsEntity[i] = tool.ToEntity()
	}

	return toolsEntity, nil
}

func (r *pgToolRepository) FindToolByID(ctx context.Context, id int) (*entity.Tool, error) {
	query := `
		SELECT 
			id, uuid, name,
			version, description, engine_interface,
			provider_interface, created_at, updated_at
		FROM tools
		WHERE id = $1
	`

	var tool entity.ToolRow
	if err := pgxscan.Get(ctx, r.db, &tool, query, id); err != nil {
		return nil, err
	}

	return tool.ToEntity(), nil
}

func (r *pgToolRepository) FindToolByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Tool, error) {
	query := `
		SELECT 
			id, uuid, name,
			version, description, engine_interface,
			provider_interface, created_at, updated_at
		FROM tools
		WHERE uuid = $1
	`

	var tool entity.ToolRow
	if err := pgxscan.Get(ctx, r.db, &tool, query, uuid); err != nil {
		return nil, err
	}

	return tool.ToEntity(), nil
}

func (r *pgToolRepository) FindAllToolsByClientID(
	ctx context.Context, clientID int, permissionLevel valueobject.ToolClientPermissionLevel,
) ([]*entity.Tool, error) {
	query := `
		SELECT
			t.id, t.uuid, t.name,
			t.version, t.description, t.engine_interface,
			t.provider_interface, t.created_at, t.updated_at
		FROM tools t
		JOIN tool_client_permissions tcp ON t.id = tcp.tool_id
		WHERE tcp.client_id = $1 AND tcp.permission_level = $2
	`

	var tools []*entity.ToolRow
	if err := pgxscan.Select(ctx, r.db, &tools, query, clientID, permissionLevel); err != nil {
		return nil, err
	}

	if len(tools) == 0 {
		return []*entity.Tool{}, nil
	}

	toolsEntity := make([]*entity.Tool, len(tools))
	for i, tool := range tools {
		toolsEntity[i] = tool.ToEntity()
	}

	return toolsEntity, nil
}

func (r *pgToolRepository) CreateTool(ctx context.Context, tool *entity.Tool) (*entity.Tool, error) {
	query := `
		INSERT INTO tools (uuid, name, version, description, engine_interface, provider_interface)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING 
			id, uuid, name, 
			version, description, engine_interface, 
			provider_interface, created_at, updated_at
	`

	toolRaw := tool.ToRow()

	createdTool := &entity.ToolRow{}
	if err := r.db.QueryRow(ctx, query,
		toolRaw.UUID, toolRaw.Name, toolRaw.Version,
		toolRaw.Description, toolRaw.EngineInterface, toolRaw.ProviderInterface,
	).Scan(
		&createdTool.ID,
		&createdTool.UUID,
		&createdTool.Name,
		&createdTool.Version,
		&createdTool.Description,
		&createdTool.EngineInterface,
		&createdTool.ProviderInterface,
		&createdTool.CreatedAt,
		&createdTool.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return createdTool.ToEntity(), nil
}

func (r *pgToolRepository) UpdateTool(ctx context.Context, tool *entity.Tool) error {
	query := `
		UPDATE tools
		SET 
			name = $1, version = $2, description = $3, 
			engine_interface = $4, provider_interface = $5, 
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`

	toolRaw := tool.ToRow()

	_, err := r.db.Exec(ctx, query,
		toolRaw.Name, toolRaw.Version, toolRaw.Description,
		toolRaw.EngineInterface, toolRaw.ProviderInterface, toolRaw.ID,
	)

	return err
}

func (r *pgToolRepository) DeleteTool(ctx context.Context, id int) error {
	query := `
		DELETE FROM tools
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *pgToolRepository) FindAllToolClientPermissionsByToolID(
	ctx context.Context, toolID int,
) ([]*entity.ToolClientPermission, error) {
	query := `
		SELECT id, tool_id, client_id, permission_level, created_at, updated_at
		FROM tool_client_permissions
		WHERE tool_id = $1
	`

	var permissions []*entity.ToolClientPermissionRow
	if err := pgxscan.Select(ctx, r.db, &permissions, query, toolID); err != nil {
		return nil, err
	}

	if len(permissions) == 0 {
		return []*entity.ToolClientPermission{}, nil
	}

	permissionsEntity := make([]*entity.ToolClientPermission, len(permissions))
	for i, permission := range permissions {
		permissionsEntity[i] = permission.ToEntity()
	}

	return permissionsEntity, nil
}

func (r *pgToolRepository) FindAllToolClientPermissionsByClientID(
	ctx context.Context, clientID int,
) ([]*entity.ToolClientPermission, error) {
	query := `
		SELECT id, tool_id, client_id, permission_level, created_at, updated_at
		FROM tool_client_permissions
		WHERE client_id = $1
	`

	var permissions []*entity.ToolClientPermissionRow
	if err := pgxscan.Select(ctx, r.db, &permissions, query, clientID); err != nil {
		return nil, err
	}

	if len(permissions) == 0 {
		return []*entity.ToolClientPermission{}, nil
	}

	permissionsEntity := make([]*entity.ToolClientPermission, len(permissions))
	for i, permission := range permissions {
		permissionsEntity[i] = permission.ToEntity()
	}

	return permissionsEntity, nil
}

func (r *pgToolRepository) GetToolClientPermissionByToolIDAndClientID(
	ctx context.Context, toolID int, clientID int,
) (*entity.ToolClientPermission, error) {
	query := `
		SELECT id, tool_id, client_id, permission_level, created_at, updated_at
		FROM tool_client_permissions
		WHERE tool_id = $1 AND client_id = $2
	`

	var permission entity.ToolClientPermissionRow
	if err := pgxscan.Get(ctx, r.db, &permission, query, toolID, clientID); err != nil {
		return nil, err
	}

	return permission.ToEntity(), nil
}

func (r *pgToolRepository) CreateToolClientPermission(
	ctx context.Context, permission *entity.ToolClientPermission,
) (*entity.ToolClientPermission, error) {
	query := `
		INSERT INTO tool_client_permissions (tool_id, client_id, permission_level)
		VALUES ($1, $2, $3)
		RETURNING id, tool_id, client_id, permission_level, created_at, updated_at
	`

	createdPermission := &entity.ToolClientPermissionRow{}
	if err := r.db.QueryRow(ctx, query,
		permission.ToolID, permission.ClientID, permission.PermissionLevel,
	).Scan(
		&createdPermission.ID,
		&createdPermission.ToolID,
		&createdPermission.ClientID,
		&createdPermission.PermissionLevel,
		&createdPermission.CreatedAt,
		&createdPermission.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return createdPermission.ToEntity(), nil
}

func (r *pgToolRepository) UpdateToolClientPermission(
	ctx context.Context, permission *entity.ToolClientPermission,
) error {
	query := `
		UPDATE tool_client_permissions
		SET permission_level = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, permission.PermissionLevel, permission.ID)
	return err
}

func (r *pgToolRepository) DeleteToolClientPermission(ctx context.Context, id int) error {
	query := `
		DELETE FROM tool_client_permissions
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *pgToolRepository) FindToolRequestByID(
	ctx context.Context, id int,
) (*entity.ToolRequest, error) {
	query := `
		SELECT 
			tr.id, 
			tr.tool_id, 
			t.name as tool_name,
			tr.client_id, 
			tr.request_data, 
			tr.response_data, 
			tr.status, 
			tr.created_at, 
			tr.updated_at
		FROM tool_requests tr
		JOIN tools t ON tr.tool_id = t.id
		WHERE tr.id = $1
	`

	var request entity.ToolRequestRow
	if err := pgxscan.Get(ctx, r.db, &request, query, id); err != nil {
		return nil, err
	}

	return request.ToEntity(), nil
}

func (r *pgToolRepository) FindAllToolRequestsByToolID(
	ctx context.Context, toolID int,
) ([]*entity.ToolRequest, error) {
	query := `
		SELECT 
			tr.id, 
			tr.tool_id, 
			t.name as tool_name,
			tr.client_id, 
			tr.request_data, 
			tr.response_data, 
			tr.status, 
			tr.created_at, 
			tr.updated_at
		FROM tool_requests tr
		JOIN tools t ON tr.tool_id = t.id
		WHERE tool_id = $1
	`

	var requests []*entity.ToolRequestRow
	if err := pgxscan.Select(ctx, r.db, &requests, query, toolID); err != nil {
		return nil, err
	}

	if len(requests) == 0 {
		return []*entity.ToolRequest{}, nil
	}

	requestsEntity := make([]*entity.ToolRequest, len(requests))
	for i, request := range requests {
		requestsEntity[i] = request.ToEntity()
	}

	return requestsEntity, nil
}

func (r *pgToolRepository) FindAllToolRequestsByClientID(
	ctx context.Context, clientID int,
) ([]*entity.ToolRequest, error) {
	query := `
		SELECT 
			tr.id, 
			tr.tool_id, 
			t.name as tool_name,
			tr.client_id, 
			tr.request_data, 
			tr.response_data, 
			tr.status, 
			tr.created_at, 
			tr.updated_at
		FROM tool_requests tr
		JOIN tools t ON tr.tool_id = t.id
		WHERE client_id = $1
	`

	var requests []*entity.ToolRequestRow
	if err := pgxscan.Select(ctx, r.db, &requests, query, clientID); err != nil {
		return nil, err
	}

	if len(requests) == 0 {
		return []*entity.ToolRequest{}, nil
	}

	requestsEntity := make([]*entity.ToolRequest, len(requests))
	for i, request := range requests {
		requestsEntity[i] = request.ToEntity()
	}

	return requestsEntity, nil
}

func (r *pgToolRepository) CreateToolRequest(
	ctx context.Context, request *entity.ToolRequest,
) (*entity.ToolRequest, error) {
	query := `
		INSERT INTO tool_requests (tool_id, client_id, request_data, response_data, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id, tool_id, client_id,
			request_data, response_data, status, 
			created_at, updated_at
	`

	requestRaw := request.ToRow()

	createdRequest := &entity.ToolRequestRow{}
	if err := r.db.QueryRow(ctx, query,
		requestRaw.ToolID, requestRaw.ClientID, requestRaw.RequestData, requestRaw.ResponseData, requestRaw.Status,
	).Scan(
		&createdRequest.ID,
		&createdRequest.ToolID,
		&createdRequest.ClientID,
		&createdRequest.RequestData,
		&createdRequest.ResponseData,
		&createdRequest.Status,
		&createdRequest.CreatedAt,
		&createdRequest.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return createdRequest.ToEntity(), nil
}

func (r *pgToolRepository) UpdateToolRequest(
	ctx context.Context, request *entity.ToolRequest,
) error {
	query := `
		UPDATE tool_requests
		SET request_data = $1, response_data = $2, status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`

	requestRaw := request.ToRow()

	_, err := r.db.Exec(ctx, query,
		requestRaw.RequestData, requestRaw.ResponseData, requestRaw.Status, requestRaw.ID,
	)

	return err
}

func (r *pgToolRepository) DeleteToolRequest(ctx context.Context, id int) error {
	query := `
		DELETE FROM tool_requests
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
