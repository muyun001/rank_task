package services

import (
	"errors"
	"rank-task/databases"
	"rank-task/services/request_out/download_center"
	"rank-task/services/request_out/rank_util_api"
	"rank-task/structs/models"
	"rank-task/structs/rank_util"
)

type ParserService struct {
}

// BuildDcRequest: 构建单个下载中心任务结构
func (parserService *ParserService) BuildDcRequest(task *models.Task) (*download_center.DcRequest, error) {
	var keyword models.Keyword
	databases.Db.Model(task).Association("Keyword").Find(&keyword)

	engine, ok := rank_util.MapEngine[keyword.Engine]
	if !(ok) {
		return nil, errors.New("BuildDcRequest get engine error")
	}
	rankUrlApi := rank_util_api.NewRankUtilApi()

	DcPriority := download_center.DC_PRIORITY_中
	if keyword.NeedCapture == true {
		DcPriority = download_center.DC_PRIORITY_高
	}

	return rankUrlApi.PostRequestBuilder(engine, keyword.Word, task.SearchedPage, false, task.SearchCycle, DcPriority)
}

// ParseRanks: 解析某平台排名
func (parserService *ParserService) ParseRanks(checkMatch string, siteName string, html string, engine string, page int) ([]int, error) {
	startRank := (page - 1) * 10

	rankUrlApi := rank_util_api.NewRankUtilApi()
	res, err := rankUrlApi.PostRankExtractor(checkMatch, siteName, html, engine, startRank)
	if err != nil {
		return []int{}, err
	}

	return res.Ranks, nil
}

// BuildCaptureDcRequest: 构建单个下载中心任务结构 -- 截图任务
func (parserService *ParserService) BuildCaptureDcRequest(task *models.CaptureTask) (*download_center.DcRequest, error) {
	var keyword models.Keyword
	databases.Db.Model(task).Association("Keyword").Find(&keyword)

	engine, ok := rank_util.MapEngine[keyword.Engine]
	if !(ok) {
		return nil, errors.New("BuildDcRequest get engine error")
	}
	rankUrlApi := rank_util_api.NewRankUtilApi()
	return rankUrlApi.PostRequestBuilder(engine, keyword.Word, task.SearchedPage, true, task.SearchCycle, download_center.DC_PRIORITY_中)
}
