package adapters

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

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

	driveSvc, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
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
		"Затраты на участие в выставках": []*skelRange{
			{Start: 13, End: 47},
			{Start: 72, End: 73},
		},
	}

	_ = skelNamedRanges
	spreadsheetIDs = spreadsheetIDs[519:]

	// spreadsheetIDs := []string{"1jcpcuLBCvYZ_aWi70-nz4a40afZTArz0CXpdNh86uGA"}

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
					if tmpSheet.Properties.Title != "⏳ (ожидайте) "+skelSheetName {
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

func (s *SpreadsheetServiceGoogle) ActivateBlockedRanges(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	spreadsheetIDs := []string{
		"1c1yJLIt3bxbJc8ykv12RkUSjngGOI-Z_8qoaTygwhqo", "18Pc0NTmkr9kPExWtCFmlWtRggySMzRl5PGeAbaAprSo", "1JkK7renOvJ-mDTojLP_k1bjsu8sTRpr9MyERZHq1Jt0", "1UvEbSrUtN-EXY7Qc4gAXWtmzQyestRlw2H97DyCNpjw",
		"1znzdQsKIHqVEx9J_9n8iOzLs67OyPP8uCDXzljmDdT0", "1ZUPZSzElxz50owyBRBunxpi2AhoBwKmWEa3xWfHWZGg", "1Xf8bvGBt63w6TS2G3k3EN-R7c8ujSipcUSdy2fHU2MY", "1cesfjKCRezQoGSzX_fj_y8H2rO_5BQWaPG0JQ9h7hp8",
		"1rw1rqoTWHGEjTNkQt0iL1OZGyTtXsVarAPV3J16uwQo", "1nVTgSJCalLDfR2prFgdYV2IydMXcGiY4o9Fv3aCaSYg", "1r7EzdxwCfcCGD9MQzekx0yEnghofppMDQQtXXj_KSzs", "162EpE6JVTkkY4pE9XuhjTuNuC_bz6LHvTN7I5ReXQnQ",
		"18K_1Ls2UthlndfHK4YkQQtvzsvNEKAnhYw70qkN_xog", "1HsEdbx-BxLn-k5fAiOmY5dl8xazDbE4fe0E6dIx2mpA", "1eLCHxSZW-YZwgy_erhAd-kodEfLoOmHOzgHj1CHgEus", "14btaJqWcSWqV5PSaXED6kkSDXibvurGialcwkyZAGCU",
		"1WfMOr_QQhxZm6lpDSjIee-p99euhIjCBKkEbVTw4RbY", "1X1C3YS9BfPAmldenCs4ij3qPVZo_fYt9LwcM9sDKLrc", "1Ta1y1660bzM8YpXiFly_cVW9AGQgQCD25TCew_vkfyU", "1lYpH3QpJOAio9jD7HQ1GA5cIgBRYNuIX-TH2gDmGlcM",
		"1TZwHYSI3975St4bld0yG2fngEmsXoLH0afOLlzAciQs", "1VazpGDgp7IQsXuoGzDuWfG20eZuXTmluClb4TiCbWUA", "1qzB8mLsvc0dl9yC3gTskpmnubQwn__1ipT89W_s32YA", "10ZwKF9DB0bD3qEjJdKjGjM5ZIo3uc1NDYQBmYOPwlhc",
		"1AisI_Xz0PgE0OuuMWFF2WPK73RJzPRdG3BFRi7bHCJk", "1jHhf9zF-nPBh-rnQSpoXk1g2n830X6aXwbjW5k4wqTU", "1tEaZczQMwLZbcY6j9p4JzrRl0vTRh3MUNWGWUbj7Uv0", "1qihHg_NZ6G70140KMfi3FwDLBi5B2rQm1SfGdM5wPgs",
		"1sAPLUcbNqZtw7W1eoOu82oxDcqnM5uAxIx_-YHIopEU", "1Jp4sTkvvEumPFuZA7fpa1Z1U3R9i2bI9V3-z60duNWg", "1pfUsE5MWHMMocAqKHTo3wYHYlgIoluoN5DZsAs9k9vk", "12r7wG6tT5prPVuHuk92cax9kICFMRK1t165clZv18W4",
		"1kEnY1tcV3mEjfBXAFOypqTvwmbIYgnECzQVRZkAgDso", "1fvfFht0bNsKkHHQztWQKNq-AjEVu8Sj5HVdn5irwfnM", "1V0p80kryN9FpW6Hpw4f6RfqXLsBPaF9dsXCVV2Ljm4U", "1oGUL6JYfoMEo6DJNukyl6neFVvYcbZwVt6GfDhrcwjs",
		"1vnac2RXonPGLvaM697WfeZXjTjBXY0SNx8jm14T7jE0", "1y-N0SKgI5L-qJCQ9k4EhJN-PAQ61uJSOjhCtUKrH94g", "1k9Ray-9O9cvZullkBq0rLR0r8xw_mwlu6AVP7ieOiZM", "1zzu4rB20EjQbQXzhYgB1iGso59H-szZ8l8i4IENOgm8",
		"1Kc6M4K6R-PGO0tiki4LpHDY20S4TvjODDfG-39L8_vI", "1QzIHgTSS6ZFJGo_cN_JRCnJVi1m0xjxFoYE_8ux6QaI", "1xyulTOre1Tik4bbtTQAk-Jk36y-JeRSGhf8xX0Zs1-k", "1p50HLUkBpACGIsiBCGKKUpTr21tXHShYSuYQlo-Hs-k",
		"1XezdKLE4WBgmFGMzYwhIxbZEWjRYekOhdeeMDBS8s9U", "19FGj2bC1UtqcnZy3VCwIRaAmbqEHkVQMhXIoy1RonVY", "1vbnYblfnyJ8OJ7bzywTeVXDZQoiBXHjD1WM5WXfFKPw", "1GGrPEs5gbl1r-MErgJM6Sx5FCVrSp6U5rmD8ILJk_PE",
		"1g2Rl-4tCfhtn7_HpXrPkLiC60qRYISuBr9XFWOUvWYc", "115px4QhFpsSRW40rtHkLK5IxerzoSB7fCgpbYnN7bcs", "1fTLQrAmiHHADHnyc2xRREFjZlaaHLs4OQ8BY38axJgI", "1Qddmwvv0CHmgVxP3N3KdLNvxDMA8yEypFu7oRlS8N_0",
		"1MIxbo89FzVwy14VLhK7N0OGxofoo9_rIUNVeDI9F9NI", "1n0WapvYle1aSCTCg5C3LH9snIxhJL5zvm4BvwcOkrnU", "12FKiNdd8cudtAEdNCsoNIE0xlCoYNqj7tnqUAv9lWME", "1hyjli8bQB_hDervJ796blv7PrRVE5DLsX6E_QRmSp0U",
		"1k1qv_oX7EicScjWRvYe3-eKjLMmkyyz3-g3D83VWGe4", "1mQTnR9Itbyf6VX8DGj0zvR4F9QXOXOVVkRp9ODhIngE", "1hxAvB6t5wYf7GZYlBDhUeg_zZQ0RPTGR9jZblMp7maY", "1-yaMCDLWfrZPciiJlBzKOcAhQZlS9TEP7w8SnhtSC8U",
		"1dnxGao23OFWdT0v_ry5AcGn9Rrk7DQP62KO6Kg5Xymw", "1KMyJp-JVjkNPJ1vNzkg18dNusvYRj6T4XZTTAJ6Eiuo", "1t8W2lQEz64jDta5XGK3LSb-c36GOmZ-arMCXx6fBzRU", "1J0dNd4yxxa-DWmUIWUNDdnJjVtYIYsRpL98JZdZLNUM",
		"1pa0rVaX3sK4pZSu9o9RCrQwqHZpDXDjrlgztjjLAOfg", "1X67Olsph3fEqxsDpwETWj5UvzWxpLuUlrrC1mww1HZI", "1OdfMSikfzh6NvAq2rs1bFVWgs9PWVtUBCN4gCvZgW_I", "15ez9Hs8SF0RMhcAPRnzs7ptcU0j0_fYesJLOvZA3MNo",
		"1baKhMgeoNEH5TEe4H0fE5gnJPCgn_8wcJw_lCxmP8CU", "1mPvZCp4RrtMjGQprmSjNCmr9GiklrFwVJOfZNWYiDaI", "1ZICQd7hfJdVtiVsuSAew2dcPy9NXzf0zi0Xx1ulJ8SA", "119ygHqSgwjPwOekTQm2RayiCPQ00LnUtZQXjUHgQEEA",
		"1yqqhxZHfujAkUrnRW11MetPmBpzp-bMaN7Okk6y60-U", "1qGoWvPjN5w435Zm_lkt1GbcUM2VZZ5w_b5Gl_6i_ou8", "1MdqjufVjUA6pJomm1FjZ5g6s4RHoxSlbYX7N0vzswfo", "1dK6trumK_tf1Amy3jEBnHPD9Oabmv87J9arsPIgEvtY",
		"1k5yP6XW11nUVnS-LioJE18FbxH1S_GvhNJ8j9DSdsaI", "1bchRvC25IvRL12EZ_a0Y5ZIXtBm3gUIrxtlYr1Pqzlw", "1Y2nfJqPrPA7J-7d1cxHfRTE5wdkEQ2wg4oK-uH-kVLY", "1nHk5OJbeIGs2l4MHJcPI7fMYtNFddxSqbEuwBDHtwJ0",
		"1IPQx2bCbJ6TckvqI1LVgwUugb4Vj4HVLW53QT82fdTg", "19OoAnsYIInnu3WdByKb0oEQ7zC1i2BvqTSAbdlnp43Y", "1xfKRzSYLtN0GeLyIoPbu5xaI8foFpFbEZWIdX68_VkI", "1apEzPvVSgxPZshTz1TGF2QNgosZE2C1jBw0j_wjO3lQ",
		"1JHavaD9ZrrLDYZotJVd9N1B9eAMIa9XzGnZVay2hmJU", "10mytod7TsGLZJ6vASAid3dnpa6OFoZ4ZujQW5NpqUdc", "1wE4JGBT7uaZvt4J4VOucv9pnDsUxw0OwQK__IgQcQ0k", "1sL0iNHrZ4BKnXuSWAkQImz0hBwQPLLl7xda8aDHS888",
		"14ipQ1hxSoyTpkHQGtg4DK-KMCJJV9sDbR7rOeLbkJyI", "1wJvOZB215ijgALSi1NlMZRnQnWwkmQETMlh3X3tfxVg", "1Rt047s1JCdRqnlaBZ9af5jobcBK-V9QSZjqGyt3cv6E", "1cjMRYw1v5MxRoAuldhEXGjAiNhQr36VmttjhGSECjs0",
		"1CYF8AOL2dI2yY-JhqxZzEnT50BK5Op_i228By6NsRWI", "1uzkWstoM-sDWNzaNaaQecOg7DjpqEs3A2LHRl6R3Qv8", "1d5Ndtdptd_Jcstm3jEeVsx31h61qP0EghJygL2EaFzk", "1VHyaTWbvltfshMoPemkTjZEuYdbl95Yd_0-tont6c7c",
		"13BNidTuWqh4WvPsv-hFZhLtZWrITqw1T2lbG3UjRmQQ", "1IE8e7FM7CUD_gskwbGqyka43KFWfpSAvPX2vllUQmDI", "1ZPft2_E-4eXgdtByOc2k1YP2lf-7NWzPKjRW0-ZVCEI", "10TA83SY50q13Wp7Ip0-Y8j4ugk3ZEhxFLy8rbra3Eq4",
		"1BurOnFFSOCFcMFSpjnZzzBJVNlhPkszeH43erAv1Ajg", "1245aNoMJBqb4wTVgXe3Qo0EuD2f1vTxOIrLctCetrsQ", "1rrD7WSY9pSr4lsjpnJ_Qo1dRxkqyy33coVQyXI5cvro", "1cD-eh5fy2wY6zdSzN9bLm-rNJ8IhVTHJcqseY9b8eR0",
		"1uBG1zftgHn_uhYmXPgk0ClyAn7zmHA8PCeIXYe3JxaE", "190U6SzY8nd1gpRcH5nhlkwmyq5A7W0rLABuBIZhOVyI", "1KQ0Ye4_BSnPdjD5mODiTFOo2Kt7oZHzPBSJ5TEA9y9c", "1OznazxS69I6nwxOgYWi3R7-vTIngt-2Thymy8J_XEgg",
		"1ZgGtPd3kjD6N_Q-Cmks4p_wxhHVTKvge2hzxP_i0riA", "1l7ACAWjJUY9NALboiOxu6NenS22dp1F3bIH0CmkW2dM", "1V_u3mLrxNMGHLGrw-6DwbG4NBXEdYiBnUH_VaxnpT7M", "1Gs3l1o5cqN49p8-2t7iZaYQ0CV3zp_5P5C57DliGoEA",
		"1y_B0iBNYPjauPkbqU7M64mwR7KB6_nIvYHYYp51jn04", "1lZ4c-8xQc1Cz9snqFghloxx0qKNXq9nrLs9LucNIkGs", "1JLVHgocleGmsHX9y_Q0e675M5Wcw3zvkT6a0nHoxlp4", "1LstA-TDRHjJ9BL9_xmR1vBM4LUCWr3qj5PF_TfUn7ko",
		"16EEyi7Dv8qI6fgdWZUyTN70MafiyecWFNbmxVXJJBXs", "1n0oxcNlKey1v-0ymRdRfSn6Cx7IrfYbU25nRPUEUmtQ", "1nDWIR0splwvay39TAQklprOHY02XettoSMz3NE174m0", "12-DvqvJSohPsCB3iFDprXGcCRYqL6F7hU5oPBOB7hDs",
		"1XCeSvM2JAUXaVH4nTbMYPfPjAMb3Y_JCd4qH54M-r3k", "1SSe-zjdSITh3QNAMXiRoDth12feOY2qipLXB_DzFQxI", "1CoVi36PW6H3Xqu4XZV-ONvIOlP9Dhm4DY2UHWx-AGRk", "1wxppiBVNo2r2dRGfdN7oJYicWePbbHrhU-spQCf9U7I",
		"1q6QaHEziPecGDd7tDgvvCR4bg_J6VENI9ehGrUylx1U", "1dqHHiR40Ref4QkqhJxU1n1ttXeHzyXgEl6VSVIoUbbA", "1842Go0HUVGEV44J_qvnzbN1qiX707NrIqA0gjJx6yrg", "1oeBw7XAMlz_gR57E3kDvJYV3-xy6akgULdntKj5IyOg",
		"1nFl47XN3seedUOzFlo3hUCrDFHLNpcVph0H3nk-yaak", "18Y19f4d3WHYYAq1IMkhUZbS-C4v2MmKkJv5pabXhAWM", "1TggB6guUTWfe6PhVoxHmLiSgaybpEZM0XFUuo5psBkE", "1QQ7_7PUWaT5cR1Af0CGrNa0c_JPTzNE1y5EG8ciZ8YE",
		"17iJ0SPpGSVQKr-RigpJZa9EE7DWs72__c7tha2PsytA", "1WdAM3ZGOJJszy1C0XFEtyG7i3yOGuL1fC97p-cTvZAU", "1MD-UcWfw8Z-lhXBSkrBPkl7-6Q4GBW1KEL4gUDd6pOw", "16vwyTal7ikUlxq5HDlKvyOz1CStaNjh-DKGp79Al6NE",
		"16z3zl-ik8WYPbMgQrBAn7tz1tgtK48adKWxem4PmsoE", "1CDiLmBMh64_nUczpIGb0DRvEwN7dK0v-AQ1KXTm6J2w", "1sCGWsrBVgFbUOjm5JA4u425d5-yPuU2r1aaDGgr-f24", "18Meaq87Ck1Zm_4wOyVrBRWIZUdKaYe4Taj7Gcddpj34",
		"1ewEk5Vi0wFHTIHfsgL2rjljP5Qe2LK--V9x4vwOgtwI", "1vp8Z74I9y_OH0DKLu6FEY2_hs3TvFOqSvKJaKAiv87E", "1H136nbYJZge9GcCJ5m86g4JcGLou-SCfPTQvGeT_McA", "1DusAkoStlVZ-Zu4wGFAuzFJNIHN47ohOK4hVDwxMzbw",
		"1CIDQgAOWT4_aIcCwq4P1AOv7W9koAcupQ90H1o76pVo", "1ygHt7sghAG4GNJuvKuUEcueO1BrtfiUFN-QuyO79ad4", "1dgAaZ3HSq5-Kj4mdPEW-ovas7P4VOUYkXQIIacoTHoo", "1Y_lw9YmehaWLn-IDs9fXgpdLgbEGeHL3FOyeTxl86ew",
		"1KNwqVDR7qQ6sqpP_vJFN110th7QKsWPRA8B3JpZZ99Q", "18mxTBDTFxSdmBM7xir1S_hhPnFezRWB0zLYa9eh_OVg", "1bt7YLhw49THgy748zVLdigO-UKeOr5CyC91ZSvWvAVM", "1yS_4R86jW1cwaHUp23UkMLIjGrMb6MN42vNBbHl-3CE",
		"1NiF677nPlXmWl3u5ZJuxJjkArQbavvTrvZiOl6u2eHg", "1rh5yp4EB1Nfk-YOQpvjPazWn4Wnubn2eUa0vn_K91q4", "1ArP38u-L_VnaUSBdRxleGt7rlU8PSWLJhzJ9hRCkf8U", "1ujXrHDjrGFnD5-uFv3KD4sIgy_ZAlN5TkbkjdoaO3VY",
		"1yo0EejD1g6Np92ad8B2ZwRPjGiW0lSbIAwqW-ZPXwh0", "109sWUkfqZHt0NAW8tt4FQCIflEN-usIBgA4EZiIIUqo", "1_QyLKrEiLid9QupWle-Cwx_AhF8Uu5kz_xzYFosVPZg", "1ZGoTqYnL1n9C_8sRzhXVoccU6iUYPS-S3AWNuFPbJyg",
		"1Si1V6sqqPcDAsRIsBd-seH4AQHnI8lwPLFx4OXFi7K8", "1xy0b090uVuDizYgpXe3Z-LUIdTZ8AlVe9KkJmXoB3tE", "1R4E9xsIHiORRdUBD7qqvnWw2Y7NGRETaN4gUGryz4aE", "1jAVi-klxRf7tTCchV4271MzghdJltP6bb_EVL1fyN-I",
		"1_77bLWWtHUHKkiwYoIqJGoOW7rTpGpCUCMdhj1eCERI", "1gUQZue8ISmwGBbawh4auQ2zb6-O2smWRT-bSGy4oxrA", "1rWU5QFMW_cVmsSvDSRuWh3piWk7hINn-dR11Bv_0-pk", "1_BYswiuoWd2YLiJ3bYLCbIpxWlbYkQfhfZev4_ieU-M",
		"1ZsgpGvX_1iFAzzvCjX5dEBiZA_l1dZW-08j65DYTm0c", "1YPPJ3Qv4frkBYLZzI3oocF5R-YEfuly-bbQcBOWKXrE", "1qBfZ3_eU1bdIjxqyWgct-hvQvkLwB9g6EarLMT-RhCA", "1295iuFaxht8HfiaglgS7UscKAwkgpjwpjldr8tkKBbc",
		"1vphJ_OBLxMdKOkYtm78QnUwAMRCnsZwncb62guOFIe4", "19-d2mEOT-AvfJk4NTUyBngEd3uoaphsAOiiq3L3RYiA", "1opQb-_fLy3MpcRGCt82PY7gxT9-PtxFB5nds9Fh0yXs", "1FwoGFMlOwm4igs-soznJ3lvP9IxVWVZsVf4VrL7i88A",
		"1m6a_eT9h34RE_6aAN6_cWezqAwKhFqk0CkElJ-iX-eQ", "18ZmhgGWoUkrFsBXgG4NYgX-Z9D0_7Fa3R-8ADmgfZ_I", "1jlckU47lRAB471qmQTGNe4JMqOb83Jpw73enjz-SfO8", "1nVRK_uqKpVHYXjkozQRLuK8h3Ipdbh1SzXoIt92Vzic",
		"1opCfdDDlYTxZX0mVh-0UyKeKH7vF3jLeW98WdUA6t1Y", "1TBnkuHM3_qktELXVnzVptNii6QJvN5oCN9H4LbAHRkc", "1OUV2p1jWN9mkZhsVAVVjUR5tIAoffJajIV8MEajX8hs", "1xonCsWjavJATyjX5lWCvoczl4wxRq4apvo9i1sS8k3U",
		"1LYulgQlh8eKBvcXfplbL9grYyYtpWm4ql7dOldvJ_Zc", "1ppdeJBSHDjOnnGsxl1jZLyOBTj8StqqINMhOAgEzKh8", "1RQDGWYYGwCmIa8RNfCd9LmIkPN0aarDZTzy9cAsM6L8", "1OBjbY1lSe0QeXA4mWi8eBITCBjWx0DLh1lX3NvwMwDM",
		"1VYLhBd1xh77kdhxgilxGRE2g5XpszIVM1VP-Sab-7Yw", "1ruS7wHSBeSs6wSnIzjVoHl7ISzqaw3txq7ODvpzz38k", "1Zju_PpNRNy3qjKs-_brr28OfmFDcicynjCKHAPoVBxE", "1aaHEpGSP8ZtRlNiStwttDnyUv12FOh6y86KfLIOjIQY",
		"1PYApk2Yxw63i6y3R210XbZ9fXLKKkX1g5d0MxWyTW3o", "17GOo7ZIE7UgGcRlMYXQu1zZcPskyYqUTkpWctA_b6ZU", "1rb7Ys5InInp-TYcFIwmWZ-_C8kqOE5lALgJZst1p8UM", "1hzJHbloNkbKRHySu56k9fSEmpFFeIhHJlUdHRgCROoI",
		"1irHuX1ZKsYpsjjaDxuUfi_o4JVL4TAM35rufG-0UjGQ", "10yrR9wzDAUdRzUsy7wJ8hS4uJV3z-hGla5NRNBgl6oE", "1oSOSPKBOcIJv_bA8IY9wzx-Cp8FCRcJH-ckNNg8SKw8", "1yaGojYoxNFShLOTcrUzcduVlRcIBD1hh-V8LpMElUWg",
		"1xf96tk8XQ5gco6k8evxxrdbg_b9_EbmFtAVYc-LuhVg", "1KRj3v-2wBeO9IJYaqufF9TjzWQ1-Si8fGaWlnQP9BiQ", "1vIq5_YBiZF9p4TNCa7fi_3JGJ4116CrnWmZNB5uqEbw", "1TgXtWUNIcuQa_QYcThJpW-0vAxB82eOzbz1qGAJKymQ",
		"1TLIwBboooSG8qECGACqUnSa3J63VciE9ndCC-o2tbZ8", "1ACH_viO5Jw-as_5857ssmHp8VeJfTlb5O_8M-7Ai0A0", "1G19LO6V429E2X7O_VFYlqOSTOvZY1ilUq8atECIJHcw", "13Xs5PBwK7u_Zi7hrlfqEon4szBM6HTUMeNXEfZBGsSE",
		"17emaF0JmQxx-l6X4sEDvk2dL-lhTpZaoHvcGirhcizs", "1Khycip6LvADcClgGmsdpXi42wGUsLvVVL5Uf6VgTQdE", "1Xqfv2MNOdCh4QiaJQlhRnSGW6FKSgPP09qkRyzMFxZo", "1pQRb9QupZtdQSPYpSGJbKYZBOgnPe83nL5MWT4HqHxc",
		"1R7v3i3TYWKNHFPRiA3bryLalVn50QCzOi_f8mI67NA8", "1C5iYcDMHCfajIUevaAIAD6ILDUzd7ps9Ezju3j9fpsc", "1NPC3JesGwXi4vizmQIWdGulZDBBPUEnbPlzvxzWjWlU", "1fylgaRE1yF_DpXKq-t_2Z6KFAZ2apwyjc1czckmk8-g",
		"1FP_XEA-RPuoFx_wCtr3-XRLKfsN1ZRWPIGE9iY8TjEU", "1jsSmsI-EC9-FAYr89tTi6MJPVWqrQLyVdN6-cX6_SU4", "1Fa431lAssVU4aOy5yu3T6B2JjynWSLUugD1qE-sPoCE", "1ynG1XRmlNvZLRhMo48uJafHIBLDcCa_i-VeXWLrvQvY",
	}

	spreadsheetIDs = spreadsheetIDs[165:]

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).Do()
		if err != nil {
			log.Printf("Spreadsheets.Get: %v", err)
			continue
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
			return fmt.Errorf("batch.Do: %w", err)
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
	}

	return nil
}

