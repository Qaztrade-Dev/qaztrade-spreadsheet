package adapters

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetClient struct {
	sheetsService *sheets.Service
	driveService  *drive.Service
	credentials   *google.Credentials
}

var _ domain.SpreadsheetRepository = (*SpreadsheetClient)(nil)

func NewSpreadsheetClient(ctx context.Context, credentialsJson []byte) (*SpreadsheetClient, error) {
	sheetsService, err := sheets.NewService(
		ctx,
		option.WithCredentialsJSON(credentialsJson),
	)
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(
		ctx,
		option.WithCredentialsJSON(credentialsJson),
	)
	if err != nil {
		return nil, err
	}

	credentials, err := google.CredentialsFromJSON(
		ctx,
		credentialsJson,
		"https://www.googleapis.com/auth/spreadsheets",
	)
	if err != nil {
		return nil, err
	}

	return &SpreadsheetClient{
		sheetsService: sheetsService,
		driveService:  driveService,
		credentials:   credentials,
	}, nil
}

func (c *SpreadsheetClient) GetApplication(ctx context.Context, spreadsheetID string) (*domain.Application, error) {
	var result domain.Application

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

	batchDataValues, err := c.getDataFromRanges(ctx, spreadsheetID, strRanges)
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

func (c *SpreadsheetClient) getDataFromRanges(ctx context.Context, spreadsheetID string, ranges []string) ([][][]interface{}, error) {
	resp, err := c.sheetsService.Spreadsheets.Values.BatchGet(spreadsheetID).Ranges(ranges...).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	datas := make([][][]interface{}, len(resp.ValueRanges))
	for i := range resp.ValueRanges {
		datas[i] = resp.ValueRanges[i].Values
	}
	return datas, nil
}

func (c *SpreadsheetClient) GetExpenseValues(ctx context.Context, spreadsheetID string, expensesTitles []string) ([]float64, error) {
	strRanges := make([]string, 0, len(expensesTitles))
	for i := range expensesTitles {
		strRanges = append(strRanges, fmt.Sprintf("'%s'!%s_expense_value", expensesTitles[i], strings.ReplaceAll(expensesTitles[i], " ", "_")))
	}

	batchDataValues, err := c.getDataFromRanges(ctx, spreadsheetID, strRanges)
	if err != nil {
		return nil, err
	}

	expenseValues := make([]float64, len(expensesTitles))
	for i := range batchDataValues {
		var value string
		if len(batchDataValues[i]) > 0 && len(batchDataValues[i][0]) > 0 {
			value = strings.TrimSpace(batchDataValues[i][0][0].(string))
		}

		value = strings.ReplaceAll(value, ",", ".")
		value = strings.ReplaceAll(value, " ", "")
		value = strings.ReplaceAll(value, "\u00a0", "")

		expenseValues[i], err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
	}

	return expenseValues, nil
}

func (c *SpreadsheetClient) GetExpensesSheetTitles(ctx context.Context, spreadsheetID string) ([]string, error) {
	spreadsheet, err := c.sheetsService.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	sheetTitles := make([]string, 0)
	for _, sheet := range spreadsheet.Sheets {
		switch sheet.Properties.Title {
		case "Заявление", "ТНВЭД", "ОКВЭД":
			continue
		}
		sheetTitles = append(sheetTitles, sheet.Properties.Title)
	}
	return sheetTitles, nil
}

func (c *SpreadsheetClient) GetAttachments(ctx context.Context, spreadsheetID string, expensesTitles []string) ([]io.ReadSeeker, error) {
	spreadsheet, err := c.sheetsService.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Ranges(expensesTitles...).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	sheets := spreadsheet.Sheets

	exportRequests := make([]*exportRequest, 0, len(sheets))
	for _, sheet := range sheets {
		nonEmptyRange := getNonEmptyRange(sheet)
		if nonEmptyRange.RowEnd <= 2 {
			continue
		}

		exportRequests = append(exportRequests, &exportRequest{
			RowStart:      nonEmptyRange.RowStart,
			RowEnd:        nonEmptyRange.RowEnd,
			ColumnStart:   nonEmptyRange.ColumnStart,
			ColumnEnd:     nonEmptyRange.ColumnEnd,
			sheetID:       sheet.Properties.SheetId,
			spreadsheetID: spreadsheetID,
		})
	}

	attachments := make([]io.ReadSeeker, 0, len(sheets))
	for _, exportReq := range exportRequests {
		req, err := http.NewRequestWithContext(ctx, "GET", exportReq.ExportURL(), nil)
		if err != nil {
			return nil, err
		}

		token, err := c.credentials.TokenSource.Token()
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+token.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		attachments = append(attachments, bytes.NewReader(body))
	}

	return attachments, nil
}

type getNonEmptyRangeResponse struct {
	RowStart    int64
	RowEnd      int64 // inclusive
	ColumnStart int64
	ColumnEnd   int64
}

func getNonEmptyRange(sheet *sheets.Sheet) *getNonEmptyRangeResponse {
	var (
		sheetLength = len(sheet.Data[0].RowData)
		rowEnd      = sheetLength - 1
		columnEnd   = -1
	)

	for i := sheetLength - 1; i >= 0; i-- {
		var (
			row                = sheet.Data[0].RowData[i]
			nonEmptyCellsCount = 0
		)

		for _, cell := range row.Values {
			switch {
			case cell.UserEnteredValue == nil:
				continue
			case cell.UserEnteredValue.StringValue != nil:
				value := strings.TrimSpace(*cell.UserEnteredValue.StringValue)
				if len(value) > 0 {
					nonEmptyCellsCount += 1
				}
			case cell.UserEnteredValue.BoolValue != nil,
				cell.UserEnteredValue.NumberValue != nil,
				cell.Hyperlink != "":
				nonEmptyCellsCount += 1
			}
		}

		if nonEmptyCellsCount > 0 {
			rowEnd = i
			break
		}
	}

	for i := 0; i < sheetLength; i++ {
		var (
			row = sheet.Data[0].RowData[i]
		)

		for j, cell := range row.Values {
			if cell.UserEnteredValue == nil {
				continue
			}
			if cell.UserEnteredValue.StringValue == nil {
				continue
			}
			value := strings.TrimSpace(*cell.UserEnteredValue.StringValue)
			if len(value) == 0 {
				continue
			}
			if value == "Затраты заявленные Заявителем (по докум. заявки) без НДС и акцизы РК" {
				columnEnd = j
				break
			}
		}

		if columnEnd != -1 {
			break
		}
	}

	return &getNonEmptyRangeResponse{
		RowStart:    0,
		RowEnd:      int64(rowEnd),
		ColumnStart: 0,
		ColumnEnd:   int64(columnEnd),
	}
}

type exportRequest struct {
	spreadsheetID string
	sheetID       int64
	RowStart      int64
	RowEnd        int64 // inclusive
	ColumnStart   int64
	ColumnEnd     int64 // inclusive
}

func (p *exportRequest) getQueryParams() *url.Values {
	queryParams := &url.Values{
		"exportFormat":  {"pdf"},
		"format":        {"pdf"},
		"size":          {"A4"},
		"portrait":      {"false"},
		"fitw":          {"false"},
		"sheetnames":    {"true"},
		"printtitle":    {"false"},
		"pagenumbers":   {"false"},
		"gridlines":     {"true"},
		"fzr":           {"true"},
		"top_margin":    {"0"},
		"bottom_margin": {"0"},
		"left_margin":   {"0"},
		"right_margin":  {"0"},
		"gid":           {fmt.Sprintf("%d", p.sheetID)},
		"range":         {p.getRange()},
	}

	return queryParams
}

func (p *exportRequest) ExportURL() string {
	var (
		queryParams    = p.getQueryParams()
		queryParamsStr = queryParams.Encode()
		urlStr         = fmt.Sprintf(
			"https://docs.google.com/spreadsheets/d/%s/export?%s",
			p.spreadsheetID,
			queryParamsStr,
		)
	)

	return urlStr
}

func (p *exportRequest) getRange() string {
	var (
		columnStartA1 = numToColName(p.ColumnStart)
		columnEndA1   = numToColName(p.ColumnEnd)
		fromRangeA1   = fmt.Sprintf("%s%d", columnStartA1, p.RowStart+1)
		toRangeA1     = fmt.Sprintf("%s%d", columnEndA1, p.RowEnd+1)
		rangeA1       = fmt.Sprintf("%s:%s", fromRangeA1, toRangeA1)
	)

	return rangeA1
}

func (c *SpreadsheetClient) UpdateSigningTime(ctx context.Context, spreadsheetID, signingTime string) error {
	var mappings = []struct {
		Range string
		Value string
	}{
		{Range: "signing_time", Value: signingTime},
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

	_, err := c.sheetsService.Spreadsheets.Values.BatchUpdate(spreadsheetID, updateValuesRequest).Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadsheetClient) SwitchModeRead(ctx context.Context, spreadsheetID string) error {
	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	_, err := s.driveService.Permissions.Insert(spreadsheetID, permission).Do()
	if err != nil {
		return err
	}
	return nil
}

// numToColName converts a zero-indexed number to Google Sheets column name.
func numToColName(num int64) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var colName strings.Builder

	for num >= 0 {
		remainder := num % 26
		colName.WriteString(string(chars[remainder]))
		num = (num / 26) - 1
	}

	return reverseString(colName.String())
}

// reverseString is a helper function to reverseString a string.
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
