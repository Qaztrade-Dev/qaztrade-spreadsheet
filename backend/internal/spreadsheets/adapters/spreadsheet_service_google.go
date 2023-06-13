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

// func (s *SpreadsheetServiceGoogle) Test(ctx context.Context) error {
// 	httpClient, err := s.oauth2.GetClient(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	driveSvc, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
// 	if err != nil {
// 		return err
// 	}

// 	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
// 	if err != nil {
// 		return err
// 	}

// 	query := fmt.Sprintf("mimeType!='application/vnd.google-apps.folder' and trashed = false and '%s' in parents", s.destinationFolderID)
// 	fileListCall := driveSvc.Files.List().Q(query).Fields("nextPageToken, files(id, name)").OrderBy("createdTime asc")

// 	spreadsheetIDs := make([]string, 0)
// 	err = fileListCall.Pages(ctx, func(filesList *drive.FileList) error {
// 		for _, file := range filesList.Files {
// 			spreadsheetIDs = append(spreadsheetIDs, file.Id)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	// spreadsheetIDs := []string{"1NoaQE9dWeMHF368EheuhFL2ttZueDqBc0If_Sil40bg"}

// 	replaceValues := []string{
// 		"40210", "40320", "40390", "40410", "40490", "40510", "40520", "40610", "40620", "40630", "40640",
// 		"40690", "20110", "20120", "20130", "20210", "20220", "20230", "20311", "20312", "20319", "20321",
// 		"20322", "20329", "20410", "20421", "20422", "20423", "20430", "20441", "20442", "20443", "20450",
// 		"20500", "20610", "20621", "20622", "20629", "20630", "20641", "20649", "20680", "20690", "20711",
// 		"20712", "20713", "20714", "20724", "20725", "20726", "20727", "20741", "20742", "20743", "20744",
// 		"20745", "20751", "20752", "20753", "20754", "20755", "20760", "20810", "20830", "20840", "20850",
// 		"20860", "20890", "20910", "20990", "21011", "21012", "21019", "21020", "21091", "21092", "21093",
// 		"21099", "30291", "30311", "30314", "30319", "30323", "30324", "30325", "30326", "30329", "30331",
// 		"30332", "30333", "30334", "30339", "30341", "30342", "30343", "30344", "30345", "30346", "30349",
// 		"30351", "30353", "30354", "30355", "30356", "30359", "30363", "30364", "30365", "30366", "30367",
// 		"30368", "30369", "30381", "30382", "30383", "30384", "30389", "30391", "30392", "30399", "30431",
// 		"30432", "30433", "30439", "30441", "30442", "30443", "30444", "30445", "30446", "30447", "30448",
// 		"30449", "30451", "30452", "30453", "30454", "30455", "30456", "30457", "30459", "30461", "30462",
// 		"30463", "30469", "30471", "30472", "30473", "30474", "30475", "30479", "30481", "30482", "30483",
// 		"30484", "30485", "30486", "30487", "30488", "30489", "30491", "30492", "30493", "30494", "30495",
// 		"30496", "30497", "30499", "30520", "30531", "30532", "30539", "30541", "30542", "30543", "30544",
// 		"30549", "30551", "30552", "30553", "30554", "30559", "30561", "30562", "30563", "30564", "30569",
// 		"30571", "30572", "30579", "30611", "30612", "30614", "30615", "30616", "30617", "30619", "30719",
// 		"30722", "30729", "30732", "30739", "30743", "30749", "30752", "30759", "30760", "30772", "30779",
// 		"30783", "30784", "30787", "30788", "30792", "30799", "30812", "30819", "30822", "30829", "30910",
// 		"30990", "40110", "40120", "40140", "40150", "40221", "40229", "40291", "40299", "40590", "40790",
// 		"40811", "40819", "40891", "40899", "50210", "50290", "50400", "50510", "50590", "50610", "50690",
// 		"50710", "50790", "51000", "51191", "51199", "71010", "71021", "71022", "71029", "71030", "71040",
// 		"71080", "71090", "71120", "71140", "71151", "71159", "71190", "71220", "71231", "71232", "71233",
// 		"71234", "71239", "71290", "80620", "81110", "81120", "81190", "81210", "81290", "81310", "81320",
// 		"81330", "81340", "81350", "81400", "90112", "90121", "90122", "90190", "90210", "90230", "90412",
// 		"90422", "90520", "90620", "90720", "90812", "90822", "90832", "90922", "90932", "90962", "91012",
// 	}
// 	tnvedRegex := strings.Join(replaceValues, "|")

// 	for i, spreadsheetID := range spreadsheetIDs {
// 		spreadsheetID := spreadsheetID

// 		// Define the new format
// 		cellFormat := &sheets.CellFormat{
// 			NumberFormat: &sheets.NumberFormat{
// 				Type: "TEXT",
// 			},
// 		}

// 		batch := NewBatchUpdate(spreadsheetsSvc)
// 		batch.WithRequest(&sheets.Request{
// 			RepeatCell: &sheets.RepeatCellRequest{
// 				Range: &sheets.GridRange{
// 					SheetId:          292380577,
// 					StartRowIndex:    0,
// 					StartColumnIndex: 2,
// 					EndColumnIndex:   3,
// 				},
// 				Cell:   &sheets.CellData{UserEnteredFormat: cellFormat},
// 				Fields: "userEnteredFormat(numberFormat)",
// 			},
// 		})

// 		batch.WithRequest(&sheets.Request{
// 			FindReplace: &sheets.FindReplaceRequest{
// 				Find:          "^(\\d{5})$",
// 				Replacement:   "0$1",
// 				SheetId:       292380577,
// 				SearchByRegex: true,
// 				MatchCase:     true,
// 			},
// 		})

// 		batch.WithRequest(&sheets.Request{
// 			FindReplace: &sheets.FindReplaceRequest{
// 				Find:          "^('*)(\\d{6})$",
// 				Replacement:   "'$2",
// 				SheetId:       292380577,
// 				SearchByRegex: true,
// 				MatchCase:     true,
// 			},
// 		})

// 		batch.WithRequest(&sheets.Request{
// 			FindReplace: &sheets.FindReplaceRequest{
// 				Find:          "^('+)(\\d{6})$",
// 				Replacement:   "$2",
// 				SheetId:       292380577,
// 				SearchByRegex: true,
// 				MatchCase:     true,
// 			},
// 		})

// 		batch.WithRequest(&sheets.Request{
// 			FindReplace: &sheets.FindReplaceRequest{
// 				Find:          fmt.Sprintf("^(%s)\\s-", tnvedRegex),
// 				Replacement:   "0$1 -",
// 				AllSheets:     true,
// 				SearchByRegex: true,
// 				MatchCase:     true,
// 			},
// 		})

// 		if err := batch.Do(ctx, spreadsheetID); err != nil {
// 			return err
// 		}

// 		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
// 	}

// 	return nil
// }
