package entity

import (
	"encoding/json"
	"time"

	"aigendrug.com/router-core/internal/tool/application/dto"
	"aigendrug.com/router-core/internal/tool/domain/shared_type"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Tool struct {
	ID                int                           `json:"id" db:"id"`
	UUID              uuid.UUID                     `json:"uuid" db:"uuid"`
	Name              string                        `json:"name" db:"name"`
	Version           string                        `json:"version" db:"version"`
	Description       string                        `json:"description" db:"description"`
	EngineInterface   shared_type.EngineInterface   `json:"engine_interface" db:"engine_interface"`
	ProviderInterface shared_type.ProviderInterface `json:"provider_interface" db:"provider_interface"`
	CreatedAt         time.Time                     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time                     `json:"updated_at" db:"updated_at"`
}

type ToolRow struct {
	ID                int                `json:"id" db:"id"`
	UUID              pgtype.UUID        `json:"uuid" db:"uuid"`
	Name              string             `json:"name" db:"name"`
	Version           string             `json:"version" db:"version"`
	Description       string             `json:"description" db:"description"`
	EngineInterface   string             `json:"engine_interface" db:"engine_interface"`
	ProviderInterface string             `json:"provider_interface" db:"provider_interface"`
	CreatedAt         pgtype.Timestamptz `json:"created_at" db:"created_at"`
	UpdatedAt         pgtype.Timestamptz `json:"updated_at" db:"updated_at"`
}

func (t *Tool) ToRow() *ToolRow {
	engineInterface, err := json.Marshal(t.EngineInterface)
	if err != nil {
		return nil
	}
	providerInterface, err := json.Marshal(t.ProviderInterface)
	if err != nil {
		return nil
	}

	uuidBin, err := t.UUID.MarshalBinary()
	if err != nil {
		return nil
	}
	var uuidBytes [16]byte
	copy(uuidBytes[:], uuidBin)
	uuid := pgtype.UUID{Bytes: uuidBytes, Valid: true}

	return &ToolRow{
		ID:                t.ID,
		UUID:              uuid,
		Name:              t.Name,
		Version:           t.Version,
		Description:       t.Description,
		EngineInterface:   string(engineInterface),
		ProviderInterface: string(providerInterface),
		CreatedAt:         pgtype.Timestamptz{Time: t.CreatedAt},
		UpdatedAt:         pgtype.Timestamptz{Time: t.UpdatedAt},
	}
}

func (tr *ToolRow) ToEntity() *Tool {
	engineInterface := shared_type.EngineInterface{}
	if err := json.Unmarshal([]byte(tr.EngineInterface), &engineInterface); err != nil {
		return nil
	}
	providerInterface := shared_type.ProviderInterface{}
	if err := json.Unmarshal([]byte(tr.ProviderInterface), &providerInterface); err != nil {
		return nil
	}

	uuid, err := uuid.FromBytes(tr.UUID.Bytes[:])
	if err != nil {
		return nil
	}

	return &Tool{
		ID:                tr.ID,
		UUID:              uuid,
		Name:              tr.Name,
		Version:           tr.Version,
		Description:       tr.Description,
		EngineInterface:   engineInterface,
		ProviderInterface: providerInterface,
		CreatedAt:         tr.CreatedAt.Time,
		UpdatedAt:         tr.UpdatedAt.Time,
	}
}

func (t *Tool) ToDTO() *dto.ReadToolDTO {
	return &dto.ReadToolDTO{
		ID:                t.ID,
		UUID:              t.UUID,
		Name:              t.Name,
		Version:           t.Version,
		Description:       t.Description,
		EngineInterface:   t.EngineInterface,
		ProviderInterface: t.ProviderInterface,
	}
}
