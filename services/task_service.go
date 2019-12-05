package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"rank-task/databases"
	"rank-task/services/keyword_service"
	"rank-task/settings"
	"strconv"
	"time"
)

type TaskService struct {
}

// AddSomeKeywordsToTasks: 按序添加一些关键词到任务
func (taskService *TaskService) AddSomeKeywordsToTasks() int {
	keywordIds := keyword_service.GetSomeUnTaskedKeywordIds()
	if len(keywordIds) != 0 {
		databases.Db.Exec("INSERT INTO tasks (`keyword_id`, `status`, `search_cycle`, `created_at`, `updated_at`) (SELECT id, 1, `search_cycle`, NOW(), NOW() FROM keywords WHERE id in (?))", keywordIds)
		databases.Db.Exec("UPDATE keywords SET searched_at = NOW(), searched_cycle = search_cycle WHERE id IN (?)", keywordIds)
	}
	return len(keywordIds)
}

// DaysApartToday: 某个时间点距离今天相隔的天数
func DaysApartToday(formerDay time.Time) int {
	t, _ := time.ParseInLocation("2006-01-02", formerDay.Format("2006-01-02"), time.Local)
	duration := time.Now().Sub(t)
	return int(duration.Hours() / 24)
}

// uploadCaptureToCos: 上传截图到腾讯云
func UploadCaptureToCos(capture string, keywordId int, engine string, ip string) (string, error) {
	cosUrl := fmt.Sprintf("http://%s.cos.%s.myqcloud.com", settings.QcloudCosBucket, settings.QcloudCosRegion)
	u, _ := url.Parse(cosUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  settings.QcloudCosSecretId,
			SecretKey: settings.QcloudCosSecretKey,
		},
	})

	name := fmt.Sprintf("%s/rank_imgs/%s/%s/%s/%s_%s_%s.png", settings.QcloudCosPrefix, time.Now().Format("2006-01-02"), engine, os.Getenv("DB_DATABASE"), strconv.Itoa(keywordId), ip, time.Now().Format("150405"))
	captureBin, err := base64.StdEncoding.DecodeString(capture)

	f := bytes.NewReader(captureBin)
	res, err := c.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		return "", err
	}

	defer res.Response.Body.Close()

	return cosUrl + "/" + name, nil
}
