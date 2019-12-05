package rank_util

type RequestBuilderRequest struct {
	SearchWord  string `json:"search_word" bson:"search_word"`
	Page        int    `json:"page" bson:"page"`
	Capture     bool   `json:"capture" bson:"capture"`
	SearchCycle int    `json:"search_cycle"`
	Priority    string `json:"priority"`
}

type ParseRankRequest struct {
	StartRank  int    `json:"start_rank" bson:"start_rank"`
	Body       string `json:"body" bson:"body"`
	CheckMatch string `json:"check_match" bson:"check_match"`
	SiteName   string `json:"site_name" bson:"site_name"`
}

type ParseRankResponse struct {
	Ranks []int `json:"ranks"`
}
