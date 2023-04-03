package service

import (
	"context"

	"golang.org/x/oauth2"
)

func (s *service) GetRedirectLink(ctx context.Context) (string, error) {
	authURL := s.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL, nil
}
