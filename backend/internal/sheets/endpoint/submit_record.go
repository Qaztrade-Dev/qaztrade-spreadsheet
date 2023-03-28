package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	"github.com/go-kit/kit/endpoint"
)

type SubmitRecordRequest struct {
	SpreadsheetID string
	Payload       *domain.Payload
}

type SubmitRecordResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *SubmitRecordResponse) Error() error { return r.Err }

func MakeSubmitRecordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SubmitRecordRequest)
		err := s.SubmitRecord(ctx, &service.SubmitRecordRequest{
			SpreadsheetID: req.SpreadsheetID,
			Payload:       req.Payload,
		})
		return SubmitRecordResponse{Err: err}, nil
	}
}
