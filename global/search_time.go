package global

import (
	"rank-task/settings"
	"time"
)

func IsBetweenSearchTime() bool {
	now := time.Now()
	if now.Sub(settings.SearchEndTime) > 0 {
		settings.SearchStartTime = settings.SearchStartTime.Add(time.Hour * 24)
		settings.SearchEndTime = settings.SearchEndTime.Add(time.Hour * 24)
	}
	return now.Sub(settings.SearchStartTime) > 0 && now.Sub(settings.SearchEndTime) < 0
}
