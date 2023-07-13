package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

func (s *service) CreateBatch(ctx context.Context) error {
	/*

		1. Lock signed applications into batch
		2. Get sheets for dostavka
		3. Get managers with role=digital
		4. Call simple distribute
		5. Create assignments
		6. Update batch step=1

	*/

	managerIDs, err := s.assignmentRepo.GetManagerIDs(ctx, domain.TypeDigital)
	if err != nil {
		return fmt.Errorf("GetManagerIDs: %w", err)
	}

	if len(managerIDs) == 0 {
		return domain.ErrorEmptyManagers
	}

	batchID, err := s.assignmentRepo.LockApplications(ctx)
	if err != nil {
		return fmt.Errorf("LockApplications: %w", err)
	}

	sheets, err := s.assignmentRepo.GetSheets(ctx, batchID, domain.SheetsЗатратыНаДоставкуТранспортом)
	if err != nil {
		return fmt.Errorf("GetSheets: %w", err)
	}

	if len(sheets) == 0 {
		return domain.ErrorEmptySheets
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(managerIDs), func(i, j int) { managerIDs[i], managerIDs[j] = managerIDs[j], managerIDs[i] })

	assignees := domain.DistributeSimple(len(managerIDs), sheets)

	assignments := make([]*domain.AssignmentInput, 0, len(sheets))
	for i, assignee := range assignees.Managers {
		for _, sheet := range assignee.Sheets {
			assignments = append(assignments, &domain.AssignmentInput{
				ApplicationID:  sheet.ApplicationID,
				SheetTitle:     sheet.SheetTitle,
				SheetID:        sheet.SheetID,
				AssignmentType: domain.TypeDigital,
				ManagerID:      managerIDs[i],
				TotalRows:      sheet.TotalRows,
				TotalSum:       sheet.TotalSum,
			})
		}
	}

	if err := s.assignmentRepo.CreateAssignments(ctx, assignments); err != nil {
		return fmt.Errorf("CreateAssignments: %w", err)
	}

	if err := s.assignmentRepo.UpdateBatchStep(ctx, batchID, 1); err != nil {
		return fmt.Errorf("UpdateBatchStep: %w", err)
	}

	return nil
}
