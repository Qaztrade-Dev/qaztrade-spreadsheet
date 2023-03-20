package domain

import "context"

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

type SheetsRepository interface {
	InsertRecord(ctx context.Context, spreadsheetID string, payload *Payload) error
}