func (s *SpreadsheetServiceGoogle) GetExpoCount(ctx context.Context) error {
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

	count := 18
	spreadsheetIDs = spreadsheetIDs[268:]
	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).Do()
		if err != nil {
			return err
		}

		for _, sheet := range spreadsheet.Sheets {
			if sheet.Properties.Title == "Затраты на участие в выставках" {
				count++
			}
		}

		fmt.Printf("%v/%v, count: %v\n", i+1, len(spreadsheetIDs), count)
		time.Sleep(time.Second)
	}

	return nil
}

func (s *SpreadsheetServiceGoogle) getSpreadsheets(ctx context.Context, httpClient *http.Client) ([]string, error) {
	driveSvc, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return spreadsheetIDs, nil
}

func (s *SpreadsheetServiceGoogle) AddTotalSumCells(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	spreadsheetIDs, err := s.getSpreadsheets(ctx, httpClient)
	if err != nil {
		return err
	}
	spreadsheetIDs = []string{"1oMJFttuiPxoBdejx3Ul3D2nscE0x45oeNfo-upVe1gE"}

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).Do()
		if err != nil {
			return err
		}

		batch := NewBatchUpdate(spreadsheetsSvc)
		for _, sheet := range spreadsheet.Sheets {
			sheetID := sheet.Properties.SheetId
			title := strings.ReplaceAll(sheet.Properties.Title, "⏳ (ожидайте) ", "")
			switch title {
			case "Затраты на доставку транспортом":
				totalSum_Затраты_на_доставку_транспортом(batch, sheetID)
			case "Затраты на сертификацию предприятия",
				"Затраты на рекламу ИКУ за рубежом",
				"Затраты на перевод каталога ИКУ",
				"Затраты на аренду помещения ИКУ",
				"Затраты на сертификацию ИКУ",
				"Затраты на демонстрацию ИКУ",
				"Затраты на франчайзинг",
				"Затраты на регистрацию товарных знаков",
				"Затраты на аренду",
				"Затраты на перевод",
				"Затраты на рекламу товаров за рубежом",
				"Затраты на участие в выставках",
				"Затраты на участие в выставках ИКУ",
				"Затраты на соответствие товаров требованиям":
				totalSum_Затраты_на_сертификацию_предприятия(batch, sheetID)
			}
		}

		if err := batch.Do(ctx, spreadsheetID); err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
		time.Sleep(time.Second)
	}

	return nil
}

