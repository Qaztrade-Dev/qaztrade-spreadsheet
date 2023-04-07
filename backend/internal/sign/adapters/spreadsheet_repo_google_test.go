package adapters

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"testing"
)

//go:embed credentials.json
var credentials []byte

func TestGetApplication(t *testing.T) {
	var (
		ctx           = context.Background()
		spreadsheetID = "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"
	)

	cli, err := NewSpreadsheetClient(ctx, credentials)
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
		spreadsheetID = "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"
	)

	cli, err := NewSpreadsheetClient(ctx, credentials)
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
