package entity

import (
	"time"

	"aigendrug.com/router-core/internal/client/application/dto"
	"github.com/jackc/pgx/v5/pgtype"
)

type Client struct {
	ID               int       `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	ClientIdentifier string    `json:"client_identifier" db:"client_identifier"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type ClientRow struct {
	ID               int                `json:"id" db:"id"`
	Name             string             `json:"name" db:"name"`
	ClientIdentifier string             `json:"client_identifier" db:"client_identifier"`
	IsActive         bool               `json:"is_active" db:"is_active"`
	CreatedAt        pgtype.Timestamptz `json:"created_at" db:"created_at"`
}

func (c *Client) ToDTO() *dto.ReadClientDTO {
	return &dto.ReadClientDTO{
		ID:               c.ID,
		Name:             c.Name,
		ClientIdentifier: c.ClientIdentifier,
		IsActive:         c.IsActive,
		CreatedAt:        c.CreatedAt,
	}
}

func (c *ClientRow) ToEntity() *Client {
	return &Client{
		ID:               c.ID,
		Name:             c.Name,
		ClientIdentifier: c.ClientIdentifier,
		IsActive:         c.IsActive,
		CreatedAt:        c.CreatedAt.Time,
	}
}

func (c *Client) ToRow() *ClientRow {
	return &ClientRow{
		ID:               c.ID,
		Name:             c.Name,
		ClientIdentifier: c.ClientIdentifier,
		IsActive:         c.IsActive,
		CreatedAt:        pgtype.Timestamptz{Time: c.CreatedAt},
	}
}

type ClientApiKey struct {
	ID        int       `json:"id" db:"id"`
	ClientID  int       `json:"client_id" db:"client_id"`
	ApiKey    string    `json:"api_key" db:"api_key"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ClientApiKeyRow struct {
	ID        int                `json:"id" db:"id"`
	ClientID  int                `json:"client_id" db:"client_id"`
	ApiKey    string             `json:"api_key" db:"api_key"`
	CreatedAt pgtype.Timestamptz `json:"created_at" db:"created_at"`
}

func (c *ClientApiKey) ToDTO() *dto.ReadClientApiKeyDTO {
	return &dto.ReadClientApiKeyDTO{
		ID:        c.ID,
		ClientID:  c.ClientID,
		ApiKey:    c.ApiKey,
		CreatedAt: c.CreatedAt,
	}
}

func (c *ClientApiKeyRow) ToEntity() *ClientApiKey {
	return &ClientApiKey{
		ID:        c.ID,
		ClientID:  c.ClientID,
		ApiKey:    c.ApiKey,
		CreatedAt: c.CreatedAt.Time,
	}
}

func (c *ClientApiKey) ToRow() *ClientApiKeyRow {
	return &ClientApiKeyRow{
		ID:        c.ID,
		ClientID:  c.ClientID,
		ApiKey:    c.ApiKey,
		CreatedAt: pgtype.Timestamptz{Time: c.CreatedAt},
	}
}
