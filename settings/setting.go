package settings

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	DbConnection string
	DbUsername   string
	DbPassword   string
	DbHost       string
	DbPort       string
	DbDatabase   string
)

var (
	ReachRank        int
	CheckRank        int
	SearchCycleLimit int
	SearchStartTime  time.Time
	SearchEndTime    time.Time
)

var (
	RankUtilApi    string
	DcWrapperApi   string
	RankArchiveApi string
)

var (
	QcloudCosSecretId  string
	QcloudCosSecretKey string
	QcloudCosRegion    string
	QcloudCosScheme    string
	QcloudCosBucket    string
	QcloudCosPrefix    string
)

var Debug bool

func init() {
	checkEnv()
	LoadSetting()
}

func checkEnv() {
	_ = godotenv.Load()
	needChecks := []string{
		"DB_CONNECTION", "DB_HOST", "DB_PORT", "DB_DATABASE", "DB_USERNAME", "DB_PASSWORD",
		"CHECK_RANK", "REACH_RANK", "SEARCH_CYCLE_LIMIT", "SEARCH_START_TIME", "SEARCH_END_TIME",
		"RANK_UTIL_API", "DC_WRAPPER_API", "SEND_RANKS_TO_ARCHIVE_API",
		"QCLOUD_COS_SECRET_ID", "QCLOUD_COS_SECRET_KEY", "QCLOUD_COS_REGION", "QCLOUD_COS_SCHEME", "QCLOUD_COS_BUCKET", "QCLOUD_COS_PREFIX",
	}

	for _, envKey := range needChecks {
		if os.Getenv(envKey) == "" {
			log.Fatalf("env %s missed", envKey)
		}
	}
}

func LoadSetting() {
	DbConnection = os.Getenv("DB_CONNECTION")
	DbUsername = os.Getenv("DB_USERNAME")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbDatabase = os.Getenv("DB_DATABASE")

	debug := os.Getenv("DEBUG")
	if debug != "" && debug != "false" && debug != "0" {
		Debug = true
	}

	CheckRank = loadIntFatal("CHECK_RANK")
	ReachRank = loadIntFatal("REACH_RANK")
	SearchCycleLimit = loadIntFatal("SEARCH_CYCLE_LIMIT")
	now := time.Now()
	SearchStartTime, _ = time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02")+" "+os.Getenv("SEARCH_START_TIME"), now.Location())
	SearchEndTime, _ = time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02")+" "+os.Getenv("SEARCH_END_TIME"), now.Location())

	RankUtilApi = os.Getenv("RANK_UTIL_API")
	DcWrapperApi = os.Getenv("DC_WRAPPER_API")
	RankArchiveApi = os.Getenv("SEND_RANKS_TO_ARCHIVE_API")

	QcloudCosSecretId = os.Getenv("QCLOUD_COS_SECRET_ID")
	QcloudCosSecretKey = os.Getenv("QCLOUD_COS_SECRET_KEY")
	QcloudCosRegion = os.Getenv("QCLOUD_COS_REGION")
	QcloudCosScheme = os.Getenv("QCLOUD_COS_SCHEME")
	QcloudCosBucket = os.Getenv("QCLOUD_COS_BUCKET")
	QcloudCosPrefix = os.Getenv("QCLOUD_COS_PREFIX")
}

func loadIntFatal(e string) int {
	intVar, err := strconv.Atoi(os.Getenv(e))
	if err != nil {
		log.Fatalf("env %s invalid\n", e)
	}

	return intVar
}
