package domain

import "context"

type UserClaims struct {
	UserID string `json:"uid"`
	Role   string `json:"r,omitempty"`
}

const (
	RoleUser    = "user"
	RoleManager = "manager"
)

type User struct {
	ID   string
	Role string
}

type Credentials struct {
	AccessToken string `json:"access_token"`
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
	SignUp(ctx context.Context, input *SignUpInput) (*User, error)
	SignIn(ctx context.Context, input *SignInInput) (*User, error)
	UpdatePassword(ctx context.Context, userID, password string) error
	GetOne(ctx context.Context, input *GetQuery) (*User, error)
}

type CredentialsRepository interface {
	Create(ctx context.Context, user *User) (*Credentials, error)
}

type EmailService interface {
	Send(ctx context.Context, toEmail, mailName string, payload interface{}) error
}
