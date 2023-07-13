package adapters

import (
	"context"
	"fmt"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/doodocs/qaztrade/backend/pkg/qaztradeoauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetServiceGoogle struct {
	oauth2                *qaztradeoauth2.Client
	svcAccount            string
	reviewerAccount       string
	jwtcli                *jwt.Client
	originSpreadsheetID   string
	templateSpreadsheetID string
	destinationFolderID   string
}

var _ domain.SpreadsheetService = (*SpreadsheetServiceGoogle)(nil)

func NewSpreadsheetServiceGoogle(
	oauth2 *qaztradeoauth2.Client,
	svcAccount string,
	reviewerAccount string,
	jwtcli *jwt.Client,
	originSpreadsheetID string,
	templateSpreadsheetID string,
	destinationFolderID string,
) *SpreadsheetServiceGoogle {
	client := &SpreadsheetServiceGoogle{
		oauth2:                oauth2,
		svcAccount:            svcAccount,
		reviewerAccount:       reviewerAccount,
		jwtcli:                jwtcli,
		originSpreadsheetID:   originSpreadsheetID,
		templateSpreadsheetID: templateSpreadsheetID,
		destinationFolderID:   destinationFolderID,
	}

	return client
}

func (s *SpreadsheetServiceGoogle) Create(ctx context.Context, user *domain.User) (string, error) {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return "", err
	}

	driveSvc, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return "", err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return "", err
	}

	spreadsheetID, err := s.copyFile(ctx, driveSvc, user)
	if err != nil {
		return "", err
	}

	if err := s.initSpreadsheet(ctx, spreadsheetsSvc, spreadsheetID); err != nil {
		return "", err
	}

	if err := s.setPublic(ctx, driveSvc, spreadsheetID); err != nil {
		return "", err
	}

	if err := s.setReviewer(ctx, driveSvc, spreadsheetID); err != nil {
		return "", err
	}

	return spreadsheetID, nil
}

func (s *SpreadsheetServiceGoogle) copyFile(ctx context.Context, svc *drive.Service, user *domain.User) (string, error) {
	newFileName, err := domain.CreateSpreadsheetName(user)
	if err != nil {
		return "", err
	}

	copy := &drive.File{
		Name:    newFileName,
		Parents: []string{s.destinationFolderID},
	}
	copiedFile, err := svc.Files.Copy(s.templateSpreadsheetID, copy).Context(ctx).Do()
	if err != nil {
		return "", err
	}

	return copiedFile.Id, nil
}

func (s *SpreadsheetServiceGoogle) initSpreadsheet(ctx context.Context, svc *sheets.Service, spreadsheetID string) error {
	templateSpreadsheet, err := svc.Spreadsheets.Get(s.templateSpreadsheetID).Context(ctx).Do()

	if err != nil {
		return err
	}

	spreadsheet, err := svc.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return err
	}

	batch := NewBatchUpdate(svc)

	for _, sheet := range spreadsheet.Sheets {
		for _, protectedRange := range sheet.ProtectedRanges {
			protectedRange := protectedRange
			s.deleteProtectedRange(protectedRange, batch)
		}
	}

	for _, sheet := range templateSpreadsheet.Sheets {
		for _, protectedRange := range sheet.ProtectedRanges {
			protectedRange := protectedRange
			s.addProtectedRange(protectedRange, batch)
		}
	}

	if err := s.setMetadata(spreadsheetID, batch); err != nil {
		return err
	}

	if err := batch.Do(ctx, spreadsheetID); err != nil {
		return err
	}

	return nil
}

func (s *SpreadsheetServiceGoogle) deleteProtectedRange(protectedRange *sheets.ProtectedRange, batch *BatchUpdate) {
	batch.WithRequest(&sheets.Request{
		DeleteProtectedRange: &sheets.DeleteProtectedRangeRequest{
			ProtectedRangeId: protectedRange.ProtectedRangeId,
		},
	})
}

