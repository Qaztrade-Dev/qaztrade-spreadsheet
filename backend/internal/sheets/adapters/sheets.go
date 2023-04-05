package adapters

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetClient struct {
	service             *sheets.Service
	originSpreadsheetID string
}

var _ domain.SheetsRepository = (*SpreadsheetClient)(nil)

func NewSpreadsheetClient(ctx context.Context, credentialsJson []byte, originSpreadsheetID string) (*SpreadsheetClient, error) {
	service, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, err
	}

	return &SpreadsheetClient{
		service:             service,
		originSpreadsheetID: originSpreadsheetID,
	}, nil
}

func (c *SpreadsheetClient) InsertRecord(ctx context.Context, spreadsheetID, sheetName string, sheetID int64, payload *domain.Payload) error {
	sheetClient, err := c.NewSheetClient(ctx, spreadsheetID, sheetName, sheetID)
	if err != nil {
		return err
	}
	return sheetClient.InsertRecord(ctx, payload)
}

func (c *SpreadsheetClient) RemoveRecord(ctx context.Context, spreadsheetID, sheetName string, sheetID int64, input *domain.RemoveInput) error {
	sheetClient, err := c.NewSheetClient(ctx, spreadsheetID, sheetName, sheetID)
	if err != nil {
		return err
	}
	return sheetClient.RemoveParent(ctx, input)
}

func (c *SpreadsheetClient) UpdateApplication(ctx context.Context, spreadsheetID string, application *domain.Application) error {
	var mappings = []struct {
		Range string
		Value string
	}{
		{Range: "from", Value: application.From},
		{Range: "gov_reg", Value: application.GovReg},
		{Range: "fact_addr", Value: application.FactAddr},
		{Range: "bin", Value: application.Bin},
		{Range: "industry", Value: application.Industry},
		{Range: "activity", Value: application.Activity},
		{Range: "emp_count", Value: application.EmpCount},
		{Range: "manufacturer", Value: application.Manufacturer},
		{Range: "item", Value: application.Item},
		{Range: "item_volume", Value: application.ItemVolume},
		{Range: "fact_volume_earnings", Value: application.FactVolumeEarnings},
		{Range: "fact_workload", Value: application.FactWorkload},
		{Range: "chief_lastname", Value: application.ChiefLastname},
		{Range: "chief_firstname", Value: application.ChiefFirstname},
		{Range: "chief_middlename", Value: application.ChiefMiddlename},
		{Range: "chief_position", Value: application.ChiefPosition},
		{Range: "chief_phone", Value: application.ChiefPhone},
		{Range: "cont_lastname", Value: application.ContLastname},
		{Range: "cont_firstname", Value: application.ContFirstname},
		{Range: "cont_middlename", Value: application.ContMiddlename},
		{Range: "cont_position", Value: application.ContPosition},
		{Range: "cont_phone", Value: application.ContPhone},
		{Range: "cont_email", Value: application.ContEmail},
		{Range: "country", Value: application.Country},
		{Range: "code_tnved", Value: application.CodeTnved},
	}

	data := make([]*sheets.ValueRange, 0, len(mappings))
	for i := range mappings {
		data = append(data, &sheets.ValueRange{
			Range:  mappings[i].Range,
			Values: [][]interface{}{{mappings[i].Value}},
		})
	}

	updateValuesRequest := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "RAW",
		Data:             data,
	}

	_, err := c.service.Spreadsheets.Values.BatchUpdate(spreadsheetID, updateValuesRequest).Do()
	if err != nil {
		return err
	}

	return nil
}

func (c *SpreadsheetClient) AddSheet(ctx context.Context, spreadsheetID string, sheetName string) error {
	var (
		mappings = map[string]int64{ // sheetName:sheetID
			"Доставка ЖД транспортом": 0,
		}
		sourceSheetID = mappings[sheetName]
	)

	containsSheet, err := c.containsSheet(ctx, spreadsheetID, sheetName)
	if err != nil {
		return err
	}

	if containsSheet {
		return domain.ErrorSheetPresent
	}

	sheetID, err := c.copySheet(ctx, c.originSpreadsheetID, spreadsheetID, sourceSheetID)
	if err != nil {
		return err
	}

	protectedRanges, err := c.getProtectedRanges(ctx, c.originSpreadsheetID, sourceSheetID)
	if err != nil {
		return err
	}

	batchUpdate := NewBatchUpdate(c.service)
	{
		batchUpdate.WithProtectedRange(sheetID, protectedRanges)
		batchUpdate.WithSheetName(sheetID, sheetName)
	}

	if err := batchUpdate.Do(ctx, spreadsheetID); err != nil {
		return err
	}

	return nil
}

