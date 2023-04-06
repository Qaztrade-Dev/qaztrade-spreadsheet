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
		ctx                       = context.Background()
		originSpreadsheetID       = "1YvRrTIVWz1kigSke6pN8Uz87r0fWl-kyarogwAjKx5c"
		spreadsheetID             = "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"
		sheetName                 = "Доставка ЖД транспортом"
		sheetID             int64 = 932754288
	)

	cli, err := NewSpreadsheetClient(ctx, credentials, originSpreadsheetID)
	if err != nil {
		t.Fatal("NewSheetsClient error:", err)
	}

	err = cli.InsertRecord(ctx, spreadsheetID, sheetName, sheetID, &domain.Payload{
		ParentID: "null",
		ChildKey: "№",
		Value: domain.PayloadValue{
			"№": "2",
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
		ctx                       = context.Background()
		originSpreadsheetID       = "1YvRrTIVWz1kigSke6pN8Uz87r0fWl-kyarogwAjKx5c"
		spreadsheetID             = "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"
		sheetName                 = "Доставка ЖД транспортом"
		sheetID             int64 = 932754288
	)

	cli, err := NewSpreadsheetClient(ctx, credentials, originSpreadsheetID)
	if err != nil {
		t.Fatal("NewSheetsClient error:", err)
	}

	err = cli.InsertRecord(ctx, spreadsheetID, sheetName, sheetID, &domain.Payload{
		ParentID: "1",
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

func TestUpdateApplication(t *testing.T) {
	var (
		ctx                 = context.Background()
		originSpreadsheetID = "1YvRrTIVWz1kigSke6pN8Uz87r0fWl-kyarogwAjKx5c"
		spreadsheetID       = "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"
		appl                = &domain.Application{
			From:                  "Kaspi Ltd.",
			GovReg:                "11111111",
			FactAddr:              "Алматы",
			Bin:                   "950223347566",
			Industry:              "Информационные технологии",
			Activity:              "ЭДО",
			EmpCount:              "1-5",
			Manufacturer:          "Doodocs",
			Item:                  "Подписка",
			ItemVolume:            "200",
			FactVolumeEarnings:    "35000000",
			FactWorkload:          "100",
			ChiefLastname:         "Давлетов",
			ChiefFirstname:        "Дагар",
			ChiefMiddlename:       "Гусманович",
			ChiefPosition:         "Директор",
			ChiefPhone:            "+77777777774",
			ContLastname:          "Тлекбаи",
			ContFirstname:         "Али",
			ContMiddlename:        "Кайратулы",
			ContPosition:          "Разработчик",
			ContPhone:             "+77777777777",
			ContEmail:             "a@gmail.com",
			InfoManufacturedGoods: "Казахстан",
			NameOfGoods:           "123123",
		}
	)

	cli, err := NewSpreadsheetClient(ctx, credentials, originSpreadsheetID)
	if err != nil {
		t.Fatal("NewSheetsClient error:", err)
	}

	err = cli.UpdateApplication(ctx, spreadsheetID, appl)
	if err != nil {
		t.Fatal("UpdateApplication error:", err)
	}
}

func TestAddSheet(t *testing.T) {
	var (
		ctx                 = context.Background()
		originSpreadsheetID = "1BY6-dstDDWP1k6Xv-HzmZ6q4SJ3i088z26gBgkfwXow"
		spreadsheetID       = "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"
		sheetName           = "Затраты на продвижение"
	)

	cli, err := NewSpreadsheetClient(ctx, credentials, originSpreadsheetID)
	if err != nil {
		t.Fatal("NewSheetsClient error:", err)
	}

	err = cli.AddSheet(ctx, spreadsheetID, sheetName)
	if err != nil {
		t.Fatal("AddSheet error:", err)
	}
}

func TestRemoveParent(t *testing.T) {
	var (
		ctx                 = context.Background()
		originSpreadsheetID = "1YvRrTIVWz1kigSke6pN8Uz87r0fWl-kyarogwAjKx5c"
		spreadsheetID       = "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"
		sheetName           = "Доставка ЖД транспортом"
		sheetID             = int64(1974041431)
	)

	cli, err := NewSpreadsheetClient(ctx, credentials, originSpreadsheetID)
	if err != nil {
		t.Fatal("NewSheetsClient error:", err)
	}

	err = cli.RemoveRecord(ctx, spreadsheetID, sheetName, sheetID, &domain.RemoveInput{
		Value: "21",
		Name:  "Дистрибьюторский договор",
	})
	if err != nil {
		t.Fatal("RemoveRecord error ", err)
	}
}
