package adapters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/doodocs/qaztrade/backend/pkg/qaztradeoauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type SpreadsheetServiceGoogle struct {
	oauth2                *qaztradeoauth2.Client
	svcAccount            string
	reviewerAccount       string
	originSpreadsheetID   string
	templateSpreadsheetID string
	destinationFolderID   string
}

func NewSpreadsheetServiceGoogle(
	oauth2 *qaztradeoauth2.Client,
	svcAccount string,
	reviewerAccount string,
	originSpreadsheetID string,
	templateSpreadsheetID string,
	destinationFolderID string,
) *SpreadsheetServiceGoogle {
	client := &SpreadsheetServiceGoogle{
		oauth2:                oauth2,
		svcAccount:            svcAccount,
		reviewerAccount:       reviewerAccount,
		originSpreadsheetID:   originSpreadsheetID,
		templateSpreadsheetID: templateSpreadsheetID,
		destinationFolderID:   destinationFolderID,
	}

	return client
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
