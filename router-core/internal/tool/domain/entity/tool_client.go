package entity

import (
	"time"

	"aigendrug.com/router-core/internal/tool/application/dto"
	"aigendrug.com/router-core/internal/tool/domain/valueobject"
	"github.com/jackc/pgx/v5/pgtype"
)

type ToolClientPermission struct {
	ID              int                                   `json:"id" db:"id"`
	ToolID          int                                   `json:"tool_id" db:"tool_id"`
	ClientID        int                                   `json:"client_id" db:"client_id"`
	PermissionLevel valueobject.ToolClientPermissionLevel `json:"permission_level" db:"permission_level"`
	CreatedAt       time.Time                             `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time                             `json:"updated_at" db:"updated_at"`
}

type ToolClientPermissionRow struct {
	ID              int                `json:"id" db:"id"`
	ToolID          int                `json:"tool_id" db:"tool_id"`
	ClientID        int                `json:"client_id" db:"client_id"`
	PermissionLevel int                `json:"permission_level" db:"permission_level"`
	CreatedAt       pgtype.Timestamptz `json:"created_at" db:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at" db:"updated_at"`
}

func (t *ToolClientPermission) ToRow() *ToolClientPermissionRow {
	return &ToolClientPermissionRow{
		ID:              t.ID,
		ToolID:          t.ToolID,
		ClientID:        t.ClientID,
		PermissionLevel: t.PermissionLevel.Int(),
		CreatedAt:       pgtype.Timestamptz{Time: t.CreatedAt},
		UpdatedAt:       pgtype.Timestamptz{Time: t.UpdatedAt},
	}
}

func (t *ToolClientPermissionRow) ToEntity() *ToolClientPermission {
	return &ToolClientPermission{
		ID:              t.ID,
		ToolID:          t.ToolID,
		ClientID:        t.ClientID,
		PermissionLevel: valueobject.ToolClientPermissionLevel(t.PermissionLevel),
		CreatedAt:       t.CreatedAt.Time,
		UpdatedAt:       t.UpdatedAt.Time,
	}
}

func (t *ToolClientPermission) ToDTO() *dto.ReadToolClientPermissionDTO {
	return &dto.ReadToolClientPermissionDTO{
		ID:              t.ID,
		ToolID:          t.ToolID,
		ClientID:        t.ClientID,
		PermissionLevel: t.PermissionLevel,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}
