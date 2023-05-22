package adapters

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed credentials_sa.json
var credentialsSA []byte

func TestComments(t *testing.T) {
	var (
		ctx          = context.Background()
		adminAccount = getenv("ADMIN_ACCOUNT")
		// spreadsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
		spreadsheetID = "15Y8kld4d3PmFdNEXjjLHgFoTC4TtDeJjwPCAe1aLXD4"
	)

	svc, err := NewSpreadsheetService(ctx, credentialsSA, adminAccount)
	require.Nil(t, err)

	err = svc.Comments(ctx, spreadsheetID)
	require.Nil(t, err)
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
