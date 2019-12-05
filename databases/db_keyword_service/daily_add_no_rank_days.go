package db_keyword_service

import (
	"github.com/jinzhu/gorm"
	"rank-task/databases"
	"rank-task/structs/models"
)

// AddNoRankDays 每日执行增加NoRankDays
// 需要在DailyResetKeywords和DownPriority之前操作
func AddNoRankDays() {
	databases.Db.Model(&models.Keyword{}).
		Where("rank = 0").
		Update("no_rank_days", gorm.Expr("no_rank_days + ?", 1))
}
