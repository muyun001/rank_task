package routers

import (
	"github.com/gin-gonic/gin"
	"rank-task/actions"
)

var r *gin.Engine

func Load() *gin.Engine {
	r.POST("migrate", actions.MigratePost)
	r.PUT("keywords", actions.KeywordsPut)
	r.GET("keywords/:check-match/:engine", actions.GroupKeywordsGet)
	r.PUT("keywords/:check-match/:engine", actions.GroupKeywordsPut)
	r.PUT("captured-keywords/:check-match/:engine", actions.CapturedGroupKeywordsPut)
	r.POST("ranks/get/:check-match/:engine/:request-hash", actions.RanksGet)
	r.POST("captured-ranks/get/:check-match/:engine/:request-hash", actions.CapturedRanksGet)
	r.PUT("ranks/confirmed/:request-hash", actions.RanksGetConfirmed)
	r.GET("rank/test/:check-match/:engine/:word/:cycle", actions.RankTestGet)
	r.POST("ranks/back", actions.RanksBackPost)
	r.PUT("daily-reset", actions.DailyResetPut)
	r.PUT("try-re-search", actions.TryReSearch)
	r.GET("tasks/stat", actions.TasksStatGet)
	r.POST("re-check-site-name", actions.ReCheckSiteName)
	r.PUT("ranks/archive", actions.RankArchivePut)
	r.GET("domains/today-searched-count", actions.DomainsTodaySearchedCountGet)
	r.GET("site-name/:check-match", actions.SiteNameGet)

	return r
}

func init() {
	r = gin.Default()
}
