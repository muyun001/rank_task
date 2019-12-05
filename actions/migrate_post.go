package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/databases"
)

func MigratePost(c *gin.Context) {
	databases.AutoMigrate()

	c.JSON(http.StatusOK, gin.H{
		"msg": "migrated",
	})
}
