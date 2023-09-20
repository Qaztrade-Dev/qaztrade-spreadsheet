package adapters

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/api/sheets/v4"
)

type BatchUpdate struct {
	service  *sheets.Service
	requests []*sheets.Request
}

func NewBatchUpdate(service *sheets.Service) *BatchUpdate {
	return &BatchUpdate{
		service:  service,
		requests: make([]*sheets.Request, 0),
	}
}

func (b *BatchUpdate) WithRequest(requests ...*sheets.Request) {
	b.requests = append(b.requests, requests...)
}

func (b *BatchUpdate) Do(ctx context.Context, spreadsheetID string) error {
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{Requests: b.requests}
	_, err := b.service.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Context(ctx).Do()
	return err
}

type Cell struct {
	Col int64
	Row int64
}

type Range struct {
	From *Cell
	To   *Cell
}

func columnToNumber(column string) int64 {
	column = strings.ToUpper(column)
	var num int64 = 0
	for _, char := range column {
		num = num*26 + int64(char-'A') + 1
	}
	return num
}

func A1ToCell(a1 string) *Cell {
	for i, r := range a1 {
		if r >= '0' && r <= '9' {
			col := a1[:i]
			row, _ := strconv.Atoi(a1[i:])
			return &Cell{
				Col: columnToNumber(col) - 1,
				Row: int64(row - 1),
			}
		}
	}
	return nil
}

func numberToColumn(num int64) string {
	column := ""
	for num >= 0 {
		column = string(rune((num%26)+'A')) + column
		num = num/26 - 1
	}
	return column
}

func (cell *Cell) Equals(b *Cell) bool {
	return cell.Col == b.Col && cell.Row == b.Row
}

func (cell *Cell) ToA1() string {
	return numberToColumn(cell.Col) + strconv.Itoa(int(cell.Row)+1)
}

func (r *Range) ToA1() string {
	return r.From.ToA1() + ":" + r.To.ToA1()
}

func A1ToRange(fromToA1 string) *Range {
	splitted := strings.Split(fromToA1, ":")
	if len(splitted) != 2 {
		return nil
	}

	var (
		fromCell = A1ToCell(splitted[0])
		toCell   = A1ToCell(splitted[1])
	)

	return &Range{
		From: fromCell,
		To:   toCell,
	}
}
