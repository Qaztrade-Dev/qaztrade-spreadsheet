package domain

import "testing"

func TestDistributeSimpleWork(t *testing.T) {
	sheets := []*Sheet{
		{"", "some title 1", 1, 100, 10000},
		{"", "some title 2", 2, 30, 5000},
	}

	managersCount := 2
	pq := DistributeSimple(managersCount, sheets)

	// Print out manager assignments
	for _, manager := range pq.Managers {
		println("Total rows: ", manager.TotalRows)
	}
}
