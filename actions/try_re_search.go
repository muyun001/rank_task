package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rank-task/services/keyword_service"
)

func TryReSearch(c *gin.Context) {
	c.JSON(http.StatusOK, keyword_service.TryReSearch())
}
