package task_service_test

import (
	"rank-task/services/task_service"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"testing"
	"time"
)

func TestNextSearchPageTask(t *testing.T) {
	task := &models.Task{
		ID:           1,
		KeywordId:    1,
		Status:       logics.TASK_STATUS_查询不达标,
		UniqueKey:    "abcdefgh",
		SearchedPage: 2,
		SearchCycle:  1,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}

	nextPageTask := task_service.NextSearchPageTask(task)
	if nextPageTask.ID != 0 {
		t.Errorf("expect ID 0, got %d", nextPageTask.ID)
	}
	if nextPageTask.Status != logics.TASK_STATUS_未查询 {
		t.Errorf("expect Status %d, got %d", logics.TASK_STATUS_未查询, nextPageTask.Status)
	}
	if nextPageTask.UniqueKey != "" {
		t.Errorf("expect UniqueKey empty, got %s", nextPageTask.UniqueKey)
	}
	if nextPageTask.SearchedPage != 3 {
		t.Errorf("expect SearchedPage 3, got %d", nextPageTask.SearchedPage)
	}
	if nextPageTask.SearchCycle != 1 {
		t.Errorf("expect SearchCycle 1, got %d", nextPageTask.SearchCycle)
	}
}
