package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

type UpdateAssignmentStatusRequest struct {
	AssignmentID uint64
	StatusName   string
}

func (s *service) UpdateAssignmentStatus(ctx context.Context, input *UpdateAssignmentStatusRequest) error {
	return s.assignmentRepo.UpdateStatus(ctx, &domain.UpdateStatusInput{
		AssignmentID: input.AssignmentID,
		StatusName:   input.StatusName,
	})
}
