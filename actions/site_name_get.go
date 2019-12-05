package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/databases"
	"rank-task/structs/models"
)

func SiteNameGet(c *gin.Context) {
	checkMatch := c.Param("check-match")

	siteName := models.SiteName{}
	databases.Db.Where("site_domain = ?", checkMatch).First(&siteName)

	c.JSON(http.StatusOK, siteName.SiteName)
}
