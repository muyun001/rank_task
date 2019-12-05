package rank_archive

type HistoryRank struct {
	Keyword    string `json:"keyword"`
	Engine     string `json:"engine"`
	CheckMatch string `json:"check_match"`
	TopRank    int    `json:"top_rank"`
	Ranks      string `json:"ranks"`
	Date       string `json:"date"`
	CaptureUrl string `json:"capture_url"`
	Ip         string `json:"ip"`
}
