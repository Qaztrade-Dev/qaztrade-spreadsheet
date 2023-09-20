package adapters

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	spreadsheetsAdapters "github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (s *SpreadsheetServiceGoogle) AddParityCheckColumn(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	// driveSvc, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
	// if err != nil {
	// 	return err
	// }

	// query := fmt.Sprintf("mimeType!='application/vnd.google-apps.folder' and trashed = false and '%s' in parents", s.destinationFolderID)
	// fileListCall := driveSvc.Files.List().Q(query).Fields("nextPageToken, files(id, name)").OrderBy("createdTime asc")

	// spreadsheetIDs := make([]string, 0)
	// err = fileListCall.Pages(ctx, func(filesList *drive.FileList) error {
	// 	for _, file := range filesList.Files {
	// 		spreadsheetIDs = append(spreadsheetIDs, file.Id)
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	return err
	// }

	// Учитывать текущий оффсет
	// spreadsheetIDs = spreadsheetIDs[345:]

	// spreadsheetIDs := []string{"1VHyaTWbvltfshMoPemkTjZEuYdbl95Yd_0-tont6c7c"}

	// // master spreadsheet
	spreadsheetIDs := []string{"1oMJFttuiPxoBdejx3Ul3D2nscE0x45oeNfo-upVe1gE"}

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Do()
		if err != nil {
			return err
		}

		var sheet *sheets.Sheet = nil
		for i, tmpSheet := range spreadsheet.Sheets {
			// if tmpSheet.Properties.Title != "Затраты на доставку транспортом" {
			if tmpSheet.Properties.Title != "⏳ (ожидайте) Затраты на доставку транспортом" {
				continue
			}
			sheet = spreadsheet.Sheets[i]
			break
		}

		if sheet == nil {
			continue
		}

		batch := spreadsheetsAdapters.NewBatchUpdate(spreadsheetsSvc)
		batch.WithRequest(&sheets.Request{
			InsertDimension: &sheets.InsertDimensionRequest{
				Range: &sheets.DimensionRange{
					SheetId:    sheet.Properties.SheetId,
					Dimension:  "COLUMNS",
					StartIndex: 82,
					EndIndex:   83,
				},
				InheritFromBefore: true,
			},
		})

		batch.WithRequest(&sheets.Request{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          sheet.Properties.SheetId,
					StartRowIndex:    3,
					StartColumnIndex: 82,
					EndColumnIndex:   83,
				},
				Fields: "*",
				Cell: &sheets.CellData{
					UserEnteredValue: &sheets.ExtendedValue{
						StringValue: nil,
					},
					UserEnteredFormat: &sheets.CellFormat{
						NumberFormat: &sheets.NumberFormat{
							Type: "TEXT",
						},
						BackgroundColor: &sheets.Color{
							Red:   243.0 / 255.0,
							Green: 243.0 / 255.0,
							Blue:  243.0 / 255.0,
							Alpha: 1.0,
						},
					},
					DataValidation: nil,
				},
			},
		})

		batch.WithRequest(
			&sheets.Request{
				UpdateCells: &sheets.UpdateCellsRequest{
					Range: &sheets.GridRange{
						SheetId:          sheet.Properties.SheetId,
						StartRowIndex:    3, // 4th row
						EndRowIndex:      4,
						StartColumnIndex: 82,
						EndColumnIndex:   83,
					},
					Fields: "userEnteredValue",
					Rows: []*sheets.RowData{
						{
							Values: []*sheets.CellData{
								{
									UserEnteredValue: &sheets.ExtendedValue{
										FormulaValue: aws.String(`=ARRAYFORMULA(IFS(CC4:CC=""; ""; CU4:CU=""; ""; (CC4:CC=CU4:CU)*(CD4:CD=CV4:CV)=1; "✓ да"; (CC4:CC<>CU4:CU)+(CD4:CD<>CV4:CV)>=1; "✗ нет"))`),
									},
								},
							},
						},
					},
				},
			},
			&sheets.Request{
				UpdateCells: &sheets.UpdateCellsRequest{
					Range: &sheets.GridRange{
						SheetId:          sheet.Properties.SheetId,
						StartRowIndex:    2, // 3th row
						EndRowIndex:      3,
						StartColumnIndex: 82,
						EndColumnIndex:   83,
					},
					Fields: "userEnteredValue",
					Rows: []*sheets.RowData{
						{
							Values: []*sheets.CellData{
								{
									UserEnteredValue: &sheets.ExtendedValue{
										StringValue: aws.String(`соответствие № и даты контракта поставки данным столбца контрактов поставки`),
									},
								},
							},
						},
					},
				},
			},
			&sheets.Request{
				AddProtectedRange: &sheets.AddProtectedRangeRequest{
					ProtectedRange: &sheets.ProtectedRange{
						Description: "parity_delivery_contract",
						Editors: &sheets.Editors{
							Users: []string{
								"qaztrade.export@gmail.com",
								s.svcAccount,
							},
						},
						Range: &sheets.GridRange{
							SheetId:          sheet.Properties.SheetId,
							StartColumnIndex: 82,
							EndColumnIndex:   83,
						},
					},
				},
			},
		)
		if err := batch.Do(ctx, spreadsheetID); err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
	}

	return nil
}
