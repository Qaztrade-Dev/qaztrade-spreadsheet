package adapters

import (
	"context"
	"fmt"

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

func (s *SpreadsheetServiceGoogle) GrantAdminPermissions(ctx context.Context, spreadsheetID, email string) error {
	if err := s.grantWritePermission(ctx, spreadsheetID, email); err != nil {
		return err
	}

	if err := s.grantPermissionToProtectedRanges(ctx, spreadsheetID, email); err != nil {
		return err
	}

	return nil
}

func (s *SpreadsheetServiceGoogle) grantWritePermission(ctx context.Context, spreadsheetID, email string) error {
	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer",
		EmailAddress: email,
	}
	_, err := s.driveSvc.Permissions.Create(spreadsheetID, permission).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *SpreadsheetServiceGoogle) grantPermissionToProtectedRanges(ctx context.Context, spreadsheetID, email string) error {
	spreadsheet, err := s.sheetsSvc.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return err
	}

	updateRanges := make([]*sheets.Request, 0)

	for _, sheet := range spreadsheet.Sheets {
		for _, protectedRange := range sheet.ProtectedRanges {
			protectedRange := protectedRange

			if s.editorsContains(email, protectedRange.Editors.Users) {
				continue
			}

			protectedRange.Editors.Users = append(protectedRange.Editors.Users, email)

			updateRanges = append(updateRanges, &sheets.Request{
				UpdateProtectedRange: &sheets.UpdateProtectedRangeRequest{
					ProtectedRange: protectedRange,
					Fields:         "editors",
				},
			})
		}
	}

	batch := NewBatchUpdate(s.sheetsSvc)
	batch.WithRequest(updateRanges...)
	if err := batch.Do(ctx, spreadsheetID); err != nil {
		return err
	}

	return nil
}

func (s *SpreadsheetServiceGoogle) editorsContains(email string, editors []string) bool {
	for _, editor := range editors {
		if editor == email {
			return true
		}
	}
	return false
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
		Link:           s.GetPublicLink(ctx, application.SpreadsheetID),
		BIN:            applicationMap["bin"],
		Manufactor:     applicationMap["manufacturer"],
		To:             applicationMap["from"],
		ApplicantEmail: applicationMap["cont_email"],
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
				var (
					x, _, _            = excelize.CellNameToCoordinates(j.Cell)
					column_cell, _     = excelize.CoordinatesToCellName(x, 2)
					column, _          = file_xlsx.GetCellValue(i, column_cell)
					column_add_cell, _ = excelize.CoordinatesToCellName(x, 3)
					column_add, _      = file_xlsx.GetCellValue(i, column_add_cell)
				)

				summary.Remarks += fmt.Sprintf("%d) %s", cnt, column)
				if column != column_add {
					summary.Remarks += fmt.Sprintf(" - %s", column_add)
				}
				summary.Remarks += fmt.Sprintf(" (Клетка-%s), Замечания: %s\n", j.Cell, j.Text)
			}
		}
	}
	if err != nil {
		return summary, err
	}

	return summary, nil
}
