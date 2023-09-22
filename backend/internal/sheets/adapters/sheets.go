package adapters

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetClient struct {
	service *sheets.Service
}

var _ domain.SheetsRepository = (*SpreadsheetClient)(nil)

func NewSpreadsheetClient(ctx context.Context, credentialsJson []byte) (*SpreadsheetClient, error) {
	service, err := sheets.NewService(
		ctx,
		option.WithCredentialsJSON(credentialsJson),
	)
	if err != nil {
		return nil, err
	}

	return &SpreadsheetClient{
		service: service,
	}, nil
}

func (c *SpreadsheetClient) NewSpreadsheetServiceMetadata() *sheets.SpreadsheetsDeveloperMetadataService {
	return sheets.NewSpreadsheetsDeveloperMetadataService(c.service)
}

func (c *SpreadsheetClient) UpdateApplication(ctx context.Context, spreadsheetID string, application *domain.Application) error {
	var mappings = []struct {
		Range string
		Value string
	}{
		{Range: "from", Value: application.From},
		{Range: "gov_reg", Value: application.GovReg},
		{Range: "fact_addr", Value: application.FactAddr},
		{Range: "bin", Value: application.Bin},
		{Range: "industry", Value: application.Industry},
		{Range: "industry_other", Value: application.IndustryOther},
		{Range: "activity", Value: application.Activity},
		{Range: "emp_count", Value: application.EmpCount},
		{Range: "tax_sum", Value: application.TaxSum},
		{Range: "product_capacity", Value: application.ProductCapacity},
		{Range: "manufacturer", Value: application.Manufacturer},
		{Range: "item", Value: application.Item},
		{Range: "item_volume", Value: application.ItemVolume},
		{Range: "fact_volume_earnings", Value: application.FactVolumeEarnings},
		{Range: "fact_workload", Value: application.FactWorkload},
		{Range: "chief_lastname", Value: application.ChiefLastname},
		{Range: "chief_firstname", Value: application.ChiefFirstname},
		{Range: "chief_middlename", Value: application.ChiefMiddlename},
		{Range: "chief_position", Value: application.ChiefPosition},
		{Range: "chief_phone", Value: application.ChiefPhone},
		{Range: "cont_lastname", Value: application.ContLastname},
		{Range: "cont_firstname", Value: application.ContFirstname},
		{Range: "cont_middlename", Value: application.ContMiddlename},
		{Range: "cont_position", Value: application.ContPosition},
		{Range: "cont_phone", Value: application.ContPhone},
		{Range: "cont_email", Value: application.ContEmail},
		{Range: "info_manufactured_goods", Value: application.InfoManufacturedGoods},
		{Range: "name_of_goods", Value: application.NameOfGoods},
		{Range: "spend_plan", Value: application.SpendPlan},
		{Range: "spend_plan_other", Value: application.SpendPlanOther},
		{Range: "metrics_2022", Value: application.Metrics2022},
		{Range: "metrics_2023", Value: application.Metrics2023},
		{Range: "metrics_2024", Value: application.Metrics2024},
		{Range: "metrics_2025", Value: application.Metrics2025},
		{Range: "has_agreement", Value: application.HasAgreement},
		{Range: "agreement_file", Value: application.AgreementFile},
	}

	data := make([]*sheets.ValueRange, 0, len(mappings))
	for i := range mappings {
		data = append(data, &sheets.ValueRange{
			Range:  mappings[i].Range,
			Values: [][]interface{}{{mappings[i].Value}},
		})
	}

	updateValuesRequest := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
		Data:             data,
	}

	_, err := c.service.Spreadsheets.Values.BatchUpdate(spreadsheetID, updateValuesRequest).Do()
	if err != nil {
		return fmt.Errorf("BatchUpdate: %w", err)
	}

	return nil
}

type SheetClient struct {
	service       *sheets.Service
	spreadsheetID string
	sheetName     string
	sheetID       int64
	data          [][]string
}

func (c *SpreadsheetClient) NewSheetClient(ctx context.Context, spreadsheetID, sheetName string, sheetID int64) (*SheetClient, error) {
	sheetClient := &SheetClient{
		service:       c.service,
		spreadsheetID: spreadsheetID,
		sheetName:     sheetName,
		sheetID:       sheetID,
	}

	var (
		dataRangeName = getSheetRangeData(sheetName)
	)

	t1 := time.Now()
	batchDataValues, err := c.getDataFromRanges(ctx, spreadsheetID, []string{dataRangeName})
	fmt.Println("time: getDataFromRanges", time.Since(t1))
	if err != nil {
		return nil, err
	}

	var (
		dataValues = batchDataValues[0]
	)

	t1 = time.Now()
	data, err := sheetClient.getData(ctx, sheetName, dataValues)
	fmt.Println("time: getData", time.Since(t1))
	if err != nil {
		return nil, err
	}

	sheetClient.data = data

	return sheetClient, nil
}

func (c *SpreadsheetClient) getDataFromRanges(ctx context.Context, spreadsheetID string, ranges []string) ([][][]interface{}, error) {
	resp, err := c.service.Spreadsheets.Values.BatchGet(spreadsheetID).Ranges(ranges...).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	datas := make([][][]interface{}, len(resp.ValueRanges))
	for i := range resp.ValueRanges {
		datas[i] = resp.ValueRanges[i].Values
	}
	return datas, nil
}

func getSheetRangeData(sheetName string) string {
	rangeName := fmt.Sprintf("'%s'!%s_%s", sheetName, strings.ReplaceAll(sheetName, " ", "_"), "data")

	return rangeName
}

