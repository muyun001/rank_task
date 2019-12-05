package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/databases"
	"rank-task/structs/models"
)

// GroupKeywordsGet 分组查看关键词
func GroupKeywordsGet(c *gin.Context) {
	checkMatch := c.Param("check-match")
	engine := c.Param("engine")
	keywords := make([]models.Keyword, 0)

	databases.Db.Model(models.Keyword{}).
		Where("check_match = ? and engine = ?", checkMatch, engine).
		Find(&keywords)

	words := make([]string, 0)
	for i := range keywords {
		words = append(words, keywords[i].Word)
	}

	siteName := models.SiteName{}
	databases.Db.Where("site_domain = ?", checkMatch).First(&siteName)

	c.JSON(http.StatusOK, gin.H{
		"site_name": siteName.SiteName,
		"words":     words,
	})
}
