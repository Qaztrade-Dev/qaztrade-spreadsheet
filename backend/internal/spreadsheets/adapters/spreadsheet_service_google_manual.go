package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (s *SpreadsheetServiceGoogle) FormatTextCells(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	driveSvc, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	query := fmt.Sprintf("mimeType!='application/vnd.google-apps.folder' and trashed = false and '%s' in parents", s.destinationFolderID)
	fileListCall := driveSvc.Files.List().Q(query).Fields("nextPageToken, files(id, name)").OrderBy("createdTime asc")

	spreadsheetIDs := make([]string, 0)
	err = fileListCall.Pages(ctx, func(filesList *drive.FileList) error {
		for _, file := range filesList.Files {
			spreadsheetIDs = append(spreadsheetIDs, file.Id)
		}
		return nil
	})
	if err != nil {
		return err
	}
	// spreadsheetIDs := []string{"1NoaQE9dWeMHF368EheuhFL2ttZueDqBc0If_Sil40bg"}

	replaceValues := []string{
		"40210", "40320", "40390", "40410", "40490", "40510", "40520", "40610", "40620", "40630", "40640",
		"40690", "20110", "20120", "20130", "20210", "20220", "20230", "20311", "20312", "20319", "20321",
		"20322", "20329", "20410", "20421", "20422", "20423", "20430", "20441", "20442", "20443", "20450",
		"20500", "20610", "20621", "20622", "20629", "20630", "20641", "20649", "20680", "20690", "20711",
		"20712", "20713", "20714", "20724", "20725", "20726", "20727", "20741", "20742", "20743", "20744",
		"20745", "20751", "20752", "20753", "20754", "20755", "20760", "20810", "20830", "20840", "20850",
		"20860", "20890", "20910", "20990", "21011", "21012", "21019", "21020", "21091", "21092", "21093",
		"21099", "30291", "30311", "30314", "30319", "30323", "30324", "30325", "30326", "30329", "30331",
		"30332", "30333", "30334", "30339", "30341", "30342", "30343", "30344", "30345", "30346", "30349",
		"30351", "30353", "30354", "30355", "30356", "30359", "30363", "30364", "30365", "30366", "30367",
		"30368", "30369", "30381", "30382", "30383", "30384", "30389", "30391", "30392", "30399", "30431",
		"30432", "30433", "30439", "30441", "30442", "30443", "30444", "30445", "30446", "30447", "30448",
		"30449", "30451", "30452", "30453", "30454", "30455", "30456", "30457", "30459", "30461", "30462",
		"30463", "30469", "30471", "30472", "30473", "30474", "30475", "30479", "30481", "30482", "30483",
		"30484", "30485", "30486", "30487", "30488", "30489", "30491", "30492", "30493", "30494", "30495",
		"30496", "30497", "30499", "30520", "30531", "30532", "30539", "30541", "30542", "30543", "30544",
		"30549", "30551", "30552", "30553", "30554", "30559", "30561", "30562", "30563", "30564", "30569",
		"30571", "30572", "30579", "30611", "30612", "30614", "30615", "30616", "30617", "30619", "30719",
		"30722", "30729", "30732", "30739", "30743", "30749", "30752", "30759", "30760", "30772", "30779",
		"30783", "30784", "30787", "30788", "30792", "30799", "30812", "30819", "30822", "30829", "30910",
		"30990", "40110", "40120", "40140", "40150", "40221", "40229", "40291", "40299", "40590", "40790",
		"40811", "40819", "40891", "40899", "50210", "50290", "50400", "50510", "50590", "50610", "50690",
		"50710", "50790", "51000", "51191", "51199", "71010", "71021", "71022", "71029", "71030", "71040",
		"71080", "71090", "71120", "71140", "71151", "71159", "71190", "71220", "71231", "71232", "71233",
		"71234", "71239", "71290", "80620", "81110", "81120", "81190", "81210", "81290", "81310", "81320",
		"81330", "81340", "81350", "81400", "90112", "90121", "90122", "90190", "90210", "90230", "90412",
		"90422", "90520", "90620", "90720", "90812", "90822", "90832", "90922", "90932", "90962", "91012",
	}
	tnvedRegex := strings.Join(replaceValues, "|")

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID

		// Define the new format
		cellFormat := &sheets.CellFormat{
			NumberFormat: &sheets.NumberFormat{
				Type: "TEXT",
			},
		}

		batch := NewBatchUpdate(spreadsheetsSvc)
		batch.WithRequest(&sheets.Request{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          292380577,
					StartRowIndex:    0,
					StartColumnIndex: 2,
					EndColumnIndex:   3,
				},
				Cell:   &sheets.CellData{UserEnteredFormat: cellFormat},
				Fields: "userEnteredFormat(numberFormat)",
			},
		})

		batch.WithRequest(&sheets.Request{
			FindReplace: &sheets.FindReplaceRequest{
				Find:          "^(\\d{5})$",
				Replacement:   "0$1",
				SheetId:       292380577,
				SearchByRegex: true,
				MatchCase:     true,
			},
		})

		batch.WithRequest(&sheets.Request{
			FindReplace: &sheets.FindReplaceRequest{
				Find:          "^('*)(\\d{6})$",
				Replacement:   "'$2",
				SheetId:       292380577,
				SearchByRegex: true,
				MatchCase:     true,
			},
		})

		batch.WithRequest(&sheets.Request{
			FindReplace: &sheets.FindReplaceRequest{
				Find:          "^('+)(\\d{6})$",
				Replacement:   "$2",
				SheetId:       292380577,
				SearchByRegex: true,
				MatchCase:     true,
			},
		})

		batch.WithRequest(&sheets.Request{
			FindReplace: &sheets.FindReplaceRequest{
				Find:          fmt.Sprintf("^(%s)\\s-", tnvedRegex),
				Replacement:   "0$1 -",
				AllSheets:     true,
				SearchByRegex: true,
				MatchCase:     true,
			},
		})

		if err := batch.Do(ctx, spreadsheetID); err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
	}

	return nil
}

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

		batch := NewBatchUpdate(spreadsheetsSvc)
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

