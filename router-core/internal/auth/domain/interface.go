package domain

import (
	"context"
)

type SessionRepository interface {
	GetClientIDByAPIKey(ctx context.Context, apiKey string) (int, error)
}
