package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgSessionRepository struct {
	db *pgxpool.Pool
}

func NewPgSessionRepository(db *pgxpool.Pool) *PgSessionRepository {
	return &PgSessionRepository{db: db}
}

func (r *PgSessionRepository) GetClientIDByAPIKey(ctx context.Context, apiKey string) (int, error) {
	var clientID int
	err := r.db.QueryRow(ctx, `
		SELECT client_id FROM client_api_keys WHERE api_key = $1
	`, apiKey).Scan(&clientID)
	if err != nil {
		return 0, err
	}
	return clientID, nil
}
