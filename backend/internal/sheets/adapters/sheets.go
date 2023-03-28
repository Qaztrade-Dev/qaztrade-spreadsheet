package adapters

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsClient struct {
	service *sheets.Service
}

var _ domain.SheetsRepository = (*SheetsClient)(nil)

func NewSheetsClient(ctx context.Context, credentialsJson []byte) (*SheetsClient, error) {
	service, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, err
	}

	return &SheetsClient{
		service: service,
	}, nil
}

func (c *SheetsClient) InsertRecord(ctx context.Context, spreadsheetID string, payload *domain.Payload) error {
	spreadsheetClient, err := c.NewSpreadsheetClient(ctx, spreadsheetID)
	if err != nil {
		return err
	}

	return spreadsheetClient.InsertRecord(ctx, payload)
}

type SpreadsheetClient struct {
	service       *sheets.Service
	spreadsheetID string
	headersMap    HeaderCellMap
}

func (c *SheetsClient) NewSpreadsheetClient(ctx context.Context, spreadsheetID string) (*SpreadsheetClient, error) {
	spreadsheetClient := &SpreadsheetClient{
		service:       c.service,
		spreadsheetID: spreadsheetID,
	}

	// TODO
	// cache or reuse or hard-code
	headersMap, err := spreadsheetClient.getHeaderCells(ctx)
	if err != nil {
		return nil, err
	}

	spreadsheetClient.headersMap = headersMap

	return spreadsheetClient, nil
}

func (c *SpreadsheetClient) getHeaderCells(ctx context.Context) (HeaderCellMap, error) {
	sheetRange, err := c.service.Spreadsheets.Values.Get(c.spreadsheetID, "Header").
		ValueRenderOption("FORMATTED_VALUE").
		Context(ctx).
		Do()
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

type UpdateCellRequest struct {
	RowIndex    int64
	ColumnIndex int64
	Value       string
}

// TODO
// FillRecord construct batch
func (c *SpreadsheetClient) FillRecord(ctx context.Context, payload domain.PayloadValue, headers HeaderCellMap, rowNum int) error {
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
			var p interface{} = payload[key]
			var m map[string]interface{} = p.(map[string]interface{})
			var d domain.PayloadValue = domain.PayloadValue(m)

			if err := c.FillRecord(ctx, d, cell.Values, rowNum); err != nil {
				return err
			}
		}
	}

	if len(batch) == 0 {
		return nil
	}

	requests := make([]*sheets.Request, 0, len(batch))
	for i := range batch {
		requests = append(requests, batch[i].Encode())
	}

	_, err := c.service.Spreadsheets.BatchUpdate(
		c.spreadsheetID,
		&sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		},
	).Context(ctx).Do()

	return err
}

func (r *UpdateCellRequest) Encode() *sheets.Request {
	return &sheets.Request{
		UpdateCells: &sheets.UpdateCellsRequest{
			Fields: "*",
			Start: &sheets.GridCoordinate{
				RowIndex:    int64(r.RowIndex),
				ColumnIndex: int64(r.ColumnIndex),
				SheetId:     0,
			},
			Rows: []*sheets.RowData{
				{
					Values: []*sheets.CellData{
						{
							UserEnteredValue: &sheets.ExtendedValue{
								StringValue: &r.Value,
							},
						},
					},
				},
			},
		},
	}
}

const ParentKeyRoot = "root"

func (c *SpreadsheetClient) GetNodeBounds(ctx context.Context, nodeKey, nodeID string, parentBound ...*Bound) (*Bound, error) {
	fmt.Printf("GetNodeBounds. nodeKey=%v, nodeID=%v\n", nodeKey, nodeID)
	if nodeKey == ParentKeyRoot {
		sheetRange, err := c.service.Spreadsheets.Values.Get(c.spreadsheetID, "A6:A").
			ValueRenderOption("FORMATTED_VALUE").
			Context(ctx).
			Do()
		if err != nil {
			return nil, err
		}

		return &Bound{Top: 5, Bottom: 5 + len(sheetRange.Values)}, nil
	}

	fmt.Printf("%#v\n", nodeKey)
	var (
		parentHeaderCell = c.GetHeaderCell(nodeKey)
		columnA1         = EncodeColumn(parentHeaderCell.Range.Left)
		range_           = columnA1 + ":" + columnA1

		left   = 0
		right  = 0
		offset = 0
	)

	if len(parentBound) > 0 {
		range_ = parentBound[0].EncodeRange(parentHeaderCell.Range.Left)
		offset = parentBound[0].Top
	}

	sheetRange, err := c.service.Spreadsheets.Values.Get(c.spreadsheetID, range_).
		Context(ctx).
		Do()
	fmt.Println("Get err", err)
	if err != nil {
		return nil, err
	}

	for i := range sheetRange.Values {
		fmt.Printf("left: i=%v, sheetRange.Values[i]=%#v\n", i, sheetRange.Values[i])
		value := DecodeRow(sheetRange.Values[i])

		if value == nodeID {
			left = i
			right = i
			break
		}
	}

	for i := left + 1; i < len(sheetRange.Values); i++ {
		fmt.Printf("right: i=%v, sheetRange.Values[i]=%#v\n", i, sheetRange.Values[i])
		value := DecodeRow(sheetRange.Values[i])
		if value != "" {
			break
		}

		right = i
	}

	result := &Bound{Top: left + offset, Bottom: right + offset}
	fmt.Printf("%#v\n", result)

	return result, nil
}

