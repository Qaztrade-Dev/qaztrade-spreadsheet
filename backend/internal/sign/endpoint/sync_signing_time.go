package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sign/service"
	"github.com/go-kit/kit/endpoint"
)

type SyncSigningTimeRequest struct {
	DocumentID string
}

type SyncSigningTimeResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *SyncSigningTimeResponse) Error() error { return r.Err }

func MakeSyncSigningTimeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SyncSigningTimeRequest)

		err := s.SyncSigningTime(ctx, &service.SyncSigningTimeRequest{
			SignDocumentID: req.DocumentID,
		})

		return &SyncSigningTimeResponse{
			Err: err,
		}, nil
	}
}
