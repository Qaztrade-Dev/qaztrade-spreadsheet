package adapters

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed credentials_sa.json
var credentialsSA []byte

func TestGetApplication(t *testing.T) {
	var (
		ctx           = context.Background()
		adminAccount  = os.Getenv("ADMIN_ACCOUNT")
		svcAccount    = os.Getenv("SERVICE_ACCOUNT")
		spreadsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
	)

	cli, err := NewSpreadsheetClient(ctx, credentialsSA, adminAccount, svcAccount)
	require.Nil(t, err)

	application, err := cli.GetApplication(ctx, spreadsheetID)
	require.Nil(t, err)

	fmt.Printf("%#v\n", application)
}

func TestGetAttachments(t *testing.T) {
	var (
		ctx           = context.Background()
		adminAccount  = os.Getenv("ADMIN_ACCOUNT")
		svcAccount    = os.Getenv("SERVICE_ACCOUNT")
		spreadsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
	)

	cli, err := NewSpreadsheetClient(ctx, credentialsSA, adminAccount, svcAccount)
	require.Nil(t, err)

	expensesTitles, err := cli.GetExpensesSheetTitles(ctx, spreadsheetID)
	require.Nil(t, err)

	fmt.Println(expensesTitles)
	if len(expensesTitles) == 0 {
		return
	}

	attachments, err := cli.GetAttachments(ctx, spreadsheetID, expensesTitles)
	if err != nil {
		t.Fatal("GetAttachments error:", err)
	}

	for i, attachment := range attachments {
		name := fmt.Sprintf("%v.pdf", i)
		body, _ := io.ReadAll(attachment)
		os.WriteFile(name, body, 0644)
	}
}

func TestGetExpensesData(t *testing.T) {
	var (
		ctx           = context.Background()
		adminAccount  = os.Getenv("ADMIN_ACCOUNT")
		svcAccount    = os.Getenv("SERVICE_ACCOUNT")
		spreadsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
	)

	cli, err := NewSpreadsheetClient(ctx, credentialsSA, adminAccount, svcAccount)
	require.Nil(t, err)

	expensesTitles, err := cli.GetExpensesSheetTitles(ctx, spreadsheetID)
	require.Nil(t, err)
	fmt.Println(expensesTitles)

	expensesValues, err := cli.GetExpenseValues(ctx, spreadsheetID, expensesTitles)
	require.Nil(t, err)
	fmt.Println(expensesValues)
}

func TestHasMergedCells(t *testing.T) {
	var (
		ctx           = context.Background()
		adminAccount  = os.Getenv("ADMIN_ACCOUNT")
		svcAccount    = os.Getenv("SERVICE_ACCOUNT")
		spreadsheetID = "15wAKoZVRz1FbayCTA9SjYvIs3v_vbTZB2_mgIlhJL0g"
	)

	cli, err := NewSpreadsheetClient(ctx, credentialsSA, adminAccount, svcAccount)
	require.Nil(t, err)

	expensesTitles, err := cli.GetExpensesSheetTitles(ctx, spreadsheetID)
	require.Nil(t, err)

	hasMergedCells, err := cli.HasMergedCells(ctx, spreadsheetID, expensesTitles)
	require.Nil(t, err)
	require.False(t, hasMergedCells)
}

func TestBlockImportantRanges(t *testing.T) {
	var (
		ctx           = context.Background()
		adminAccount  = os.Getenv("ADMIN_ACCOUNT")
		svcAccount    = os.Getenv("SERVICE_ACCOUNT")
		spreadsheetID = "13mbVQopQziZO4lEok42y-ThyVdX8Fk8AHF5qbrVR0nw"
	)

	cli, err := NewSpreadsheetClient(ctx, credentialsSA, adminAccount, svcAccount)
	require.Nil(t, err)

	err = cli.BlockImportantRanges(ctx, spreadsheetID)
	require.Nil(t, err)
}
