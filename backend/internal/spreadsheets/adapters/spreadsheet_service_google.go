package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetServiceGoogle struct {
	pg         *pgxpool.Pool
	config     *oauth2.Config
	svcAccount string // sheets@secret-beacon-380907.iam.gserviceaccount.com
	jwtcli     *jwt.Client
}

var _ domain.SpreadsheetService = (*SpreadsheetServiceGoogle)(nil)

func NewSpreadsheetServiceGoogle(clientSecretBytes []byte, svcAccount string, jwtcli *jwt.Client, pg *pgxpool.Pool) (*SpreadsheetServiceGoogle, error) {
	config, err := google.ConfigFromJSON(clientSecretBytes, drive.DriveScope, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	return &SpreadsheetServiceGoogle{
		pg:         pg,
		config:     config,
		svcAccount: svcAccount,
		jwtcli:     jwtcli,
	}, err
}

func (s *SpreadsheetServiceGoogle) Create(ctx context.Context, user *domain.User) (string, error) {
	httpClient, err := s.getOauth2Client(ctx)
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

	return spreadsheetID, nil
}

func (s *SpreadsheetServiceGoogle) copyFile(ctx context.Context, svc *drive.Service, user *domain.User) (string, error) {
	var (
		// TODO
		// move to struct
		sourceSpreadsheetID = "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"
		destinationFolderID = "1c04RznMaAumXl9OfVkstH4ZAIG3ULOgR"
	)

	newFileName, err := domain.CreateSpreadsheetName(user)
	if err != nil {
		return "", err
	}

	copy := &drive.File{
		Title:   newFileName,
		Parents: []*drive.ParentReference{{Id: destinationFolderID}},
	}
	copiedFile, err := svc.Files.Copy(sourceSpreadsheetID, copy).Context(ctx).Do()
	if err != nil {
		return "", err
	}

	return copiedFile.Id, nil
}

func (s *SpreadsheetServiceGoogle) initSpreadsheet(ctx context.Context, svc *sheets.Service, spreadsheetID string) error {
	spreadsheet, err := svc.Spreadsheets.Get(spreadsheetID).IncludeGridData(false).Context(ctx).Do()
	if err != nil {
		return err
	}

	batch := NewBatchUpdate(svc)
	for _, sheet := range spreadsheet.Sheets {
		for _, protectedRange := range sheet.ProtectedRanges {
			s.setProtectedRange(spreadsheetID, protectedRange, batch)
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

func (s *SpreadsheetServiceGoogle) setProtectedRange(
	spreadsheetID string,
	protectedRange *sheets.ProtectedRange,
	batch *BatchUpdate,
) {
	editors := protectedRange.Editors
	if stringSliceIndex(editors.Users, s.svcAccount) >= 0 {
		return
	}

	editors.Users = append(editors.Users, s.svcAccount)

	batch.WithRequest(&sheets.Request{
		UpdateProtectedRange: &sheets.UpdateProtectedRangeRequest{
			ProtectedRange: &sheets.ProtectedRange{
				ProtectedRangeId: protectedRange.ProtectedRangeId,
				Editors:          editors,
			},
			Fields: "editors",
		},
	})
}

func (s *SpreadsheetServiceGoogle) setMetadata(spreadsheetID string, batch *BatchUpdate) error {
	tokenStr, err := jwt.NewTokenString(s.jwtcli, &domain.SpreadsheetClaims{
		SpreadsheetID: spreadsheetID,
	})
	if err != nil {
		return err
	}

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
	_, err := svc.Permissions.Insert(spreadsheetID, permission).Do()
	if err != nil {
		return err
	}
	return nil
}

func stringSliceIndex(arr []string, str string) int {
	for i, v := range arr {
		if v == str {
			return i
		}
	}
	return -1
}

func (s *SpreadsheetServiceGoogle) getOauth2Client(ctx context.Context) (*http.Client, error) {
	tokenStr, err := s.getOauthToken(ctx)
	if err != nil {
		return nil, err
	}

	tokenOauth2, err := tokenFromStr(tokenStr)
	if err != nil {
		return nil, err
	}

	httpCli := clientWithToken(s.config, tokenOauth2)
	return httpCli, nil

}

func (s *SpreadsheetServiceGoogle) getOauthToken(ctx context.Context) (string, error) {
	const sql = `
		select 
			token
		from "oauth2_tokens"
		where id = 1
	`

	var (
		token string
	)

	err := s.pg.QueryRow(ctx, sql).Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func tokenFromStr(tokenStr string) (*oauth2.Token, error) {
	tok := &oauth2.Token{}
	err := json.Unmarshal([]byte(tokenStr), &tok)
	return tok, err
}

func clientWithToken(config *oauth2.Config, token *oauth2.Token) *http.Client {
	return config.Client(context.Background(), token)
}

func (s *SpreadsheetServiceGoogle) GetPublicLink(_ context.Context, spreadsheetID string) string {
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit?usp=sharing", spreadsheetID)
	return url
}
