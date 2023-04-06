package adapters

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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

	service, err := sheets.NewService(
		ctx,
		option.WithCredentialsJSON(credentialsJson),
	)
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
		{Range: "product_capacity", Value: application.ProductCapacity},
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
		{Range: "info_manufactured_goods", Value: application.InfoManufacturedGoods},
		{Range: "name_of_goods", Value: application.NameOfGoods},
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
			"Затраты на продвижение":  1156025711,
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
	sheetParents  domain.Relations
	sheetChildren domain.Relations
	offset        int
	headersMap    HeadersMap
	data          [][]string
}

func (c *SpreadsheetClient) NewSheetClient(ctx context.Context, spreadsheetID, sheetName string, sheetID int64) (*SheetClient, error) {
	sheetClient := &SheetClient{
		service:       c.service,
		spreadsheetID: spreadsheetID,
		sheetName:     sheetName,
		sheetID:       sheetID,
		sheetParents:  domain.SheetParents[sheetName],
		sheetChildren: domain.SheetChildren[sheetName],
	}

	var (
		headerRangeName = getSheetRangeHeader(sheetName)
		dataRangeName   = getSheetRangeData(sheetName)
	)

	t1 := time.Now()
	batchDataValues, err := c.getDataFromRanges(ctx, spreadsheetID, []string{headerRangeName, dataRangeName})
	fmt.Println("time: getDataFromRanges", time.Since(t1))
	if err != nil {
		return nil, err
	}

	var (
		headerValues = batchDataValues[0]
		dataValues   = batchDataValues[1]
	)

	t1 = time.Now()
	headersMap, err := sheetClient.getHeadersMap(ctx, sheetName, headerValues)
	fmt.Println("time: getHeaderCells", time.Since(t1))
	if err != nil {
		return nil, err
	}

	t1 = time.Now()
	data, err := sheetClient.getData(ctx, sheetName, dataValues)
	fmt.Println("time: getData", time.Since(t1))
	if err != nil {
		return nil, err
	}

	sheetClient.headersMap = headersMap
	sheetClient.data = data
	sheetClient.offset = 2

	return sheetClient, nil
}

func (c *SpreadsheetClient) getDataFromRanges(ctx context.Context, spreadsheetID string, ranges []string) ([][][]interface{}, error) {
	resp, err := c.service.Spreadsheets.Values.BatchGet(spreadsheetID).Ranges(ranges...).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	datas := make([][][]interface{}, len(resp.ValueRanges))
	for i := range resp.ValueRanges {
		datas[i] = resp.ValueRanges[i].Values
	}
	return datas, nil
}

func getSheetRangeData(sheetName string) string {
	rangeName := fmt.Sprintf("'%s'!%s_%s", sheetName, strings.ReplaceAll(sheetName, " ", "_"), "data")

	return rangeName
}

func (c *SheetClient) getData(ctx context.Context, sheetName string, values [][]interface{}) ([][]string, error) {
	data := make([][]string, len(values))
	for i, row := range values {
		data[i] = make([]string, len(row))
		for j := range row {
			data[i][j] = strings.TrimSpace(row[j].(string))
		}
	}
	// for i := range data {
	// 	for j := range data[i] {
	// 		fmt.Printf("%v, ", data[i][j])
	// 	}
	// 	fmt.Println()
	// }
	return data, nil
}

func getSheetRangeHeader(sheetName string) string {
	rangeName := fmt.Sprintf("'%s'!%s_%s", sheetName, strings.ReplaceAll(sheetName, " ", "_"), "header")

	return rangeName
}

func (c *SheetClient) getHeadersMap(ctx context.Context, sheetName string, values [][]interface{}) (HeadersMap, error) {
	var (
		topLevel   = values[0]
		lowLevel   = values[1]
		headersMap = make(HeadersMap)
	)

	for i := range topLevel {
		if topLevel[i] == "" {
			continue
		}

		var (
			topLevelValue   = strings.TrimSpace(topLevel[i].(string))
			innerHeadersMap = make(HeadersMap)
			rangeR          = i
		)

		for j := i; j < len(lowLevel); j++ {
			lowLevelValue := strings.TrimSpace(lowLevel[j].(string))

			if lowLevelValue == "" {
				break
			}
			if !(topLevel[j] == "" || i == j) {
				break
			}

			innerHeadersMap[lowLevelValue] = NewHeader(lowLevelValue, topLevelValue, j, j, nil)
			rangeR = j
		}

		headersMap[topLevelValue] = NewHeader(topLevelValue, topLevelValue, i, rangeR, innerHeadersMap)
	}

	// for k := range hcellMap {
	// 	fmt.Println(k)
	// 	fmt.Println(hcellMap[k].Values)
	// 	for v := range hcellMap[k].Values {
	// 		fmt.Println("\t", v)
	// 		fmt.Println("\t", hcellMap[k].Values[v].Values)
	// 	}
	// }

	return headersMap, nil
}

type UpdateCellRequest struct {
	RowIndex    int64
	ColumnIndex int64
	Value       string
}

