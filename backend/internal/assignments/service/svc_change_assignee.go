package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

type ChangeAssigneeRequest struct {
	UserID       string
	AssignmentID uint64
}

func (s *service) ChangeAssignee(ctx context.Context, input *ChangeAssigneeRequest) error {
	if err := s.assignmentRepo.ChangeAssignee(ctx, &domain.ChangeAssigneeInput{
		UserID:       input.UserID,
		AssignmentID: input.AssignmentID,
	}); err != nil {
		return err
	}

	return nil
}
