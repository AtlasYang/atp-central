package entity

import (
	"encoding/json"
	"time"

	"aigendrug.com/router-core/internal/tool/application/dto"
	"aigendrug.com/router-core/internal/tool/domain/shared_type"
	"aigendrug.com/router-core/internal/tool/domain/valueobject"
	"github.com/jackc/pgx/v5/pgtype"
)

type ToolRequest struct {
	ID           int                                 `json:"id" db:"id"`
	ToolID       int                                 `json:"tool_id" db:"tool_id"`
	ToolName     string                              `json:"tool_name" db:"tool_name"`
	ClientID     int                                 `json:"client_id" db:"client_id"`
	RequestData  shared_type.ToolRequestData         `json:"request_data" db:"request_data"`
	ResponseData shared_type.ToolRequestResponseData `json:"response_data" db:"response_data"`
	Status       valueobject.ToolRequestStatus       `json:"status" db:"status"`
	CreatedAt    time.Time                           `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time                           `json:"updated_at" db:"updated_at"`
}

type ToolRequestRow struct {
	ID           int                `json:"id" db:"id"`
	ToolID       int                `json:"tool_id" db:"tool_id"`
	ToolName     string             `json:"tool_name" db:"tool_name"`
	ClientID     int                `json:"client_id" db:"client_id"`
	RequestData  string             `json:"request_data" db:"request_data"`
	ResponseData string             `json:"response_data" db:"response_data"`
	Status       string             `json:"status" db:"status"`
	CreatedAt    pgtype.Timestamptz `json:"created_at" db:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at" db:"updated_at"`
}

func (t *ToolRequest) ToRow() *ToolRequestRow {
	requestData, err := json.Marshal(t.RequestData)
	if err != nil {
		return nil
	}
	responseData, err := json.Marshal(t.ResponseData)
	if err != nil {
		return nil
	}

	return &ToolRequestRow{
		ID:           t.ID,
		ToolID:       t.ToolID,
		ToolName:     t.ToolName,
		ClientID:     t.ClientID,
		RequestData:  string(requestData),
		ResponseData: string(responseData),
		Status:       t.Status.String(),
		CreatedAt:    pgtype.Timestamptz{Time: t.CreatedAt},
		UpdatedAt:    pgtype.Timestamptz{Time: t.UpdatedAt},
	}
}

func (t *ToolRequestRow) ToEntity() *ToolRequest {
	requestData := shared_type.ToolRequestData{}
	if err := json.Unmarshal([]byte(t.RequestData), &requestData); err != nil {
		return nil
	}
	responseData := shared_type.ToolRequestResponseData{}
	if err := json.Unmarshal([]byte(t.ResponseData), &responseData); err != nil {
		return nil
	}

	return &ToolRequest{
		ID:           t.ID,
		ToolID:       t.ToolID,
		ToolName:     t.ToolName,
		ClientID:     t.ClientID,
		RequestData:  requestData,
		ResponseData: responseData,
		Status:       valueobject.ToolRequestStatus(t.Status),
		CreatedAt:    t.CreatedAt.Time,
		UpdatedAt:    t.UpdatedAt.Time,
	}
}

func (t *ToolRequest) ToDTO() *dto.ReadToolRequestDTO {
	return &dto.ReadToolRequestDTO{
		ID:           t.ID,
		ToolID:       t.ToolID,
		ToolName:     t.ToolName,
		ClientID:     t.ClientID,
		RequestData:  t.RequestData,
		ResponseData: t.ResponseData,
		Status:       t.Status,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}
