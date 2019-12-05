package databases

import "rank-task/structs/models"

func AutoMigrate() {
	Db.AutoMigrate(&models.Keyword{}, &models.Task{}, &models.CaptureTask{}, &models.SearchedRank{}, &models.SiteName{}, &models.AppKey{}, &models.GetRank{})
}
