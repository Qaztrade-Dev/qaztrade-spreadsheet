package adapters

import (
	"strconv"
	"strings"

	"google.golang.org/api/sheets/v4"
)

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

func UnmergeRequest(sheetID int64, fromToA1 string) *sheets.Request {
	cellRange := A1ToRange(fromToA1)

	return &sheets.Request{
		UnmergeCells: &sheets.UnmergeCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    cellRange.From.Row,
				EndRowIndex:      cellRange.To.Row + 1,
				StartColumnIndex: cellRange.From.Col,
				EndColumnIndex:   cellRange.To.Col + 1,
			},
		},
	}
}

func MergeRequest(sheetID int64, fromToA1 string) *sheets.Request {
	cellRange := A1ToRange(fromToA1)

	return &sheets.Request{
		MergeCells: &sheets.MergeCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    cellRange.From.Row,
				EndRowIndex:      cellRange.To.Row + 1,
				StartColumnIndex: cellRange.From.Col,
				EndColumnIndex:   cellRange.To.Col + 1,
			},
			MergeType: "MERGE_ALL",
		},
	}
}