func totalSum(sheetID int64, parentA1, targetA1 string) []*sheets.Request {
	var (
		result = make([]*sheets.Request, 0)

		a1Range = A1ToRange(parentA1)
		cell    = A1ToCell(targetA1)
	)

	beforeTargetRange := &Range{
		From: a1Range.From,
		To: &Cell{
			Col: cell.Col - 1,
			Row: cell.Row,
		},
	}

	afterTargetRange := &Range{
		From: &Cell{
			Col: cell.Col + 1,
			Row: cell.Row,
		},
		To: a1Range.To,
	}

	result = append(result,
		UnmergeRequest(sheetID, parentA1),
		SetCellFormula(sheetID, targetA1, fmt.Sprintf("=sum(%[1]s4:%[1]s)", numberToColumn(cell.Col))),
	)

	// fmt.Println("--------")

	// fmt.Println(parentA1, targetA1)
	// fmt.Printf("=sum(%[1]s4:%[1]s)\n", numberToColumn(cell.Col))

	if beforeTargetRange.From.Col < cell.Col && !beforeTargetRange.From.Equals(beforeTargetRange.To) {
		result = append(result, MergeRequest(sheetID, beforeTargetRange.ToA1()))
		// fmt.Printf("merge %s\n", beforeTargetRange.ToA1())
	}

	if afterTargetRange.To.Col > cell.Col && !afterTargetRange.From.Equals(afterTargetRange.To) {
		result = append(result, MergeRequest(sheetID, afterTargetRange.ToA1()))
		// fmt.Printf("merge %s\n", afterTargetRange.ToA1())
	}

	// fmt.Println("--------")

	return result
}

