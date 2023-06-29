package adapters

import (
	"context"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
)

type CredentialsRepositoryJWT struct {
	jwtcli *jwt.Client
}

var _ domain.CredentialsRepository = (*CredentialsRepositoryJWT)(nil)

func NewCredentialsRepositoryJWT(jwtcli *jwt.Client) *CredentialsRepositoryJWT {
	return &CredentialsRepositoryJWT{
		jwtcli: jwtcli,
	}
}

func (r *CredentialsRepositoryJWT) Create(ctx context.Context, claims *domain.UserClaims) (*domain.Credentials, error) {
	var (
		expireAt = time.Now().Add(time.Duration(72 * time.Hour)) // 3 days
	)

	accessToken, err := jwt.NewTokenString(r.jwtcli, claims, jwt.WithExpire(expireAt))
	if err != nil {
		return nil, err
	}

	creds := &domain.Credentials{
		AccessToken: accessToken,
	}

	return creds, nil
}
