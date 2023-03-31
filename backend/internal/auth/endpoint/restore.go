package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/auth/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
)

type RestoreRequest struct {
	AccessToken string
	Password    string
}

type RestoreResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *RestoreResponse) Error() error { return r.Err }

func MakeRestoreEndpoint(s service.Service, j *jwt.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RestoreRequest)

		claims, err := jwt.Parse[domain.UserClaims](j, req.AccessToken)
		if err != nil {
			return nil, err
		}

		err = s.Restore(ctx, &service.RestoreRequest{
			UserID:   claims.UserID,
			Password: req.Password,
		})
		return RestoreResponse{Err: err}, nil
	}
}
