package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/service"
	"github.com/go-kit/kit/endpoint"
)

type ForgotRequest struct {
	Email string
}

type ForgotResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *ForgotResponse) Error() error { return r.Err }

func MakeForgotEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ForgotRequest)

		err := s.Forgot(ctx, &service.ForgotRequest{
			Email: req.Email,
		})
		return ForgotResponse{Err: err}, nil
	}
}
