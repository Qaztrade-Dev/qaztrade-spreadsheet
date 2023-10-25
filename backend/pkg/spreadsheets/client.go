package spreadsheets

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetClient struct {
	driveSvc  *drive.Service
	sheetsSvc *sheets.Service
}

var _ SpreadsheetService = (*SpreadsheetClient)(nil)

func NewSpreadsheetClient(ctx context.Context, credentialsJson []byte) (*SpreadsheetClient, error) {
	driveSvc, err := drive.NewService(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, err
	}

	sheetsSvc, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, err
	}

	return &SpreadsheetClient{
		driveSvc:  driveSvc,
		sheetsSvc: sheetsSvc,
	}, nil
}

func (s *SpreadsheetClient) SwitchModeRead(ctx context.Context, spreadsheetID string) error {
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

func (s *SpreadsheetClient) SwitchModeEdit(ctx context.Context, spreadsheetID string) error {
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

func (s *SpreadsheetClient) LockSheets(ctx context.Context, spreadsheetID string) error {
	spreadsheet, err := s.sheetsSvc.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return err
	}

	batch := NewBatchUpdate(s.sheetsSvc)

	for _, sheet := range spreadsheet.Sheets {
		sheet := sheet

		if sheet.Properties.Title == "Заявление" ||
			sheet.Properties.Title == "ТНВЭД" ||
			sheet.Properties.Title == "ОКВЭД" {
			continue
		}

		for _, protectedRange := range sheet.ProtectedRanges {
			protectedRange := protectedRange

			if !strings.HasPrefix(protectedRange.Description, "Protecting entire sheet") {
				continue
			}

			batch.WithRequest(
				&sheets.Request{
					DeleteProtectedRange: &sheets.DeleteProtectedRangeRequest{
						ProtectedRangeId: protectedRange.ProtectedRangeId,
					},
				},
			)
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

func (s *SpreadsheetClient) GrantAdminPermissions(ctx context.Context, spreadsheetID, email string) error {
	if err := s.grantWritePermission(ctx, spreadsheetID, email); err != nil {
		return err
	}

	if err := s.grantPermissionToProtectedRanges(ctx, spreadsheetID, email); err != nil {
		return err
	}

	return nil
}

func (s *SpreadsheetClient) grantWritePermission(ctx context.Context, spreadsheetID, email string) error {
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

func (s *SpreadsheetClient) grantPermissionToProtectedRanges(ctx context.Context, spreadsheetID, email string) error {
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

func (s *SpreadsheetClient) GetApplicationAttrs(ctx context.Context, spreadsheetID string) (*ApplicationAttrs, error) {
	var result ApplicationAttrs

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

func (s *SpreadsheetClient) getDataFromRanges(ctx context.Context, spreadsheetID string, ranges []string) ([][][]interface{}, error) {
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

func (s *SpreadsheetClient) editorsContains(email string, editors []string) bool {
	for _, editor := range editors {
		if editor == email {
			return true
		}
	}
	return false
}

func (c *SpreadsheetClient) GetSheetData(ctx context.Context, spreadsheetID string, sheetTitle string) ([][]string, error) {
	spreadsheet, err := c.sheetsSvc.Spreadsheets.Get(spreadsheetID).
		IncludeGridData(true).
		Ranges(sheetTitle).
		Context(ctx).
		Do()
	if err != nil {
		return nil, err
	}

	if len(spreadsheet.Sheets) == 0 {
		return nil, ErrorSheetNotFound
	}

	var (
		sheet  = spreadsheet.Sheets[0]
		data   = sheet.Data[0].RowData
		result = make([][]string, len(data))
	)

	for i := range data {
		row := data[i]
		result[i] = make([]string, len(row.Values))

		for j := range row.Values {
			cell := row.Values[j]
			result[i][j] = getValue(cell.UserEnteredValue)
		}
	}

	return result, nil
}

func getValue(input *sheets.ExtendedValue) string {
	result := ""
	if input == nil {
		return result
	}

	switch {
	case input.BoolValue != nil:
		result = fmt.Sprintf("%v", *input.BoolValue)
	case input.NumberValue != nil:
		result = fmt.Sprintf("%f", *input.NumberValue)
	case input.StringValue != nil:
		result = *input.StringValue
	}

	return result
}
