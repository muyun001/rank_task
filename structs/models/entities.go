package models

import "time"

type Keyword struct {
	ID            int       `gorm:"primary_key" json:"id"`
	Word          string    `gorm:"type:varchar(255);not null;unique_index:get_rank" json:"word"`
	Engine        string    `gorm:"type:varchar(255);unique_index:get_rank" json:"engine"`
	CheckMatch    string    `gorm:"type:varchar(255);unique_index:get_rank;index" json:"check_match"`
	NeedCapture   bool      `gorm:"type:tinyint;default:0" json:"need_capture"`
	SearchedCycle int       `gorm:"type:int;default:0" json:"searched_cycle"`
	SearchCycle   int       `gorm:"type:int;default:1" json:"search_cycle"`
	Priority      int       `gorm:"type:int;default:2;index:get_tasks,down_priority" json:"priority"`
	NoRankDays    int       `gorm:"type:int;index:down_priority" json:"no_rank_days"`
	HasNewRank    bool      `gorm:"type:tinyint;index:get_results" json:"has_new_rank"`
	TopRank       int       `gorm:"type:int;" json:"top_rank"`
	CaptureUrl    string    `gorm:"type:varchar(255)" json:"capture_url"`
	CreatedAt     time.Time `json:"created_at"`
	SearchedAt    time.Time `gorm:"index:get_tasks" json:"searched_at"`
	ProcessedAt   time.Time `gorm:"index:get_results" json:"processed_at"`
	SiteName      SiteName  `gorm:"foreignkey:SiteDomain;association_foreignkey:CheckMatch"`
}

type Task struct {
	ID           int       `gorm:"primary_key" json:"id"`
	KeywordId    int       `gorm:"type:int;index:task_keyword_id" json:"keyword_id"`
	Status       int       `gorm:"type:int;not null;index:task_status" json:"status"`
	UniqueKey    string    `gorm:"type:varchar(255);index:keyword_unique_key;" json:"unique_key"`
	SearchedPage int       `gorm:"type:int;not null;default:1" json:"searching_page"`
	SearchCycle  int       `gorm:"type:int;not null;default:0" json:"search_cycle"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Keyword      Keyword
	SearchedRank SearchedRank
}

type CaptureTask struct {
	ID           int       `gorm:"primary_key" json:"id"`
	KeywordId    int       `gorm:"type:int;index:task_keyword_id" json:"keyword_id"`
	Status       int       `gorm:"type:int;not null;index:task_status" json:"status"`
	UniqueKey    string    `gorm:"type:varchar(255);index:capture_kw_unique_key;" json:"unique_key"`
	SearchedPage int       `gorm:"type:int;not null;default:1" json:"searching_page"`
	SearchCycle  int       `gorm:"type:int;not null;default:0" json:"search_cycle"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Keyword      Keyword
	SearchedRank SearchedRank
}

type SearchedRank struct {
	ID         int       `gorm:"primary_key" json:"id"`
	KeywordId  int       `gorm:"type:int;index:keyword_id_created_at" json:"keyword_id"`
	TopRank    int       `gorm:"type:int;" json:"top_rank"`
	Ranks      string    `gorm:"type:varchar(64);" json:"ranks"`
	CaptureUrl string    `gorm:"type:longtext;" json:"capture_url"`
	Ip         string    `gorm:"type:varchar(64)" json:"ip"`
	IsSend     bool      `gorm:"type:tinyint;index;default:0" json:"is_send"`
	CreatedAt  time.Time `gorm:"index:keyword_id_created_at" json:"created_at"`
	Keyword    Keyword
}

type SiteName struct {
	ID         int    `gorm:"primary_key" json:"id"`
	SiteDomain string `gorm:"type:varchar(128);unique_index:site_name_site_domain" json:"site_domain"`
	SiteName   string `gorm:"type:varchar(128);unique_index:site_name_site_domain" json:"site_name"`
}

type AppKey struct {
	ID         int    `gorm:"primary_key" json:"id"`
	SiteDomain string `gorm:"type:varchar(128);index:app_keys_site_domain" json:"site_domain"`
	AppKey     string `gorm:"type:varchar(255);" json:"app_key"`
}

type GetRank struct {
	CheckMatch  string    `gorm:"type:varchar(255);" json:"check_match"`
	Engine      string    `gorm:"type:varchar(255);" json:"engine"`
	Words       string    `gorm:"type:longtext;not null;" json:"words"`
	RequestHash string    `gorm:"type:varchar(255);not null;unique_index:request_hash" json:"request_hash"`
	CreatedAt   time.Time `json:"created_at"`
}
