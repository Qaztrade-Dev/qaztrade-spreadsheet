package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

func (s *service) RedistributeAssignments(ctx context.Context, assignmentType string) error {
	list, err := s.assignmentRepo.GetMany(ctx, &domain.GetManyInput{
		AssignmentType: &assignmentType,
	})
	if err != nil {
		return err
	}

	managerIDs, err := s.assignmentRepo.GetManagerIDs(ctx, assignmentType)
	if err != nil {
		return err
	}

	var (
		managersMap = make(map[string]*domain.Manager)
		sheets      = make([]*domain.Sheet, 0, len(list.Objects))
	)

	for _, assignment := range list.Objects {
		if assignment.RowsCompleted == 0 {
			sheets = append(sheets, &domain.Sheet{
				AssignmentID: assignment.AssignmentID,
				TotalRows:    uint64(assignment.TotalRows),
				TotalSum:     float64(assignment.TotalSum),
			})
			continue
		}

		manager, ok := managersMap[assignment.AssigneeID]
		if !ok {
			manager = &domain.Manager{
				UserID: assignment.AssigneeID,
			}
			managersMap[assignment.AssigneeID] = manager
		}

		manager.TotalRows += uint64(assignment.TotalRows)
		manager.TotalSum += float64(assignment.TotalSum)

		manager.Sheets = append(manager.Sheets, &domain.Sheet{
			AssignmentID: assignment.AssignmentID,
			TotalRows:    uint64(assignment.TotalRows),
			TotalSum:     float64(assignment.TotalSum),
		})
	}

	for _, managerID := range managerIDs {
		if _, ok := managersMap[managerID]; ok {
			continue
		}

		managersMap[managerID] = &domain.Manager{
			UserID: managerID,
		}
	}

	managers := make([]*domain.Manager, 0, len(managersMap))
	for _, manager := range managersMap {
		manager := manager
		managers = append(managers, manager)
	}

	var assignees domain.Queue
	if assignmentType == domain.TypeDigital {
		assignees = domain.DistributeSimpleWithData(managers, sheets)
	} else {
		assignees = domain.DistributeAdvancedWithData(managers, sheets)
	}

	assignments := make([]*domain.AssignmentInput, 0, len(sheets))
	for i, assignee := range assignees.GetManagers() {
		for _, sheet := range assignee.Sheets {
			assignments = append(assignments, &domain.AssignmentInput{
				AssignmentID: sheet.AssignmentID,
				ManagerID:    managerIDs[i],
			})
		}
	}

	if err := s.assignmentRepo.UpdateAssignees(ctx, assignments); err != nil {
		return err
	}

	return nil
}
