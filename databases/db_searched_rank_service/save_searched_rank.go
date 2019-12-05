package db_searched_rank_service

import (
	"rank-task/common/ints"
	"rank-task/databases"
	"rank-task/structs/models"
)

// SaveSearchedRank 保存排名到SearchedRank
func SaveSearchedRank(keywordId, topRank int, ranks []int, captureUrl string, ip string) {
	ranksJoined := ints.Join(ranks, ",")
	searchedRank := models.SearchedRank{
		Ip:         ip,
		KeywordId:  keywordId,
		TopRank:    topRank,
		Ranks:      ranksJoined,
		CaptureUrl: captureUrl,
	}
	databases.Db.Save(&searchedRank)
}
