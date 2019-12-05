package main

import (
	"rank-task/databases"
	"rank-task/jobs"
	"time"
)

func foreverGo(run func(), routineLimits int) {
	for i := 0; i < routineLimits; i++ {
		go func() {
			for {
				run()
			}
		}()
	}
}

func main() {
	databases.AutoMigrate()

	foreverGo(jobs.SendTasksToChan, 1)
	foreverGo(jobs.SendDcRequestsFromChan, 20)
	foreverGo(jobs.FetchQueryParsedResult, 1)
	foreverGo(jobs.SendCaptureTasksToChan, 1)
	foreverGo(jobs.SendCaptureDcRequestsFromChan, 20)
	foreverGo(jobs.FetchCaptureQueryParsedResult, 1)
	foreverGo(jobs.ReSearch, 1)
	foreverGo(jobs.AddKeywordsToTasks, 1)
	foreverGo(jobs.SyncBeforeQueriedCount, 1)
	foreverGo(jobs.RankArchive, 1)

	for {
		time.Sleep(time.Minute)
	}
}
