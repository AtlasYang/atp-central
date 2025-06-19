package domain

import (
	"context"

	"aigendrug.com/router-core/internal/client/domain/entity"
	"github.com/jackc/pgx/v5"
)

type ClientRepository interface {
	WithTx(ctx context.Context, tx pgx.Tx) ClientRepository

	FindAllClients(ctx context.Context) ([]*entity.Client, error)
	FindClientByID(ctx context.Context, id int) (*entity.Client, error)
	FindClientByClientIdentifier(ctx context.Context, clientIdentifier string) (*entity.Client, error)
	CreateClient(ctx context.Context, client *entity.Client) (*entity.Client, error)
	UpdateClient(ctx context.Context, client *entity.Client) error
	DeleteClient(ctx context.Context, id int) error

	FindClientApiKeyByClientID(ctx context.Context, clientID int) ([]*entity.ClientApiKey, error)
	CreateClientApiKey(ctx context.Context, clientApiKey *entity.ClientApiKey) (*entity.ClientApiKey, error)
	DeleteClientApiKey(ctx context.Context, id int) error
}
