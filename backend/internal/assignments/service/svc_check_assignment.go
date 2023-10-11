package service

import (
	"fmt"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

func (s *service) CheckAssignment(ctx context.Context, assignmentID uint64) error {
	assignment, err := s.assignmentRepo.GetOne(ctx, &domain.GetManyInput{AssignmentID: &assignmentID})
	if err != nil {
		return fmt.Errorf("assignmentRepo.GetOne: %w", err)
	}

	var (
		spreadsheetID = assignment.SpreadsheetID
		asgnType      = assignment.AssignmentType
	)

	sheetTitles := strings.Split(assignment.SheetTitle, ", ")
	for _, sheetTitle := range sheetTitles {
		// На оцифровку только "Затраты на доставку транспортом"
		if asgnType == domain.TypeDigital && sheetTitle != domain.TitleЗатратыНаДоставкуТранспортом {
			return nil
		}

		sheetData, err := s.spreadsheetRepo.GetSheetData(ctx, spreadsheetID, sheetTitle)
		if err != nil {
			return fmt.Errorf("spreadsheetRepo.GetSheetData: %w", err)
		}

		var (
			totalCompleted uint64 = 0
			headerOffset          = 3
			col                   = columnMapping[sheetTitle][asgnType]
		)

		for _, row := range sheetData[headerOffset:] {
			if row[col] != "" {
				totalCompleted++
			}
		}

		if err := s.assignmentRepo.InsertAssignmentResult(ctx, assignmentID, totalCompleted); err != nil {
			return fmt.Errorf("assignmentRepo.InsertAssignmentResult: %w", err)
		}
	}

	return nil
}

var (
	columnMapping = map[string]map[string]int{
		domain.TitleЗатратыНаДоставкуТранспортом: {
			domain.TypeDigital: 136,
			domain.TypeFinance: 130,
			domain.TypeLegal:   135,
		},
		domain.TitleЗатратыНаСертификациюПредприятия: {
			domain.TypeFinance: 60,
			domain.TypeLegal:   64,
		},
		domain.TitleЗатратыНаРекламуИкуЗаРубежом: {
			domain.TypeFinance: 61,
			domain.TypeLegal:   65,
		},
		domain.TitleЗатратыНаПереводКаталогаИку: {
			domain.TypeFinance: 59,
			domain.TypeLegal:   63,
		},
		domain.TitleЗатратыНаАрендуПомещенияИку: {
			domain.TypeFinance: 56,
			domain.TypeLegal:   60,
		},
		domain.TitleЗатратыНаСертификациюИку: {
			domain.TypeFinance: 62,
			domain.TypeLegal:   66,
		},
		domain.TitleЗатратыНаДемонстрациюИку: {
			domain.TypeFinance: 56,
			domain.TypeLegal:   61,
		},
		domain.TitleЗатратыНаФранчайзинг: {
			domain.TypeFinance: 56,
			domain.TypeLegal:   60,
		},
		domain.TitleЗатратыНаРегистрациюТоварныхЗнаков: {
			domain.TypeFinance: 68,
			domain.TypeLegal:   72,
		},
		domain.TitleЗатратыНаАренду: {
			domain.TypeFinance: 56,
			domain.TypeLegal:   60,
		},
		domain.TitleЗатратыНаПеревод: {
			domain.TypeFinance: 59,
			domain.TypeLegal:   63,
		},
		domain.TitleЗатратыНаРекламуТоваровЗаРубежом: {
			domain.TypeFinance: 61,
			domain.TypeLegal:   65,
		},
		domain.TitleЗатратыНаУчастиеВВыставках: {
			domain.TypeFinance: 75,
			domain.TypeLegal:   79,
		},
		domain.TitleЗатратыНаУчастиеВВыставкахИку: {
			domain.TypeFinance: 75,
			domain.TypeLegal:   79,
		},
		domain.TitleЗатратыНаСоответствиеТоваровТребованиям: {
			domain.TypeFinance: 70,
			domain.TypeLegal:   74,
		},
	}
)
