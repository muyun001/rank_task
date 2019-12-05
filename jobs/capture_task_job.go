package jobs

import (
	"rank-task/channels"
	"rank-task/common/ints"
	"rank-task/databases"
	"rank-task/databases/db_keyword_service"
	"rank-task/databases/db_searched_rank_service"
	"rank-task/databases/scopes/task_scope"
	"rank-task/global"
	"rank-task/services"
	"rank-task/services/request_out/download_center"
	"rank-task/services/task_service"
	"rank-task/settings"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"time"
)

// SendCaptureTasksToChan: 从Chan发送截图任务
func SendCaptureTasksToChan() {
	var tasks []*models.CaptureTask
	limit := cap(channels.CaptureTaskSendingChan) - len(channels.CaptureTaskSendingChan)
	if limit == 0 {
		time.Sleep(time.Second)
		return
	}
	databases.Db.Model(&models.CaptureTask{}).Scopes(task_scope.UnQueried).Limit(limit).Scan(&tasks)

	if len(tasks) != 0 {
		task_service.SendCaptureTasksToChan(tasks, channels.CaptureTaskSendingChan, func(task *models.CaptureTask) {
			task.Status = logics.TASK_STATUS_查询中
			databases.Db.Save(&task)
		})
	} else {
		time.Sleep(time.Second * 5)
	}
}

// SendCaptureDcRequestsFromChan: 将截图任务发送到下载中心
func SendCaptureDcRequestsFromChan() {
	task := <-channels.CaptureTaskSendingChan
	parserService := &services.ParserService{}
	dcRequest, err := parserService.BuildCaptureDcRequest(task)
	if err != nil {
		time.Sleep(time.Second * 5)
		databases.Db.Model(&task).Update(models.Task{Status: logics.TASK_STATUS_未查询})
		return
	}

	dc := download_center.NewDownloadCenter()
	err = dc.PutRequest(dcRequest)
	if err != nil {
		time.Sleep(time.Second * 5)
		databases.Db.Model(&task).Update(models.Task{Status: logics.TASK_STATUS_未查询})
		return
	}

	task.UniqueKey = dcRequest.UniqueKey
	databases.Db.Save(&task)
}

type UniqueKeyCaptureTaskGroup struct {
	UniqueKey    string
	CaptureTasks []*models.CaptureTask
}

// FetchCaptureQueryParsedResult: 获取查询和解析结果
func FetchCaptureQueryParsedResult() {
	var beforeQueriedTasks []*models.CaptureTask
	databases.Db.
		Preload("Keyword").
		Scopes(task_scope.Querying).
		Order("updated_at").
		Find(&beforeQueriedTasks)

	if len(beforeQueriedTasks) == 0 {
		time.Sleep(time.Second * 2)
		return
	}

	fetchCaptureTasksResults(beforeQueriedTasks)
	time.Sleep(time.Second)
}

func fetchCaptureTasksResults(captureTasks []*models.CaptureTask) {
	var uniqueKeys []string
	uniqueKeyIndexedTasks := make(map[string]*models.CaptureTask)
	for _, captureTask := range captureTasks {
		uniqueKeys = append(uniqueKeys, captureTask.UniqueKey)
		uniqueKeyIndexedTasks[captureTask.UniqueKey] = captureTask
	}

	dc := download_center.NewDownloadCenter()
	finishedUniqueKeys, err := dc.PostResponsesCheck(uniqueKeys)
	if err != nil {
		time.Sleep(time.Second * 5)
		return
	}
	if len(finishedUniqueKeys) == 0 {
		time.Sleep(time.Second * 30)
		return
	}

	var finishedTasks []*models.CaptureTask
	databases.Db.
		Preload("Keyword").
		Preload("Keyword.SiteName").
		Scopes(task_scope.UniqueKeysIn(finishedUniqueKeys)).
		Find(&finishedTasks)

	uniqueKeyCaptureTasksMap := task_service.UniqueKeyMappedCaptureTasks(finishedTasks)

	for _, uniqueKey := range finishedUniqueKeys {
		uniqueKeyCaptureTaskGroup := UniqueKeyCaptureTaskGroup{
			UniqueKey:    uniqueKey,
			CaptureTasks: uniqueKeyCaptureTasksMap[uniqueKey],
		}
		tryFinishCaptureTask(uniqueKeyCaptureTaskGroup)
	}
}

