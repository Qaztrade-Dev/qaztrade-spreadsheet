package service

import (
	"fmt"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

func (s *service) EnqueueAssignments(ctx context.Context) error {
	isCompleted := false

	assignments, err := s.assignmentRepo.GetMany(ctx, &domain.GetManyInput{
		IsCompleted: &isCompleted,
	})
	if err != nil {
		return fmt.Errorf("assignmentRepo.GetMany: %w", err)
	}

	assignmentIDs := make([]uint64, 0, assignments.Total)
	for _, assignment := range assignments.Objects {
		assignmentIDs = append(assignmentIDs, assignment.ID)
	}

	if err := s.publisher.Publish(ctx, assignmentIDs...); err != nil {
		return fmt.Errorf("publisher.Publish: %w", err)
	}

	return nil
}
