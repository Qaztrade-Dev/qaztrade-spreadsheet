package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

func (s *service) RedistributeAssignments(ctx context.Context, assignmentType string) error {
	legalAssignments, err := s.distributeAssignments(ctx, domain.TypeLegal)
	if err != nil {
		return err
	}

	financeAssignments, err := s.distributeAssignments(ctx, domain.TypeFinance)
	if err != nil {
		return err
	}

	digitalAssignment, err := s.distributeDigital(ctx, legalAssignments, financeAssignments)
	if err != nil {
		return err
	}

	if err := s.assignmentRepo.UpdateAssignees(ctx, legalAssignments); err != nil {
		return err
	}

	if err := s.assignmentRepo.UpdateAssignees(ctx, financeAssignments); err != nil {
		return err
	}

	if err := s.assignmentRepo.UpdateAssignees(ctx, digitalAssignment); err != nil {
		return err
	}

	return nil
}

func (s *service) distributeDigital(ctx context.Context, legal, finance []*domain.AssignmentInput) ([]*domain.AssignmentInput, error) {
	assignmentType := domain.TypeDigital

	digital, err := s.assignmentRepo.GetMany(ctx, &domain.GetManyInput{
		AssignmentType: &assignmentType,
	})
	if err != nil {
		return nil, err
	}

	var (
		legalMap   = toMap(legal)
		financeMap = toMap(finance)
	)

	managerID := ""
	result := make([]*domain.AssignmentInput, 0, len(digital.Objects))
	for i := 0; i < len(digital.Objects); i++ {
		assignment := digital.Objects[i]

		if _, ok := legalMap[assignment.ApplicationID]; !ok {
			continue
		}

		if i%2 == 0 {
			managerID = legalMap[assignment.ApplicationID].ManagerID
		} else {
			managerID = financeMap[assignment.ApplicationID].ManagerID
		}

		result = append(result, &domain.AssignmentInput{
			AssignmentID: assignment.AssignmentID,
			ManagerID:    managerID,
		})
	}

	return result, nil
}

func toMap(assignments []*domain.AssignmentInput) map[string]*domain.AssignmentInput {
	assignmentsMap := make(map[string]*domain.AssignmentInput, len(assignments))

	for _, assignment := range assignments {
		assignment := assignment
		assignmentsMap[assignment.ApplicationID] = assignment
	}

	return assignmentsMap
}

func (s *service) distributeAssignments(ctx context.Context, assignmentType string) ([]*domain.AssignmentInput, error) {
	list, err := s.assignmentRepo.GetMany(ctx, &domain.GetManyInput{
		AssignmentType: &assignmentType,
	})
	if err != nil {
		return nil, err
	}

	managerIDs, err := s.assignmentRepo.GetManagerIDs(ctx, assignmentType)
	if err != nil {
		return nil, err
	}

	var (
		managersMap = make(map[string]*domain.Manager)
		sheets      = make([]*domain.Sheet, 0, len(list.Objects))
	)

	for _, assignment := range list.Objects {
		if assignment.RowsCompleted == 0 {
			sheets = append(sheets, &domain.Sheet{
				AssignmentID:  assignment.AssignmentID,
				ApplicationID: assignment.ApplicationID,
				TotalRows:     uint64(assignment.TotalRows),
				TotalSum:      float64(assignment.TotalSum),
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
			AssignmentID:  assignment.AssignmentID,
			ApplicationID: assignment.ApplicationID,
			TotalRows:     uint64(assignment.TotalRows),
			TotalSum:      float64(assignment.TotalSum),
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
				AssignmentID:  sheet.AssignmentID,
				ApplicationID: sheet.ApplicationID,
				ManagerID:     managerIDs[i],
			})
		}
	}

	return assignments, nil
}