func (s *SpreadsheetServiceGoogle) AddLockNamedRanges(ctx context.Context) error {
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

	type skelRange struct {
		Start int
		End   int
	}

	skelNamedRanges := map[string]([]*skelRange){
		"Затраты на доставку транспортом": []*skelRange{
			{Start: 20, End: 49},
			{Start: 123, End: 124},
		},
		"Затраты на сертификацию предприятия": []*skelRange{
			{Start: 13, End: 47},
			{Start: 57, End: 58},
		},
		"Затраты на рекламу ИКУ за рубежом": []*skelRange{
			{Start: 13, End: 47},
			{Start: 58, End: 59},
		},
		"Затраты на перевод каталога ИКУ": []*skelRange{
			{Start: 13, End: 47},
			{Start: 56, End: 57},
		},
		"Затраты на аренду помещения ИКУ": []*skelRange{
			{Start: 13, End: 47},
			{Start: 53, End: 54},
		},
		"Затраты на сертификацию ИКУ": []*skelRange{
			{Start: 13, End: 47},
			{Start: 59, End: 60},
		},
		"Затраты на демонстрацию ИКУ": []*skelRange{
			{Start: 13, End: 47},
			{Start: 53, End: 54},
		},
		"Затраты на франчайзинг": []*skelRange{
			{Start: 13, End: 47},
			{Start: 53, End: 54},
		},
		"Затраты на соответствие товаров требованиям": []*skelRange{
			{Start: 13, End: 47},
			{Start: 63, End: 64},
		},
		"Затраты на регистрацию товарных знаков": []*skelRange{
			{Start: 13, End: 47},
			{Start: 65, End: 66},
		},
		"Затраты на аренду": []*skelRange{
			{Start: 13, End: 47},
			{Start: 53, End: 54},
		},
		"Затраты на перевод": []*skelRange{
			{Start: 13, End: 47},
			{Start: 56, End: 57},
		},
		"Затраты на рекламу товаров за рубежом": []*skelRange{
			{Start: 13, End: 47},
			{Start: 58, End: 59},
		},
		"Затраты на участие в выставках": []*skelRange{ // TODO
			{Start: 13, End: 47},
			{Start: 72, End: 73},
		},
	}

	_ = skelNamedRanges

	spreadsheetIDs := []string{"1jcpcuLBCvYZ_aWi70-nz4a40afZTArz0CXpdNh86uGA"}

	// // master spreadsheet
	// spreadsheetIDs := []string{"1oMJFttuiPxoBdejx3Ul3D2nscE0x45oeNfo-upVe1gE"}

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).Do()
		if err != nil {
			return err
		}

		batch := NewBatchUpdate(spreadsheetsSvc)

		// for _, namedRange := range spreadsheet.NamedRanges {
		// 	if strings.Contains(namedRange.Name, "_blocked_") {
		// 		fmt.Printf("%#v\n", namedRange.NamedRangeId)
		// 		fmt.Printf("%#v\n", namedRange.Range)
		// 		batch.WithRequest(
		// 			&sheets.Request{
		// 				DeleteNamedRange: &sheets.DeleteNamedRangeRequest{
		// 					NamedRangeId: namedRange.NamedRangeId,
		// 				},
		// 			},
		// 		)
		// 	}
		// }

		for skelSheetName, skelRanges := range skelNamedRanges {
			skelSheetName := skelSheetName

			for j := range skelRanges {
				var (
					skelRange                    = skelRanges[j]
					namedRangeName               = fmt.Sprintf("%s_blocked_%d", strings.ReplaceAll(skelSheetName, " ", "_"), j)
					sheet          *sheets.Sheet = nil
				)

				for k, tmpSheet := range spreadsheet.Sheets {
					tmpSheet := tmpSheet
					if tmpSheet.Properties.Title != skelSheetName {
						continue
					}

					sheet = spreadsheet.Sheets[k]
					break
				}

				if sheet == nil {
					continue
				}

				_ = skelRange
				_ = namedRangeName

				batch.WithRequest(
					&sheets.Request{
						AddNamedRange: &sheets.AddNamedRangeRequest{
							NamedRange: &sheets.NamedRange{
								Name: namedRangeName,
								Range: &sheets.GridRange{
									SheetId:          sheet.Properties.SheetId,
									StartColumnIndex: int64(skelRange.Start),
									EndColumnIndex:   int64(skelRange.End + 1),
									EndRowIndex:      1,
								},
							},
						},
					},
				)
			}
		}

		if err := batch.Do(ctx, spreadsheetID); err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
	}

	return nil
}

func (s *SpreadsheetServiceGoogle) Test(ctx context.Context) error {
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

	spreadsheetIDs := []string{"1jcpcuLBCvYZ_aWi70-nz4a40afZTArz0CXpdNh86uGA"}

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).Do()
		if err != nil {
			return err
		}

		batch := NewBatchUpdate(spreadsheetsSvc)

		for _, namedRange := range spreadsheet.NamedRanges {
			namedRange := namedRange
			if strings.Contains(namedRange.Name, "_blocked_") {
				batch.WithRequest(
					&sheets.Request{
						AddProtectedRange: &sheets.AddProtectedRangeRequest{
							ProtectedRange: &sheets.ProtectedRange{
								Description: namedRange.Name,
								Editors: &sheets.Editors{
									Users: []string{
										"qaztrade.export@gmail.com",
									},
								},
								Range: &sheets.GridRange{
									SheetId:          namedRange.Range.SheetId,
									StartColumnIndex: namedRange.Range.StartColumnIndex,
									EndColumnIndex:   namedRange.Range.EndColumnIndex,
								},
							},
						},
					},
				)
			}
		}

		if err := batch.Do(ctx, spreadsheetID); err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
	}

	return nil
}
