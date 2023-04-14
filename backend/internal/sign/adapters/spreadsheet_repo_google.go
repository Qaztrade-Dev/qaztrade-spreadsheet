package adapters

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
		{Range: "activity", Value: &result.Activity},
		{Range: "emp_count", Value: &result.EmpCount},
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

func (c *SpreadsheetClient) getExpensesSheetTitles(ctx context.Context, spreadsheetID string) ([]string, error) {
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

func (c *SpreadsheetClient) GetAttachments(ctx context.Context, spreadsheetID string) ([]io.ReadSeeker, error) {
	expensesSheetTitles, err := c.getExpensesSheetTitles(ctx, spreadsheetID)
	if err != nil {
		return nil, err
	}

	if len(expensesSheetTitles) == 0 {
		return nil, errors.New("no expenses")
	}

	spreadsheet, err := c.sheetsService.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Ranges(expensesSheetTitles...).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	sheets := spreadsheet.Sheets

	exportRequests := make([]*exportRequest, 0, len(sheets))
	for _, sheet := range sheets {
		nonEmptyRange := getNonEmptyRange(sheet)
		if nonEmptyRange.RangeEnd <= 2 {
			continue
		}

		exportRequests = append(exportRequests, &exportRequest{
			RangeStart:    nonEmptyRange.RangeStart,
			RangeEnd:      nonEmptyRange.RangeEnd,
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
	RangeStart int64
	RangeEnd   int64 // inclusive
}

func getNonEmptyRange(sheet *sheets.Sheet) *getNonEmptyRangeResponse {
	var (
		sheetLength = len(sheet.Data[0].RowData)
		rangeEnd    = sheetLength - 1
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
			rangeEnd = i
			break
		}
	}

	return &getNonEmptyRangeResponse{
		RangeStart: 0,
		RangeEnd:   int64(rangeEnd),
	}
}

type exportRequest struct {
	spreadsheetID string
	sheetID       int64
	RangeStart    int64
	RangeEnd      int64 // inclusive
}

func (p *exportRequest) getQueryParams() *url.Values {
	var (
		rangeStart = p.RangeStart + 1
		rangeEnd   = p.RangeEnd + 1
	)

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
		"range":         {fmt.Sprintf("%d:%d", rangeStart, rangeEnd)},
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
		ValueInputOption: "RAW",
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
