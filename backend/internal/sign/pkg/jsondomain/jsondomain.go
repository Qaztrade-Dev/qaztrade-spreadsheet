package jsondomain

import "github.com/doodocs/qaztrade/backend/internal/sign/domain"

type Sheet struct {
	Title    string     `json:"title,omitempty"`
	SheetID  int64      `json:"sheet_id,omitempty"`
	Expenses float64    `json:"expenses"`
	Rows     int64      `json:"rows"`
	Data     [][]string `json:"data"`
	Header   [][]string `json:"header"`
}

type Application struct {
	From                  string `json:"from"`
	GovReg                string `json:"gov_reg"`
	FactAddr              string `json:"fact_addr"`
	Bin                   string `json:"bin"`
	Industry              string `json:"industry"`
	IndustryOther         string `json:"industry_other"`
	Activity              string `json:"activity"`
	EmpCount              string `json:"emp_count"`
	TaxSum                string `json:"tax_sum"`
	ProductCapacity       string `json:"product_capacity"`
	Manufacturer          string `json:"manufacturer"`
	Item                  string `json:"item"`
	ItemVolume            string `json:"item_volume"`
	FactVolumeEarnings    string `json:"fact_volume_earnings"`
	FactWorkload          string `json:"fact_workload"`
	ChiefLastname         string `json:"chief_lastname"`
	ChiefFirstname        string `json:"chief_firstname"`
	ChiefMiddlename       string `json:"chief_middlename"`
	ChiefPosition         string `json:"chief_position"`
	ChiefPhone            string `json:"chief_phone"`
	ContLastname          string `json:"cont_lastname"`
	ContFirstname         string `json:"cont_firstname"`
	ContMiddlename        string `json:"cont_middlename"`
	ContPosition          string `json:"cont_position"`
	ContPhone             string `json:"cont_phone"`
	ContEmail             string `json:"cont_email"`
	InfoManufacturedGoods string `json:"info_manufactured_goods"`
	NameOfGoods           string `json:"name_of_goods"`
	HasAgreement          string `json:"has_agreement"`
	SpendPlan             string `json:"spend_plan"`
	SpendPlanOther        string `json:"spend_plan_other"`
	Metrics2022           string `json:"metrics_2022"`
	Metrics2023           string `json:"metrics_2023"`
	Metrics2024           string `json:"metrics_2024"`
	Metrics2025           string `json:"metrics_2025"`
	AgreementFile         string `json:"agreement_file"`
	ExpensesSum           string `json:"expenses_sum"`
	ExpensesList          string `json:"expenses_list"`
	ApplicationDate       string `json:"application_date"`
}

type ApplicationAttrs struct {
	Application *Application `json:"application"`
	SheetsAgg   *Sheet       `json:"sheets_agg"`
	Sheets      []*Sheet     `json:"sheets"`
}

func EncodeSheet(input *domain.Sheet) *Sheet {
	return &Sheet{
		Title:    input.Title,
		SheetID:  input.SheetID,
		Expenses: input.Expenses,
		Rows:     input.Rows,
		Data:     input.Data,
		Header:   input.Header,
	}
}

func EncodeApplication(input *domain.Application) *Application {
	return &Application{
		From:                  input.From,
		GovReg:                input.GovReg,
		FactAddr:              input.FactAddr,
		Bin:                   input.Bin,
		Industry:              input.Industry,
		IndustryOther:         input.IndustryOther,
		Activity:              input.Activity,
		EmpCount:              input.EmpCount,
		TaxSum:                input.TaxSum,
		ProductCapacity:       input.ProductCapacity,
		Manufacturer:          input.Manufacturer,
		Item:                  input.Item,
		ItemVolume:            input.ItemVolume,
		FactVolumeEarnings:    input.FactVolumeEarnings,
		FactWorkload:          input.FactWorkload,
		ChiefLastname:         input.ChiefLastname,
		ChiefFirstname:        input.ChiefFirstname,
		ChiefMiddlename:       input.ChiefMiddlename,
		ChiefPosition:         input.ChiefPosition,
		ChiefPhone:            input.ChiefPhone,
		ContLastname:          input.ContLastname,
		ContFirstname:         input.ContFirstname,
		ContMiddlename:        input.ContMiddlename,
		ContPosition:          input.ContPosition,
		ContPhone:             input.ContPhone,
		ContEmail:             input.ContEmail,
		InfoManufacturedGoods: input.InfoManufacturedGoods,
		NameOfGoods:           input.NameOfGoods,
		HasAgreement:          input.HasAgreement,
		SpendPlan:             input.SpendPlan,
		SpendPlanOther:        input.SpendPlanOther,
		Metrics2022:           input.Metrics2022,
		Metrics2023:           input.Metrics2023,
		Metrics2024:           input.Metrics2024,
		Metrics2025:           input.Metrics2025,
		AgreementFile:         input.AgreementFile,
		ExpensesSum:           input.ExpensesSum,
		ExpensesList:          input.ExpensesList,
		ApplicationDate:       input.ApplicationDate,
	}
}

func EncodeApplicationAttrs(input *domain.ApplicationAttrs) *ApplicationAttrs {
	return &ApplicationAttrs{
		Application: EncodeApplication(input.Application),
		SheetsAgg:   EncodeSheet(input.SheetsAgg),
		Sheets:      EncodeSlice(input.Sheets, EncodeSheet),
	}
}

// encoder func for encoding from D (domain) to J (json)
type encoder[D any, J any] func(*D) *J

// EncodeSlice encodes slice of domain objects to a slice of json objects.
func EncodeSlice[D any, J any](slice []*D, encode encoder[D, J]) []*J {
	result := make([]*J, len(slice))
	for i, v := range slice {
		result[i] = encode(v)
	}
	return result
}
