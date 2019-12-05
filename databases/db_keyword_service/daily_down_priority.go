package db_keyword_service

import (
	"github.com/jinzhu/gorm"
	"rank-task/databases"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"time"
)

// DownToNormalPriority 优先级降级到中
func DownToNormalPriority() *gorm.DB {
	return databases.Db.Model(&models.Keyword{}).
		Where("priority = ? ", logics.KEYWORD_PRIORITY_高).
		Where("no_rank_days >= ?", logics.KEYWORD_PRIORITY_调至中级连续无排名天数).
		Update(map[string]interface{}{
			"priority": logics.KEYWORD_PRIORITY_中,
		})
}

// DownToLowPriority 优先级降级到低
func DownToLowPriority() *gorm.DB {
	return databases.Db.Model(&models.Keyword{}).
		Where("priority = ?", logics.KEYWORD_PRIORITY_中).
		Where("no_rank_days >= ?", logics.KEYWORD_PRIORITY_调至低级连续无排名天数).
		Update(map[string]interface{}{
			"priority": logics.KEYWORD_PRIORITY_低,
		})
}

// DownToNoSearchPriority: 优先级降到不查
func DownToNoSearchPriority() *gorm.DB {
	return databases.Db.Model(&models.Keyword{}).
		Where("processed_at <= ?", time.Now().Add(-time.Hour*24*2)).
		Update(map[string]interface{}{
			"priority": logics.KEYWORD_PRIORITY_不查,
		})
}
