package domain

import (
	"context"
	"errors"
	"io"
)

type (
	Relations map[string]*Node
	Node      struct {
		Name string
		Key  string
	}
)

var (
	SheetParents = map[string]Relations{
		"Доставка ЖД транспортом": TrainDeliveryParents,
		"Затраты на продвижение":  AdvertisementExpensesParents,
	}

	SheetChildren = map[string]Relations{
		"Доставка ЖД транспортом": TrainDeliveryChildren,
		"Затраты на продвижение":  AdvertisementExpensesChildren,
	}
)

// Доставка ЖД транспортом
var (
	TrainDeliveryParents = Relations{
		"№": &Node{
			Name: "root",
			Key:  "root",
		},
		"Дистрибьюторский договор": &Node{
			Name: "№",
			Key:  "№",
		},
		"контракт на поставку": &Node{
			Name: "Дистрибьюторский договор",
			Key:  "Дистрибьюторский договор|№",
		},
	}

	TrainDeliveryChildren = Relations{
		"root": &Node{
			Name: "№",
			Key:  "№",
		},
		"№": &Node{
			Name: "Дистрибьюторский договор",
			Key:  "Дистрибьюторский договор|№",
		},
		"Дистрибьюторский договор": &Node{
			Name: "контракт на поставку",
			Key:  "контракт на поставку|№",
		},
	}
)

// Затраты на продвижение
var (
	AdvertisementExpensesParents = Relations{
		"№": &Node{
			Name: "root",
			Key:  "root",
		},
	}

	AdvertisementExpensesChildren = Relations{
		"root": &Node{
			Name: "№",
			Key:  "№",
		},
	}
)

type (
	PayloadValue map[string]interface{}

	Payload struct {
		RowNumber int
		ParentID  string
		ChildKey  string
		Value     PayloadValue
	}
)

var ErrorSheetPresent = errors.New("sheet already present")

type (
	RemoveInput struct {
		Value  string
		RowNum int
		Name   string
	}

	SheetsRepository interface {
		InsertRecord(ctx context.Context, spreadsheetID, sheetName string, sheetID int64, payload *Payload) error
		UpdateApplication(ctx context.Context, spreadsheetID string, application *Application) error
		AddSheet(ctx context.Context, spreadsheetID string, sheetName string) error
	}

	Storage interface {
		Upload(ctx context.Context, folderName, fileName string, fileSize int64, fileReader io.Reader) (string, error)
	}
)

type Claims struct {
	SpreadsheetID string `json:"sid"`
}
