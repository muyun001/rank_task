package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/databases"
	"rank-task/services/keyword_service"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
)

type TaskCounts struct {
	UnQueried int `json:"un_queried"`
	Querying  int `json:"querying"`
	Reached   int `json:"reached"`
	UnReached int `json:"un_reached"`
	Failed    int `json:"failed"`
}

type CapturedTaskCounts TaskCounts

type StatResponse struct {
	TaskCounts            TaskCounts         `json:"task"`
	CapturedTaskCounts    CapturedTaskCounts `json:"captured_task"`
	NextReSearchInSeconds float64            `json:"next_re_search_in_seconds"`
}

func TasksStatGet(c *gin.Context) {
	taskCounts := TaskCounts{}
	capturedTaskCounts := CapturedTaskCounts{}
	type StatusCount struct {
		Status int
		Count  int
	}
	taskStatusCounts := make([]StatusCount, 0)
	databases.Db.Model(&models.Task{}).Select("status, count(status) as count").Group("status").Scan(&taskStatusCounts)
	for _, statusCount := range taskStatusCounts {
		switch statusCount.Status {
		case logics.TASK_STATUS_未查询:
			taskCounts.UnQueried = statusCount.Count
		case logics.TASK_STATUS_查询中:
			taskCounts.Querying = statusCount.Count
		case logics.TASK_STATUS_查询达标:
			taskCounts.Reached = statusCount.Count
		case logics.TASK_STATUS_查询不达标:
			taskCounts.UnReached = statusCount.Count
		case logics.TASK_STATUS_查询失败:
			taskCounts.Failed = statusCount.Count
		}
	}

	capturedTaskStatusCounts := make([]StatusCount, 0)
	databases.Db.Model(&models.CaptureTask{}).Select("status, count(status) as count").Group("status").Scan(&capturedTaskStatusCounts)
	for _, statusCount := range capturedTaskStatusCounts {
		switch statusCount.Status {
		case logics.TASK_STATUS_未查询:
			capturedTaskCounts.UnQueried = statusCount.Count
		case logics.TASK_STATUS_查询中:
			capturedTaskCounts.Querying = statusCount.Count
		case logics.TASK_STATUS_查询达标:
			capturedTaskCounts.Reached = statusCount.Count
		case logics.TASK_STATUS_查询不达标:
			capturedTaskCounts.UnReached = statusCount.Count
		case logics.TASK_STATUS_查询失败:
			capturedTaskCounts.Failed = statusCount.Count
		}
	}

	c.JSON(http.StatusOK, StatResponse{
		TaskCounts:            taskCounts,
		CapturedTaskCounts:    capturedTaskCounts,
		NextReSearchInSeconds: keyword_service.DurationToNextReSearch().Seconds(),
	})
}
