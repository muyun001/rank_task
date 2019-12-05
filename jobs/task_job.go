package jobs

import (
	"fmt"
	"rank-task/channels"
	"rank-task/common/debug_log"
	"rank-task/common/ints"
	"rank-task/databases"
	"rank-task/databases/db_keyword_service"
	"rank-task/databases/db_searched_rank_service"
	"rank-task/databases/scopes/task_scope"
	"rank-task/global"
	"rank-task/services"
	"rank-task/services/request_out/baidu_site_name"
	"rank-task/services/request_out/download_center"
	"rank-task/services/task_service"
	"rank-task/settings"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"sync"
	"sync/atomic"
	"time"
)

// AddKeywordsToTasks: 按序添加关键词到任务
func AddKeywordsToTasks() {
	debug_log.Info(fmt.Sprintf("BeforeQueriedTasksCount: %d", global.BeforeQueriedTasksCount), "COUNT")
	if global.IsBetweenSearchTime() && global.BeforeQueriedTasksCount < logics.TASK_查询完成前任务数量限制 {
		service := &services.TaskService{}
		keywordsCount := service.AddSomeKeywordsToTasks()
		if keywordsCount > 0 {
			atomic.AddInt64(&global.BeforeQueriedTasksCount, int64(keywordsCount))
			return
		}
	}
	time.Sleep(time.Second * 3)
}

// SendTasksToChan: 从Chan发送任务
func SendTasksToChan() {
	var tasks []*models.Task
	limit := cap(channels.TaskSendingChan) - len(channels.TaskSendingChan)
	if limit == 0 {
		time.Sleep(time.Second)
		return
	}
	databases.Db.Model(&models.Task{}).Scopes(task_scope.UnQueried).Limit(limit).Scan(&tasks)

	if len(tasks) != 0 {
		task_service.SendTasksToChan(tasks, channels.TaskSendingChan, func(task *models.Task) {
			task.Status = logics.TASK_STATUS_查询中
			databases.Db.Save(&task)
		})
	} else {
		time.Sleep(time.Second * 5)
	}
}

// SendDcRequestsFromChan: 发送下载中心
func SendDcRequestsFromChan() {
	task := <-channels.TaskSendingChan
	parserService := &services.ParserService{}
	dcRequest, err := parserService.BuildDcRequest(task)
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

type UniqueKeyTaskGroup struct {
	UniqueKey string
	Tasks     []*models.Task
}

// FetchQueryParsedResult: 获取查询和解析结果
func FetchQueryParsedResult() {
	var beforeQueriedTasks []*models.Task
	databases.Db.
		Preload("Keyword").
		Scopes(task_scope.Querying).
		Order("updated_at").
		Find(&beforeQueriedTasks)

	if len(beforeQueriedTasks) == 0 {
		time.Sleep(time.Second * 2)
		return
	}

	fetchTasksResults(beforeQueriedTasks)
	time.Sleep(time.Second)
}

func fetchTasksResults(tasks []*models.Task) {
	var uniqueKeys []string
	uniqueKeyIndexedTasks := make(map[string]*models.Task)
	for _, task := range tasks {
		uniqueKeys = append(uniqueKeys, task.UniqueKey)
		uniqueKeyIndexedTasks[task.UniqueKey] = task
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

	var finishedTasks []*models.Task
	databases.Db.
		Preload("Keyword").
		Preload("Keyword.SiteName").
		Scopes(task_scope.UniqueKeysIn(finishedUniqueKeys)).
		Find(&finishedTasks)

	uniqueKeyTasksMap := task_service.UniqueKeyMappedTasks(finishedTasks)

	uniqueKeyTaskGroupChan := make(chan UniqueKeyTaskGroup)
	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			for {
				uniqueKeyTaskGroup, ok := <-uniqueKeyTaskGroupChan
				if !ok {
					break
				}
				tryFinishTask(uniqueKeyTaskGroup)
			}
			wg.Done()
		}()
	}
	for _, uniqueKey := range finishedUniqueKeys {
		uniqueKeyTaskGroup := UniqueKeyTaskGroup{
			UniqueKey: uniqueKey,
			Tasks:     uniqueKeyTasksMap[uniqueKey],
		}
		uniqueKeyTaskGroupChan <- uniqueKeyTaskGroup
		time.Sleep(time.Microsecond * 10)
	}
	close(uniqueKeyTaskGroupChan)

	wg.Wait()
}

