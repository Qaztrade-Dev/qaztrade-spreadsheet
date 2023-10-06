package domain

import "testing"

func TestDistributeAdvancedWork(t *testing.T) {
	sheets := []*Sheet{
		{0, "", "some title 1", 1, 100, 10000},
		{0, "", "some title 2", 2, 30, 5000},
	}

	managersCount := 2
	pq := DistributeAdvanced(managersCount, sheets)

	// Print out manager assignments
	for _, manager := range pq.Managers {
		println("Total rows: ", manager.TotalRows)
		println("Total sum: ", manager.TotalSum)
	}
}
