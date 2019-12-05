package task_service

import (
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"time"
)

// NextSearchPageTask 生成搜索下一页任务
func NextSearchPageTask(task *models.Task) *models.Task {
	return &models.Task{
		KeywordId:    task.KeywordId,
		Status:       logics.TASK_STATUS_未查询,
		SearchedPage: task.SearchedPage + 1,
		SearchCycle:  task.SearchCycle,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}
}