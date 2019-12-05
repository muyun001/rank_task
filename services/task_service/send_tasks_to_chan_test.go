package task_service_test

import (
	"rank-task/services/task_service"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"testing"
)

func TestSendTasksToChan(t *testing.T) {
	var tasks []*models.Task
	task1 := &models.Task{KeywordId: 1}
	task2 := &models.Task{KeywordId: 2}
	tasks = append(tasks, task1, task2)
	channel := make(chan *models.Task, logics.TASK_发送下载缓冲区大小)

	task_service.SendTasksToChan(tasks, channel, func(task *models.Task) {
		task.Status = logics.TASK_STATUS_查询中
	})

	verifyTask1 := <-channel
	verifyTask2 := <-channel

	if verifyTask1.KeywordId != verifyTask1.KeywordId {
		t.Error("keywordId not match")
	}
	if verifyTask2.KeywordId != verifyTask2.KeywordId {
		t.Error("keywordId not match")
	}
	if verifyTask1.Status != logics.TASK_STATUS_查询中 {
		t.Error("status not changed")
	}
	if verifyTask2.Status != logics.TASK_STATUS_查询中 {
		t.Error("status not changed")
	}
}