func tryFinishTask(group UniqueKeyTaskGroup) {
	dc := download_center.NewDownloadCenter()
	dcResponse, err := dc.GetResponse(group.UniqueKey)
	if err != nil {
		taskCount := int64(len(group.Tasks))
		atomic.AddInt64(&global.BeforeQueriedTasksCount, -taskCount)
		databases.Db.Model(&models.Task{}).Where(models.Task{UniqueKey: group.UniqueKey}).Updates(models.Task{Status: logics.TASK_STATUS_查询失败})
		return
	}
	if dcResponse.Body == "" {
		_ = download_center.NewDownloadCenter().ResetRequest(group.UniqueKey)
		atomic.AddInt64(&global.BeforeQueriedTasksCount, 1)
		return
	}
	for _, finishedTask := range group.Tasks {
		atomic.AddInt64(&global.BeforeQueriedTasksCount, -1)
		parserService := services.ParserService{}
		ranks, err := parserService.ParseRanks(finishedTask.Keyword.CheckMatch, finishedTask.Keyword.SiteName.SiteName, dcResponse.Body, finishedTask.Keyword.Engine, finishedTask.SearchedPage)
		if err != nil {
			databases.Db.Model(&finishedTask).Updates(models.Task{Status: logics.TASK_STATUS_查询失败})
			continue
		}

		topRank := ints.Min(ranks...)
		isRankReached := topRank > 0 && topRank <= settings.ReachRank
		captureUrl := ""

		db_searched_rank_service.SaveSearchedRank(finishedTask.KeywordId, topRank, ranks, captureUrl, dcResponse.Ip)

		if isRankReached {
			databases.Db.Model(&finishedTask).Updates(models.Task{Status: logics.TASK_STATUS_查询达标})
		} else {
			databases.Db.Model(&finishedTask).Updates(models.Task{Status: logics.TASK_STATUS_查询不达标})
		}

		if topRank > 0 {
			if finishedTask.Keyword.NeedCapture {
				if isRankReached {
					captureTask := models.CaptureTask{
						KeywordId:    finishedTask.KeywordId,
						Status:       logics.TASK_STATUS_未查询,
						SearchedPage: 1,
						SearchCycle:  finishedTask.SearchCycle,
						CreatedAt:    time.Time{},
						UpdatedAt:    time.Time{},
					}
					databases.Db.Save(&captureTask)
				}
				continue
			}
			// 更新"has_new_rank","top_rank","no_rank_days"的三种情况:
			// 1.processed_at是昨天或者更早;
			// 2.processed_at是今天,top_rank是0;
			// 3.processed_at是今天,但现在排名比今天历史排名靠前(todayHistoryTopRank > topRank).
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

			if task_service.HasNextSearchPage(finishedTask.SearchedPage, settings.CheckRank) {
				if global.IsBetweenSearchTime() {
					nextPageTask := task_service.NextSearchPageTask(finishedTask)
					databases.Db.Save(nextPageTask)
					atomic.AddInt64(&global.BeforeQueriedTasksCount, 1)
				}
			}
		}
	}
}

// CheckSiteName: 检查未查过的SiteName
func CheckSiteName() {
	var checkMatches []string
	databases.Db.Raw("SELECT DISTINCT check_match FROM keywords LEFT JOIN site_names ON keywords.check_match = site_names.site_domain WHERE site_names.site_domain IS NULL").Pluck("check_match", &checkMatches)
	if len(checkMatches) == 0 {
		time.Sleep(time.Second * 5)
		return
	}
	for _, domain := range checkMatches {
		siteName := models.SiteName{
			SiteDomain: domain,
			SiteName:   baidu_site_name.BaiduPcSiteName(domain),
		}
		databases.Db.Save(&siteName)
	}
}

// ReCheckSiteName: 重新检查已查过的SiteName
func ReCheckSiteName() {
	var siteDomains []string
	databases.Db.Raw("SELECT site_domain FROM site_names").Pluck("site_domain", &siteDomains)
	for i := range siteDomains {
		siteName := models.SiteName{
			SiteDomain: siteDomains[i],
		}
		databases.Db.Model(&siteName).Where(siteName).Update(models.SiteName{SiteName: baidu_site_name.BaiduPcSiteName(siteDomains[i])})
	}
}
