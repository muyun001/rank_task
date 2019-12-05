package searched_rank_scope

import (
	"github.com/jinzhu/gorm"
)

func KeywordId(keywordId int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("keyword_id = ?", keywordId)
	}
}

func TopRankLT(rank int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("top_rank < ?", rank)
	}
}
