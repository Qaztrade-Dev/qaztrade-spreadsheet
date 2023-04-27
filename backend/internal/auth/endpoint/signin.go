package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/auth/service"
	"github.com/go-kit/kit/endpoint"
)

type SignInRequest struct {
	Email    string
	Password string
}

type SignInResponse struct {
	Err         error               `json:"err,omitempty"`
	Credentials *domain.Credentials `json:"credentials,omitempty"`
}

func (r *SignInResponse) Error() error { return r.Err }

func MakeSignInEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignInRequest)

		credentials, err := s.SignIn(ctx, &service.SignInRequest{
			Email:    req.Email,
			Password: req.Password,
		})

		response := &SignInResponse{Credentials: credentials, Err: err}

		return response, nil
	}
}
