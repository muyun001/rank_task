package actions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/databases"
	"rank-task/services/request_out/download_center"
	"rank-task/services/request_out/rank_util_api"
	"rank-task/settings"
	"rank-task/structs/models"
	"rank-task/structs/rank_util"
	"strconv"
)

func RankTestGet(c *gin.Context) {
	checkMatch := c.Param("check-match")
	engine := c.Param("engine")
	word := c.Param("word")
	cycle, err := strconv.Atoi(c.Param("cycle"))
	checkPageCount := settings.CheckRank / 10
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "cycle必须为数字",
		})
		return
	}

	keyword := models.Keyword{}
	if databases.Db.Where(&models.Keyword{Word: word, Engine: engine, CheckMatch: checkMatch}).First(&keyword).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "未找到关键词记录",
		})
		return
	}

	for i := 0; i < checkPageCount; i++ {
		dcRequest, err := rank_util_api.NewRankUtilApi().PostRequestBuilder(rank_util.MapEngine[engine], keyword.Word, 1, keyword.NeedCapture, cycle, download_center.DC_PRIORITY_中)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"msg": "模拟构建请求数据失败",
			})
			return
		}

		fmt.Println(dcRequest.UniqueKey)
		dcResponse, err := download_center.NewDownloadCenter().GetResponse(dcRequest.UniqueKey)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"msg": "获取搜索页面数据失败",
			})
			return
		}

		if dcResponse.Body == "" {
			c.JSON(http.StatusOK, gin.H{
				"msg": "暂无数据",
			})
			return
		}

		siteName := models.SiteName{}
		siteNameStr := ""
		if databases.Db.Where(models.SiteName{SiteDomain: checkMatch}).First(&siteName).RecordNotFound() == false {
			siteNameStr = siteName.SiteName
		}
		parseRankResponse, err := rank_util_api.NewRankUtilApi().PostRankExtractor(checkMatch, siteNameStr, dcResponse.Body, engine, 0)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"msg":  "解析排名失败",
				"html": dcResponse.Body,
			})
			return
		}

		if len(parseRankResponse.Ranks) > 0 {
			c.JSON(http.StatusOK, parseRankResponse.Ranks)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("未查到%d以内排名", settings.CheckRank),
	})
}
