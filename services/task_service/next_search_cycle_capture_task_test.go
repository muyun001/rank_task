package task_service_test

import (
	"rank-task/services/task_service"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"testing"
	"time"
)

func TestNextSearchCycleCaptureTask(t *testing.T) {
	captureTask := &models.CaptureTask{
		ID:           1,
		KeywordId:    1,
		Status:       logics.TASK_STATUS_查询不达标,
		UniqueKey:    "abcdefgh",
		SearchedPage: 2,
		SearchCycle:  2,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}

	nextPageCaptureTask := task_service.NextSearchCycleCaptureTask(captureTask)
	if nextPageCaptureTask.ID != 0 {
		t.Errorf("expect ID 0, got %d", nextPageCaptureTask.ID)
	}
	if nextPageCaptureTask.Status != logics.TASK_STATUS_未查询 {
		t.Errorf("expect Status %d, got %d", logics.TASK_STATUS_未查询, nextPageCaptureTask.Status)
	}
	if nextPageCaptureTask.UniqueKey != "" {
		t.Errorf("expect UniqueKey empty, got %s", nextPageCaptureTask.UniqueKey)
	}
	if nextPageCaptureTask.SearchedPage != 1 {
		t.Errorf("expect SearchedPage 1, got %d", nextPageCaptureTask.SearchedPage)
	}
	if nextPageCaptureTask.SearchCycle != 3 {
		t.Errorf("expect SearchCycle 3, got %d", nextPageCaptureTask.SearchCycle)
	}
}
