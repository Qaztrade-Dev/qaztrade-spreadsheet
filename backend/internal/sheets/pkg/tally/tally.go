package tally

import (
	"encoding/json"
	"strconv"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

type tallyFieldOption struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type tallyField struct {
	Label   string             `json:"label"`
	Type    string             `json:"type"`
	Value   interface{}        `json:"value"`
	Options []tallyFieldOption `json:"options"`
}

type tallyResponseData struct {
	Fields []tallyField `json:"fields"`
}

type tallyResponse struct {
	Data tallyResponseData `json:"data"`
}

func Decode(jsonBytes []byte) (*domain.Application, error) {
	var (
		response tallyResponse
		appl     domain.Application
	)

	err := json.Unmarshal(jsonBytes, &response)
	if err != nil {
		return nil, err
	}

	mapping := map[string]*string{
		"От кого": &appl.From,
		"Государственная регистрация/перерегистрация": &appl.GovReg,
		"Фактический адрес":                           &appl.FactAddr,
		"БИН/ИИН":                                     &appl.Bin,
		"Наименование отрасли":                        &appl.Industry,
		"Индустрия, если выбрано \"Иное\"":            &appl.IndustryOther,
		"Вид деятельности":                            &appl.Activity,
		"Численность сотрудников":                     &appl.EmpCount,
		"Сумма уплаченных налогов без учета косвенных налогов (налог на добавленную стоимость и акцизы) в году, предшествующем году подачи заявки": &appl.TaxSum,
		"Производственная мощность, возможности увеличения": &appl.ProductCapacity,
		"Производитель": &appl.Manufacturer,
		"Товар":         &appl.Item,
		"Объем товара":  &appl.ItemVolume,
		"Объем фактической валютной выручки за полугодие предшествующей дате подачи заявки": &appl.FactVolumeEarnings,
		"Фактическая загруженность производства":                                            &appl.FactWorkload,
		"Фамилия руководителя":        &appl.ChiefLastname,
		"Имя руководителя":            &appl.ChiefFirstname,
		"Отчество руководителя":       &appl.ChiefMiddlename,
		"Должность руководителя":      &appl.ChiefPosition,
		"Номер телефона руководителя": &appl.ChiefPhone,
		"Фамилия конт. лица":          &appl.ContLastname,
		"Имя конт. лица":              &appl.ContFirstname,
		"Отчество конт. лица":         &appl.ContMiddlename,
		"Должность конт. лица":        &appl.ContPosition,
		"Телефона конт. лица":         &appl.ContPhone,
		"Эл. адрес конт. лица":        &appl.ContEmail,
		"Сведения о реализуемых отечественных товарах обрабатывающей промышленности и/или о предоставляемых ИКУ":                                &appl.InfoManufacturedGoods,
		"Наименование товаров с указанием товарной позиции на уровне 6 и более знаков ЕТН ВЭД ЕАЭС и/или ИКУ на уровне не менее 4 знаков ОКВЭД": &appl.NameOfGoods,
		"Наличие соглашения о промышленной сборке": &appl.HasAgreement,
		"Файл соглашения": &appl.AgreementFile,
		"token":           &appl.Token,
	}

	for _, field := range response.Data.Fields {
		ptr, ok := mapping[field.Label]
		if !ok {
			continue
		}
		*ptr = extractValue(&field)
	}
	return &appl, nil
}

func extractValue(field *tallyField) string {
	switch field.Type {
	case "DROPDOWN":
		return extractDropdown(field)
	case "FILE_UPLOAD":
		return extractFileUpload(field)
	case "INPUT_NUMBER":
		return extractNumber(field)
	default:
		return extractText(field)
	}
}

func extractText(field *tallyField) string {
	if field.Value == nil {
		return ""
	}
	return field.Value.(string)
}

func extractNumber(field *tallyField) string {
	v := int(field.Value.(float64))
	return strconv.Itoa(v)
}

func extractDropdown(field *tallyField) string {
	for _, option := range field.Options {
		if option.ID == field.Value {
			return option.Text
		}
	}
	return ""
}

func extractFileUpload(field *tallyField) string {
	if field.Value == nil {
		return ""
	}

	fileValues, ok := field.Value.([]interface{})
	if !ok {
		return ""
	}

	if len(fileValues) == 0 {
		return ""
	}

	value, ok := fileValues[0].(map[string]interface{})
	if !ok {
		return ""
	}

	return value["url"].(string)
}