func (s *SpreadsheetServiceGoogle) addProtectedRange(protectedRange *sheets.ProtectedRange, batch *BatchUpdate) {
	batch.WithRequest(&sheets.Request{
		AddProtectedRange: &sheets.AddProtectedRangeRequest{
			ProtectedRange: protectedRange,
		},
	})
}

func (s *SpreadsheetServiceGoogle) setReviewer(ctx context.Context, svc *drive.Service, spreadsheetID string) error {
	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer",
		EmailAddress: s.reviewerAccount,
	}
	_, err := svc.Permissions.Create(spreadsheetID, permission).SendNotificationEmail(false).Context(ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadsheetServiceGoogle) setMetadata(spreadsheetID string, batch *BatchUpdate) error {
	tokenStr, err := jwt.NewTokenString(s.jwtcli, &domain.SpreadsheetClaims{
		SpreadsheetID: spreadsheetID,
	})
	if err != nil {
		return err
	}

	batch.WithRequest(&sheets.Request{
		DeleteDeveloperMetadata: &sheets.DeleteDeveloperMetadataRequest{
			DataFilter: &sheets.DataFilter{
				DeveloperMetadataLookup: &sheets.DeveloperMetadataLookup{
					MetadataKey: "token",
				},
			},
		},
	})

	batch.WithRequest(&sheets.Request{
		CreateDeveloperMetadata: &sheets.CreateDeveloperMetadataRequest{
			DeveloperMetadata: &sheets.DeveloperMetadata{
				Location: &sheets.DeveloperMetadataLocation{
					Spreadsheet: true,
				},
				Visibility:    "DOCUMENT",
				MetadataKey:   "token",
				MetadataValue: tokenStr,
			},
		},
	})
	return nil
}

func (s *SpreadsheetServiceGoogle) setPublic(ctx context.Context, svc *drive.Service, spreadsheetID string) error {
	permission := &drive.Permission{
		Type: "anyone",
		Role: "writer",
	}
	_, err := svc.Permissions.Create(spreadsheetID, permission).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *SpreadsheetServiceGoogle) GetPublicLink(_ context.Context, spreadsheetID string) string {
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit?usp=sharing", spreadsheetID)
	return url
}

func (c *SpreadsheetServiceGoogle) AddSheet(ctx context.Context, spreadsheetID string, sheetName string) error {
	var (
		mappings = map[string]int64{ // sheetName:sheetID
			"Затраты на доставку транспортом":             928848876,
			"Затраты на сертификацию предприятия":         693636717,
			"Затраты на рекламу ИКУ за рубежом":           1717840340,
			"Затраты на перевод каталога ИКУ":             784543090,
			"Затраты на аренду помещения ИКУ":             699998073,
			"Затраты на сертификацию ИКУ":                 830264833,
			"Затраты на демонстрацию ИКУ":                 826367543,
			"Затраты на франчайзинг":                      808421585,
			"Затраты на соответствие товаров требованиям": 564452334,
			"Затраты на регистрацию товарных знаков":      719601032,
			"Затраты на аренду":                           1638330977,
			"Затраты на перевод":                          545808572,
			"Затраты на рекламу товаров за рубежом":       1112839101,
			"Затраты на участие в выставках":              662810845,
			"Затраты на участие в выставках ИКУ":          1197759146,
		}
		sourceSheetID = mappings[sheetName]
	)

	httpClient, err := c.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	containsSheet, err := c.containsSheet(ctx, spreadsheetsSvc, spreadsheetID, sheetName)
	if err != nil {
		return err
	}

	if containsSheet {
		return domain.ErrorSheetPresent
	}

	sheetID, err := c.copySheet(ctx, spreadsheetsSvc, c.originSpreadsheetID, spreadsheetID, sourceSheetID)
	if err != nil {
		return err
	}

	dataToCopy, err := c.getDataToCopy(ctx, spreadsheetsSvc, c.originSpreadsheetID, sheetID, sourceSheetID)
	if err != nil {
		return err
	}

	batchUpdate := NewBatchUpdate(spreadsheetsSvc)
	{
		batchUpdate.WithProtectedRange(sheetID, dataToCopy.protectedRanges)
		batchUpdate.WithSheetName(sheetID, sheetName)
		batchUpdate.WithRequest(dataToCopy.updateCellRequests...)
	}

	if err := batchUpdate.Do(ctx, spreadsheetID); err != nil {
		return err
	}

	return nil
}

