package selector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"aigendrug.com/router-core/internal/config"
)

type selectorService struct {
	selectorURL string
}

type SelectorService interface {
	Select(ctx context.Context, request SelectorRequest) (SelectorResponse, error)
}

func NewSelectorService(config *config.Config) SelectorService {
	return &selectorService{
		selectorURL: config.Selector.URL,
	}
}

func (s *selectorService) Select(ctx context.Context, request SelectorRequest) (SelectorResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return SelectorResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	response, err := http.Post(fmt.Sprintf("%s/api/v1/select", s.selectorURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return SelectorResponse{}, fmt.Errorf("failed to select tool: %w", err)
	}

	defer response.Body.Close()

	var selectorResponse SelectorResponse
	if err := json.NewDecoder(response.Body).Decode(&selectorResponse); err != nil {
		return SelectorResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return selectorResponse, nil
}
