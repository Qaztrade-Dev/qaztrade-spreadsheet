package endpoint

import (
	"context"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type GetUserAssignmentsRequest struct {
	Limit  int
	Offset int
}

type GetUserAssignmentsResponse struct {
	Err             error                   `json:"err,omitempty"`
	AssignmentsList *domain.AssignmentsList `json:"assignments_list"`
	AssignmentsInfo *domain.AssignmentsInfo `json:"assignments_info"`
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
			AssignmentsList: response.AssignmentsList,
			AssignmentsInfo: response.AssignmentsInfo,
		}, nil
	}
}
