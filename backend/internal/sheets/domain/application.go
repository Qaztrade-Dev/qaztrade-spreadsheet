package domain

const (
	StatusUserFilling      = "user_filling"
	StatusManagerReviewing = "manager_reviewing"
	StatusUserFixing       = "user_fixing"
	StatusCompleted        = "completed"
	StatusRejected         = "rejected"
)

type Application struct {
	From                  string
	GovReg                string
	FactAddr              string
	Bin                   string
	Industry              string
	IndustryOther         string
	Activity              string
	EmpCount              string
	TaxSum                string
	ProductCapacity       string
	Manufacturer          string
	Item                  string
	ItemVolume            string
	FactVolumeEarnings    string
	FactWorkload          string
	ChiefLastname         string
	ChiefFirstname        string
	ChiefMiddlename       string
	ChiefPosition         string
	ChiefPhone            string
	ContLastname          string
	ContFirstname         string
	ContMiddlename        string
	ContPosition          string
	ContPhone             string
	ContEmail             string
	InfoManufacturedGoods string
	NameOfGoods           string
	HasAgreement          string
	SpendPlan             string
	SpendPlanOther        string
	Metrics2022           string
	Metrics2023           string
	Metrics2024           string
	Metrics2025           string
	AgreementFile         string

	Token string
}

type StatusApplication struct {
	SpreadsheetID string
	Status        string
}
