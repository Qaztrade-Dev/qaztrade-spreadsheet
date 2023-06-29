package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/auth/service"
	"github.com/go-kit/kit/endpoint"
)

type SignUpRequest struct {
	Email    string
	Password string
	OrgName  string
}

type SignUpResponse struct {
	Err         error               `json:"err,omitempty"`
	Credentials *domain.Credentials `json:"credentials,omitempty"`
}

func (r *SignUpResponse) Error() error { return r.Err }

func MakeSignUpEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)

		credentials, err := s.SignUp(ctx, &service.SignUpRequest{
			Email:    req.Email,
			Password: req.Password,
			OrgName:  req.OrgName,
		})
		return &SignUpResponse{Credentials: credentials, Err: err}, nil
	}
}