func tryFinishCaptureTask(group UniqueKeyCaptureTaskGroup) {
	dc := download_center.NewDownloadCenter()
	dcResponse, err := dc.GetResponse(group.UniqueKey)
	if err != nil {
		databases.Db.Model(&models.Task{}).Where(models.Task{UniqueKey: group.UniqueKey}).Updates(models.Task{Status: logics.TASK_STATUS_查询失败})
		return
	}
	if dcResponse.Body == "" {
		_ = download_center.NewDownloadCenter().ResetRequest(group.UniqueKey)
		return
	}

	for _, finishedTask := range group.CaptureTasks {
		parserService := services.ParserService{}
		ranks, err := parserService.ParseRanks(finishedTask.Keyword.CheckMatch, finishedTask.Keyword.SiteName.SiteName, dcResponse.Body, finishedTask.Keyword.Engine, finishedTask.SearchedPage)
		if err != nil {
			databases.Db.Model(&finishedTask).Updates(models.CaptureTask{Status: logics.TASK_STATUS_查询失败})
			continue
		}
		topRank := ints.Min(ranks...)
		isRankReached := topRank > 0 && topRank <= settings.ReachRank
		captureUrl := ""

		if isRankReached && finishedTask.Keyword.NeedCapture {
			cosUrl, err := services.UploadCaptureToCos(dcResponse.Capture, finishedTask.Keyword.ID, finishedTask.Keyword.Engine, dcResponse.Ip)
			if err != nil {
				captureUrl = "data:image/png;base64," + dcResponse.Capture
			} else {
				captureUrl = cosUrl
			}
		}

		db_searched_rank_service.SaveSearchedRank(finishedTask.KeywordId, topRank, ranks, captureUrl, dcResponse.Ip)

		if isRankReached {
			databases.Db.Model(&finishedTask).Updates(models.CaptureTask{Status: logics.TASK_STATUS_查询达标})
		} else {
			databases.Db.Model(&finishedTask).Updates(models.CaptureTask{Status: logics.TASK_STATUS_查询不达标})
		}

		if topRank > 0 {
			// 更新"has_new_rank","top_rank","no_rank_days"的三种情况:
			// 1.processed_at是昨天或者更早;
			// 2.processed_at是今天,top_rank是0;
			// 3.processed_at是今天,但今天的历史排名相比现在的排名要靠后(todayHistoryTopRank > topRank)
			isProcessedBeforeToday := services.DaysApartToday(finishedTask.Keyword.ProcessedAt) > 0
			hasNewRankToday := services.DaysApartToday(finishedTask.Keyword.ProcessedAt) == 0 && (finishedTask.Keyword.TopRank == 0 || finishedTask.Keyword.TopRank > topRank)
			if isProcessedBeforeToday || hasNewRankToday {
				db_keyword_service.SetHasNewRank(&finishedTask.Keyword, topRank, captureUrl)
			}
			db_keyword_service.UpdatePriority(&finishedTask.Keyword, logics.KEYWORD_PRIORITY_高)
		} else {
			days := services.DaysApartToday(finishedTask.Keyword.ProcessedAt)
			if days >= 1 {
				db_keyword_service.UpdateNoRankDays(&finishedTask.Keyword, days)
			}

			if task_service.HasNextSearchPage(finishedTask.SearchedPage, settings.ReachRank) {
				if global.IsBetweenSearchTime() {
					nextPageCaptureTask := task_service.NextSearchPageCaptureTask(finishedTask)
					databases.Db.Save(nextPageCaptureTask)
				}
			} else {
				if global.IsBetweenSearchTime() {
					nextCycleCaptureTask := task_service.NextSearchCycleCaptureTask(finishedTask)
					databases.Db.Save(nextCycleCaptureTask)
				}
			}
		}
	}
}