func (c *SpreadsheetServiceGoogle) containsSheet(ctx context.Context, svc *sheets.Service, spreadsheetID string, sheetName string) (bool, error) {
	spreadsheet, err := svc.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return false, err
	}

	// Iterate through the sheets and check if the sheet name exists
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			return true, nil
		}
	}

	return false, nil
}

// copySheet returns sheetId
func (c *SpreadsheetServiceGoogle) copySheet(ctx context.Context, svc *sheets.Service, sourceSpreadsheetID, targetSpreadsheetID string, sheetID int64) (int64, error) {
	copyRequest := &sheets.CopySheetToAnotherSpreadsheetRequest{
		DestinationSpreadsheetId: targetSpreadsheetID,
	}

	resp, err := svc.Spreadsheets.Sheets.CopyTo(sourceSpreadsheetID, sheetID, copyRequest).Context(ctx).Do()
	if err != nil {
		return 0, err
	}

	return resp.SheetId, nil
}

type getDataToCopyResponse struct {
	protectedRanges    []*sheets.ProtectedRange
	updateCellRequests []*sheets.Request
}

func (c *SpreadsheetServiceGoogle) getDataToCopy(ctx context.Context, svc *sheets.Service, spreadsheetID string, destinationSheetID, sourceSheetID int64) (*getDataToCopyResponse, error) {
	spreadsheet, err := svc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	var sheet *sheets.Sheet

	for i, tmpSheet := range spreadsheet.Sheets {
		if tmpSheet.Properties.SheetId == sourceSheetID {
			sheet = spreadsheet.Sheets[i]
			break
		}
	}

	var (
		protectedRanges    []*sheets.ProtectedRange = sheet.ProtectedRanges
		updateCellRequests []*sheets.Request
	)

	// adapt data validations of range
	for rowIdx, row := range sheet.Data[0].RowData {
		for cellIdx, cell := range row.Values {
			if cell.DataValidation != nil && cell.DataValidation.Condition.Type == "ONE_OF_RANGE" {
				updateCellRequests = append(updateCellRequests, &sheets.Request{
					UpdateCells: &sheets.UpdateCellsRequest{
						Start: &sheets.GridCoordinate{
							RowIndex:    int64(rowIdx),
							ColumnIndex: int64(cellIdx),
							SheetId:     destinationSheetID,
						},
						Fields: "dataValidation",
						Rows: []*sheets.RowData{
							{
								Values: []*sheets.CellData{
									{
										DataValidation: sheet.Data[0].RowData[rowIdx].Values[cellIdx].DataValidation,
									},
								},
							},
						},
					},
				})
			}
		}
	}

	// adapt formulas
	for rowIdx, row := range sheet.Data[0].RowData {
		for cellIdx, cell := range row.Values {
			if cell.UserEnteredValue != nil && cell.UserEnteredValue.FormulaValue != nil {
				updateCellRequests = append(updateCellRequests, &sheets.Request{
					UpdateCells: &sheets.UpdateCellsRequest{
						Start: &sheets.GridCoordinate{
							RowIndex:    int64(rowIdx),
							ColumnIndex: int64(cellIdx),
							SheetId:     destinationSheetID,
						},
						Fields: "userEnteredValue",
						Rows: []*sheets.RowData{
							{
								Values: []*sheets.CellData{
									{
										UserEnteredValue: &sheets.ExtendedValue{
											FormulaValue: sheet.Data[0].RowData[rowIdx].Values[cellIdx].UserEnteredValue.FormulaValue,
										},
									},
								},
							},
						},
					},
				})
			}
		}
	}

	result := &getDataToCopyResponse{
		protectedRanges:    protectedRanges,
		updateCellRequests: updateCellRequests,
	}

	return result, nil
}
