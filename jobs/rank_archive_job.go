package jobs

import (
	"rank-task/services/request_out/rank_archive_api"
	"time"
)

func RankArchive() {
	rank_archive_api.SendRanks()
	time.Sleep(time.Second)
}