func (c *SheetClient) getData(ctx context.Context, sheetName string, values [][]interface{}) ([][]string, error) {
	data := make([][]string, len(values))
	for i, row := range values {
		data[i] = make([]string, len(row))
		for j := range row {
			data[i][j] = strings.TrimSpace(row[j].(string))
		}
	}
	return data, nil
}

type UpdateCellRequest struct {
	RowIndex    int64
	ColumnIndex int64
	Value       string
}

func (r *UpdateCellRequest) encode(sheetID int64, newStringVal *string, newTextFormat []*sheets.TextFormatRun) *sheets.Request {
	return &sheets.Request{
		UpdateCells: &sheets.UpdateCellsRequest{
			Fields: "UserEnteredValue,TextFormatRuns",
			Start: &sheets.GridCoordinate{
				RowIndex:    r.RowIndex,
				ColumnIndex: r.ColumnIndex,
				SheetId:     sheetID,
			},
			Rows: []*sheets.RowData{
				{
					Values: []*sheets.CellData{
						{
							UserEnteredValue: &sheets.ExtendedValue{
								StringValue: newStringVal,
							},
							TextFormatRuns: newTextFormat,
						},
					},
				},
			},
		},
	}
}

func (c *SpreadsheetClient) UpdateCell(ctx context.Context, spreadsheetID string, input *domain.UpdateCellInput) error {
	var (
		batch             = NewBatchUpdate(c.service)
		updateCellRequest = &UpdateCellRequest{
			RowIndex:    input.RowIdx - 1,
			ColumnIndex: input.ColumnIdx - 1,
			Value:       input.Value,
		}
		cell, err     = c.getCellValue(ctx, spreadsheetID, input.SheetName, input.RowIdx, input.ColumnIdx)
		curTime       = time.Now()
		timeToString  = fmt.Sprintf("(%s)", curTime.Format("02-01-2006"))
		newStringVal  = "файл " + timeToString
		newTextFormat = []*sheets.TextFormatRun{}
	)
	if err != nil {
		return err
	}
	newTextFormat = append(newTextFormat, &sheets.TextFormatRun{
		Format: &sheets.TextFormat{
			Link: &sheets.Link{
				Uri: input.Value,
			},
		},
		StartIndex: 0,
	})

	if !input.Replace && *cell.FormattedValue != "" {
		newStringVal += "\nстарый " + *cell.FormattedValue
		if len(cell.TextFormatRun) == 0 {
			newTextFormat = append(newTextFormat, &sheets.TextFormatRun{
				Format: &sheets.TextFormat{
					Link: &sheets.Link{
						Uri: *cell.Hyperlink,
					},
				},
				StartIndex: 18,
			})

		} else {
			cell.TextFormatRun[0].StartIndex += 18
			for i := 1; i < len(cell.TextFormatRun); i++ {
				cell.TextFormatRun[i].StartIndex += 25
			}
			newTextFormat = append(newTextFormat, cell.TextFormatRun...)
		}
	}

	batch.WithRequest(
		updateCellRequest.encode(input.SheetID, &newStringVal, newTextFormat),
	)

	if err := batch.Do(ctx, spreadsheetID); err != nil {
		return err
	}

	return nil
}

func (c *SpreadsheetClient) AddRows(ctx context.Context, spreadsheetID string, input *domain.AddRowsInput) error {
	var (
		batch = NewBatchUpdate(c.service)
	)

	spreadsheet, err := c.service.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return err
	}

	var sheet *sheets.Sheet

	for _, s := range spreadsheet.Sheets {
		if s.Properties.SheetId == input.SheetID {
			sheet = s
			break
		}
	}

	lastRowIndex := int(sheet.Properties.GridProperties.RowCount)

	insertDimensionRequest := &sheets.Request{
		InsertDimension: &sheets.InsertDimensionRequest{
			Range: &sheets.DimensionRange{
				SheetId:    input.SheetID,
				Dimension:  "ROWS",
				StartIndex: int64(lastRowIndex),
				EndIndex:   int64(lastRowIndex + input.RowsAmount),
			},
			InheritFromBefore: true,
		},
	}

	batch.WithRequest(insertDimensionRequest)

	if err := batch.Do(ctx, spreadsheetID); err != nil {
		return err
	}

	return nil
}

type CellValue struct {
	TextFormatRun  []*sheets.TextFormatRun
	FormattedValue *string
	Hyperlink      *string
}

func (s *SpreadsheetClient) GetHyperLink(ctx context.Context, spreadsheetID string, SheetName string, Row_idx int64, Column_idx int64) (*string, error) {
	cell, err := s.getCellValue(ctx, spreadsheetID, SheetName, Row_idx, Column_idx)
	if err != nil {
		return nil, err
	}
	return cell.Hyperlink, nil
}

func (s *SpreadsheetClient) getCellValue(ctx context.Context, spreadsheetID string, SheetName string, Row_idx int64, Column_idx int64) (*CellValue, error) {
	A1Notion := (&Cell{
		Col: Column_idx - 1,
		Row: Row_idx - 1,
	}).ToA1()

	resp, err := s.service.Spreadsheets.GetByDataFilter(spreadsheetID, &sheets.GetSpreadsheetByDataFilterRequest{
		DataFilters: []*sheets.DataFilter{
			{A1Range: SheetName + "!" + A1Notion},
		},
		IncludeGridData: true,
	}).Do()

	if err != nil {
		return nil, err
	}

	for _, s := range resp.Sheets {
		if s.Properties.Title == SheetName {
			val := s.Data[0].RowData[0].Values[0]
			return &CellValue{
				TextFormatRun:  val.TextFormatRuns,
				FormattedValue: &val.FormattedValue,
				Hyperlink:      &val.Hyperlink,
			}, nil
		}
	}
	return nil, err
}
