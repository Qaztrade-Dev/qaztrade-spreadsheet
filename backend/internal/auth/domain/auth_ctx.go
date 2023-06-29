package domain

import (
	"context"
	"fmt"

	"github.com/doodocs/qaztrade/backend/pkg/jwt"
)

type ctxKey int

const ctxKeyAuth ctxKey = iota

var (
	ErrorEmptyAuthorization = fmt.Errorf("unauthorized")
	ErrorPermissionDenied   = fmt.Errorf("permission denied")
)

func WithToken(ctx context.Context, input string) context.Context {
	return context.WithValue(ctx, ctxKeyAuth, input)
}

func ExtractToken(ctx context.Context) (string, error) {
	token, ok := ctx.Value(ctxKeyAuth).(string)
	if !ok {
		return "", ErrorEmptyAuthorization
	}
	return token, nil
}

func ExtractClaims(ctx context.Context, jwtCli *jwt.Client) (*UserClaims, error) {
	token, err := ExtractToken(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := jwt.Parse[UserClaims](jwtCli, token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
