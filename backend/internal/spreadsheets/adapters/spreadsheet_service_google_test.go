package adapters

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

//go:embed client_secret.json
var clientSecretBytes []byte

func TestSpreadsheetCreate(t *testing.T) {
	var (
		ctx               = context.Background()
		clientSecretBytes = clientSecretBytes
		svcAccount        = "sheets@secret-beacon-380907.iam.gserviceaccount.com"
		jwtcli            = jwt.NewClient("qaztradesecret")

		postgresLogin         = getenv("POSTGRES_LOGIN", "postgres")
		postgresPassword      = getenv("POSTGRES_PASSWORD", "postgres")
		postgresHost          = getenv("POSTGRES_HOST", "localhost")
		postgresDatabase      = getenv("POSTGRES_DATABASE", "qaztrade")
		templateSpreadsheetId = getenv("TEMPLATE_SPREADSHEET_ID")
		destinationFolderId   = getenv("DESTINATION_FOLDER_ID")

		postgresURL = fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", postgresLogin, postgresPassword, postgresHost, postgresDatabase)

		user = &domain.User{ID: "75455b90-edad-4281-9509-611c7cc24df8", OrgName: "Doodocs1"}
	)

	pg, err := pgxpool.Connect(ctx, postgresURL)
	require.Nil(t, err)

	svc, err := NewSpreadsheetServiceGoogle(clientSecretBytes, svcAccount, jwtcli, pg, templateSpreadsheetId, destinationFolderId)
	require.Nil(t, err)

	id, err := svc.Create(ctx, user)
	fmt.Println(id)
	fmt.Println(err)
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
