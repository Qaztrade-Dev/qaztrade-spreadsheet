package adapters

type HeadersMap map[string]*Header

type Range struct {
	Left  int
	Right int
}

type Bounds struct {
	Top    int
	Bottom int
}

func (b *Bounds) Equals(a *Bounds) bool {
	return b.Top == a.Top && b.Bottom == a.Bottom
}

type Header struct {
	Key      string
	Range    Range
	GroupKey string
	Values   HeadersMap
}

func NewHeader(key, groupKey string, rangeL, rangeR int, hcellMap HeadersMap) *Header {
	return &Header{
		Key:      key,
		GroupKey: groupKey,
		Range: Range{
			Left:  rangeL,
			Right: rangeR,
		},
		Values: hcellMap,
	}
}

func (h *Header) IsLeaf() bool {
	if h.Values == nil {
		return true
	}
	return len(h.Values) == 0
}

type Cell struct {
	Value      string
	RowNum     int
	ColumnNum  int
	HeaderCell *Header
}

func NewCell(value string, rowNum, columnNum int, headerCell *Header) *Cell {
	return &Cell{
		Value:      value,
		RowNum:     rowNum,
		ColumnNum:  columnNum,
		HeaderCell: headerCell,
	}
}
