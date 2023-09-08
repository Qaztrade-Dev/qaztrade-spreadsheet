package adapters

import (
	"context"
	"fmt"
	"strings"

	excelize "github.com/xuri/excelize/v2"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetServiceGoogle struct {
	driveSvc     *drive.Service
	sheetsSvc    *sheets.Service
	adminAccount string
	svcAccount   string
}

var _ domain.SpreadsheetService = (*SpreadsheetServiceGoogle)(nil)

func NewSpreadsheetService(ctx context.Context, credentialsJson []byte, adminAccount, svcAccount string) (*SpreadsheetServiceGoogle, error) {
	driveSvc, err := drive.NewService(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, err
	}

	sheetsSvc, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, err
	}

	return &SpreadsheetServiceGoogle{
		driveSvc:     driveSvc,
		sheetsSvc:    sheetsSvc,
		adminAccount: adminAccount,
		svcAccount:   svcAccount,
	}, err
}

func (s *SpreadsheetServiceGoogle) SwitchModeRead(ctx context.Context, spreadsheetID string) error {
	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}

	_, err := s.driveSvc.Permissions.Create(spreadsheetID, permission).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *SpreadsheetServiceGoogle) SwitchModeEdit(ctx context.Context, spreadsheetID string) error {
	permission := &drive.Permission{
		Type: "anyone",
		Role: "writer",
	}

	_, err := s.driveSvc.Permissions.Create(spreadsheetID, permission).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *SpreadsheetServiceGoogle) BlockImportantRanges(ctx context.Context, spreadsheetID string) error {
	spreadsheet, err := s.sheetsSvc.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return err
	}
	batch := NewBatchUpdate(s.sheetsSvc)

	for _, namedRange := range spreadsheet.NamedRanges {
		namedRange := namedRange
		if strings.Contains(namedRange.Name, "_blocked_") {
			batch.WithRequest(
				&sheets.Request{
					AddProtectedRange: &sheets.AddProtectedRangeRequest{
						ProtectedRange: &sheets.ProtectedRange{
							Description: namedRange.Name,
							Editors: &sheets.Editors{
								Users: []string{s.adminAccount, s.svcAccount},
							},
							Range: &sheets.GridRange{
								SheetId:          namedRange.Range.SheetId,
								StartColumnIndex: namedRange.Range.StartColumnIndex,
								EndColumnIndex:   namedRange.Range.EndColumnIndex,
							},
						},
					},
				},
			)
		}
	}

	if err := batch.Do(ctx, spreadsheetID); err != nil {
		return err
	}

	return nil
}

func (s *SpreadsheetServiceGoogle) UnlockImportantRanges(ctx context.Context, spreadsheetID string) error {
	spreadsheet, err := s.sheetsSvc.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return err
	}

	batch := NewBatchUpdate(s.sheetsSvc)

	for _, sheet := range spreadsheet.Sheets {
		for _, protectedRange := range sheet.ProtectedRanges {
			if strings.Contains(protectedRange.Description, "_blocked_") {
				batch.WithRequest(
					&sheets.Request{
						DeleteProtectedRange: &sheets.DeleteProtectedRangeRequest{
							ProtectedRangeId: protectedRange.ProtectedRangeId,
						},
					},
				)
			}
		}
	}

	if err := batch.Do(ctx, spreadsheetID); err != nil {
		return err
	}

	return nil
}

func (s *SpreadsheetServiceGoogle) GetPublicLink(_ context.Context, spreadsheetID string) string {
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit?usp=sharing", spreadsheetID)
	return url
}

func (s *SpreadsheetServiceGoogle) Comments(ctx context.Context, application *domain.Application) (*domain.Revision, error) {

	dataMap, ok := application.Attrs.(map[string]interface{})
	if !ok {
		return &domain.Revision{}, nil
	}
	applicationRawData, ok := dataMap["application"].(map[string]interface{})
	if !ok {
		return &domain.Revision{}, nil
	}
	applicationMap := map[string]string{}
	for i, j := range applicationRawData {
		applicationMap[i] = j.(string)
	}
	summary := &domain.Revision{
		ApplicationID:  application.ID,
		SpreadsheetID:  application.SpreadsheetID,
		No:             application.No,
		CreatedAt:      application.CreatedAt,
		Link:           s.GetPublicLink(ctx, application.SpreadsheetID),
		BIN:            applicationMap["bin"],
		Manufactor:     applicationMap["manufacturer"],
		To:             applicationMap["from"],
		ApplicantEmail: applicationMap["cont_email"],
		Address:        applicationMap["fact_addr"],
	}

	spreadsheetID := application.SpreadsheetID
	exportMimeType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	exportedContent, err := s.driveSvc.Files.Export(spreadsheetID, exportMimeType).Download()
	if err != nil {
		return summary, err
	}

	content := exportedContent.Body
	defer exportedContent.Body.Close()
	file_xlsx, err := excelize.OpenReader(content)
	if err != nil {
		return summary, err
	}
	sheet_list := file_xlsx.GetSheetList()

	cnt := 0
	for _, i := range sheet_list {
		comments, _ := file_xlsx.GetComments(i)
		if len(comments) != 0 {
			summary.Remarks += fmt.Sprintln("Таблица " + i + ":")
			for _, j := range comments {
				cnt++
				y2 := 3
				x2 := 1
				x, y, _ := excelize.CellNameToCoordinates(j.Cell)
				if i != "Заявление" {
					x2 = x
					y = 3
					y2 = 2
				}
				if i == "ТНВЭД" || i == "ОКВЭД" {
					y = 1
					y2 = 1
				}
				var (
					column_cell, _     = excelize.CoordinatesToCellName(x, y2)
					column, _          = file_xlsx.GetCellValue(i, column_cell)
					column_add_cell, _ = excelize.CoordinatesToCellName(x2, y)
					column_add, _      = file_xlsx.GetCellValue(i, column_add_cell)
				)

				summary.Remarks += fmt.Sprintf("%d) %s", cnt, column)
				if column != column_add {
					summary.Remarks += fmt.Sprintf(" - %s", column_add)
				}
				index := strings.LastIndex(j.Text, "-")
				summary.Remarks += fmt.Sprintf(" (Клетка-%s), Замечания: %s\n", j.Cell, j.Text[:index-2])
			}
		}
	}
	fmt.Println(summary.Remarks)
	if err != nil {
		return summary, err
	}

	return summary, nil
}
