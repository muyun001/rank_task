package rank_util_api

import (
	"encoding/json"
	"github.com/panwenbin/ghttpclient"
	"rank-task/services/request_out/download_center"
	"rank-task/settings"
	"rank-task/structs/rank_util"
	"strings"
)

type RankUtilApi struct {
	BaseUrl string
}

const (
	POST_REQUEST_BUILDER = "/request-builder/:engine"
	POST_RANK_EXTRACTOR  = "/rank-extractor/:engine"
)

// NewRankUtilApi: 创建新的排名工具api
func NewRankUtilApi() *RankUtilApi {
	rankUtilApi := &RankUtilApi{}
	rankUtilApi.BaseUrl = settings.RankUtilApi

	return rankUtilApi
}

// apiUrl: 填充API参数并返回完整API地址
func (r *RankUtilApi) apiUrl(path string, params map[string]string) string {
	for key, value := range params {
		path = strings.Replace(path, key, value, 1)
	}

	return r.BaseUrl + path
}

// PostRequestBuilder: 请求构建DcRequest对象
func (r *RankUtilApi) PostRequestBuilder(engine, searchWord string, page int, capture bool, searchCycle int, priority string) (*download_center.DcRequest, error) {
	apiUrl := r.apiUrl(POST_REQUEST_BUILDER, map[string]string{
		":engine": engine,
	})

	requestBuilderRequest := rank_util.RequestBuilderRequest{
		SearchWord:  searchWord,
		Page:        page,
		Capture:     capture,
		SearchCycle: searchCycle,
		Priority:    priority,
	}
	jsonBytes, err := json.Marshal(requestBuilderRequest)
	var dcRequest download_center.DcRequest
	err = ghttpclient.PostJson(apiUrl, jsonBytes, nil).ReadJsonClose(&dcRequest)
	if err != nil {
		return nil, err
	}

	return &dcRequest, nil
}

// PostRankExtractor: 排名解析
func (r *RankUtilApi) PostRankExtractor(checkMatch string, siteName string, html string, engine string, startRank int) (*rank_util.ParseRankResponse, error) {
	apiUrl := r.apiUrl(POST_RANK_EXTRACTOR, map[string]string{
		":engine": rank_util.MapEngine[engine],
	})
	rankExtractor := rank_util.ParseRankRequest{
		StartRank:  startRank,
		Body:       html,
		CheckMatch: checkMatch,
		SiteName:   siteName,
	}
	jsonBytes, err := json.Marshal(rankExtractor)
	rankExtractorResponse := rank_util.ParseRankResponse{}
	err = ghttpclient.PostJson(apiUrl, jsonBytes, nil).ReadJsonClose(&rankExtractorResponse)
	if err != nil {
		return nil, err
	}

	return &rankExtractorResponse, nil
}
