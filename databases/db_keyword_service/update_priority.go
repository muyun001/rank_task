package db_keyword_service

import (
	"github.com/jinzhu/gorm"
	"rank-task/databases"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"time"
)

// UpdatePriority 更新单个Keyword优先级
func UpdatePriority(keyword *models.Keyword, priority int) *gorm.DB {
	return databases.Db.Model(keyword).Update(map[string]interface{}{
		"priority": priority,
	})
}

// GroupWordsResetPriority 分组重置Keywords优先级
// 由不查重置为中
func GroupWordsResetPriority(checkMatch, engine string, words []string) *gorm.DB {
	return databases.Db.Model(models.Keyword{}).
		Where("priority = ? and check_match = ? and engine = ? and word in (?)", logics.KEYWORD_PRIORITY_不查, checkMatch, engine, words).
		Updates(models.Keyword{Priority: logics.KEYWORD_PRIORITY_中, ProcessedAt: time.Now()})
}
