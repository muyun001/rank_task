package global

import (
	"rank-task/databases"
	"rank-task/databases/scopes/task_scope"
	"rank-task/structs/models"
)

var BeforeQueriedTasksCount int64

func init() {
	ReadBeforeQueriedCount()
}

func ReadBeforeQueriedCount() {
	databases.Db.Model(&models.Task{}).Scopes(task_scope.BeforeQueried).Count(&BeforeQueriedTasksCount)
}
