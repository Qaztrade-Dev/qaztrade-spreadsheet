package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type EnqueueAssignmentsResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *EnqueueAssignmentsResponse) Error() error { return r.Err }

func MakeEnqueueAssignmentsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := s.EnqueueAssignments(ctx)

		return &EnqueueAssignmentsResponse{
			Err: err,
		}, nil
	}
}
