package domain

import "testing"

func TestDistributeWork(t *testing.T) {
	sheets := []*Sheet{
		{"", "some title 1", 1, 100, 10000},
		{"", "some title 2", 2, 30, 5000},
	}

	managersCount := 2
	pq := DistributeWork(managersCount, sheets)

	// Print out manager assignments
	for _, manager := range pq.Managers {
		println("Manager ID: ", manager.ID)
		println("Total rows: ", manager.TotalRows)
		println("Total sum: ", manager.TotalSum)
	}
}
