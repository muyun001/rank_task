package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/common/ints"
	"rank-task/databases"
	"rank-task/databases/db_keyword_service"
	"rank-task/databases/db_searched_rank_service"
	"rank-task/services"
	"rank-task/settings"
	"rank-task/structs/models"
)

type RanksBack struct {
	Ip         string `json:"ip"`
	Word       string `json:"word"`
	Engine     string `json:"engine"`
	CheckMatch string `json:"check_match"`
	Ranks      []int  `json:"ranks"`
	Capture    string `json:"capture"`
}

// RanksBackPost 回传排名数据
// {
//  "ip": "ip地址",
//  "word": "词",
//  "engine": "baidu_pc",
//  "check_match": "www.fxt.cn",
//  "ranks": [1,3,5],
//  "capture": "base64编码",
// }
func RanksBackPost(c *gin.Context) {
	ranksBack := RanksBack{}
	err := c.Bind(&ranksBack)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "返回格式不正确",
		})
		return
	}

	if len(ranksBack.Ranks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "没有排名数据",
		})
		return
	}

	keyword := models.Keyword{
		Word:       ranksBack.Word,
		Engine:     ranksBack.Engine,
		CheckMatch: ranksBack.CheckMatch,
	}
	if databases.Db.Where(&keyword).First(&keyword).RecordNotFound() {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "关键词不存在",
		})
		return
	}

	topRank := ints.Min(ranksBack.Ranks...)
	isRankReached := topRank > 0 && topRank <= settings.ReachRank
	captureUrl := ""

	if isRankReached {
		if keyword.NeedCapture {
			cosUrl, err := services.UploadCaptureToCos(ranksBack.Capture, keyword.ID, keyword.Engine, ranksBack.Ip)
			if err != nil {
				captureUrl = "data:image/png;base64," + ranksBack.Capture
			} else {
				captureUrl = cosUrl
			}
		}

		db_keyword_service.SetHasNewRank(&keyword, topRank, captureUrl)
		db_searched_rank_service.SaveSearchedRank(keyword.ID, topRank, ranksBack.Ranks, captureUrl, ranksBack.Ip)

		c.JSON(http.StatusOK, gin.H{
			"msg": "排名已保存",
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "排名未达标",
		})
		return
	}
}
