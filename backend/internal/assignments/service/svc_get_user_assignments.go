package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

type GetUserAssignmentsRequest struct {
	UserID string
	Limit  uint64
	Offset uint64
}

type GetUserAssignmentsResponse struct {
	AssignmentsList *domain.AssignmentsList
	AssignmentsInfo *domain.AssignmentsInfo
}

func (s *service) GetUserAssignments(ctx context.Context, input *GetUserAssignmentsRequest) (*GetUserAssignmentsResponse, error) {
	assignmentsList, err := s.assignmentRepo.GetMany(ctx, &domain.GetManyInput{
		UserID: &input.UserID,
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		return nil, err
	}

	assignmentsInfo, err := s.assignmentRepo.GetInfo(ctx, &domain.GetInfoInput{
		UserID: &input.UserID,
	})
	if err != nil {
		return nil, err
	}

	output := &GetUserAssignmentsResponse{
		AssignmentsList: assignmentsList,
		AssignmentsInfo: assignmentsInfo,
	}

	return output, nil
}
