package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/databases"
	"rank-task/databases/db_keyword_service"
	"rank-task/databases/scopes/task_scope"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"time"
)

func DailyResetPut(c *gin.Context) {
	go func() {
		db_keyword_service.AddNoRankDays()
		db_keyword_service.DownToNormalPriority()
		db_keyword_service.DownToLowPriority()
		db_keyword_service.DownToNoSearchPriority()
		db_keyword_service.DailyResetKeywords()

		databases.Db.Exec("TRUNCATE TABLE searched_ranks")

		var beforeQueriedTasks []*models.Task
		databases.Db.
			Scopes(task_scope.Querying).
			Order("updated_at").
			Find(&beforeQueriedTasks)

		databases.Db.Exec("TRUNCATE TABLE tasks")

		var beforeQueriedTaskIds []int
		for _, beforeQueriedTask := range beforeQueriedTasks {
			beforeQueriedTask.ID = 0
			databases.Db.Save(&beforeQueriedTask)
			beforeQueriedTaskIds = append(beforeQueriedTaskIds, beforeQueriedTask.ID)
		}

		var beforeQueriedCaptureTasks []*models.CaptureTask
		databases.Db.
			Scopes(task_scope.Querying).
			Order("updated_at").
			Find(&beforeQueriedCaptureTasks)

		databases.Db.Exec("TRUNCATE TABLE capture_tasks")

		var beforeQueriedCaptureTaskIds []int
		for _, beforeQueriedCaptureTask := range beforeQueriedCaptureTasks {
			beforeQueriedCaptureTask.ID = 0
			databases.Db.Save(&beforeQueriedCaptureTask)
			beforeQueriedCaptureTaskIds = append(beforeQueriedCaptureTaskIds, beforeQueriedCaptureTask.ID)
		}

		go func(beforeQueriedTaskIds []int, beforeQueriedCaptureTaskIds []int) {
			time.Sleep(time.Minute * 10)
			databases.Db.Where("status = ?", logics.TASK_STATUS_查询中).Where("id in (?)", beforeQueriedTaskIds).Delete(models.Task{})
			databases.Db.Where("status = ?", logics.TASK_STATUS_查询中).Where("id in (?)", beforeQueriedCaptureTaskIds).Delete(models.CaptureTask{})
		}(beforeQueriedTaskIds, beforeQueriedCaptureTaskIds)
	}()

	c.JSON(http.StatusOK, gin.H{
		"msg": "已重置",
	})
}