func (c *SpreadsheetClient) containsSheet(ctx context.Context, spreadsheetID string, sheetName string) (bool, error) {
	spreadsheet, err := c.service.Spreadsheets.Get(spreadsheetID).IncludeGridData(false).Context(ctx).Do()
	if err != nil {
		return false, err
	}

	// Iterate through the sheets and check if the sheet name exists
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			return true, nil
		}
	}

	return false, nil
}

// copySheet returns sheetId
func (c *SpreadsheetClient) copySheet(ctx context.Context, sourceSpreadsheetID, targetSpreadsheetID string, sheetID int64) (int64, error) {
	copyRequest := &sheets.CopySheetToAnotherSpreadsheetRequest{
		DestinationSpreadsheetId: targetSpreadsheetID,
	}

	resp, err := c.service.Spreadsheets.Sheets.CopyTo(sourceSpreadsheetID, sheetID, copyRequest).Context(ctx).Do()
	if err != nil {
		return 0, err
	}

	return resp.SheetId, nil
}

func (c *SpreadsheetClient) getProtectedRanges(ctx context.Context, spreadsheetID string, sheetID int64) ([]*sheets.ProtectedRange, error) {
	spreadsheet, err := c.service.Spreadsheets.Get(spreadsheetID).IncludeGridData(false).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	var protectedRanges []*sheets.ProtectedRange
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.SheetId == sheetID {
			protectedRanges = sheet.ProtectedRanges
			break
		}
	}

	return protectedRanges, nil
}

type SheetClient struct {
	service       *sheets.Service
	spreadsheetID string
	sheetName     string
	sheetID       int64
	offset        int
	headersMap    HeaderCellMap
	data          [][]string
}

func (c *SpreadsheetClient) NewSheetClient(ctx context.Context, spreadsheetID, sheetName string, sheetID int64) (*SheetClient, error) {
	sheetClient := &SheetClient{
		service:       c.service,
		spreadsheetID: spreadsheetID,
		sheetName:     sheetName,
		sheetID:       sheetID,
	}

	headersMap, err := sheetClient.getHeaderCells(ctx, sheetName)
	if err != nil {
		return nil, err
	}

	data, err := sheetClient.getData(ctx, sheetName)
	if err != nil {
		return nil, err
	}

	sheetClient.headersMap = headersMap
	sheetClient.data = data
	sheetClient.offset = 4

	return sheetClient, nil
}

func getSheetRangeData(sheetName string) string {
	rangeName := fmt.Sprintf("'%s'!%s_%s", sheetName, strings.ReplaceAll(sheetName, " ", "_"), "data")

	return rangeName
}

func (c *SheetClient) getData(ctx context.Context, sheetName string) ([][]string, error) {
	headerRangeName := getSheetRangeData(sheetName)
	sheetRange, err := c.service.Spreadsheets.Values.Get(c.spreadsheetID, headerRangeName).
		Context(ctx).
		Do()
	if err != nil {
		return nil, err
	}

	data := make([][]string, len(sheetRange.Values))
	for i, row := range sheetRange.Values {
		data[i] = make([]string, len(row))
		for j := range row {
			data[i][j] = strings.TrimSpace(row[j].(string))
		}
	}
	for i := range data {
		for j := range data[i] {
			fmt.Printf("%v, ", data[i][j])
		}
		fmt.Println()
	}
	return data, nil
}

func getSheetRangeHeader(sheetName string) string {
	rangeName := fmt.Sprintf("'%s'!%s_%s", sheetName, strings.ReplaceAll(sheetName, " ", "_"), "header")

	return rangeName
}

