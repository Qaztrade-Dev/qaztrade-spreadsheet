package service

import (
	"fmt"
	"strconv"

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

	assignmentIDs := make([][]byte, 0, assignments.Total)
	for _, assignment := range assignments.Objects {
		payload := strconv.FormatUint(assignment.AssignmentID, 10)
		assignmentIDs = append(assignmentIDs, []byte(payload))
	}

	if err := s.publisher.Publish(ctx, assignmentIDs...); err != nil {
		return fmt.Errorf("publisher.Publish: %w", err)
	}

	return nil
}
