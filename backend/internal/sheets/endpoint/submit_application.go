package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	"github.com/go-kit/kit/endpoint"
)

type SubmitApplicationRequest struct {
	SpreadsheetID string
	Application   *domain.Application
}

type SubmitApplicationResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *SubmitApplicationResponse) Error() error { return r.Err }

func MakeSubmitApplicationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SubmitApplicationRequest)
		err := s.SubmitApplication(ctx, &service.SubmitApplicationRequest{
			SpreadsheetID: req.SpreadsheetID,
			Application:   req.Application,
		})
		return SubmitApplicationResponse{Err: err}, nil
	}
}
