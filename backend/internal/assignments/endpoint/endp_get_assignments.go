package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/pkg/jsondomain"
	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type GetAssignmentsRequest struct {
	Limit  uint64
	Offset uint64
	UserID *string
}

type GetAssignmentsResponse struct {
	Err             error                       `json:"err,omitempty"`
	AssignmentsList *jsondomain.AssignmentsList `json:"assignments_list"`
	AssignmentsInfo *jsondomain.AssignmentsInfo `json:"assignments_info"`
}

func (r *GetAssignmentsResponse) Error() error { return r.Err }

func MakeGetAssignmentsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(GetAssignmentsRequest)

		response, err := s.GetAssignments(ctx, &service.GetAssignmentsRequest{
			UserID: input.UserID,
			Limit:  input.Limit,
			Offset: input.Offset,
		})

		return &GetAssignmentsResponse{
			Err:             err,
			AssignmentsList: jsondomain.EncodeAssignmentsList(response.AssignmentsList),
			AssignmentsInfo: jsondomain.EncodeAssignmentsInfo(response.AssignmentsInfo),
		}, nil
	}
}
