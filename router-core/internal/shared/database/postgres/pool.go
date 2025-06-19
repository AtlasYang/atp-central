package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(c context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.NewWithConfig(c, config)

	if err != nil {
		return nil, err
	}

	return dbpool, nil
}

func NewPostgresPoolFromConnectionString(c context.Context, connectionString string) (*pgxpool.Pool, error) {
	config, err := ParsePostgresConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return NewPostgresPool(c, config)
}
