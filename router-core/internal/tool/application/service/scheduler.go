package service

import (
	"aigendrug.com/router-core/internal/tool/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ToolRequestScheduler interface {
}

type tooloolRequestScheduler struct {
	db       *pgxpool.Pool
	toolRepo domain.ToolRepository
}

func NewToolRequestScheduler(
	db *pgxpool.Pool,
	toolRepo domain.ToolRepository,
) ToolRequestScheduler {
	return &tooloolRequestScheduler{
		db:       db,
		toolRepo: toolRepo,
	}
}

// func (s *tooloolRequestScheduler) reftechToolRequest(
// 	ctx context.Context,
// 	toolRequest *entity.ToolRequest,
// ) error {
// 	tool, err := s.toolRepo.FindToolByID(ctx, toolRequest.ToolID)
// 	if err != nil {
// 		return err
// 	}
// }
