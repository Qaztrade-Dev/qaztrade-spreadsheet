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
	metaDataSvc  *sheets.SpreadsheetsDeveloperMetadataService
}
type MetaDataCommentsPack struct {
	sheetName string
	rowIdx    int64
	colIdx    int64
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

	metaDataSvc := sheets.NewSpreadsheetsDeveloperMetadataService(sheetsSvc)
	return &SpreadsheetServiceGoogle{
		driveSvc:     driveSvc,
		sheetsSvc:    sheetsSvc,
		adminAccount: adminAccount,
		svcAccount:   svcAccount,
		metaDataSvc:  metaDataSvc,
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

func (s *SpreadsheetServiceGoogle) LockSheets(ctx context.Context, spreadsheetID string) error {
	spreadsheet, err := s.sheetsSvc.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return err
	}

	batch := NewBatchUpdate(s.sheetsSvc)

	for _, sheet := range spreadsheet.Sheets {
		sheet := sheet

		if !(sheet.Properties.Title == "Заявление" ||
			sheet.Properties.Title == "ТНВЭД" ||
			sheet.Properties.Title == "ОКВЭД") {
			continue
		}

		batch.WithRequest(
			&sheets.Request{
				AddProtectedRange: &sheets.AddProtectedRangeRequest{
					ProtectedRange: &sheets.ProtectedRange{
						Range: &sheets.GridRange{
							SheetId: sheet.Properties.SheetId,
						},
						Description: "Protecting entire sheet" + " " + sheet.Properties.Title,
						WarningOnly: false,
					},
				},
			},
		)
	}

	if err := batch.Do(ctx, spreadsheetID); err != nil {
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

func (s *SpreadsheetServiceGoogle) GetApplication(ctx context.Context, spreadsheetID string) (*domain.ApplicationAttrs, error) {
	var result domain.ApplicationAttrs

	var mappings = []struct {
		Range string
		Value *string
	}{
		{Range: "from", Value: &result.From},
		{Range: "gov_reg", Value: &result.GovReg},
		{Range: "fact_addr", Value: &result.FactAddr},
		{Range: "bin", Value: &result.Bin},
		{Range: "industry", Value: &result.Industry},
		{Range: "industry_other", Value: &result.IndustryOther},
		{Range: "activity", Value: &result.Activity},
		{Range: "emp_count", Value: &result.EmpCount},
		{Range: "tax_sum", Value: &result.TaxSum},
		{Range: "product_capacity", Value: &result.ProductCapacity},
		{Range: "manufacturer", Value: &result.Manufacturer},
		{Range: "item", Value: &result.Item},
		{Range: "item_volume", Value: &result.ItemVolume},
		{Range: "fact_volume_earnings", Value: &result.FactVolumeEarnings},
		{Range: "fact_workload", Value: &result.FactWorkload},
		{Range: "chief_lastname", Value: &result.ChiefLastname},
		{Range: "chief_firstname", Value: &result.ChiefFirstname},
		{Range: "chief_middlename", Value: &result.ChiefMiddlename},
		{Range: "chief_position", Value: &result.ChiefPosition},
		{Range: "chief_phone", Value: &result.ChiefPhone},
		{Range: "cont_lastname", Value: &result.ContLastname},
		{Range: "cont_firstname", Value: &result.ContFirstname},
		{Range: "cont_middlename", Value: &result.ContMiddlename},
		{Range: "cont_position", Value: &result.ContPosition},
		{Range: "cont_phone", Value: &result.ContPhone},
		{Range: "cont_email", Value: &result.ContEmail},
		{Range: "info_manufactured_goods", Value: &result.InfoManufacturedGoods},
		{Range: "name_of_goods", Value: &result.NameOfGoods},
		{Range: "spend_plan", Value: &result.SpendPlan},
		{Range: "spend_plan_other", Value: &result.SpendPlanOther},
		{Range: "metrics_2022", Value: &result.Metrics2022},
		{Range: "metrics_2023", Value: &result.Metrics2023},
		{Range: "metrics_2024", Value: &result.Metrics2024},
		{Range: "metrics_2025", Value: &result.Metrics2025},
		{Range: "has_agreement", Value: &result.HasAgreement},
	}

	strRanges := make([]string, 0, len(mappings))
	for i := range mappings {
		strRanges = append(strRanges, mappings[i].Range)
	}

	batchDataValues, err := s.getDataFromRanges(ctx, spreadsheetID, strRanges)
	if err != nil {
		return nil, err
	}

	for i := range batchDataValues {
		var value string
		if len(batchDataValues[i]) > 0 && len(batchDataValues[i][0]) > 0 {
			value = strings.TrimSpace(batchDataValues[i][0][0].(string))
		}
		*mappings[i].Value = value
	}

	return &result, nil
}

func (s *SpreadsheetServiceGoogle) getDataFromRanges(ctx context.Context, spreadsheetID string, ranges []string) ([][][]interface{}, error) {
	resp, err := s.sheetsSvc.Spreadsheets.Values.BatchGet(spreadsheetID).Ranges(ranges...).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	datas := make([][][]interface{}, len(resp.ValueRanges))
	for i := range resp.ValueRanges {
		datas[i] = resp.ValueRanges[i].Values
	}
	return datas, nil
}

func (s *SpreadsheetServiceGoogle) Comments(ctx context.Context, application *domain.Application) (*domain.Revision, error) {
	applicationAttr, err := s.GetApplication(ctx, application.SpreadsheetID)
	if err != nil {
		return nil, err
	}

	summary := &domain.Revision{
		ApplicationID:  application.ID,
		SpreadsheetID:  application.SpreadsheetID,
		No:             application.No,
		SignedAt:       application.SignedAt,
		Link:           s.GetPublicLink(ctx, application.SpreadsheetID),
		BIN:            applicationAttr.Bin,
		Manufactor:     applicationAttr.Manufacturer,
		To:             applicationAttr.From,
		ApplicantEmail: applicationAttr.ContEmail,
		Address:        applicationAttr.FactAddr,
	}

	var (
		spreadsheetID  = application.SpreadsheetID
		exportMimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	)

	exportedContent, err := s.driveSvc.Files.Export(spreadsheetID, exportMimeType).Download()
	if err != nil {
		return summary, err
	}
	defer exportedContent.Body.Close()

	fileXLSX, err := excelize.OpenReader(exportedContent.Body)
	if err != nil {
		return summary, err
	}

	sheetList := fileXLSX.GetSheetList()

	if err := s.deleteMetadata(ctx, spreadsheetID); err != nil {
		return nil, err
	}

	var (
		arrMetadata = []*MetaDataCommentsPack{}
		cnt         = 0
	)

	for _, i := range sheetList {
		comments, _ := fileXLSX.GetComments(i)
		if len(comments) != 0 {
			summary.Remarks += fmt.Sprint("\u200b         Таблица " + i + ":\n")

			for _, j := range comments {
				cnt++
				y2 := 3
				x2 := 1
				x, y, _ := excelize.CellNameToCoordinates(j.Cell)
				arrMetadata = append(arrMetadata, &MetaDataCommentsPack{
					sheetName: i,
					rowIdx:    int64(y),
					colIdx:    int64(x),
				})
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
					column, _          = fileXLSX.GetCellValue(i, column_cell)
					column_add_cell, _ = excelize.CoordinatesToCellName(x2, y)
					column_add, _      = fileXLSX.GetCellValue(i, column_add_cell)
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
	if err := s.SetMetadata(ctx, spreadsheetID, arrMetadata); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return summary, nil
}
func (s *SpreadsheetServiceGoogle) deleteMetadataByKey(ctx context.Context, spreadsheetID, key string, batch *BatchUpdate) {
	batch.WithRequest(&sheets.Request{
		DeleteDeveloperMetadata: &sheets.DeleteDeveloperMetadataRequest{
			DataFilter: &sheets.DataFilter{
				DeveloperMetadataLookup: &sheets.DeveloperMetadataLookup{
					Visibility:  "DOCUMENT",
					MetadataKey: key,
				},
			},
		},
	})
}

func (s *SpreadsheetServiceGoogle) deleteMetadata(ctx context.Context, spreadsheetID string) error {
	batch := NewBatchUpdate(s.sheetsSvc)

	var (
		filter []*sheets.DataFilter
	)

	filter = append(filter, &sheets.DataFilter{
		DeveloperMetadataLookup: &sheets.DeveloperMetadataLookup{
			Visibility: "DOCUMENT",
		},
	})
	reqMeta := &sheets.SearchDeveloperMetadataRequest{
		DataFilters: filter,
	}

	response, err := s.metaDataSvc.Search(spreadsheetID, reqMeta).Do()
	if err != nil {
		return err
	}
	for _, i := range response.MatchedDeveloperMetadata {
		if i.DeveloperMetadata.MetadataKey[0] == '!' {
			s.deleteMetadataByKey(ctx, spreadsheetID, i.DeveloperMetadata.MetadataKey, batch)
		}
	}

	if err := batch.Do(ctx, spreadsheetID); err != nil {
		return err
	}
	return nil
}

func (s *SpreadsheetServiceGoogle) SetMetadata(ctx context.Context, spreadsheetID string, arrMetadata []*MetaDataCommentsPack) error {
	if err := s.deleteMetadata(ctx, spreadsheetID); err != nil {
		return err
	}

	batch := NewBatchUpdate(s.sheetsSvc)
	for _, i := range arrMetadata {
		batch.WithRequest(&sheets.Request{
			CreateDeveloperMetadata: &sheets.CreateDeveloperMetadataRequest{
				DeveloperMetadata: &sheets.DeveloperMetadata{
					Location: &sheets.DeveloperMetadataLocation{
						Spreadsheet: true,
					},
					Visibility:    "DOCUMENT",
					MetadataKey:   fmt.Sprintf("!%s-%d:%d", i.sheetName, i.rowIdx, i.colIdx),
					MetadataValue: "true",
				},
			},
		})
	}

	if err := batch.Do(ctx, spreadsheetID); err != nil {
		return err
	}

	return nil
}
