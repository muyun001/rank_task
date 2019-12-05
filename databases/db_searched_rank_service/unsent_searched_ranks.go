package db_searched_rank_service

import (
	"rank-task/databases"
	"rank-task/structs/models"
)

// UnsentSearchedRanks 获取没有发送到存档的SearchedRank记录
func UnsentSearchedRanks() []models.SearchedRank {
	searchedRanks := make([]models.SearchedRank, 0)
	databases.Db.
		Preload("Keyword").
		Where("searched_ranks.is_send = ?", 0).
		Limit(1000).
		Find(&searchedRanks)

	return searchedRanks
}
