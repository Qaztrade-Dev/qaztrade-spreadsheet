package endpoint

import (
	"context"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"

	"github.com/doodocs/qaztrade/backend/internal/assignments/pkg/jsondomain"
	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type GetUserAssignmentsRequest struct {
	Limit  uint64
	Offset uint64
}

type GetUserAssignmentsResponse struct {
	Err             error                       `json:"err,omitempty"`
	AssignmentsList *jsondomain.AssignmentsList `json:"assignments_list"`
	AssignmentsInfo *jsondomain.AssignmentsInfo `json:"assignments_info"`
}

func (r *GetUserAssignmentsResponse) Error() error { return r.Err }

func MakeGetUserAssignmentsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(GetUserAssignmentsRequest)

		claims, err := authDomain.ExtractClaims[authDomain.UserClaims](ctx)
		if err != nil {
			return nil, err
		}

		response, err := s.GetUserAssignments(ctx, &service.GetUserAssignmentsRequest{
			UserID: claims.UserID,
			Limit:  input.Limit,
			Offset: input.Offset,
		})

		return &GetUserAssignmentsResponse{
			Err:             err,
			AssignmentsList: jsondomain.EncodeAssignmentsList(response.AssignmentsList),
			AssignmentsInfo: jsondomain.EncodeAssignmentsInfo(response.AssignmentsInfo),
		}, nil
	}
}
