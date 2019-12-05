package db_keyword_service

import (
	"github.com/jinzhu/gorm"
	"rank-task/databases"
	"rank-task/structs/models"
)

// DailyResetKeywords 每日重置所有关键词
func DailyResetKeywords() *gorm.DB {
	return databases.Db.Model(&models.Keyword{}).Updates(map[string]interface{}{
		"search_cycle":   1,
		"searched_cycle": 0,
		"has_new_rank":   false,
		"top_rank":       0,
	})
}
