package adapters

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed credentials_sa.json
var credentialsSA []byte

func TestComments(t *testing.T) {
	var (
		ctx = context.Background()
		// spreadsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
		spreadsheetID = "15Y8kld4d3PmFdNEXjjLHgFoTC4TtDeJjwPCAe1aLXD4"
	)

	svc, err := NewSpreadsheetService(ctx, credentialsSA)
	require.Nil(t, err)

	err = svc.Comments(ctx, spreadsheetID)
	require.Nil(t, err)
}
