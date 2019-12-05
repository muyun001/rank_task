package task_service_test

import (
	"rank-task/services/task_service"
	"testing"
)

func TestHasNextSearchPage(t *testing.T) {
	type checkValue struct {
		SearchedPage int
		CheckRank    int
		Result       bool
	}

	checkValues := []checkValue{
		{1, 50, true},
		{2, 50, true},
		{3, 50, true},
		{4, 50, true},
		{5, 50, false},
		{6, 50, false},
	}

	for _, c := range checkValues {
		if task_service.HasNextSearchPage(c.SearchedPage, c.CheckRank) != c.Result {
			t.Errorf("page %d next page to check rank %d should be %t", c.SearchedPage, c.CheckRank, c.Result)
		}
	}
}
