package domain

import (
	"container/heap"
	"sort"
)

type SimplePriorityQueue struct {
	Managers []*Manager
}

func (pq SimplePriorityQueue) Len() int { return len(pq.Managers) }

func (pq SimplePriorityQueue) Less(i, j int) bool {
	return pq.Managers[i].TotalRows < pq.Managers[j].TotalRows
}

func (pq SimplePriorityQueue) Swap(i, j int) {
	pq.Managers[i], pq.Managers[j] = pq.Managers[j], pq.Managers[i]
	pq.Managers[i].index = i
	pq.Managers[j].index = j
}

func (pq *SimplePriorityQueue) Push(x interface{}) {
	n := len(pq.Managers)
	item := x.(*Manager)
	item.index = n
	pq.Managers = append(pq.Managers, item)
}

func (pq *SimplePriorityQueue) Pop() interface{} {
	old := pq.Managers
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	pq.Managers = old[0 : n-1]
	return item
}

// DistributeSimple distributes sheets based on total rows of a sheet
func DistributeSimple(managersCount int, sheets []*Sheet) SimplePriorityQueue {
	pq := SimplePriorityQueue{Managers: make([]*Manager, managersCount)}

	for i := 0; i < managersCount; i++ {
		pq.Managers[i] = &Manager{
			TotalRows: 0,
			TotalSum:  0,
			index:     i,
		}
	}

	heap.Init(&pq)

	sort.Slice(sheets, func(i, j int) bool {
		return sheets[i].TotalRows > sheets[j].TotalRows
	})

	for _, sheet := range sheets {
		sheet := sheet
		manager := heap.Pop(&pq).(*Manager)
		manager.Sheets = append(manager.Sheets, sheet)
		manager.TotalRows += sheet.TotalRows
		heap.Push(&pq, manager)
	}

	return pq
}
