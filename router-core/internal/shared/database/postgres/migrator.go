package postgres

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed sql/init.sql
var initial_sql string

func AutoMigrateFromConnectionString(ctx context.Context, connectionString string) (bool, error) {
	config, err := ParsePostgresConnectionString(connectionString)
	if err != nil {
		return false, err
	}

	dbName := config.ConnConfig.Database
	config.ConnConfig.Database = "postgres"

	db, err := NewPostgresPool(ctx, config)
	if err != nil {
		return false, err
	}

	_, err = db.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return true, err
	}

	db.Close()

	config.ConnConfig.Database = dbName

	db, err = NewPostgresPool(ctx, config)
	if err != nil {
		return false, err
	}

	_, err = db.Exec(ctx, initial_sql)
	if err != nil {
		return false, err
	}

	return true, nil
}
