package channels

import (
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
)

var CaptureTaskSendingChan chan *models.CaptureTask

func init() {
	CaptureTaskSendingChan = make(chan *models.CaptureTask, logics.TASK_发送下载缓冲区大小)
}
