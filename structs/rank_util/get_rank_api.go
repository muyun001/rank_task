package rank_util

type PutKeywords struct {
	Keyword     string `json:"keyword"`
	Engine      string `json:"engine"`
	NeedCapture bool   `json:"need_capture"`
	CheckMatch  string `json:"check_match"`
}
