package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
)

func MakeClaimsMiddleware[T any](jc *jwt.Client) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			token, err := domain.ExtractToken(ctx)
			if err != nil {
				return nil, err
			}

			claims, err := jwt.Parse[T](jc, token)
			if err != nil {
				return nil, err
			}

			newCtx := domain.WithClaims(ctx, claims)

			return next(newCtx, request)
		}
	}
}

func MakeAuthMiddleware(requiredRoles ...string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			claims, err := domain.ExtractClaims[domain.UserClaims](ctx)
			if err != nil {
				return nil, err
			}

			if !mustContainRoles(claims.Roles, requiredRoles) {
				return nil, domain.ErrorPermissionDenied
			}

			return next(ctx, request)
		}
	}
}

func mustContainRoles(checkRoles, requiredRoles []string) bool {
	for _, requiredRole := range requiredRoles {
		found := false
		for _, checkRole := range checkRoles {
			if checkRole == requiredRole {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
