package actions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/common/stringse"
	"rank-task/databases"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"strings"
	"time"
)

// GroupKeywordsPut: 分组接收关键词
func GroupKeywordsPut(c *gin.Context) {
	checkMatch := c.Param("check-match")
	engine := c.Param("engine")
	words := make([]string, 0)
	err := c.BindJSON(&words)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求格式不正确",
		})
		return
	}

	err = databases.Db.Model(models.Keyword{}).
		Where("priority = ? and check_match = ? and engine = ? and word in (?)", logics.KEYWORD_PRIORITY_不查, checkMatch, engine, words).
		Updates(models.Keyword{Priority: logics.KEYWORD_PRIORITY_中, ProcessedAt: time.Now()}).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "更新优先级错误",
		})
		return
	}

	existsWords := make([]string, 0)
	databases.Db.Model(models.Keyword{}).Where("word in (?)", words).Select("word").Scan(&existsWords)
	nonExistsWords := stringse.Diff(words, existsWords)
	insertColumns := "(word, engine, check_match, need_capture, priority, processed_at)"
	insertTempl := fmt.Sprintf("(?,\"%s\",\"%s\",%d,%d,\"%s\")", engine, checkMatch, 0, 2, time.Now().Format("2006-01-02 13:04:05"))
	insertTempls := make([]string, 0)
	insertBinds := make([]interface{}, 0)
	for _, word := range nonExistsWords {
		insertTempls = append(insertTempls, insertTempl)
		insertBinds = append(insertBinds, word)
	}
	err = databases.Db.Exec(fmt.Sprintf("INSERT IGNORE INTO keywords%s VALUES %s", insertColumns, strings.Join(insertTempls, ",")), insertBinds...).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "批量插入错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "数据保存成功",
	})
	return
}
