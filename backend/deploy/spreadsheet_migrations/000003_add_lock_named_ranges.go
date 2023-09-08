package adapters

import (
	"context"
	"fmt"
	"strings"

	spreadsheetsAdapters "github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

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

		batch := spreadsheetsAdapters.NewBatchUpdate(spreadsheetsSvc)

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
