package service

import (
	"strings"

	"aigendrug.com/router-core/internal/auth/domain"
	"github.com/gin-gonic/gin"
)

type SessionService struct {
	Repo domain.SessionRepository
}

func NewSessionService(repo domain.SessionRepository) *SessionService {
	return &SessionService{Repo: repo}
}

func (s *SessionService) Authorize(c *gin.Context, apiKey string) (int, error) {
	clientID, err := s.Repo.GetClientIDByAPIKey(c.Request.Context(), apiKey)
	if err != nil {
		return 0, err
	}
	return clientID, nil
}

// route /console is not exposed to network, so validate admin by host 'localhost'
func (s *SessionService) AuthorizeAdmin(c *gin.Context) (bool, error) {
	host := c.Request.Host

	if !strings.HasPrefix(host, "localhost") {
		return false, nil
	}

	return true, nil
}
