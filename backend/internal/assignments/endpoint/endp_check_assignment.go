package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type CheckAssignmentRequest struct {
	AssignmentID uint64
}

type CheckAssignmentResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *CheckAssignmentResponse) Error() error { return r.Err }

func MakeCheckAssignmentEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(CheckAssignmentRequest)

		err := s.CheckAssignment(ctx, input.AssignmentID)

		return &CheckAssignmentResponse{
			Err: err,
		}, nil
	}
}
