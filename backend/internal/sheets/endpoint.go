package sheets

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
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

func makeSubmitRecordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitRecordRequest)
		err := s.SubmitRecord(ctx, &SubmitRecordRequest{
			SpreadsheetID: req.SpreadsheetID,
			Payload:       req.Payload,
		})
		return submitRecordResponse{Err: err}, nil
	}
}
