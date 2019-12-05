package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/services/get_rank_api"
)

// RanksGetConfirmed: 确认排名已被成功获取
func RanksGetConfirmed(c *gin.Context) {
	requestHash := c.Param("request-hash")

	err := get_rank_api.ConfirmRanks(requestHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "确认排名失败",
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "确认排名成功",
	})
}
