package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type CreateBatchResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *CreateBatchResponse) Error() error { return r.Err }

func MakeCreateBatchEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.CreateBatch(ctx)

		return &CreateBatchResponse{
			Err: err,
		}, nil
	}
}
