package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

//go:embed credentials.json
var credentials []byte

func main() {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentials))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	var (
		spreadsheetID = "1bv_mj8-xnNzBGYmF2YqbEwNPz2IyOuZVaD4E4203trc"
		cli           = NewSheetsClient(srv)
	)

	headerCells, err := cli.GetHeaderCells(ctx, spreadsheetID)
	if err != nil {
		log.Fatalf("GetHeaderCells error: %v", err)
	}

	payload := &Payload{
		ParentID: "null",
		ChildKey: "№",
		Value: PayloadValue{
			"№": "13",
			"Производитель/дочерняя компания/дистрибьютор/СПК": "Aviata1",
			"подтверждающий документ": PayloadValue{
				"производитель":       "Doodocs",
				"наименование":        "Дудокс",
				"№":                   "3",
				"наименование товара": "Подписи",
				"ТН ВЭД (6 знаков)":   "120934",
				"дата":                "12.09.2019",
				"срок":                "123",
				"подтверждение на сайте уполномоченного органа": "http://google.com",
			},
		},
	}
	err = cli.SubmitRow(ctx, spreadsheetID, payload, headerCells)
	if err != nil {
		log.Fatalf("SubmitRow error: %v", err)
	}
}

type SheetsClient struct {
	service *sheets.Service
}

func NewSheetsClient(service *sheets.Service) *SheetsClient {
	return &SheetsClient{
		service: service,
	}
}

type HeaderCellMap map[string]*HeaderCell

type Range struct {
	Left  int
	Right int
}

type Bound struct {
	Top    int
	Bottom int
}

type HeaderCell struct {
	Key      string
	Range    Range
	Values   HeaderCellMap
	GroupKey string
}

func NewHeaderCell(key, groupKey string, rangeL, rangeR int, hcellMap HeaderCellMap) *HeaderCell {
	return &HeaderCell{
		Key:      key,
		GroupKey: groupKey,
		Range: Range{
			Left:  rangeL,
			Right: rangeR,
		},
		Values: hcellMap,
	}
}

func (h *HeaderCell) IsLeaf() bool {
	if h.Values == nil {
		return true
	}
	return len(h.Values) == 0
}

type Cell struct {
	Value      string
	RowNum     int
	ColumnNum  int
	HeaderCell *HeaderCell
}

func NewCell(value string, rowNum, columnNum int, headerCell *HeaderCell) *Cell {
	return &Cell{
		Value:      value,
		RowNum:     rowNum,
		ColumnNum:  columnNum,
		HeaderCell: headerCell,
	}
}

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

func (c *SheetsClient) GetHeaderCells(ctx context.Context, spreadsheetID string) (HeaderCellMap, error) {
	sheetRange, err := c.service.Spreadsheets.Values.Get(spreadsheetID, "Header").ValueRenderOption("FORMATTED_VALUE").Do()
	if err != nil {
		return nil, err
	}
	var (
		topLevel = sheetRange.Values[0]
		lowLevel = sheetRange.Values[1]

		hcellMap = make(HeaderCellMap)
	)

	for i := range topLevel {
		if topLevel[i] == "" {
			continue
		}

		var (
			topLevelValue = strings.TrimSpace(topLevel[i].(string))
			innerHcellMap = make(HeaderCellMap)
			rangeR        = i
		)

		for j := i; j < len(lowLevel); j++ {
			lowLevelValue := strings.TrimSpace(lowLevel[j].(string))

			if lowLevelValue == "" {
				break
			}
			if !(topLevel[j] == "" || i == j) {
				break
			}

			innerHcellMap[lowLevelValue] = NewHeaderCell(lowLevelValue, topLevelValue, j, j, nil)
			rangeR = j
		}

		hcellMap[topLevelValue] = NewHeaderCell(topLevelValue, topLevelValue, i, rangeR, innerHcellMap)
	}

	return hcellMap, nil
}

type Payload struct {
	ParentID string
	ChildKey string
	Value    PayloadValue
}

type PayloadValue map[string]interface{}

