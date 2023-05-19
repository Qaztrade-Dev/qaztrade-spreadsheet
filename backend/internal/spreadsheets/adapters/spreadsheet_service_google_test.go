package adapters

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/doodocs/qaztrade/backend/pkg/qaztradeoauth2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

//go:embed credentials_oauth.json
var credentialsOAuth []byte

func TestSpreadsheetCreate(t *testing.T) {
	var (
		ctx              = context.Background()
		credentialsOAuth = credentialsOAuth
		jwtcli           = jwt.NewClient("qaztradesecret")

		postgresLogin         = getenv("POSTGRES_LOGIN", "postgres")
		postgresPassword      = getenv("POSTGRES_PASSWORD", "postgres")
		postgresHost          = getenv("POSTGRES_HOST", "localhost")
		postgresDatabase      = getenv("POSTGRES_DATABASE", "qaztrade")
		originSpreadsheetID   = getenv("ORIGIN_SPREADSHEET_ID")
		templateSpreadsheetId = getenv("TEMPLATE_SPREADSHEET_ID")
		destinationFolderId   = getenv("DESTINATION_FOLDER_ID")
		reviewerAccount       = getenv("REVIEWER_ACCOUNT")
		svcAccount            = getenv("SERVICE_ACCOUNT")

		postgresURL = fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", postgresLogin, postgresPassword, postgresHost, postgresDatabase)

		user = &domain.User{ID: "75455b90-edad-4281-9509-611c7cc24df8", OrgName: "Doodocs1"}
	)

	pg, err := pgxpool.Connect(ctx, postgresURL)
	require.Nil(t, err)

	oauth2, err := qaztradeoauth2.NewClient(credentialsOAuth, pg)
	require.Nil(t, err)

	svc := NewSpreadsheetServiceGoogle(
		oauth2,
		svcAccount,
		reviewerAccount,
		jwtcli,
		originSpreadsheetID,
		templateSpreadsheetId,
		destinationFolderId,
	)

	id, err := svc.Create(ctx, user)
	fmt.Println(id)
	fmt.Println(err)
}

func TestAddSheet(t *testing.T) {
	var (
		ctx              = context.Background()
		credentialsOAuth = credentialsOAuth
		jwtcli           = jwt.NewClient("qaztradesecret")

		postgresLogin         = getenv("POSTGRES_LOGIN", "postgres")
		postgresPassword      = getenv("POSTGRES_PASSWORD", "postgres")
		postgresHost          = getenv("POSTGRES_HOST", "localhost")
		postgresDatabase      = getenv("POSTGRES_DATABASE", "qaztrade")
		templateSpreadsheetId = getenv("TEMPLATE_SPREADSHEET_ID")
		destinationFolderId   = getenv("DESTINATION_FOLDER_ID")
		reviewerAccount       = getenv("REVIEWER_ACCOUNT")
		svcAccount            = getenv("SERVICE_ACCOUNT")

		postgresURL = fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", postgresLogin, postgresPassword, postgresHost, postgresDatabase)

		originSpreadsheetID = os.Getenv("ORIGIN_SPREADSHEET_ID")
		// spreadsheetID       = os.Getenv("TEMPLATE_SPREADSHEET_ID")
		spreadsheetID = "1ZGkQGlUFc6Fgj2yYCitE0KUCPwxbsxnwEhTrBxHeImY"
		sheetName     = "Затраты на доставку транспортом"
	)

	pg, err := pgxpool.Connect(ctx, postgresURL)
	require.Nil(t, err)

	oauth2, err := qaztradeoauth2.NewClient(credentialsOAuth, pg)
	require.Nil(t, err)

	svc := NewSpreadsheetServiceGoogle(
		oauth2,
		svcAccount,
		reviewerAccount,
		jwtcli,
		originSpreadsheetID,
		templateSpreadsheetId,
		destinationFolderId,
	)

	err = svc.AddSheet(ctx, spreadsheetID, sheetName)
	if err != nil {
		t.Fatal("AddSheet error:", err)
	}
}

func TestTest(t *testing.T) {
	var (
		ctx              = context.Background()
		credentialsOAuth = credentialsOAuth
		jwtcli           = jwt.NewClient("qaztradesecret")

		postgresLogin         = getenv("POSTGRES_LOGIN", "postgres")
		postgresPassword      = getenv("POSTGRES_PASSWORD", "postgres")
		postgresHost          = getenv("POSTGRES_HOST", "localhost")
		postgresDatabase      = getenv("POSTGRES_DATABASE", "qaztrade")
		templateSpreadsheetId = getenv("TEMPLATE_SPREADSHEET_ID")
		destinationFolderId   = getenv("DESTINATION_FOLDER_ID")
		reviewerAccount       = getenv("REVIEWER_ACCOUNT")
		svcAccount            = getenv("SERVICE_ACCOUNT")

		postgresURL = fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", postgresLogin, postgresPassword, postgresHost, postgresDatabase)

		originSpreadsheetID = os.Getenv("ORIGIN_SPREADSHEET_ID")
	)

	pg, err := pgxpool.Connect(ctx, postgresURL)
	require.Nil(t, err)

	oauth2, err := qaztradeoauth2.NewClient(credentialsOAuth, pg)
	require.Nil(t, err)

	svc := NewSpreadsheetServiceGoogle(
		oauth2,
		svcAccount,
		reviewerAccount,
		jwtcli,
		originSpreadsheetID,
		templateSpreadsheetId,
		destinationFolderId,
	)

	err = svc.Test(ctx)
	if err != nil {
		t.Fatal("AddSheet error:", err)
	}
}

func getenv(env string, fallback ...string) string {
	e := os.Getenv(env)
	if e == "" {
		value := ""
		if len(fallback) > 0 {
			value = fallback[0]
		}
		return value
	}
	return e
}
