package jobs

import (
	"rank-task/global"
	"rank-task/services/keyword_service"
	"time"
)

// ReSearch: 根据searchCycleLimit重查
func ReSearch() {
	if global.IsBetweenSearchTime() {
		time.Sleep(keyword_service.DurationToNextReSearch())
		keyword_service.TryReSearch()
	} else {
		time.Sleep(time.Minute)
	}
}
