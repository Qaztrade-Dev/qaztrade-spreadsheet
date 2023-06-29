package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

type GetAssignmentsInput struct {
	Limit  int
	Offset int
}

type GetAssignmentsOutput struct {
	AssignmentsList *domain.AssignmentsList
	AssignmentsInfo *domain.AssignmentsInfo
}

func (s *service) GetAssignments(ctx context.Context, input *GetAssignmentsInput) (*GetAssignmentsOutput, error) {
	assignmentsList, err := s.assignmentRepo.GetMany(ctx, &domain.GetManyInput{
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		return nil, err
	}

	assignmentsInfo, err := s.assignmentRepo.GetInfo(ctx, &domain.GetInfoInput{})
	if err != nil {
		return nil, err
	}

	output := &GetAssignmentsOutput{
		AssignmentsList: assignmentsList,
		AssignmentsInfo: assignmentsInfo,
	}

	return output, nil
}
