package task_service

import "rank-task/structs/models"

func UniqueKeyMappedCaptureTasks(finishedCaptureTasks []*models.CaptureTask) map[string][]*models.CaptureTask {
	uniqueKeyCaptureTasksMap := make(map[string][]*models.CaptureTask)
	for i, _ := range finishedCaptureTasks {
		uniqueKeyCaptureTasksMap[finishedCaptureTasks[i].UniqueKey] = append(uniqueKeyCaptureTasksMap[finishedCaptureTasks[i].UniqueKey], finishedCaptureTasks[i])
	}

	return uniqueKeyCaptureTasksMap
}