// Затраты на доставку транспортом
func totalSum_Затраты_на_доставку_транспортом(batch *BatchUpdate, sheetID int64) {
	args := []struct {
		parentA1 string
		targetA1 string
	}{
		{"U2:Y2", "X2"},
		{"Z2:AF2", "AF2"},
		{"AG2:AL2", "AJ2"},
		{"AM2:AP2", "AP2"},
		{"AQ2:AT2", "AT2"},
		{"AU2:AX2", "AX2"},
	}

	for _, arg := range args {
		batch.WithRequest(
			totalSum(sheetID, arg.parentA1, arg.targetA1)...,
		)
	}
}

// Затраты на сертификацию предприятия
func totalSum_Затраты_на_сертификацию_предприятия(batch *BatchUpdate, sheetID int64) {
	args := []struct {
		parentA1 string
		targetA1 string
	}{
		{"N2:V2", "U2"},
		{"W2:AB2", "AA2"},
		{"AC2:AI2", "AG2"},
		{"AJ2:AN2", "AM2"},
		{"AO2:AR2", "AR2"},
		{"AS2:AV2", "AV2"},
	}

	for _, arg := range args {
		batch.WithRequest(
			totalSum(sheetID, arg.parentA1, arg.targetA1)...,
		)
	}
}

func (s *SpreadsheetServiceGoogle) AddСоответствиеЖДН(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	// spreadsheetIDs, err := s.getSpreadsheets(ctx, httpClient)
	// if err != nil {
	// 	return err
	// }
	spreadsheetIDs := []string{"1oMJFttuiPxoBdejx3Ul3D2nscE0x45oeNfo-upVe1gE"}

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).Do()
		if err != nil {
			return err
		}

		batch := NewBatchUpdate(spreadsheetsSvc)

		for _, sheet := range spreadsheet.Sheets {
			sheetID := sheet.Properties.SheetId
			title := strings.ReplaceAll(sheet.Properties.Title, "⏳ (ожидайте) ", "")
			switch title {
			case "Затраты на доставку транспортом":
				СоответствиеЖДН(batch, sheetID)
			}
		}

		if err := batch.Do(ctx, spreadsheetID); err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
		time.Sleep(time.Second)
	}

	return nil
}

