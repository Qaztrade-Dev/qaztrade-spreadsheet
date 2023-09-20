package adapters

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	qaztradeSheets "github.com/doodocs/qaztrade/backend/internal/sheets/adapters"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func (s *SpreadsheetServiceGoogle) BackupSpreadsheets(ctx context.Context) error {
	var (
		accessKey  = os.Getenv("S3_ACCESS_KEY")
		secretKey  = os.Getenv("S3_SECRET_KEY")
		endpoint   = os.Getenv("S3_ENDPOINT")
		bucketName = "spreadsheets"
	)

	storage, err := qaztradeSheets.NewStorageS3(ctx, accessKey, secretKey, bucketName, endpoint)
	if err != nil {
		return err
	}

	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	driveSvc, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	spreadsheetIDs, err := s.getSpreadsheets(ctx, httpClient)
	if err != nil {
		return err
	}

	for i := 0; i < len(spreadsheetIDs); i++ {
		spreadsheetID := spreadsheetIDs[i]
		fmt.Println(spreadsheetID)

		httpResp, err := driveSvc.Files.Export(
			spreadsheetID, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		).Download()
		if err != nil {
			return err
		}

		defer httpResp.Body.Close()

		b, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		_, err = storage.Upload(ctx, fmt.Sprintf("%s.xlsx", spreadsheetID), int64(len(b)), bytes.NewReader(b))
		if err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
	}

	return nil
}
