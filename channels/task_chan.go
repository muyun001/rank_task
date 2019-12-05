package channels

import (
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
)

var TaskSendingChan chan *models.Task

func init() {
	TaskSendingChan = make(chan *models.Task, logics.TASK_发送下载缓冲区大小)
}