func (c *SheetClient) fillRecord(
	ctx context.Context,
	payload domain.PayloadValue,
	headers HeadersMap,
	rowNum int,
	batchUpdate *BatchUpdate,
) error {
	batch := make([]*UpdateCellRequest, 0)

	for key := range payload {
		header := headers[key]
		if header.IsLeaf() {
			batch = append(batch, &UpdateCellRequest{
				RowIndex:    int64(rowNum),
				ColumnIndex: int64(header.Range.Left),
				Value:       payload[key].(string),
			})
		} else {
			var p interface{} = payload[key]
			var m map[string]interface{} = p.(map[string]interface{})
			var d domain.PayloadValue = domain.PayloadValue(m)

			if err := c.fillRecord(ctx, d, header.Values, rowNum, batchUpdate); err != nil {
				return err
			}
		}
	}

	if len(batch) == 0 {
		return nil
	}

	for i := range batch {
		batchUpdate.WithRequest(batch[i].encode(c.sheetID))
	}

	return nil
}

func (r *UpdateCellRequest) encode(sheetID int64) *sheets.Request {
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

func (c *SheetClient) getBounds(nodeKey, nodeID string, parentBound ...*Bounds) *Bounds {
	fmt.Printf("getBounds. nodeKey=%v, nodeID=%v\n", nodeKey, nodeID)
	if nodeKey == ParentKeyRoot {
		return &Bounds{Top: 0, Bottom: len(c.data) - 1}
	}

	fmt.Printf("nodeKey: %#v\n", nodeKey)
	var (
		parentHeader = c.getHeader(nodeKey)
		columnIdx    = parentHeader.Range.Left
		fromRow      = 0
		toRow        = len(c.data) - 1

		left  = 0
		right = 0
	)

	fmt.Printf("parentHeader: %+v\n", parentHeader)

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
		value := strings.Join(c.data[i][:columnIdx+1], "")
		fmt.Printf("right: i=%v, c.data[i][j]=%v\n", i, c.data[i][columnIdx])
		if value != "" {
			break
		}
		right = i
	}

	result := &Bounds{Top: left, Bottom: right}
	fmt.Printf("result %#v\n", result)

	return result
}

func (c *SheetClient) getHeader(parentKey string) *Header {
	fmt.Println("getHeader:", parentKey)
	keys := strings.Split(parentKey, "|")
	var header *Header
	// fmt.Println(keys)
	for _, key := range keys {
		if header == nil {
			header = c.headersMap[key]
		} else {
			header = header.Values[key]
		}
		// fmt.Printf("key=%#v,cell=%#v\n", key, cell)
	}

	return header
}

var (
	ErrorChildNotFound = errors.New("child not found")
	ErrorChildNotMatch = errors.New("child name doesn't match")
)

func (c *SheetClient) getChildNode(parentKey, childName string) (*domain.Node, error) {
	fmt.Printf("GetChildNode. parentKey=%v, childName=%v\n", parentKey, childName)

	if _, ok := c.sheetChildren[parentKey]; !ok {
		return nil, ErrorChildNotFound
	}

	child := c.sheetChildren[parentKey]
	if child.Name != childName {
		return nil, ErrorChildNotMatch
	}

	return c.sheetChildren[parentKey], nil
}

func (c *SheetClient) getLastChild(parentBounds *Bounds, childHeader *Header) *Cell {
	fmt.Printf("GetLastChildCell. parentBounds=%#v, childHeaderCell=%#v\n", parentBounds, childHeader)
	var (
		fromRow   = parentBounds.Top
		toRow     = parentBounds.Bottom
		columnIdx = childHeader.Range.Left
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
			continue
		}
		value := c.data[i][columnIdx]
		if i == 0 && value == "" {
			return nil
		}
		if value != "" {
			lastIdx = i
			lastValue = value
		}
	}

	if lastValue == "" {
		return nil
	}

	return NewCell(lastValue, fromRow+lastIdx, columnIdx, childHeader)
}

