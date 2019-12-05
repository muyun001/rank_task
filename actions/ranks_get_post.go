package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/services/get_rank_api"
)

// RanksGet: 获取排名
func RanksGet(c *gin.Context) {
	checkMatch := c.Param("check-match")
	engine := c.Param("engine")
	requestHash := c.Param("request-hash")

	keywords := make([]string, 0)
	err := c.BindJSON(&keywords)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求格式不正确",
		})
		return
	}

	rankResultsResponse, err := get_rank_api.RankResultsResponse(checkMatch, engine, requestHash, keywords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "排名获取失败",
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rankResultsResponse)
}
