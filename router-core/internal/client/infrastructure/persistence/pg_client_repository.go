package persistence

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"aigendrug.com/router-core/internal/client/domain"
	"aigendrug.com/router-core/internal/client/domain/entity"
	"aigendrug.com/router-core/internal/shared/database/postgres"
)

type pgClientRepository struct {
	db postgres.DbExecutor
}

func NewPgClientRepository(dbPool *pgxpool.Pool) domain.ClientRepository {
	return &pgClientRepository{db: dbPool}
}

func (r *pgClientRepository) WithTx(ctx context.Context, tx pgx.Tx) domain.ClientRepository {
	return &pgClientRepository{db: tx}
}

func (r *pgClientRepository) FindAllClients(ctx context.Context) ([]*entity.Client, error) {
	query := `
		SELECT id, name, client_identifier, is_active, created_at
		FROM clients
	`

	var clients []*entity.ClientRow
	if err := pgxscan.Select(ctx, r.db, &clients, query); err != nil {
		return nil, err
	}

	if len(clients) == 0 {
		return []*entity.Client{}, nil
	}

	clientsEntity := make([]*entity.Client, len(clients))
	for i, client := range clients {
		clientsEntity[i] = client.ToEntity()
	}

	return clientsEntity, nil
}

func (r *pgClientRepository) FindClientByID(ctx context.Context, id int) (*entity.Client, error) {
	query := `
		SELECT id, name, client_identifier, is_active, created_at
		FROM clients
		WHERE id = $1
	`

	var client entity.ClientRow
	if err := pgxscan.Get(ctx, r.db, &client, query, id); err != nil {
		return nil, err
	}

	return client.ToEntity(), nil
}

func (r *pgClientRepository) FindClientByClientIdentifier(
	ctx context.Context, clientIdentifier string,
) (*entity.Client, error) {
	query := `
		SELECT id, name, client_identifier, is_active, created_at
		FROM clients
		WHERE client_identifier = $1
	`

	var client entity.ClientRow
	if err := pgxscan.Get(ctx, r.db, &client, query, clientIdentifier); err != nil {
		return nil, err
	}

	return client.ToEntity(), nil
}

func (r *pgClientRepository) CreateClient(
	ctx context.Context, client *entity.Client,
) (*entity.Client, error) {
	query := `
		INSERT INTO clients (name, client_identifier, is_active)
		VALUES ($1, $2, $3)
		RETURNING id, name, client_identifier, is_active, created_at
	`

	createdClient := &entity.ClientRow{}
	if err := r.db.QueryRow(ctx, query,
		client.Name, client.ClientIdentifier, client.IsActive,
	).Scan(
		&createdClient.ID,
		&createdClient.Name,
		&createdClient.ClientIdentifier,
		&createdClient.IsActive,
		&createdClient.CreatedAt,
	); err != nil {
		return nil, err
	}

	return createdClient.ToEntity(), nil
}

func (r *pgClientRepository) UpdateClient(ctx context.Context, client *entity.Client) error {
	query := `
		UPDATE clients
		SET name = $1, client_identifier = $2, is_active = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(ctx, query,
		client.Name, client.ClientIdentifier, client.IsActive, client.ID,
	)

	return err
}

func (r *pgClientRepository) DeleteClient(ctx context.Context, id int) error {
	query := `
		DELETE FROM clients
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *pgClientRepository) FindClientApiKeyByClientID(ctx context.Context, clientID int) ([]*entity.ClientApiKey, error) {
	query := `
		SELECT id, client_id, api_key, created_at
		FROM client_api_keys
		WHERE client_id = $1
		`

	var clientApiKey []*entity.ClientApiKeyRow
	if err := pgxscan.Select(ctx, r.db, &clientApiKey, query, clientID); err != nil {
		return nil, err
	}

	if len(clientApiKey) == 0 {
		return nil, nil
	}

	clientApiKeyEntity := make([]*entity.ClientApiKey, len(clientApiKey))
	for i, clientApiKey := range clientApiKey {
		clientApiKeyEntity[i] = clientApiKey.ToEntity()
	}

	return clientApiKeyEntity, nil
}

func (r *pgClientRepository) CreateClientApiKey(ctx context.Context, clientApiKey *entity.ClientApiKey) (*entity.ClientApiKey, error) {
	query := `
		INSERT INTO client_api_keys (client_id, api_key)
		VALUES ($1, $2)
		RETURNING id, client_id, api_key, created_at
	`

	createdClientApiKey := &entity.ClientApiKeyRow{}
	if err := r.db.QueryRow(ctx, query, clientApiKey.ClientID, clientApiKey.ApiKey).Scan(
		&createdClientApiKey.ID,
		&createdClientApiKey.ClientID,
		&createdClientApiKey.ApiKey,
		&createdClientApiKey.CreatedAt,
	); err != nil {
		return nil, err
	}

	return createdClientApiKey.ToEntity(), nil
}

func (r *pgClientRepository) DeleteClientApiKey(ctx context.Context, id int) error {
	query := `
		DELETE FROM client_api_keys
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
