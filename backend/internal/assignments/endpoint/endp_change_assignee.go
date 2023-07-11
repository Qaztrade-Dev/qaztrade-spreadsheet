package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type ChangeAssigneeRequest struct {
	UserID       string
	AssignmentID uint64
}

type ChangeAssigneeResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *ChangeAssigneeResponse) Error() error { return r.Err }

func MakeChangeAssigneeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(ChangeAssigneeRequest)

		err := s.ChangeAssignee(ctx, &service.ChangeAssigneeRequest{
			UserID:       input.UserID,
			AssignmentID: input.AssignmentID,
		})

		return &ChangeAssigneeResponse{
			Err: err,
		}, nil
	}
}
