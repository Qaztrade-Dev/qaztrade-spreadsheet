package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/auth/service"
	"github.com/go-kit/kit/endpoint"
)

type RestoreRequest struct {
	Password string
}

type RestoreResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *RestoreResponse) Error() error { return r.Err }

func MakeRestoreEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RestoreRequest)

		claims, err := domain.ExtractClaims[domain.UserClaims](ctx)
		if err != nil {
			return nil, err
		}

		err = svc.Restore(ctx, &service.RestoreRequest{
			UserID:   claims.UserID,
			Password: req.Password,
		})
		return &RestoreResponse{Err: err}, nil
	}
}
