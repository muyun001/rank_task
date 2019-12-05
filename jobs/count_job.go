package jobs

import (
	"rank-task/global"
	"time"
)

func SyncBeforeQueriedCount() {
	global.ReadBeforeQueriedCount()
	time.Sleep(time.Minute * 10)
}
