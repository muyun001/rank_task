package keyword_service

import (
	"rank-task/databases"
	"rank-task/structs/models/logics"
)

// GetSomeUnTaskedKeywordIds: 按序获取一些未加入任务的关键词ID
func GetSomeUnTaskedKeywordIds() []int {
	var keywordIds []int

	keywordIds = append(keywordIds, getKeywordsIds(logics.KEYWORD_PRIORITY_高)...)
	if len(keywordIds) == logics.TASK_单次放入任务数量限制 {
		return keywordIds
	}

	keywordIds = append(keywordIds, getKeywordsIds(logics.KEYWORD_PRIORITY_中)...)
	if len(keywordIds) >= logics.TASK_单次放入任务数量限制 {
		return keywordIds[0:logics.TASK_单次放入任务数量限制]
	}

	keywordIds = append(keywordIds, getKeywordsIds(logics.KEYWORD_PRIORITY_低)...)
	if len(keywordIds) >= logics.TASK_单次放入任务数量限制 {
		return keywordIds[0:logics.TASK_单次放入任务数量限制]
	}

	return keywordIds
}

// getKeywordsIds 按照优先级和searchedCycle获取未加入任务的关键词ID
func getKeywordsIds(priority int) []int {
	var keywordIds []int
	databases.Db.Table("keywords").
		Where("priority = ?", priority).
		Where("searched_cycle < search_cycle").
		Order("searched_at").
		Limit(logics.TASK_单次放入任务数量限制).
		Pluck("id", &keywordIds)

	return keywordIds
}
