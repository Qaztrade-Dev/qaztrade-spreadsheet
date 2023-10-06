package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type RedistributeAssignmentsRequest struct {
	AssignmentType string
}

type RedistributeAssignmentsResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *RedistributeAssignmentsResponse) Error() error { return r.Err }

func MakeRedistributeAssignmentsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(RedistributeAssignmentsRequest)

		err := s.RedistributeAssignments(ctx, input.AssignmentType)

		return &RedistributeAssignmentsResponse{
			Err: err,
		}, nil
	}
}
