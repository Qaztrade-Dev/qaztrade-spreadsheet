package service

import (
	"context"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
)

type SignInRequest struct {
	Email    string
	Password string
}

func (s *service) SignIn(ctx context.Context, input *SignInRequest) (*domain.Credentials, error) {
	email := strings.TrimSpace(input.Email)

	user, err := s.authRepo.SignIn(ctx, &domain.SignInInput{
		Email:    email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	roles, err := s.authRepo.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	creds, err := s.credsRepo.Create(ctx, &domain.UserClaims{
		UserID: user.ID,
		Roles:  roles,
	})
	if err != nil {
		return nil, err
	}

	return creds, nil
}
