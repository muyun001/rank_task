package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/jobs"
)

// 重新检查SiteName
func ReCheckSiteName(c *gin.Context) {
	go func() {
		jobs.ReCheckSiteName()
	}()

	c.JSON(http.StatusOK, gin.H{
		"msg": "Re-Check Site Name triggered",
	})
}
