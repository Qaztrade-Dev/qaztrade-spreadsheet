package sheets

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	"github.com/go-kit/kit/endpoint"
)

type submitRecordRequest struct {
	SpreadsheetID string
	Payload       *domain.Payload
}

type submitRecordResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *submitRecordResponse) error() error { return r.Err }

func makeSubmitRecordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitRecordRequest)
		err := s.SubmitRecord(ctx, &service.SubmitRecordRequest{
			SpreadsheetID: req.SpreadsheetID,
			Payload:       req.Payload,
		})
		return submitRecordResponse{Err: err}, nil
	}
}