func (c *SheetsClient) FillRecord(spreadsheetID string, payload PayloadValue, headers HeaderCellMap, rowNum int) error {
	type UpdateCellRequest struct {
		RowIndex    int64
		ColumnIndex int64
		Value       string
	}
	batch := make([]*UpdateCellRequest, 0)

	for key := range payload {
		cell := headers[key]
		if cell.IsLeaf() {
			batch = append(batch, &UpdateCellRequest{
				RowIndex:    int64(rowNum),
				ColumnIndex: int64(cell.Range.Left),
				Value:       payload[key].(string),
			})
		} else {
			c.FillRecord(spreadsheetID, payload[key].(PayloadValue), cell.Values, rowNum)
		}
	}

	requests := make([]*sheets.Request, 0, len(batch))
	for i := range batch {
		elem := batch[i]
		requests = append(requests, &sheets.Request{
			UpdateCells: &sheets.UpdateCellsRequest{
				Fields: "*",
				Start: &sheets.GridCoordinate{
					RowIndex:    int64(elem.RowIndex),
					ColumnIndex: int64(elem.ColumnIndex),
					SheetId:     0,
				},
				Rows: []*sheets.RowData{
					{
						Values: []*sheets.CellData{
							{
								UserEnteredValue: &sheets.ExtendedValue{
									StringValue: &elem.Value,
								},
							},
						},
					},
				},
			},
		})
	}

	_, err := c.service.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()

	fmt.Println(err)

	return err
}

const ParentKeyRoot = "root"

func (c *SheetsClient) GetNodeBounds(spreadsheetID, nodeKey, nodeID string, headerCellMap HeaderCellMap) (*Bound, error) {
	fmt.Printf("GetNodeBounds. nodeKey=%v, nodeID=%v\n", nodeKey, nodeID)
	if nodeKey == ParentKeyRoot {
		sheetRange, err := c.service.Spreadsheets.Values.Get(spreadsheetID, "A6:A30").ValueRenderOption("FORMATTED_VALUE").Do()
		if err != nil {
			return nil, err
		}

		return &Bound{Top: 5, Bottom: 5 + len(sheetRange.Values)}, nil
	}

	fmt.Printf("%#v\n", nodeKey)
	var (
		parentHeaderCell = c.GetHeaderCell(nodeKey, headerCellMap)
		columnA1         = EncodeAlphabet(parentHeaderCell.Range.Left + 1)
	)

	fmt.Println(columnA1 + ":" + columnA1)

	sheetRange, err := c.service.Spreadsheets.Values.Get(spreadsheetID, columnA1+":"+columnA1).Do()
	fmt.Println("Get err", err)
	if err != nil {
		return nil, err
	}

	var (
		left  = 0
		right = 0
	)

	for i := range sheetRange.Values {
		fmt.Printf("i=%v, sheetRange.Values[i]=%#v\n", i, sheetRange.Values[i])
		value := RowToString(sheetRange.Values[i])

		if value == nodeID {
			left = i
			right = i
			break
		}
	}

	for i := left + 1; i < len(sheetRange.Values); i++ {
		value := RowToString(sheetRange.Values[i])
		if value != "" {
			break
		}

		right = i
	}

	result := &Bound{Top: left, Bottom: right}
	fmt.Printf("%#v\n", result)

	return result, nil
}

func (c *SheetsClient) GetHeaderCell(parentKey string, headerCellMap HeaderCellMap) *HeaderCell {
	keys := strings.Split(parentKey, ".")
	var cell *HeaderCell
	for _, key := range keys {
		if cell == nil {
			cell = headerCellMap[key]
		} else {
			cell = cell.Values[key]
		}
		fmt.Printf("key=%#v,cell=%#v\n", key, cell)
	}

	return cell
}

var (
	ErrorChildNotFound = errors.New("child not found")
	ErrorChildNotMatch = errors.New("child name doesn't match")
)

func (c *SheetsClient) GetChildNode(parentKey, childName string) (*Node, error) {
	if _, ok := Children[parentKey]; !ok {
		return nil, ErrorChildNotFound
	}

	child := Children[parentKey]
	if child.Name != childName {
		return nil, ErrorChildNotMatch
	}

	return Children[parentKey], nil
}

