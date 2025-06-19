package dto

import (
	"time"
)

type ReadClientDTO struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	ClientIdentifier string    `json:"client_identifier"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}

type CreateClientDTO struct {
	Name             string `json:"name" example:"Client Name"`
	ClientIdentifier string `json:"client_identifier" example:"client_name"`
	IsActive         bool   `json:"is_active" example:"true"`
}

type UpdateClientDTO struct {
	Name             string `json:"name"`
	ClientIdentifier string `json:"client_identifier"`
	IsActive         bool   `json:"is_active"`
}

type ReadClientApiKeyDTO struct {
	ID        int       `json:"id"`
	ClientID  int       `json:"client_id"`
	ApiKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}
