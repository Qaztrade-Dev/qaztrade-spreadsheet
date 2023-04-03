package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/google/domain"
	"golang.org/x/oauth2"
)

type Service interface {
	GetRedirectLink(ctx context.Context) (string, error)
	UpdateToken(ctx context.Context, authCode string) error
}

type service struct {
	config *oauth2.Config
	repo   domain.TokenRepository
}

func NewService(config *oauth2.Config, repo domain.TokenRepository) Service {
	return &service{
		config: config,
		repo:   repo,
	}
}
