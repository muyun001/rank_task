package task_service

import (
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"time"
)

// NextSearchPageCaptureTask 生成搜索下一页截图任务
func NextSearchPageCaptureTask(task *models.CaptureTask) *models.CaptureTask {
	return &models.CaptureTask{
		KeywordId:    task.KeywordId,
		Status:       logics.TASK_STATUS_未查询,
		SearchedPage: task.SearchedPage + 1,
		SearchCycle:  task.SearchCycle,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}
}