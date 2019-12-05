package task_service

import (
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"time"
)

// NextSearchCycleCaptureTask 生成搜索下一循环截图任务
func NextSearchCycleCaptureTask(task *models.CaptureTask) *models.CaptureTask {
	return &models.CaptureTask{
		KeywordId:    task.KeywordId,
		Status:       logics.TASK_STATUS_未查询,
		SearchedPage: 1,
		SearchCycle:  task.SearchCycle + 1,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}
}