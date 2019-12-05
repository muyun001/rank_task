package task_service

import "rank-task/structs/models"

func UniqueKeyMappedTasks(finishedTasks []*models.Task) map[string][]*models.Task {
	uniqueKeyTasksMap := make(map[string][]*models.Task)
	for i, _ := range finishedTasks {
		uniqueKeyTasksMap[finishedTasks[i].UniqueKey] = append(uniqueKeyTasksMap[finishedTasks[i].UniqueKey], finishedTasks[i])
	}

	return uniqueKeyTasksMap
}