func (c *SheetClient) getHeaderCells(ctx context.Context, sheetName string) (HeaderCellMap, error) {
	headerRangeName := getSheetRangeHeader(sheetName)
	sheetRange, err := c.service.Spreadsheets.Values.Get(c.spreadsheetID, headerRangeName).
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

func (c *SheetClient) fillRecord(
	ctx context.Context,
	payload domain.PayloadValue,
	headers HeaderCellMap,
	rowNum int,
	batchUpdate *BatchUpdate,
) error {
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

			if err := c.fillRecord(ctx, d, cell.Values, rowNum, batchUpdate); err != nil {
				return err
			}
		}
	}

	if len(batch) == 0 {
		return nil
	}

	for i := range batch {
		batchUpdate.WithRequest(batch[i].Encode(c.sheetID))
	}

	return nil
}

func (r *UpdateCellRequest) Encode(sheetID int64) *sheets.Request {
	var (
		stringValue  *string
		formulaValue *string
	)

	if strings.HasPrefix(r.Value, "=") {
		formulaValue = &r.Value
	} else {
		stringValue = &r.Value
	}

	return &sheets.Request{
		UpdateCells: &sheets.UpdateCellsRequest{
			Fields: "*",
			Start: &sheets.GridCoordinate{
				RowIndex:    int64(r.RowIndex),
				ColumnIndex: int64(r.ColumnIndex),
				SheetId:     sheetID,
			},
			Rows: []*sheets.RowData{
				{
					Values: []*sheets.CellData{
						{
							UserEnteredValue: &sheets.ExtendedValue{
								StringValue:  stringValue,
								FormulaValue: formulaValue,
							},
						},
					},
				},
			},
		},
	}
}

const ParentKeyRoot = "root"

func (c *SheetClient) getNodeBounds(ctx context.Context, nodeKey, nodeID string, parentBound ...*Bound) (*Bound, error) {
	fmt.Printf("GetNodeBounds. nodeKey=%v, nodeID=%v\n", nodeKey, nodeID)
	if nodeKey == ParentKeyRoot {
		return &Bound{Top: 0, Bottom: len(c.data) - 1}, nil
	}

	fmt.Printf("nodeKey: %#v\n", nodeKey)
	var (
		parentHeaderCell = c.getHeaderCell(nodeKey)
		fromRow          = 0
		toRow            = len(c.data) - 1
		columnIdx        = parentHeaderCell.Range.Left

		left  = 0
		right = 0
	)

	fmt.Printf("parentHeaderCell: %+v\n", parentHeaderCell)

	if len(parentBound) > 0 {
		fromRow = parentBound[0].Top
		toRow = parentBound[0].Bottom
	}

	for i := fromRow; i <= toRow; i++ {
		fmt.Println(i)
		if i >= len(c.data) {
			fmt.Println("i >= len(c.data): break")
			break
		}
		if columnIdx >= len(c.data[i]) {
			fmt.Println("columnIdx >= len(c.data[i]): break")
			continue
		}
		value := c.data[i][columnIdx]
		fmt.Printf("left: i=%v, c.data[i][j]=%v\n", i, c.data[i][columnIdx])
		if value == nodeID {
			left = i
			right = i
			break
		}
	}

	for i := left + 1; i <= toRow; i++ {
		if i >= len(c.data) {
			break
		}
		if columnIdx >= len(c.data[i]) {
			continue
		}
		value := c.data[i][columnIdx]
		fmt.Printf("right: i=%v, c.data[i][j]=%v\n", i, c.data[i][columnIdx])
		if value != "" {
			break
		}
		right = i
	}

	result := &Bound{Top: left, Bottom: right}

	return result, nil
}

