package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

type GetUserAssignmentsResponse struct {
	AssignmentsList *domain.AssignmentsList
	AssignmentsInfo *domain.AssignmentsInfo
}

func (s *service) GetUserAssignments(ctx context.Context, input *domain.GetManyInput) (*GetUserAssignmentsResponse, error) {
	assignmentsList, err := s.assignmentRepo.GetMany(ctx, input)
	if err != nil {
		return nil, err
	}

	assignmentsInfo, err := s.assignmentRepo.GetInfo(ctx, &domain.GetInfoInput{
		UserID: input.AssigneeID,
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
