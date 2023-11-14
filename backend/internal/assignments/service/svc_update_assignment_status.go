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
	if input.StatusName == domain.ResolutionStatusRejected {
		assignment, err := s.assignmentRepo.GetOne(ctx, &domain.GetManyInput{
			AssignmentID: &input.AssignmentID,
		})
		if err != nil {
			return err
		}

		if err := s.applicationRepo.EditStatus(ctx, assignment.ApplicationID, domain.ResolutionStatusRejected); err != nil {
			return err
		}
	}

	return s.assignmentRepo.UpdateStatus(ctx, &domain.UpdateStatusInput{
		AssignmentID: input.AssignmentID,
		StatusName:   input.StatusName,
	})
}
