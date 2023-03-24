package adapters

import (
	"context"
	_ "embed"
	"testing"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

//go:embed credentials.json
var credentials []byte

func TestSubmit(t *testing.T) {
	var (
		ctx           = context.Background()
		spreadsheetID = "1KL-lrhs-Wu9kRAppBxAHUUFr7OCfNYla8Z7W-0tX4Mo"
	)

	cli, err := NewSheetsClient(ctx, credentials)
	if err != nil {
		t.Fatal("NewSheetsClient error:", err)
	}

	err = cli.InsertRecord(ctx, spreadsheetID, &domain.Payload{
		ParentID: "null",
		ChildKey: "№",
		Value: domain.PayloadValue{
			"№": "3",
			"Производитель/дочерняя компания/дистрибьютор/СПК": "Doodocs",
			"подтверждающий документ": map[string]interface{}{
				"производитель":       "производитель",
				"наименование":        "наименование",
				"№":                   "№",
				"наименование товара": "наименование товара",
				"ТН ВЭД (6 знаков)":   "ТН ВЭД (6 знаков)",
				"дата":                "дата",
				"срок":                "срок",
				"подтверждение на сайте уполномоченного органа": "подтверждение на сайте уполномоченного органа",
			},
		},
	})
	if err != nil {
		t.Fatal("InsertRecord error:", err)
	}
}

func TestSubmitChild(t *testing.T) {
	var (
		ctx           = context.Background()
		spreadsheetID = "1KL-lrhs-Wu9kRAppBxAHUUFr7OCfNYla8Z7W-0tX4Mo"
	)

	cli, err := NewSheetsClient(ctx, credentials)
	if err != nil {
		t.Fatal("NewSheetsClient error:", err)
	}

	err = cli.InsertRecord(ctx, spreadsheetID, &domain.Payload{
		ParentID: "2",
		ChildKey: "Дистрибьюторский договор",
		Value: domain.PayloadValue{
			"Дистрибьюторский договор": map[string]interface{}{
				"№":       "3",
				"дата":    "дата",
				"условия": "условия",
			},
		},
	})
	if err != nil {
		t.Fatal("InsertRecord error:", err)
	}
}
