package tally

import (
	"encoding/json"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

type tallyFieldOption struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type tallyField struct {
	Label   string             `json:"label"`
	Type    string             `json:"type"`
	Value   string             `json:"value"`
	Options []tallyFieldOption `json:"options"`
}

type tallyResponseData struct {
	Fields []tallyField `json:"fields"`
}

type tallyResponse struct {
	Data tallyResponseData `json:"data"`
}

func Encode(jsonBytes []byte) (*domain.Application, error) {
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
		"Вид деятельности":                            &appl.Activity,
		"Численность сотрудников":                     &appl.EmpCount,
		"Производитель":                               &appl.Manufacturer,
		"Товар":                                       &appl.Item,
		"Объем товара":                                &appl.ItemVolume,
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
		"Страна происхождения товара": &appl.Country,
		"Код ТНВЭД (6 знаков)":        &appl.CodeTnved,
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
	default:
		return extractText(field)
	}
}

func extractText(field *tallyField) string {
	return field.Value
}

func extractDropdown(field *tallyField) string {
	for _, option := range field.Options {
		if option.ID == field.Value {
			return option.Text
		}
	}
	return ""
}