func (c *SheetClient) getHeaderCell(parentKey string) *HeaderCell {
	keys := strings.Split(parentKey, ".")
	var cell *HeaderCell
	fmt.Println(keys)
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

func (c *SheetClient) getChildNode(parentKey, childName string) (*domain.Node, error) {
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

func (c *SheetClient) getLastChildCell(ctx context.Context, parentBounds *Bound, childHeaderCell *HeaderCell) (*Cell, error) {
	fmt.Printf("GetLastChildCell. parentBounds=%#v, childHeaderCell=%#v\n", parentBounds, childHeaderCell)
	var (
		fromRow   = parentBounds.Top
		toRow     = parentBounds.Bottom
		columnIdx = childHeaderCell.Range.Left
	)

	var (
		lastIdx   = 0
		lastValue = ""
	)

	for i := fromRow; i <= toRow; i++ {
		if i >= len(c.data) {
			break
		}
		if columnIdx >= len(c.data[i]) {
			break
		}
		value := c.data[i][columnIdx]
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

	return NewCell(lastValue, fromRow+lastIdx, columnIdx, childHeaderCell), nil
}

// 1. get parent row
// 2. get last child of the parent, e.g. neighbor
// 3. get last row of the farthest descendent
func (c *SheetClient) getRowNum(ctx context.Context, parentID, childName string, bounds ...*Bound) (int, bool, error) {
	fmt.Printf("GetRowNum. parentID:%v, childName:%v\n", parentID, childName)
	var (
		parentKey  = domain.Parents[childName].Key
		parentName = domain.Parents[childName].Name
	)

	parentBounds, err := c.getNodeBounds(ctx, parentKey, parentID, bounds...)
	fmt.Printf("parentBounds: %#v\n", parentBounds)
	if err != nil {
		return 0, false, err
	}

	var upperBound = parentBounds.Top

	child, err := c.getChildNode(parentName, childName)
	fmt.Printf("child: %#v\n", child)
	if err != nil {
		return 0, false, err
	}

	childHeaderCell := c.getHeaderCell(child.Key)

	lastChildCell, err := c.getLastChildCell(ctx, parentBounds, childHeaderCell)
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

	rowNum, _, err := c.getRowNum(ctx, lastChildCell.Value, grandChildNode.Name, parentBounds)
	if err != nil {
		return 0, false, err
	}

	return rowNum, true, nil
}

func (c *SheetClient) InsertRecord(ctx context.Context, payload *domain.Payload) error {
	var (
		err           error
		rowNum        = payload.RowNumber
		mustInsertRow = false
	)

	if rowNum == 0 {
		rowNum, mustInsertRow, err = c.getRowNum(ctx, payload.ParentID, payload.ChildKey)
		if err != nil {
			return err
		}
	}

	fmt.Printf("SubmitRow. rowNum=%v mustInsertRow=%v\n", rowNum, mustInsertRow)

	if mustInsertRow {
		if err := c.insertRowAfter(ctx, c.offset+rowNum); err != nil {
			return err
		}
		rowNum += 1
	}

	batchUpdate := NewBatchUpdate(c.service)
	if err = c.fillRecord(ctx, payload.Value, c.headersMap, c.offset+rowNum, batchUpdate); err != nil {
		return err
	}

	if err := batchUpdate.Do(ctx, c.spreadsheetID); err != nil {
		return err
	}

	return nil
}

func (c *SheetClient) insertRowAfter(ctx context.Context, rowIndex int) error {
	request := &sheets.Request{
		InsertDimension: &sheets.InsertDimensionRequest{
			Range: &sheets.DimensionRange{
				SheetId:    c.sheetID,
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

func (c *SheetClient) RemoveParent(ctx context.Context, input *domain.RemoveInput) error {
	var (
		parent     = domain.Parents[input.Name]
		child      = domain.Children[parent.Name]
		headerCell = c.getHeaderCell(child.Key)
	)

	bound, err := c.getNodeBounds(ctx, child.Key, input.Value)
	if err != nil {
		return err
	}

	batchUpdate := NewBatchUpdate(c.service)

	batchUpdate.WithRequest(&sheets.Request{
		DeleteRange: &sheets.DeleteRangeRequest{
			Range: &sheets.GridRange{
				SheetId:          c.sheetID,
				StartRowIndex:    int64(c.offset + bound.Top),
				EndRowIndex:      int64(c.offset + bound.Bottom + 1),
				StartColumnIndex: int64(headerCell.Range.Left),
			},
			ShiftDimension: "ROWS",
		},
	})

	if err := batchUpdate.Do(ctx, c.spreadsheetID); err != nil {
		return err
	}

	return nil
}
