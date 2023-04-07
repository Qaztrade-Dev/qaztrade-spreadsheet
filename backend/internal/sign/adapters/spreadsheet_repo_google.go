package adapters

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetClient struct {
	service     *sheets.Service
	credentials *google.Credentials
}

var _ domain.SpreadsheetRepository = (*SpreadsheetClient)(nil)

func NewSpreadsheetClient(ctx context.Context, credentialsJson []byte) (*SpreadsheetClient, error) {
	service, err := sheets.NewService(
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
		service:     service,
		credentials: credentials,
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
		value := strings.TrimSpace(batchDataValues[i][0][0].(string))
		*mappings[i].Value = value
	}

	return &result, nil
}

func (c *SpreadsheetClient) GetAttachments(ctx context.Context, spreadsheetID string) ([]io.ReadSeeker, error) {
	spreadsheet, err := c.service.Spreadsheets.Get(spreadsheetID).IncludeGridData(false).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	sheetIDs := make([]int64, 0)
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == "Заявление" {
			continue
		}
		sheetIDs = append(sheetIDs, sheet.Properties.SheetId)
	}

	attachments := make([]io.ReadSeeker, 0, len(sheetIDs))
	for _, sheetID := range sheetIDs {
		exportURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?exportFormat=pdf&format=pdf&size=A4&portrait=false&fitw=false&sheetnames=false&printtitle=false&pagenumbers=false&gridlines=true&fzr=true&top_margin=0&bottom_margin=0&left_margin=0&right_margin=0&gid=%d", spreadsheetID, sheetID)

		req, err := http.NewRequestWithContext(ctx, "GET", exportURL, nil)
		if err != nil {
			return nil, err
		}

		token, err := c.credentials.TokenSource.Token()
		if err != nil {
			return nil, err
		}

		req.Header.Add("Authorization", "Bearer "+token.AccessToken)

		// Make request
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
