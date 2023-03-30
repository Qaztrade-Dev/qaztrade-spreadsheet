package domain

import (
	"context"
	"errors"
)

type Relations map[string]*Node

type Node struct {
	Name string
	Key  string
}

var (
	Parents = Relations{
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
			Key:  "Дистрибьюторский договор.№",
		},
	}

	Children = Relations{
		"root": &Node{
			Name: "№",
			Key:  "№",
		},
		"№": &Node{
			Name: "Дистрибьюторский договор",
			Key:  "Дистрибьюторский договор.№",
		},
		"Дистрибьюторский договор": &Node{
			Name: "контракт на поставку",
			Key:  "контракт на поставку.№",
		},
	}
)

type Payload struct {
	ParentID string
	ChildKey string
	Value    PayloadValue
}

type PayloadValue map[string]interface{}

var (
	ErrorSheetPresent = errors.New("sheet already present")
)

type SheetsRepository interface {
	InsertRecord(ctx context.Context, spreadsheetID, sheetName string, sheetID int64, payload *Payload) error
	UpdateApplication(ctx context.Context, spreadsheetID string, application *Application) error
	AddSheet(ctx context.Context, spreadsheetID string, sheetName string) error
}
