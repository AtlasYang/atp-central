package service

import (
	"context"

	"aigendrug.com/router-core/internal/client/application/dto"
	"aigendrug.com/router-core/internal/client/domain"
	"aigendrug.com/router-core/internal/client/domain/entity"
	"aigendrug.com/router-core/internal/shared/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientService interface {
	GetAllClients(ctx context.Context) ([]*dto.ReadClientDTO, error)
	GetClientByID(ctx context.Context, id int) (*dto.ReadClientDTO, error)
	GetClientByClientIdentifier(ctx context.Context, clientIdentifier string) (*dto.ReadClientDTO, error)
	CreateClient(ctx context.Context, client *dto.CreateClientDTO) (*dto.ReadClientDTO, error)
	UpdateClient(ctx context.Context, id int, client *dto.UpdateClientDTO) error
	DeleteClient(ctx context.Context, id int) error

	FindClientApiKeyByClientID(ctx context.Context, clientID int) ([]*dto.ReadClientApiKeyDTO, error)
	CreateClientApiKey(ctx context.Context, clientID int) (*dto.ReadClientApiKeyDTO, error)
	DeleteClientApiKey(ctx context.Context, id int) error
}

type clientService struct {
	db         *pgxpool.Pool
	clientRepo domain.ClientRepository
}

func NewClientService(
	dbPool *pgxpool.Pool,
	clientRepo domain.ClientRepository,
) ClientService {
	return &clientService{
		db:         dbPool,
		clientRepo: clientRepo,
	}
}

func (s *clientService) GetAllClients(ctx context.Context) ([]*dto.ReadClientDTO, error) {
	clients, err := s.clientRepo.FindAllClients(ctx)
	if err != nil {
		return nil, err
	}

	clientsDTO := make([]*dto.ReadClientDTO, len(clients))
	for i, client := range clients {
		clientsDTO[i] = client.ToDTO()
	}
	return clientsDTO, nil
}

func (s *clientService) GetClientByID(ctx context.Context, id int) (*dto.ReadClientDTO, error) {
	client, err := s.clientRepo.FindClientByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return client.ToDTO(), nil
}

func (s *clientService) GetClientByClientIdentifier(
	ctx context.Context, clientIdentifier string,
) (*dto.ReadClientDTO, error) {
	client, err := s.clientRepo.FindClientByClientIdentifier(ctx, clientIdentifier)
	if err != nil {
		return nil, err
	}
	return client.ToDTO(), nil
}

func (s *clientService) CreateClient(
	ctx context.Context, client *dto.CreateClientDTO,
) (*dto.ReadClientDTO, error) {
	createdClient, err := s.clientRepo.CreateClient(ctx, &entity.Client{
		Name:             client.Name,
		ClientIdentifier: client.ClientIdentifier,
		IsActive:         client.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return createdClient.ToDTO(), nil
}

func (s *clientService) UpdateClient(ctx context.Context, id int, client *dto.UpdateClientDTO) error {
	return s.clientRepo.UpdateClient(ctx, &entity.Client{
		ID:               id,
		Name:             client.Name,
		ClientIdentifier: client.ClientIdentifier,
		IsActive:         client.IsActive,
	})
}

func (s *clientService) DeleteClient(ctx context.Context, id int) error {
	return s.clientRepo.DeleteClient(ctx, id)
}

func (s *clientService) FindClientApiKeyByClientID(ctx context.Context, clientID int) ([]*dto.ReadClientApiKeyDTO, error) {
	clientApiKeys, err := s.clientRepo.FindClientApiKeyByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	clientApiKeysDTO := make([]*dto.ReadClientApiKeyDTO, len(clientApiKeys))
	for i, clientApiKey := range clientApiKeys {
		clientApiKeysDTO[i] = clientApiKey.ToDTO()
	}
	return clientApiKeysDTO, nil
}

func (s *clientService) CreateClientApiKey(ctx context.Context, clientID int) (*dto.ReadClientApiKeyDTO, error) {
	apiKey := "atp-" + utils.GenerateRandomString(32)
	createdClientApiKey, err := s.clientRepo.CreateClientApiKey(ctx, &entity.ClientApiKey{
		ClientID: clientID,
		ApiKey:   apiKey,
	})
	if err != nil {
		return nil, err
	}
	return createdClientApiKey.ToDTO(), nil
}

func (s *clientService) DeleteClientApiKey(ctx context.Context, id int) error {
	return s.clientRepo.DeleteClientApiKey(ctx, id)
}
