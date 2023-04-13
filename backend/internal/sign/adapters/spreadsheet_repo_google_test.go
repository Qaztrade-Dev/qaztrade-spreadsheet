package adapters

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"testing"
)

//go:embed credentials_sa.json
var credentialsSA []byte

func TestGetApplication(t *testing.T) {
	var (
		ctx           = context.Background()
		spreadsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
	)

	cli, err := NewSpreadsheetClient(ctx, credentialsSA)
	if err != nil {
		t.Fatal("NewSpreadsheetClient error:", err)
	}

	application, err := cli.GetApplication(ctx, spreadsheetID)
	if err != nil {
		t.Fatal("GetApplication error:", err)
	}

	fmt.Printf("%#v\n", application)
}

func TestGetAttachments(t *testing.T) {
	var (
		ctx           = context.Background()
		spreadsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
	)

	cli, err := NewSpreadsheetClient(ctx, credentialsSA)
	if err != nil {
		t.Fatal("NewSpreadsheetClient error:", err)
	}

	attachments, err := cli.GetAttachments(ctx, spreadsheetID)
	if err != nil {
		t.Fatal("GetAttachments error:", err)
	}

	for i, attachment := range attachments {
		name := fmt.Sprintf("%v.pdf", i)
		body, _ := io.ReadAll(attachment)
		os.WriteFile(name, body, 0644)
	}
}
