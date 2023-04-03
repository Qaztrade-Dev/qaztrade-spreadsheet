package domain

import (
	"context"

	"golang.org/x/oauth2"
)

type TokenRepository interface {
	UpdateToken(ctx context.Context, token *oauth2.Token) error
}
