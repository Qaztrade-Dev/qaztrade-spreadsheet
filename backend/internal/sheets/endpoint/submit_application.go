package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/pkg/jwt"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	"github.com/go-kit/kit/endpoint"
)

type SubmitApplicationRequest struct {
	Application *domain.Application
}

type SubmitApplicationResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *SubmitApplicationResponse) Error() error { return r.Err }

func MakeSubmitApplicationEndpoint(s service.Service, j *jwt.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SubmitApplicationRequest)

		claims, err := j.Parse(req.Application.Token)
		if err != nil {
			return nil, err
		}

		err = s.SubmitApplication(ctx, &service.SubmitApplicationRequest{
			SpreadsheetID: claims.SpreadsheetID,
			Application:   req.Application,
		})

		return SubmitApplicationResponse{Err: err}, nil
	}
}