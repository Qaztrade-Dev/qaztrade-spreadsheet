package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/doodocs/qaztrade/backend/internal/assignments/pkg/jsondomain"
	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type GetAssignmentsResponse struct {
	Err             error                       `json:"err,omitempty"`
	AssignmentsList *jsondomain.AssignmentsList `json:"assignments_list"`
	AssignmentsInfo *jsondomain.AssignmentsInfo `json:"assignments_info"`
}

func (r *GetAssignmentsResponse) Error() error { return r.Err }

func MakeGetAssignmentsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(domain.GetManyInput)

		response, err := s.GetAssignments(ctx, &input)
		if err != nil {
			return &GetAssignmentsResponse{Err: err}, nil
		}

		return &GetAssignmentsResponse{
			Err:             err,
			AssignmentsList: jsondomain.EncodeAssignmentsList(response.AssignmentsList),
			AssignmentsInfo: jsondomain.EncodeAssignmentsInfo(response.AssignmentsInfo),
		}, nil
	}
}