func СоответствиеЖДН(batch *BatchUpdate, sheetID int64) {
	batch.WithRequest(
		InsertColumnLeft(sheetID, "EF"),
		SetCellText(sheetID, "EF2", &SetCellTextInput{"Соответствие ЖДН (оцифровка)", true, 8}),
		MergeRequest(sheetID, "EF2:EF3"),
		SetDataValidationOneOf(sheetID, "EF4", []string{"да", "нет"}),
	)
}

func (s *SpreadsheetServiceGoogle) ConvertToFinancialSumCells(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	spreadsheetIDs, err := s.getSpreadsheets(ctx, httpClient)
	if err != nil {
		return err
	}

	for i := 626; i < len(spreadsheetIDs); i++ {
		spreadsheetID := spreadsheetIDs[i]
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Do()
		if err != nil {
			return err
		}

		batch := NewBatchUpdate(spreadsheetsSvc)
		for _, sheet := range spreadsheet.Sheets {
			sheetID := sheet.Properties.SheetId
			title := strings.ReplaceAll(sheet.Properties.Title, "⏳ (ожидайте) ", "")

			if _, ok := sums[title]; ok {
				convertToFinancialSumCells_Expenses(batch, sheetID, sheet.Data[0], title)
			}
		}

		fmt.Println(len(batch.requests))

		if err := batch.Do(ctx, spreadsheetID); err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
		time.Sleep(time.Second)
	}

	return nil
}

