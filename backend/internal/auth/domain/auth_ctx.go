package domain

import (
	"context"
	"fmt"
)

type ctxKey int

const (
	ctxToken ctxKey = iota
	ctxClaims
)

var (
	ErrorEmptyToken       = fmt.Errorf("unauthorized token")
	ErrorEmptyClaims      = fmt.Errorf("unauthorized claims")
	ErrorPermissionDenied = fmt.Errorf("permission denied")
)

func WithToken(ctx context.Context, input string) context.Context {
	return context.WithValue(ctx, ctxToken, input)
}

func ExtractToken(ctx context.Context) (string, error) {
	token, ok := ctx.Value(ctxToken).(string)
	if !ok {
		return "", ErrorEmptyToken
	}
	return token, nil
}

func WithClaims[T any](ctx context.Context, claims *T) context.Context {
	return context.WithValue(ctx, ctxToken, claims)
}

func ExtractClaims[T any](ctx context.Context) (*T, error) {
	claims, ok := ctx.Value(ctxToken).(*T)
	if !ok {
		return nil, ErrorEmptyClaims
	}
	return claims, nil
}
