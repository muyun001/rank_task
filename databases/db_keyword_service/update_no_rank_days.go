package db_keyword_service

import (
	"github.com/jinzhu/gorm"
	"rank-task/databases"
	"rank-task/structs/models"
)

// UpdateNoRankDays 更新单个Keyword无排名天数
func UpdateNoRankDays(keyword *models.Keyword, days int) *gorm.DB {
	return databases.Db.Model(keyword).Updates(map[string]interface{}{
		"no_rank_days": days,
	})
}