var nonAlphanumericRegex = regexp.MustCompile(`[^0-9,]`)

func getColumnData(sheet *sheets.GridData, columnIdx, rowIdx int) [][]interface{} {
	data := make([][]interface{}, len(sheet.RowData[rowIdx:]))
	for i, row := range sheet.RowData[rowIdx:] {
		if columnIdx >= len(row.Values) {
			continue
		}

		cell := row.Values[columnIdx]

		if cell == nil || cell.UserEnteredValue == nil {
			continue
		}

		value := ""

		if cell.UserEnteredValue.StringValue != nil {
			value = *cell.UserEnteredValue.StringValue
			value = nonAlphanumericRegex.ReplaceAllString(value, "")
		} else if cell.UserEnteredValue.NumberValue != nil {
			value = fmt.Sprintf("%f", *cell.UserEnteredValue.NumberValue)
			value = strings.TrimRight(value, "0")
			if value[len(value)-1] == '.' {
				value = value[:len(value)-1]
			}
			value = strings.ReplaceAll(value, ".", ",")
		} else {
			value = ""
		}

		value = strings.ReplaceAll(value, ",", ".")
		floatValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			data[i] = append(data[i], floatValue)
		} else {
			data[i] = append(data[i], value)
		}
	}
	return data
}

func convertToFinancialSumCells(sheetID int64, targetColumn string, sheet *sheets.GridData, title string) ([]*sheets.Request, []*sheets.ValueRange) {
	var (
		requests    = make([]*sheets.Request, 0)
		valueRanges = make([]*sheets.ValueRange, 0)

		column        = columnToNumber(targetColumn) - 1
		fromRow int64 = 3
	)

	requests = append(requests,
		&sheets.Request{
			FindReplace: &sheets.FindReplaceRequest{
				Find:          "[^0-9,]",
				Replacement:   "",
				SearchByRegex: true,
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartColumnIndex: column,
					EndColumnIndex:   column + 1,
					StartRowIndex:    fromRow,
				},
			},
		},
		&sheets.Request{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartColumnIndex: column,
					EndColumnIndex:   column + 1,
					StartRowIndex:    fromRow,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						NumberFormat: &sheets.NumberFormat{
							Type:    "NUMBER",
							Pattern: "#,##0.00",
						},
					},
				},
				Fields: "userEnteredFormat",
			},
		},
	)

	valueRanges = append(valueRanges, &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          fmt.Sprintf("'%s'!", title) + targetColumn + "4:" + fmt.Sprintf("%s%d", targetColumn, int64(len(sheet.RowData))),
		Values:         getColumnData(sheet, int(column), int(fromRow)),
	})

	return requests, valueRanges
}

