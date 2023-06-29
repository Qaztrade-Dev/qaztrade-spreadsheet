package domain

import "context"

type UserClaims struct {
	UserID string   `json:"uid"`
	Roles  []string `json:"r,omitempty"`
}

const (
	RoleUser    = "user"
	RoleManager = "manager"
	RoleAdmin   = "admin"
	RoleDigital = "digital"
	RoleFinance = "finance"
	RoleLegal   = "legal"
)

type User struct {
	ID string
}

type SignUpInput struct {
	UserID   string
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
	SignUp(ctx context.Context, input *SignUpInput) error
	SignIn(ctx context.Context, input *SignInInput) (*User, error)
	UpdatePassword(ctx context.Context, userID, password string) error
	GetOne(ctx context.Context, input *GetQuery) (*User, error)
	GetRoles(ctx context.Context, userID string) ([]string, error)
}

type Credentials struct {
	AccessToken string `json:"access_token"`
}

type CredentialsRepository interface {
	Create(ctx context.Context, claims *UserClaims) (*Credentials, error)
}

type EmailService interface {
	Send(ctx context.Context, toEmail, mailName string, payload interface{}) error
}
