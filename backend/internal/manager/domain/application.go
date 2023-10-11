package domain

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	StatusUserFilling      = "user_filling"
	StatusManagerReviewing = "manager_reviewing"
	StatusUserFixing       = "user_fixing"
	StatusCompleted        = "completed"
	StatusRejected         = "rejected"
)

type Application struct {
	ID             string
	UserID         string
	No             int
	SpreadsheetID  string
	Link           string
	Status         string
	SignDocumentID string
	Attrs          interface{}
	SignedAt       time.Time
	CreatedAt      time.Time
}

type ApplicationList struct {
	OverallCount uint64
	Applications []*Application
}

type GetManyInput struct {
	Limit  uint64
	Offset uint64

	ApplicationID    string
	BIN              string
	CompensationType string
	SignedAtFrom     time.Time
	SignedAtUntil    time.Time
	CompanyName      string
	ApplicationNo    int
}

type Revision struct {
	ApplicationID  string
	SpreadsheetID  string
	No             int
	Link           string
	Address        string
	BIN            string
	Manufactor     string
	To             string
	ApplicantEmail string
	ManagerName    string
	ManagerEmail   string
	SignedAt       time.Time
	Remarks        string
}

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
type ApplicationRepository interface {
	GetMany(ctx context.Context, query *GetManyInput) (*ApplicationList, error)
	GetOne(ctx context.Context, query *GetManyInput) (*Application, error)
	EditStatus(ctx context.Context, applicationID, statusName string) error
	IsManagerAssigned(ctx context.Context, applicationID, userID string) (bool, error)
}

type SpreadsheetService interface {
	SwitchModeRead(ctx context.Context, spreadsheetID string) error
	SwitchModeEdit(ctx context.Context, spreadsheetID string) error
	LockSheets(ctx context.Context, spreadsheetID string) error
	GrantAdminPermissions(ctx context.Context, spreadsheetID, email string) error
	Comments(ctx context.Context, application *Application, managerName string) (*Revision, error)
	GetApplication(ctx context.Context, spreadsheetID string) (*ApplicationAttrs, error)
}

var (
	ErrorApplicationNotSigned      = fmt.Errorf("Заявление еще не подписано!")
	ErrorApplicationNotUnderReview = fmt.Errorf("Cтатус заявления не соответствует требованиям")
	ErrorPermissionDenied          = fmt.Errorf("Доступ запрещен")
)

type SigningService interface {
	GetDDCard(ctx context.Context, documentID string) (*http.Response, error)
}

type NoticeService interface {
	Create(revision *Revision) (*bytes.Buffer, error)
}

type Storage interface {
	Upload(ctx context.Context, folderName, fileName string, fileSize int64, fileReader io.Reader) (string, error)
	Remove(ctx context.Context, filePath string) error
}

type EmailService interface {
	SendNotice(ctx context.Context, toEmail, mailName, filename string, FileReader io.Reader) error
}
