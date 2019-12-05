package get_rank_api

import (
	"errors"
	"rank-task/databases"
	"rank-task/services/request_out/baidu_site_name"
	"rank-task/structs/models"
	"rank-task/structs/models/logics"
	"rank-task/structs/rank_util"
	"rank-task/structs/response"
	"strings"
	"time"
)

// ReceiveAndSaveKeywords: 接收并保存关键词数据
func ReceiveAndSaveKeywords(putKeywords *[]rank_util.PutKeywords) error {
	if len(*putKeywords) > 0 {
		for _, putKeyword := range *putKeywords {
			keyword := models.Keyword{}

			err := databases.Db.Where(map[string]interface{}{
				"word":         putKeyword.Keyword,
				"engine":       putKeyword.Engine,
				"need_capture": putKeyword.NeedCapture,
				"check_match":  putKeyword.CheckMatch,
			}).FirstOrInit(&keyword).Error
			if err != nil {
				return err
			}

			if keyword.Priority == logics.KEYWORD_PRIORITY_不查 {
				databases.Db.Model(&keyword).Updates(models.Keyword{
					Priority:    logics.KEYWORD_PRIORITY_中,
					ProcessedAt: time.Now(),
				})
			}

			if keyword.ID == 0 {
				keyword.CreatedAt = time.Now()
				keyword.ProcessedAt = time.Now()
				databases.Db.Save(&keyword)
			} else {
				keyword.CreatedAt = time.Now()
			}

			siteName := models.SiteName{
				SiteDomain: putKeyword.CheckMatch,
			}

			err = databases.Db.Where(siteName).FirstOrInit(&siteName).Error
			if err != nil {
				return err
			}

			if siteName.ID == 0 {
				siteName.SiteName = baidu_site_name.BaiduPcSiteName(putKeyword.CheckMatch)
				databases.Db.Save(&siteName)
			}
		}
	}

	return nil
}

// occupyRequestHash 占据RequestHash或者提示冲突
func occupyRequestHash(requestHash, engine, checkMatch, wordsJoined string) error {
	getRank := models.GetRank{}
	if databases.Db.Where(models.GetRank{RequestHash: requestHash}).First(&getRank).RecordNotFound() {
		getRank := &models.GetRank{
			CheckMatch:  checkMatch,
			Engine:      engine,
			RequestHash: requestHash,
			Words:       wordsJoined,
			CreatedAt:   time.Now(),
		}
		databases.Db.Save(&getRank)
	} else {
		if strings.Compare(getRank.Engine, engine) != 0 || strings.Compare(getRank.CheckMatch, checkMatch) != 0 || strings.Compare(getRank.Words, wordsJoined) != 0 {
			return errors.New("RequestHash冲突")
		}
	}

	return nil
}

// RankResultsResponse: 单站点单引擎的批量排名结果
func RankResultsResponse(checkMatch string, engine string, requestHash string, keywords []string) (response.RankResultsResponse, error) {
	matchCondition := map[string]interface{}{
		"need_capture": false,
		"check_match":  checkMatch,
		"engine":       engine,
	}

	wordRankResponse := make(response.RankResultsResponse, 0)
	err := databases.Db.Model(&models.Keyword{}).
		Select("word, top_rank as rank").
		Where(matchCondition).
		Where("has_new_rank = ?", true).
		Where("word in (?)", keywords).
		Scan(&wordRankResponse).
		Error

	if err != nil {
		return wordRankResponse, err
	}

	err = databases.Db.Model(&models.Keyword{}).
		Where(matchCondition).
		Where("word in (?)", keywords).
		Update(models.Keyword{
			ProcessedAt: time.Now(),
		}).
		Error

	if err != nil {
		return response.RankResultsResponse{}, err
	}

	if len(wordRankResponse) != 0 {
		var words []string
		for _, singleResponse := range wordRankResponse {
			words = append(words, singleResponse.Word)
		}
		wordsJoined := strings.Join(words, ",")
		err = occupyRequestHash(requestHash, engine, checkMatch, wordsJoined)
		if err != nil {
			return response.RankResultsResponse{}, err
		}
	}

	return wordRankResponse, nil
}

// CapturedRankResultsResponse: 单站点单引擎的批量带截图排名结果
func CapturedRankResultsResponse(checkMatch string, engine string, requestHash string, keywords []string) (response.CapturedRankResultsResponse, error) {
	matchCondition := models.Keyword{
		NeedCapture: true,
		CheckMatch:  checkMatch,
		Engine:      engine,
	}

	wordCapturedRankResponse := make(response.CapturedRankResultsResponse, 0)
	err := databases.Db.Model(&models.Keyword{}).
		Select("word, top_rank as rank, capture_url").
		Where(matchCondition).
		Where("has_new_rank = ?", true).
		Where("word in (?)", keywords).
		Scan(&wordCapturedRankResponse).
		Error

	if err != nil {
		return wordCapturedRankResponse, err
	}

	if len(wordCapturedRankResponse) != 0 {
		err := databases.Db.Model(&models.Keyword{}).
			Where(matchCondition).
			Where("word in (?)", keywords).
			Update(models.Keyword{
				ProcessedAt: time.Now(),
			}).
			Error
		if err != nil {
			return response.CapturedRankResultsResponse{}, err
		}

		var words []string
		for _, singleResponse := range wordCapturedRankResponse {
			words = append(words, singleResponse.Word)
		}
		wordsJoined := strings.Join(words, ",")
		err = occupyRequestHash(requestHash, engine, checkMatch, wordsJoined)
		if err != nil {
			return response.CapturedRankResultsResponse{}, err
		}
	}

	return wordCapturedRankResponse, nil
}

// ConfirmRanks: 确认排名已获取
func ConfirmRanks(requestHash string) error {
	getRank := &models.GetRank{}
	if databases.Db.Select("engine, check_match, words").Where("request_hash = ?", requestHash).First(&getRank).RecordNotFound() {
		return nil
	}

	keywords := strings.Split(getRank.Words, ",")
	err := databases.Db.Model(&models.Keyword{}).
		Where("has_new_rank = 1 and engine = ? and check_match = ? and word in (?)", getRank.Engine, getRank.CheckMatch, keywords).
		Updates(map[string]interface{}{
			"has_new_rank": false,
		}).
		Error
	if err != nil {
		return err
	}
	databases.Db.Delete(&getRank)
	return nil
}
