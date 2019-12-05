package baidu_site_name

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/panwenbin/ghttpclient"
	"github.com/panwenbin/ghttpclient/header"
	"rank-task/structs/rank_util"
)

// BaiduPcSiteName 通过域名获取BaiduPC的站点名称
func BaiduPcSiteName(siteDomain string) string {
	sourceUrl := fmt.Sprintf("https://www.baidu.com/s?ie=utf-8&wd=site:%s", siteDomain)
	headers := header.GHttpHeader{}
	headers.UserAgent(rank_util.RandomUserAgentForEngine("baidu_pc"))
	headers.AcceptEncodingGzip()

	body, err := ghttpclient.Get(sourceUrl, headers).ReadBodyClose()
	if err != nil {
		return ""
	}

	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return ""
	}

	siteName := ""
	dom.Find("div#content_left > div").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if i > 3 {
			return false
		}
		siteNameItem := selection.Find("a.c-showurl img.source-icon")
		if len(siteNameItem.Nodes) == 0 {
			return true
		}

		siteName = siteNameItem.Nodes[0].NextSibling.Data
		if siteName != "" {
			return false
		}
		return true
	})

	return siteName
}

