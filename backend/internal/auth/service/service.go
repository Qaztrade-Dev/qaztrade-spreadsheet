package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
)

type Service interface {
	SignUp(ctx context.Context, req *SignUpRequest) (*domain.Credentials, error)
	SignIn(ctx context.Context, req *SignInRequest) (*domain.Credentials, error)
	Forgot(ctx context.Context, req *ForgotRequest) error
	Restore(ctx context.Context, req *RestoreRequest) error
}

type service struct {
	authRepo  domain.AuthorizationRepository
	credsRepo domain.CredentialsRepository
	emailSvc  domain.EmailService
}

func NewService() Service {
	return &service{}
}
