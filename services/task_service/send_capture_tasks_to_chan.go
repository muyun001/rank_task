package task_service

import (
	"rank-task/structs/models"
)

func SendCaptureTasksToChan(captureTasks []*models.CaptureTask, captureTaskChan chan *models.CaptureTask, preFunc func(task *models.CaptureTask)) {
	for _, captureTask := range captureTasks {
		preFunc(captureTask)
		captureTaskChan <- captureTask
	}
}