func (c *SheetsClient) GetLastChildCell(spreadsheetID string, parentBounds *Bound, childHeaderCell *HeaderCell) (*Cell, error) {
	fmt.Printf("GetLastChildCell. parentBounds=%#v, childHeaderCell=%#v", parentBounds, childHeaderCell)
	var (
		columnNum = childHeaderCell.Range.Left
		rowTop    = parentBounds.Top
		rowBottom = parentBounds.Bottom

		columnA1 = EncodeAlphabet(columnNum + 1)
		fromA1   = fmt.Sprintf("%s%d", columnA1, rowTop+1)
		toA1     = fmt.Sprintf("%s%d", columnA1, rowBottom+1)
	)

	sheetRange, err := c.service.Spreadsheets.Values.Get(spreadsheetID, fromA1+":"+toA1).Do()
	if err != nil {
		return nil, err
	}

	var (
		lastIdx   = 0
		lastValue = ""
	)

	for i, row := range sheetRange.Values {
		value := RowToString(row)
		if i == 0 && value == "" {
			return nil, nil
		}
		if value != "" {
			lastIdx = i
			lastValue = value
		}
	}

	if lastValue == "" {
		return nil, nil
	}

	return NewCell(lastValue, rowTop+lastIdx, columnNum, childHeaderCell), nil
}

// 1. get parent row
// 2. get last child of the parent, e.g. neighbor
// 3. get last row of the farthest descendent
func (c *SheetsClient) GetRowNum(ctx context.Context, spreadsheetID, parentID, childName string, headerCellMap HeaderCellMap) (int, bool, error) {
	fmt.Printf("GetRowNum. parentID:%v, childName:%v\n", parentID, childName)
	var parentKey = Parents[childName].Key

	parentBounds, err := c.GetNodeBounds(spreadsheetID, parentKey, parentID, headerCellMap)
	fmt.Printf("parentBounds: %#v\n", parentBounds)
	if err != nil {
		return 0, false, err
	}

	var upperBound = parentBounds.Top

	child, err := c.GetChildNode(parentKey, childName)
	fmt.Printf("child: %#v\n", child)
	if err != nil {
		return 0, false, err
	}

	childHeaderCell := c.GetHeaderCell(child.Key, headerCellMap)

	lastChildCell, err := c.GetLastChildCell(spreadsheetID, parentBounds, childHeaderCell)
	fmt.Printf("lastChildCell: %#v\n", lastChildCell)
	if err != nil {
		return 0, false, err
	}

	if lastChildCell == nil {
		return upperBound, false, nil
	}

	grandChildNode, ok := Children[lastChildCell.HeaderCell.GroupKey]
	if !ok {
		return lastChildCell.RowNum, true, nil
	}

	rowNum, _, err := c.GetRowNum(ctx, spreadsheetID, lastChildCell.Value, grandChildNode.Name, headerCellMap)
	if err != nil {
		return 0, false, err
	}

	return rowNum, true, nil
}

func (c *SheetsClient) SubmitRow(ctx context.Context, spreadsheetID string, payload *Payload, headerCellMap HeaderCellMap) error {
	rowNum, mustInsertRow, err := c.GetRowNum(ctx, spreadsheetID, payload.ParentID, payload.ChildKey, headerCellMap)
	if err != nil {
		return err
	}

	fmt.Printf("SubmitRow. rowNum=%v mustInsertRow=%v\n", rowNum, mustInsertRow)

	if mustInsertRow {
		if err := c.insertRowAfter(ctx, spreadsheetID, rowNum); err != nil {
			return err
		}
	}

	err = c.FillRecord(spreadsheetID, payload.Value, headerCellMap, rowNum)
	if err != nil {
		return err
	}

	return nil
}

func (c *SheetsClient) insertRowAfter(ctx context.Context, spreadsheetID string, rowIndex int) error {
	request := &sheets.Request{
		InsertDimension: &sheets.InsertDimensionRequest{
			Range: &sheets.DimensionRange{
				SheetId:    0,
				Dimension:  "ROWS",
				StartIndex: int64(rowIndex) + 1, // Start index is 1-based
				EndIndex:   int64(rowIndex) + 2, // End index is exclusive
			},
			InheritFromBefore: true,
		},
	}

	batchUpdateRequest := sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{request},
	}

	_, err := c.service.Spreadsheets.BatchUpdate(spreadsheetID, &batchUpdateRequest).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to insert row: %v", err)
	}

	return err
}

func EncodeAlphabet(num int) string {
	result := ""
	for num > 0 {
		remainder := (num - 1) % 26
		result = string(65+remainder) + result
		num = (num - 1) / 26
	}
	return result
}

func GetColumnNum(column string) int {
	column = regexp.MustCompile(`\d+`).ReplaceAllString(column, "")

	result := 0

	for i := 0; i < len(column); i++ {
		result *= 26
		result += int(column[i]) - int('A') + 1
	}

	return result
}

func RowToString(row []interface{}) string {
	if len(row) == 0 {
		return ""
	}

	return strings.TrimSpace(row[0].(string))
}
