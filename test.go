package main

import (
	"fmt"
	"github.com/panwenbin/ghttpclient"
	"github.com/panwenbin/ghttpclient/header"
	"rank-task/structs/rank_util"
	"strings"
	"time"
)

func main() {
	url := "https://www.sogou.com/web?query=%E8%8B%8F%E5%B7%9E%E5%A4%A9%E6%B0%94&_ast=1572240866&_asf=www.sogou.com&w=01025001&cid=&s_from=result_up&oq=&ri=7&sourceid=sugg&suguuid=e07e5ef6-542a-4f54-bddd-df8a4ec49b85&sut=4064133&sst0=1572244932426&lkt=0%2C0%2C0&sugsuv=1568950571589337&sugtime=1572244932426"

	for i := 0; i < 101; i++ {
		headers := header.GHttpHeader{}
		headers.UserAgent(rank_util.RandomUserAgentForEngine("sogou_pc"))
		headers.AcceptEncodingGzip()
		body, err := ghttpclient.Get(url, headers).ReadBodyClose()
		if err != nil {
			fmt.Println(err)
		}
		if strings.Contains(string(body), "验证") {
			fmt.Println("failure")
			fmt.Println(string(body))
		} else {
			fmt.Println(i, "success")
			time.Sleep(time.Second)
		}
	}
}
