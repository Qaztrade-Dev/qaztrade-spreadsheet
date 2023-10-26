package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type UpdateAssignmentStatusRequest struct {
	AssignmentID uint64
	StatusName   string
}

type UpdateAssignmentStatusResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *UpdateAssignmentStatusResponse) Error() error { return r.Err }

func MakeUpdateAssignmentStatusEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(UpdateAssignmentStatusRequest)

		err := s.UpdateAssignmentStatus(ctx, &service.UpdateAssignmentStatusRequest{
			AssignmentID: input.AssignmentID,
			StatusName:   input.StatusName,
		})

		return &UpdateAssignmentStatusResponse{
			Err: err,
		}, nil
	}
}