var (
	baseCols1 = []string{"X", "AF", "AJ", "AP", "AT", "AX"}
	baseCols2 = []string{"U", "AA", "AG", "AM", "AR", "AV"}

	sums = map[string][]string{
		"Затраты на доставку транспортом":             append(baseCols1, "DU"),
		"Затраты на сертификацию предприятия":         append(baseCols2, "BG"),
		"Затраты на рекламу ИКУ за рубежом":           append(baseCols2, "BH"),
		"Затраты на перевод каталога ИКУ":             append(baseCols2, "BF"),
		"Затраты на аренду помещения ИКУ":             append(baseCols2, "BC"),
		"Затраты на сертификацию ИКУ":                 append(baseCols2, "BI"),
		"Затраты на демонстрацию ИКУ":                 append(baseCols2, "BC"),
		"Затраты на франчайзинг":                      append(baseCols2, "BC"),
		"Затраты на регистрацию товарных знаков":      append(baseCols2, "BO"),
		"Затраты на аренду":                           append(baseCols2, "BC"),
		"Затраты на перевод":                          append(baseCols2, "BF"),
		"Затраты на рекламу товаров за рубежом":       append(baseCols2, "BH"),
		"Затраты на участие в выставках":              append(baseCols2, "BV"),
		"Затраты на участие в выставках ИКУ":          append(baseCols2, "BV"),
		"Затраты на соответствие товаров требованиям": append(baseCols2, "BM"),
	}
)

func convertToFinancialSumCells_Expenses(batch *BatchUpdate, sheetID int64, sheet *sheets.GridData, title string) {
	args := sums[title]

	for _, arg := range args {
		requests, values := convertToFinancialSumCells(sheetID, arg, sheet, title)
		batch.WithRequest(requests...)
		batch.WithValueRange(values...)
	}
}
