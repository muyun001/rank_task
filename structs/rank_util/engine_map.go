package rank_util

var MapEngine map[string]string

func init() {
	MapEngine = map[string]string{
		"baidu_pc":           "baidu-pc",
		"baidu_mobile":       "baidu-mobile",
		"baidu_mini_program": "baidu-mini-program",
		"sogou_pc":           "sogou-pc",
		"sogou_mobile":       "sogou-mobile",
		"360_pc":             "360-pc",
		"sm_mobile":          "sm-mobile",
		"sug_baidu_pc":       "sug-baidu-pc",
		"sug_baidu_mobile":   "sug-baidu-mobile",
		"sug_360_pc":         "sug-360-pc",
		"sug_360_mobile":     "sug-360-mobile",
		"sug_sm_mobile":      "sug-sm-mobile",
		"sug_sogou_pc":       "sug-sogou-pc",
		"sug_sogou_mobile":   "sug-sogou-mobile",
	}
}
