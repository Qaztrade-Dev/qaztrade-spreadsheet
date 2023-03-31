package domain

import "context"

type Credentials struct {
	AccessToken string
}

type SignUpInput struct {
	Email    string
	Password string
	OrgName  string
}

type SignInInput struct {
	Email    string
	Password string
}

type GetQuery struct {
	Email string
}

type AuthorizationRepository interface {
	SignUp(ctx context.Context, input *SignUpInput) (userID string, err error)
	SignIn(ctx context.Context, input *SignInInput) (userID string, err error)
	UpdatePassword(ctx context.Context, userID, password string) error
	GetOne(ctx context.Context, input *GetQuery) (userID string, err error)
}

type CredentialsRepository interface {
	Create(ctx context.Context, userID string) (*Credentials, error)
}

type EmailService interface {
	Send(ctx context.Context, mailName string, payload interface{}) error
}
