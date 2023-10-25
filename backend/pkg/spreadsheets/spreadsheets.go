package spreadsheets

import (
	"context"
	"fmt"
)

type ApplicationAttrs struct {
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
	ExpensesSum           string
	ExpensesList          string
	ApplicationDate       string
}

type SpreadsheetService interface {
	SwitchModeRead(ctx context.Context, spreadsheetID string) error
	SwitchModeEdit(ctx context.Context, spreadsheetID string) error
	LockSheets(ctx context.Context, spreadsheetID string) error
	GrantAdminPermissions(ctx context.Context, spreadsheetID, email string) error
	GetApplicationAttrs(ctx context.Context, spreadsheetID string) (*ApplicationAttrs, error)
	GetSheetData(ctx context.Context, spreadsheetID string, sheetTitle string) ([][]string, error)
}

var (
	ErrorSheetNotFound = fmt.Errorf("sheet not found")
)
