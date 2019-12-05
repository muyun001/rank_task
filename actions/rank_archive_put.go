package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/services/request_out/rank_archive_api"
)

func RankArchivePut(c *gin.Context) {
	rank_archive_api.SendRanks()

	c.JSON(http.StatusOK, gin.H{
		"msg": "archived",
	})
}
