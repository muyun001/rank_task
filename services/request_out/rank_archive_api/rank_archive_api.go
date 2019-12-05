package rank_archive_api

import (
	"encoding/json"
	"github.com/panwenbin/ghttpclient"
	"net/http"
	"rank-task/databases"
	"rank-task/databases/db_searched_rank_service"
	"rank-task/settings"
	"rank-task/structs/rank_archive"
	"strings"
)

var baseUrl string

const POST_保存历史排名 = "/history-ranks"

type RankArchiveResult struct {
	Msg string
}

func init() {
	baseUrl = settings.RankArchiveApi
}

// apiUrl 填充API参数并返回完整API地址
func apiUrl(path string, params map[string]string) string {
	for key, value := range params {
		path = strings.Replace(path, key, value, 1)
	}

	return baseUrl + path
}

// SendRanks 把is_send为0的searchedRanks组成HistoryRank格式发送到rank_archive
func SendRanks() {
	searchedRanks := db_searched_rank_service.UnsentSearchedRanks()

	historyRanks := make([]rank_archive.HistoryRank, 0)
	for i := range searchedRanks {
		historyRank := rank_archive.HistoryRank{
			Ip:         searchedRanks[i].Ip,
			Keyword:    searchedRanks[i].Keyword.Word,
			Engine:     searchedRanks[i].Keyword.Engine,
			CheckMatch: searchedRanks[i].Keyword.CheckMatch,
			TopRank:    searchedRanks[i].TopRank,
			Ranks:      searchedRanks[i].Ranks,
			Date:       searchedRanks[i].CreatedAt.Format("2006-01-02"),
			CaptureUrl: searchedRanks[i].CaptureUrl,
		}
		historyRanks = append(historyRanks, historyRank)
	}

	jsonBytes, err := json.Marshal(historyRanks)
	if err != nil {
		return
	}

	rankArchiveResult := RankArchiveResult{}
	apiUrl := apiUrl(POST_保存历史排名, nil)
	client := ghttpclient.PostJson(apiUrl, jsonBytes, nil)
	response, err := client.Response()
	err = client.ReadJsonClose(&rankArchiveResult)
	if err != nil {
		return
	}

	if response.StatusCode == http.StatusOK {
		var searchedRankIds []int
		for i := range searchedRanks {
			searchedRankIds = append(searchedRankIds, searchedRanks[i].ID)
		}
		databases.Db.Table("searched_ranks").Where("id in (?)", searchedRankIds).Update("is_send", 1)
	}
}
