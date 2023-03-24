package adapters

type HeaderCellMap map[string]*HeaderCell

type Range struct {
	Left  int
	Right int
}

type Bound struct {
	Top    int
	Bottom int
}

func (b *Bound) EncodeRange(columnNum int) string {
	var (
		fromA1 = EncodeCoordinate(columnNum, b.Top)
		toA1   = EncodeCoordinate(columnNum, b.Bottom)
		range_ = fromA1 + ":" + toA1
	)

	return range_
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