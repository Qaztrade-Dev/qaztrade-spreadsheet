package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sign/service"
	"github.com/go-kit/kit/endpoint"
)

type SyncSpreadsheetsRequest struct {
	SpreadsheetID string
}

type SyncSpreadsheetsResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *SyncSpreadsheetsResponse) Error() error { return r.Err }

func MakeSyncSpreadsheetsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SyncSpreadsheetsRequest)

		err := s.SyncSpreadsheets(ctx, &service.SyncSpreadsheetsRequest{
			SpreadsheetID: req.SpreadsheetID,
		})

		return &SyncSpreadsheetsResponse{
			Err: err,
		}, nil
	}
}