func (c *SpreadsheetClient) GetHeaderCell(parentKey string) *HeaderCell {
	keys := strings.Split(parentKey, ".")
	var cell *HeaderCell
	for _, key := range keys {
		if cell == nil {
			cell = c.headersMap[key]
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

func (c *SpreadsheetClient) GetChildNode(parentKey, childName string) (*domain.Node, error) {
	fmt.Printf("GetChildNode. parentKey=%v, childName=%v\n", parentKey, childName)
	if _, ok := domain.Children[parentKey]; !ok {
		return nil, ErrorChildNotFound
	}

	child := domain.Children[parentKey]
	if child.Name != childName {
		return nil, ErrorChildNotMatch
	}

	return domain.Children[parentKey], nil
}

func (c *SpreadsheetClient) GetLastChildCell(ctx context.Context, parentBounds *Bound, childHeaderCell *HeaderCell) (*Cell, error) {
	fmt.Printf("GetLastChildCell. parentBounds=%#v, childHeaderCell=%#v", parentBounds, childHeaderCell)
	var (
		columnNum = childHeaderCell.Range.Left
		rowTop    = parentBounds.Top
		rowBottom = parentBounds.Bottom

		fromA1 = EncodeCoordinate(columnNum, rowTop)
		toA1   = EncodeCoordinate(columnNum, rowBottom)
		range_ = fromA1 + ":" + toA1
	)

	sheetRange, err := c.service.Spreadsheets.Values.Get(c.spreadsheetID, range_).
		Context(ctx).
		Do()
	if err != nil {
		return nil, err
	}

	var (
		lastIdx   = 0
		lastValue = ""
	)

	for i, row := range sheetRange.Values {
		value := DecodeRow(row)
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
func (c *SpreadsheetClient) GetRowNum(ctx context.Context, parentID, childName string, bounds ...*Bound) (int, bool, error) {
	fmt.Printf("GetRowNum. parentID:%v, childName:%v\n", parentID, childName)
	var (
		parentKey  = domain.Parents[childName].Key
		parentName = domain.Parents[childName].Name
	)

	parentBounds, err := c.GetNodeBounds(ctx, parentKey, parentID, bounds...)
	fmt.Printf("parentBounds: %#v\n", parentBounds)
	if err != nil {
		return 0, false, err
	}

	var upperBound = parentBounds.Top

	child, err := c.GetChildNode(parentName, childName)
	fmt.Printf("child: %#v\n", child)
	if err != nil {
		return 0, false, err
	}

	childHeaderCell := c.GetHeaderCell(child.Key)

	lastChildCell, err := c.GetLastChildCell(ctx, parentBounds, childHeaderCell)
	fmt.Printf("lastChildCell: %#v\n", lastChildCell)
	if err != nil {
		return 0, false, err
	}

	if lastChildCell == nil {
		return upperBound, false, nil
	}

	grandChildNode, ok := domain.Children[lastChildCell.HeaderCell.GroupKey]
	if !ok {
		return lastChildCell.RowNum, true, nil
	}

	rowNum, _, err := c.GetRowNum(ctx, lastChildCell.Value, grandChildNode.Name, parentBounds)
	if err != nil {
		return 0, false, err
	}

	return rowNum, true, nil
}

func (c *SpreadsheetClient) InsertRecord(ctx context.Context, payload *domain.Payload) error {
	rowNum, mustInsertRow, err := c.GetRowNum(ctx, payload.ParentID, payload.ChildKey)
	if err != nil {
		return err
	}

	fmt.Printf("SubmitRow. rowNum=%v mustInsertRow=%v\n", rowNum, mustInsertRow)

	if mustInsertRow {
		if err := c.insertRowAfter(ctx, rowNum); err != nil {
			return err
		}
		rowNum += 1
	}

	err = c.FillRecord(ctx, payload.Value, c.headersMap, rowNum)
	if err != nil {
		return err
	}

	return nil
}

func (c *SpreadsheetClient) UpdateApplication(ctx context.Context, spreadsheetID string, application *domain.Application) error {
	a := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "RAW",
		Data: []*sheets.ValueRange{
			{
				Range:  namedRange,
				Values: values,
			},
		},
	}
}

func (c *SpreadsheetClient) insertRowAfter(ctx context.Context, rowIndex int) error {
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

	_, err := c.service.Spreadsheets.BatchUpdate(c.spreadsheetID, &batchUpdateRequest).Context(ctx).Do()
	if err != nil {
		return err
	}

	return err
}

func EncodeCoordinate(columnNum, rowNum int) string {
	var (
		columnA1 = EncodeColumn(columnNum)
		coord    = fmt.Sprintf("%s%d", columnA1, rowNum+1)
	)

	return coord
}

// EncodeColumn encodes 0-indexed column to A1
func EncodeColumn(columnNum int) string {
	columnNum += 1
	result := ""
	for columnNum > 0 {
		remainder := (columnNum - 1) % 26
		result = fmt.Sprintf("%c", 65+remainder) + result
		columnNum = (columnNum - 1) / 26
	}
	return result
}

var regexpDigits = regexp.MustCompile(`\d+`)

func DecodeColumn(column string) int {
	var (
		columnWithoutDigits = regexpDigits.ReplaceAllString(column, "")
		result              = 0
	)

	for i := 0; i < len(columnWithoutDigits); i++ {
		result *= 26
		result += int(columnWithoutDigits[i]) - int('A') + 1
	}

	return result
}

func DecodeRow(row []interface{}) string {
	if len(row) == 0 {
		return ""
	}

	return strings.TrimSpace(row[0].(string))
}
