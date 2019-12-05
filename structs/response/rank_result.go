package response

// 不带截图的排名结果
type RankResult struct {
	Word string `json:"word"`
	Rank int    `json:"rank"`
}

// 带截图的排名结果
type CapturedRankResult struct {
	RankResult
	CaptureUrl string `json:"capture_url"`
}

// 不带截图的获取排名结果的回复
type RankResultsResponse []RankResult

// 带截图的获取排名结果的回复
type CapturedRankResultsResponse []CapturedRankResult