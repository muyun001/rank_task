package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/databases"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"strings"
	"time"
)

func DomainsTodaySearchedCountGet(c *gin.Context) {
	engines := c.Query("engines")

	type CheckMatchKeywordsTodaySearchedCount struct {
		CheckMatch string `json:"check_match"`
		Count      int    `json:"count"`
	}

	checkMatchKeywordsTodaySearchedCounts := make([]CheckMatchKeywordsTodaySearchedCount, 0)
	query := databases.Db.Model(models.Keyword{}).
		Where("priority in (?)", []int{logics.KEYWORD_PRIORITY_低, logics.KEYWORD_PRIORITY_中, logics.KEYWORD_PRIORITY_高}).
		Where("searched_at > ?", time.Now().Format("2006-01-02 00:00:00"))
	if engines != "" {
		query = query.Where("engine in (?)", strings.Split(engines, ","))
	}
	query.Select("check_match, count(id) as count").
		Group("check_match").
		Scan(&checkMatchKeywordsTodaySearchedCounts)

	checkMatchKeywordsTodaySearchedCountMap := make(map[string]int)
	for i := range checkMatchKeywordsTodaySearchedCounts {
		checkMatchKeywordsTodaySearchedCountMap[checkMatchKeywordsTodaySearchedCounts[i].CheckMatch] = checkMatchKeywordsTodaySearchedCounts[i].Count
	}

	c.JSON(http.StatusOK, checkMatchKeywordsTodaySearchedCountMap)
}