// 1. get parent row
// 2. get last child of the parent, e.g. neighbor
// 3. get last row of the farthest descendent
func (c *SheetClient) getRowNum(parentValue, childName string, bounds ...*Bounds) (int, bool, error) {
	fmt.Printf("GetRowNum. parentID:%v, childName:%v\n", parentValue, childName)
	var (
		parentKey  = c.sheetParents[childName].Key
		parentName = c.sheetParents[childName].Name

		parentBounds = c.getBounds(parentKey, parentValue, bounds...)
		upperBound   = parentBounds.Top
	)

	child, err := c.getChildNode(parentName, childName)
	fmt.Printf("child: %#v\n", child)
	if err != nil {
		return 0, false, err
	}

	fmt.Printf("c.getHeaderCell(%v)\n", child.Key)
	childHeaderCell := c.getHeader(child.Key)

	lastChildCell := c.getLastChild(parentBounds, childHeaderCell)
	fmt.Printf("lastChildCell: %#v\n", lastChildCell)

	if lastChildCell == nil {
		return upperBound, false, nil
	}

	grandChildNode, ok := c.sheetChildren[lastChildCell.HeaderCell.GroupKey]
	if !ok {
		return lastChildCell.RowNum, true, nil
	}

	rowNum, _, err := c.getRowNum(lastChildCell.Value, grandChildNode.Name, parentBounds)
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
		batchUpdate   = NewBatchUpdate(c.service)
	)

	if rowNum == 0 {
		t1 := time.Now()
		rowNum, mustInsertRow, err = c.getRowNum(payload.ParentID, payload.ChildKey)
		fmt.Println("time: getRowNum", time.Since(t1))
		if err != nil {
			return err
		}

	}

	fmt.Printf("SubmitRow. rowNum=%v mustInsertRow=%v\n", rowNum, mustInsertRow)

	if mustInsertRow {
		t1 := time.Now()
		c.insertRowAfter(ctx, c.offset+rowNum, batchUpdate)
		fmt.Println("time: insertRowAfter", time.Since(t1))
		rowNum += 1
	}

	t1 := time.Now()
	if err = c.fillRecord(ctx, payload.Value, c.headersMap, c.offset+rowNum, batchUpdate); err != nil {
		return err
	}
	fmt.Println("time: fillRecord", time.Since(t1))

	t1 = time.Now()
	if err := batchUpdate.Do(ctx, c.spreadsheetID); err != nil {
		return err
	}
	fmt.Println("time: batchUpdate.Do", time.Since(t1))

	return nil
}

func (c *SheetClient) insertRowAfter(ctx context.Context, rowIndex int, batchUpdate *BatchUpdate) {
	batchUpdate.WithRequest(&sheets.Request{
		InsertDimension: &sheets.InsertDimensionRequest{
			Range: &sheets.DimensionRange{
				SheetId:    c.sheetID,
				Dimension:  "ROWS",
				StartIndex: int64(rowIndex) + 1, // Start index is 1-based
				EndIndex:   int64(rowIndex) + 2, // End index is exclusive
			},
			InheritFromBefore: true,
		},
	})
}

func (c *SheetClient) getParentValue(childName, childValue string) string {
	if childName == ParentKeyRoot {
		return ""
	}

	var (
		parent = c.sheetParents[childName]
		child  = c.sheetParents[parent.Name]
	)

	fmt.Printf("parent: %#v\n", parent)

	var (
		childBounds = c.getBounds(child.Key, childValue)
		bottomIdx   = childBounds.Top
	)

	fmt.Printf("childBounds: %#v\n", childBounds)

	var (
		parentHeader = c.getHeader(parent.Key)
	)

	fmt.Printf("parentHeader: %#v\n", parentHeader)

	var (
		columnIdx = parentHeader.Range.Left
	)

	for i := bottomIdx; i >= 0; i-- {
		if columnIdx >= len(c.data[i]) {
			continue
		}

		value := c.data[i][columnIdx]
		if value != "" {
			return value
		}
	}

	return ""
}

func (c *SheetClient) RemoveParent(ctx context.Context, input *domain.RemoveInput) error {
	var (
		parent = c.sheetParents[input.Name]
		child  = c.sheetChildren[parent.Name]

		childHeader = c.getHeader(child.Key)
		childBounds = c.getBounds(child.Key, input.Value)

		parentValue          = ""
		parentBounds *Bounds = nil

		batchUpdate = NewBatchUpdate(c.service)
	)

	var (
		childStartRowIndex    = childBounds.Top
		childEndRowIndex      = childBounds.Bottom
		childStartColumnIndex = childHeader.Range.Left

		parentStartRowIndex  = -1
		parentEndRowIndex    = -1
		parentEndColumnIndex = -1
	)

	if parent.Name != ParentKeyRoot {
		parentValue = c.getParentValue(child.Name, input.Value)
		parentBounds = c.getBounds(parent.Key, parentValue)

		parentStartRowIndex = parentBounds.Top + 1
		parentEndRowIndex = parentBounds.Top + 1 + (childBounds.Bottom - childBounds.Top)
		parentEndColumnIndex = childHeader.Range.Left - 1
	}

	batchUpdate.WithRequest(
		&sheets.Request{
			DeleteRange: &sheets.DeleteRangeRequest{
				Range: &sheets.GridRange{
					SheetId:          c.sheetID,
					StartRowIndex:    int64(c.offset + childStartRowIndex),
					EndRowIndex:      int64(c.offset + childEndRowIndex + 1),
					StartColumnIndex: int64(childStartColumnIndex),
				},
				ShiftDimension: "ROWS",
			},
		},
	)

	if parent.Name != ParentKeyRoot {
		batchUpdate.WithRequest(
			&sheets.Request{
				DeleteRange: &sheets.DeleteRangeRequest{
					Range: &sheets.GridRange{
						SheetId:        c.sheetID,
						StartRowIndex:  int64(c.offset + parentStartRowIndex),
						EndRowIndex:    int64(c.offset + parentEndRowIndex + 1),
						EndColumnIndex: int64(parentEndColumnIndex + 1),
					},
					ShiftDimension: "ROWS",
				},
			},
		)
	}

	if err := batchUpdate.Do(ctx, c.spreadsheetID); err != nil {
		return err
	}

	return nil
}
