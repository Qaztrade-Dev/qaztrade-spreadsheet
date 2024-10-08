package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

type GetAssignmentsResponse struct {
	AssignmentsList *domain.AssignmentsList
	AssignmentsInfo *domain.AssignmentsInfo
}

func (s *service) GetAssignments(ctx context.Context, input *domain.GetManyInput) (*GetAssignmentsResponse, error) {
	assignmentsList, err := s.assignmentRepo.GetMany(ctx, input)
	if err != nil {
		return nil, err
	}

	assignmentsInfo, err := s.assignmentRepo.GetInfo(ctx, &domain.GetInfoInput{})
	if err != nil {
		return nil, err
	}

	output := &GetAssignmentsResponse{
		AssignmentsList: assignmentsList,
		AssignmentsInfo: assignmentsInfo,
	}

	return output, nil
}
