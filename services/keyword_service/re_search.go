package keyword_service

import (
	"github.com/jinzhu/gorm"
	"rank-task/databases"
	"rank-task/settings"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"time"
)

type TryReSearchResult struct {
	High   bool
	Normal bool
	Low    bool
}

// reSearchInterval 计算出每次重查间隔
// 按重查次数划分间隔
func reSearchInterval() time.Duration {
	return settings.SearchEndTime.Sub(settings.SearchStartTime) / time.Duration(settings.SearchCycleLimit)
}

// NextReSearchTime 下一个重查时间点
func NextReSearchTime() time.Time {
	reSearchInterval := reSearchInterval()
	now := time.Now()
	nextReSearchTime := settings.SearchStartTime
	for {
		nextReSearchTime = nextReSearchTime.Add(reSearchInterval)
		if now.Sub(nextReSearchTime) < 0 {
			break
		}
	}

	return nextReSearchTime
}

// DurationToNextReSearch 距离下次重查还有多久
func DurationToNextReSearch() time.Duration {
	return NextReSearchTime().Sub(time.Now())
}

// TryReSearch 尝试重查
func TryReSearch() TryReSearchResult {
	tryReSearchResult := TryReSearchResult{}
	if isReadyForNextCycle(logics.KEYWORD_PRIORITY_高) {
		nextSearchCycle(logics.KEYWORD_PRIORITY_高)
		tryReSearchResult.High = true
	}

	return tryReSearchResult
}

// isReadyForNextCycle: 判断是否可以重查
func isReadyForNextCycle(priority int) bool {
	var PageTaskCount int
	databases.Db.Table("keywords").Joins("left join tasks on keywords.id = tasks.keyword_id").Where("keywords.priority=?", priority).Where("tasks.status=? or tasks.status=?", logics.TASK_STATUS_未查询, logics.TASK_STATUS_查询中).Count(&PageTaskCount)

	if PageTaskCount < 100 {
		return true
	}

	return false
}

// nextSearchCycle: 更新重查次数
func nextSearchCycle(priority int) {
	databases.Db.Model(&models.Keyword{}).
		Where("priority = ?", priority).
		Where("search_cycle < ?", settings.SearchCycleLimit).
		Where("top_rank = ? or top_rank > ?", 0, settings.ReachRank).
		Update("search_cycle", gorm.Expr("search_cycle + ?", 1))
}
