package service

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
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

	batchID, err := s.assignmentRepo.LockApplications(ctx)
	if err != nil {
		return fmt.Errorf("LockApplications: %w", err)
	}

	financeLogisticsAssignments, err := s.createAssignments(ctx, batchID, domain.SheetsЗатратыНаДоставкуТранспортом, domain.TypeFinance)
	if err != nil {
		return fmt.Errorf("createAssignments: %w", err)
	}

	legalLogisticsAssignments, err := s.createAssignments(ctx, batchID, domain.SheetsЗатратыНаДоставкуТранспортом, domain.TypeLegal)
	if err != nil {
		return fmt.Errorf("createAssignments: %w", err)
	}

	digitalAssignments := s.getDigitalAssignments(ctx, financeLogisticsAssignments, legalLogisticsAssignments)

	assignments := make([]*domain.AssignmentInput, 0)
	assignments = append(assignments, financeLogisticsAssignments...)
	assignments = append(assignments, legalLogisticsAssignments...)
	assignments = append(assignments, digitalAssignments...)

	for _, sheetType := range []string{
		domain.SheetsЗатратыНаСертификациюПредприятия,
		domain.SheetsЗатратыНаРекламуИкуЗаРубежом,
		domain.SheetsЗатратыНаПереводКаталогаИку,
		domain.SheetsЗатратыНаАрендуПомещенияИку,
		domain.SheetsЗатратыНаСертификациюИку,
		domain.SheetsЗатратыНаДемонстрациюИку,
		domain.SheetsЗатратыНаФранчайзинг,
		domain.SheetsЗатратыНаРегистрациюТоварныхЗнаков,
		domain.SheetsЗатратыНаАренду,
		domain.SheetsЗатратыНаПеревод,
		domain.SheetsЗатратыНаРекламуТоваровЗаРубежом,
		domain.SheetsЗатратыНаУчастиеВВыставках,
		domain.SheetsЗатратыНаУчастиеВВыставкахИку,
		domain.SheetsЗатратыНаСоответствиеТоваровТребованиям,
	} {
		financeAssignments, err := s.createAssignments(ctx, batchID, sheetType, domain.TypeFinance)
		if err != nil {
			return fmt.Errorf("createAssignments: %w", err)
		}

		legalAssignments, err := s.createAssignments(ctx, batchID, sheetType, domain.TypeLegal)
		if err != nil {
			return fmt.Errorf("createAssignments: %w", err)
		}

		assignments = append(assignments, financeAssignments...)
		assignments = append(assignments, legalAssignments...)
	}

	if err := s.assignmentRepo.CreateAssignments(ctx, assignments); err != nil {
		return fmt.Errorf("CreateAssignments: %w", err)
	}

	return nil
}

func (s *service) createAssignments(ctx context.Context, batchID int, sheetType, assignmentType string) ([]*domain.AssignmentInput, error) {
	managerIDs, err := s.assignmentRepo.GetManagerIDs(ctx, assignmentType)
	if err != nil {
		return nil, fmt.Errorf("GetManagerIDs: %w", err)
	}

	if len(managerIDs) == 0 {
		return nil, domain.ErrorEmptyManagers
	}

	sheets, err := s.assignmentRepo.GetSheets(ctx, batchID, sheetType)
	if err != nil {
		return nil, fmt.Errorf("GetSheets: %w", err)
	}

	if len(sheets) == 0 {
		return []*domain.AssignmentInput{}, nil
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(managerIDs), func(i, j int) { managerIDs[i], managerIDs[j] = managerIDs[j], managerIDs[i] })

	assignees := domain.DistributeAdvanced(len(managerIDs), sheets)

	assignments := make([]*domain.AssignmentInput, 0, len(sheets))
	for i, assignee := range assignees.Managers {
		for _, sheet := range assignee.Sheets {
			assignments = append(assignments, &domain.AssignmentInput{
				ApplicationID:  sheet.ApplicationID,
				SheetTitle:     sheet.SheetTitle,
				SheetID:        sheet.SheetID,
				AssignmentType: assignmentType,
				ManagerID:      managerIDs[i],
				TotalRows:      sheet.TotalRows,
				TotalSum:       sheet.TotalSum,
			})
		}
	}

	sort.SliceStable(assignments, func(i, j int) bool {
		return strings.Compare(assignments[i].ApplicationID, assignments[j].ApplicationID) < 0
	})

	return assignments, nil
}

func (s *service) getDigitalAssignments(ctx context.Context, financeAssignments, legalAssignments []*domain.AssignmentInput) []*domain.AssignmentInput {
	digitalAssignments := make([]*domain.AssignmentInput, 0, len(financeAssignments))
	for i, sheet := range financeAssignments {

		managerID := sheet.ManagerID
		if i%2 == 0 {
			managerID = legalAssignments[i].ManagerID
		}

		digitalAssignments = append(digitalAssignments, &domain.AssignmentInput{
			ApplicationID:  sheet.ApplicationID,
			SheetTitle:     sheet.SheetTitle,
			SheetID:        sheet.SheetID,
			AssignmentType: domain.TypeDigital,
			ManagerID:      managerID,
			TotalRows:      sheet.TotalRows,
			TotalSum:       sheet.TotalSum,
		})
	}

	return digitalAssignments
}
