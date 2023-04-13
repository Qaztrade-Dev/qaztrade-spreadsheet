package adapters

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

//go:embed credentials_sa.json
var credentialsSA []byte

func TestUpdateApplication(t *testing.T) {
	var (
		ctx           = context.Background()
		spreadsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
		appl          = &domain.Application{
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

	cli, err := NewSpreadsheetClient(ctx, credentialsSA)
	if err != nil {
		t.Fatal("NewSheetsClient error:", err)
	}

	err = cli.UpdateApplication(ctx, spreadsheetID, appl)
	if err != nil {
		t.Fatal("UpdateApplication error:", err)
	}
}
