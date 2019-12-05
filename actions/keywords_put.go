package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/services/get_rank_api"
	"rank-task/structs/rank_util"
)

// KeywordsPut: 接收关键词
func KeywordsPut(c *gin.Context) {
	putKeywords := &[]rank_util.PutKeywords{}
	err := c.BindJSON(putKeywords)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求格式不正确",
		})
		return
	}

	err = get_rank_api.ReceiveAndSaveKeywords(putKeywords)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"mgs": "数据保存失败",
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "数据保存成功",
	})
}
