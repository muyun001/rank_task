package baidu_site_name_test

import (
	"rank-task/services/request_out/baidu_site_name"
	"testing"
)

func TestBaiduPcSiteName(t *testing.T) {
	domain := "www.hzbenyu.com"
	expectSiteName := "杭州奔宇科技有限公司"
	siteName := baidu_site_name.BaiduPcSiteName(domain)
	if siteName != expectSiteName {
		t.Errorf("expect %s, got %s", expectSiteName, siteName)
	}
}
