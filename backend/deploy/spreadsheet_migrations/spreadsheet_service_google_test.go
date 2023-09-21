package adapters

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/doodocs/qaztrade/backend/pkg/qaztradeoauth2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

//go:embed credentials_oauth.json
var credentialsOAuth []byte

func TestTest(t *testing.T) {
	var (
		ctx              = context.Background()
		credentialsOAuth = credentialsOAuth

		postgresLogin         = getenv("POSTGRES_LOGIN1", "postgres")
		postgresPassword      = getenv("POSTGRES_PASSWORD1", "postgres")
		postgresHost          = getenv("POSTGRES_HOST1", "localhost")
		postgresDatabase      = getenv("POSTGRES_DATABASE1", "qaztrade")
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
		originSpreadsheetID,
		templateSpreadsheetId,
		destinationFolderId,
	)

	err = svc.FixVLOOKUP(ctx)
	if err != nil {
		t.Fatal("Test error:", err)
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
