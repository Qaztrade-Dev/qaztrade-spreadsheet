package service

import (
	"context"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/google/uuid"
)

type SignUpRequest struct {
	Email    string
	Password string
	OrgName  string
}

func (s *service) SignUp(ctx context.Context, input *SignUpRequest) (*domain.Credentials, error) {
	var (
		userID = uuid.NewString()
		email  = strings.TrimSpace(input.Email)
	)

	err := s.authRepo.SignUp(ctx, &domain.SignUpInput{
		UserID:   userID,
		Email:    email,
		Password: input.Password,
		OrgName:  input.OrgName,
	})
	if err != nil {
		return nil, err
	}

	creds, err := s.credsRepo.Create(ctx, &domain.UserClaims{
		UserID: userID,
		Roles:  []string{domain.RoleUser},
	})
	if err != nil {
		return nil, err
	}

	return creds, nil
}
