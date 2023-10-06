package domain

import (
	"container/heap"
	"math"
	"sort"
)

type AdvancedPriorityQueue struct {
	Managers []*Manager
	meanRows float64
	meanSum  float64
	stdRows  float64
	stdSum   float64
}

func (pq AdvancedPriorityQueue) GetManagers() []*Manager {
	return pq.Managers
}

func (pq AdvancedPriorityQueue) Len() int { return len(pq.Managers) }

func (pq AdvancedPriorityQueue) Less(i, j int) bool {
	zScoreRowsI := (float64(pq.Managers[i].TotalRows) - pq.meanRows) / pq.stdRows
	zScoreSumI := (float64(pq.Managers[i].TotalSum) - pq.meanSum) / pq.stdSum
	zScoreI := zScoreRowsI + zScoreSumI

	zScoreRowsJ := (float64(pq.Managers[j].TotalRows) - pq.meanRows) / pq.stdRows
	zScoreSumJ := (float64(pq.Managers[j].TotalSum) - pq.meanSum) / pq.stdSum
	zScoreJ := zScoreRowsJ + zScoreSumJ

	return zScoreI < zScoreJ
}

func (pq AdvancedPriorityQueue) Swap(i, j int) {
	pq.Managers[i], pq.Managers[j] = pq.Managers[j], pq.Managers[i]
	pq.Managers[i].index = i
	pq.Managers[j].index = j
}

func (pq *AdvancedPriorityQueue) Push(x interface{}) {
	n := len(pq.Managers)
	item := x.(*Manager)
	item.index = n
	pq.Managers = append(pq.Managers, item)
}

func (pq *AdvancedPriorityQueue) Pop() interface{} {
	old := pq.Managers
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	pq.Managers = old[0 : n-1]
	return item
}

// DistributeAdvanced distributes sheets based on total rows and total sum of a sheet
func DistributeAdvancedWithData(managers []*Manager, sheets []*Sheet) AdvancedPriorityQueue {
	for i := 0; i < len(managers); i++ {
		managers[i].index = i
	}

	pq := AdvancedPriorityQueue{Managers: managers}
	heap.Init(&pq)

	meanRows, meanSum, varRows, varSum := calculateStats(sheets)
	pq.meanRows = meanRows
	pq.meanSum = meanSum
	pq.stdRows = math.Sqrt(varRows)
	pq.stdSum = math.Sqrt(varSum)

	sort.Slice(sheets, func(i, j int) bool {
		zScoreRowsI := (float64(sheets[i].TotalRows) - meanRows) / pq.stdRows
		zScoreSumI := (float64(sheets[i].TotalSum) - meanSum) / pq.stdSum
		zScoreI := zScoreRowsI + zScoreSumI

		zScoreRowsJ := (float64(sheets[j].TotalRows) - meanRows) / pq.stdRows
		zScoreSumJ := (float64(sheets[j].TotalSum) - meanSum) / pq.stdSum
		zScoreJ := zScoreRowsJ + zScoreSumJ

		return zScoreI > zScoreJ
	})

	for _, sheet := range sheets {
		sheet := sheet
		manager := heap.Pop(&pq).(*Manager)
		manager.Sheets = append(manager.Sheets, sheet)
		manager.TotalRows += sheet.TotalRows
		manager.TotalSum += sheet.TotalSum
		heap.Push(&pq, manager)
	}

	return pq
}

// DistributeAdvanced distributes sheets based on total rows and total sum of a sheet
func DistributeAdvanced(managersCount int, sheets []*Sheet) AdvancedPriorityQueue {
	return DistributeAdvancedWithData(make([]*Manager, managersCount), sheets)
}

func calculateStats(sheets []*Sheet) (meanRows, meanSum, varRows, varSum float64) {
	for _, sheet := range sheets {
		meanRows += float64(sheet.TotalRows)
		meanSum += float64(sheet.TotalSum)
	}

	meanRows /= float64(len(sheets))
	meanSum /= float64(len(sheets))

	for _, sheet := range sheets {
		varRows += math.Pow(float64(sheet.TotalRows)-meanRows, 2)
		varSum += math.Pow(float64(sheet.TotalSum)-meanSum, 2)
	}
	varRows /= float64(len(sheets) - 1)
	varSum /= float64(len(sheets) - 1)

	return meanRows, meanSum, varRows, varSum
}
