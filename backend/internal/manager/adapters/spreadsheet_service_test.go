package adapters

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"github.com/stretchr/testify/require"
)

//go:embed credentials_sa.json
var credentialsSA []byte

func TestComments(t *testing.T) {
	var (
		ctx          = context.Background()
		adminAccount = getenv("ADMIN_ACCOUNT")
		svcAccount   = getenv("SERVICE_ACCOUNT")
		// spreadsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
		spreadsheetID = "15Y8kld4d3PmFdNEXjjLHgFoTC4TtDeJjwPCAe1aLXD4"
	)

	svc, err := NewSpreadsheetService(ctx, credentialsSA, adminAccount, svcAccount)
	require.Nil(t, err)

	revision, err := svc.Comments(ctx, &domain.Application{
		ID:            "1",
		SpreadsheetID: spreadsheetID,
		No:            1,
		Attrs: map[string]interface{}{
			"application": map[string]interface{}{
				"bin":          "123456",
				"manufacturer": "manufacturer",
				"from":         "from",
				"cont_email":   "cont_email",
			},
		},
	})
	require.Nil(t, err)
	require.NotNil(t, revision)
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

func TestGrantAdminPermissions(t *testing.T) {
	var (
		ctx           = context.Background()
		adminAccount  = getenv("ADMIN_ACCOUNT")
		svcAccount    = getenv("SERVICE_ACCOUNT")
		spreadsheetID = "1IBbhwlDLYkD0vUT2PHn6d92xsszAJXD95yfKO8ieltQ"
		email         = "ali.tlekbai@gmail.com"
	)

	svc, err := NewSpreadsheetService(ctx, credentialsSA, adminAccount, svcAccount)
	require.Nil(t, err)

	err = svc.GrantAdminPermissions(ctx, spreadsheetID, email)
	require.Nil(t, err)
}
