package db_keyword_service

import (
	"github.com/jinzhu/gorm"
	"rank-task/databases"
	"rank-task/structs/models"
)

// SetHasNewRank 设置单个Keyword的has_new_rank
func SetHasNewRank(keyword *models.Keyword, topRank int, captureUrl string) *gorm.DB {
	return databases.Db.Model(keyword).Updates(map[string]interface{}{
		"has_new_rank": true,
		"top_rank":     topRank,
		"capture_url":  captureUrl,
		"no_rank_days": 0,
	})
}
