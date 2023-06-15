package endpoint

import (
	"context"
	"errors"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
)

type ctxKey int

const ctxKeyAuth ctxKey = iota

var (
	ErrorEmptyAuthorization = errors.New("unauthorized")
	ErrorPermissionDenied   = errors.New("permission denied")
)

func AuthManagerMiddleware(j *jwt.Client) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			token, ok := ctx.Value(ctxKeyAuth).(string)
			if !ok {
				return nil, ErrorEmptyAuthorization
			}

			claims, err := jwt.Parse[authDomain.UserClaims](j, token)
			if err != nil {
				return nil, err
			}

			if claims.Role != authDomain.RoleManager {
				return nil, ErrorPermissionDenied
			}

			return next(ctx, request)
		}
	}
}

func WithToken(ctx context.Context, input string) context.Context {
	return context.WithValue(ctx, ctxKeyAuth, input)
}
