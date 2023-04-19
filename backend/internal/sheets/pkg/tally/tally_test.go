package tally

import (
	"testing"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	var (
		jsonBytes = []byte(`
		{
			"eventId": "a5376178-ef3d-482e-aa4b-f4bd9d5a101a",
			"eventType": "FORM_RESPONSE",
			"createdAt": "2023-03-28T06:49:28.934Z",
			"data": {
			  "responseId": "zNVzZR",
			  "submissionId": "zNVzZR",
			  "respondentId": "YP0arv",
			  "formId": "m6LxZB",
			  "formName": "Возмещение затрат по экспорту",
			  "createdAt": "2023-03-28T06:49:28.000Z",
			  "fields": [
				{
				  "key": "question_woLbzO",
				  "label": "От кого",
				  "type": "INPUT_TEXT",
				  "value": "Doodocs Ltd."
				},
				{
				  "key": "question_nGAxbL",
				  "label": "Государственная регистрация/перерегистрация",
				  "type": "INPUT_TEXT",
				  "value": "120120102"
				},
				{
				  "key": "question_mOMDEY",
				  "label": "Фактический адрес",
				  "type": "INPUT_TEXT",
				  "value": "Астана"
				},
				{
				  "key": "question_mVOLNM",
				  "label": "БИН/ИИН",
				  "type": "INPUT_TEXT",
				  "value": "950223347566"
				},
				{
				  "key": "question_nPJdrB",
				  "label": "Наименование отрасли",
				  "type": "INPUT_TEXT",
				  "value": "Информационные технологии"
				},
				{
				  "key": "question_3E7k6B",
				  "label": "Вид деятельности",
				  "type": "INPUT_TEXT",
				  "value": "ЭДО"
				},
				{
				  "key": "question_nrMjY2",
				  "label": "Численность сотрудников",
				  "type": "INPUT_TEXT",
				  "value": "1-5"
				},
				{
					"key": "question_nrMjY2",
					"label": "Производственная мощность, возможности увеличения",
					"type": "INPUT_TEXT",
					"value": "900"
				  },
				{
				  "key": "question_w4V4OA",
				  "label": "Производитель",
				  "type": "INPUT_TEXT",
				  "value": "Doodocs"
				},
				{
				  "key": "question_3jAazJ",
				  "label": "Товар",
				  "type": "INPUT_TEXT",
				  "value": "Подписка"
				},
				{
				  "key": "question_w2GEoD",
				  "label": "Объем товара",
				  "type": "INPUT_TEXT",
				  "value": "200"
				},
				{
				  "key": "question_meWqzl",
				  "label": "Объем фактической валютной выручки за полугодие предшествующей дате подачи заявки",
				  "type": "INPUT_TEXT",
				  "value": "35000000"
				},
				{
				  "key": "question_nWVOqk",
				  "label": "Фактическая загруженность производства",
				  "type": "INPUT_TEXT",
				  "value": "100"
				},
				{
				  "key": "question_waqQzX",
				  "label": "Фамилия руководителя",
				  "type": "INPUT_TEXT",
				  "value": "Давлетов"
				},
				{
				  "key": "question_m6E8PY",
				  "label": "Имя руководителя",
				  "type": "INPUT_TEXT",
				  "value": "Дагар"
				},
				{
				  "key": "question_w7rRE0",
				  "label": "Отчество руководителя",
				  "type": "INPUT_TEXT",
				  "value": "Гусманович"
				},
				{
				  "key": "question_wbQ5z1",
				  "label": "Должность руководителя",
				  "type": "INPUT_TEXT",
				  "value": "Директор"
				},
				{
				  "key": "question_wAe79e",
				  "label": "Номер телефона руководителя",
				  "type": "INPUT_PHONE_NUMBER",
				  "value": "+77777777774"
				},
				{
				  "key": "question_nWW6QP",
				  "label": "Фамилия конт. лица",
				  "type": "INPUT_TEXT",
				  "value": "Тлекбаи"
				},
				{
				  "key": "question_waL1kE",
				  "label": "Имя конт. лица",
				  "type": "INPUT_TEXT",
				  "value": "Али"
				},
				{
				  "key": "question_m6lY4O",
				  "label": "Отчество конт. лица",
				  "type": "INPUT_TEXT",
				  "value": "Кайратулы"
				},
				{
				  "key": "question_w7AYj9",
				  "label": "Должность конт. лица",
				  "type": "INPUT_TEXT",
				  "value": "Разработчик"
				},
				{
				  "key": "question_wbpMg2",
				  "label": "Телефона конт. лица",
				  "type": "INPUT_PHONE_NUMBER",
				  "value": "+77777777777"
				},
				{
				  "key": "question_nppbz8",
				  "label": "Эл. адрес конт. лица",
				  "type": "INPUT_EMAIL",
				  "value": "a@gmail.com"
				},
				{
				  "key": "question_wMpeVk",
				  "label": "Сведения о реализуемых отечественных товарах обрабатывающей промышленности и/или о предоставляемых ИКУ",
				  "type": "INPUT_TEXT",
				  "value": "Information"
				},
				{
					"key": "question_wMpeVk",
					"label": "Наименование товаров с указанием товарной позиции на уровне 6 и более знаков ЕТН ВЭД ЕАЭС и/или ИКУ на уровне не менее 4 знаков ОКВЭД",
					"type": "INPUT_TEXT",
					"value": "Goods"
				},
				{
					"key": "question_wMpeVk",
					"label": "token",
					"type": "INPUT_TEXT",
					"value": "token-1"
				},
				{
					"key": "question_mKbGd7",
					"label": "Файл соглашения",
					"type": "FILE_UPLOAD",
					"value": [
						{
							"id": "jaakzE",
							"name": "sample.pdf",
							"url": "https://storage.googleapis.com/tally-response-assets/BbbVld/4c61ceca-fbaa-4e8b-bbd2-bf35321b62ae/sample.pdf",
							"mimeType": "application/pdf",
							"size": 3028
						}
					]
				}
			  ]
			}
		  }
		`)

		expAppl = &domain.Application{
			From:                  "Doodocs Ltd.",
			GovReg:                "120120102",
			FactAddr:              "Астана",
			Bin:                   "950223347566",
			Industry:              "Информационные технологии",
			Activity:              "ЭДО",
			EmpCount:              "1-5",
			ProductCapacity:       "900",
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
			InfoManufacturedGoods: "Information",
			NameOfGoods:           "Goods",
			AgreementFile:         "https://storage.googleapis.com/tally-response-assets/BbbVld/4c61ceca-fbaa-4e8b-bbd2-bf35321b62ae/sample.pdf",
			Token:                 "token-1",
		}
	)

	appl, err := Decode(jsonBytes)
	require.Nil(t, err)
	require.Equal(t, expAppl, appl)
}
