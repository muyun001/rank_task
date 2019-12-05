package task_service

import (
	"rank-task/structs/models"
)

func SendTasksToChan(tasks []*models.Task, taskChan chan *models.Task, preFunc func(task *models.Task)) {
	for _, task := range tasks {
		preFunc(task)
		taskChan <- task
	}
}
