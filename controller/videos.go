package controller

import (
	"fmt"
	"time"
	models "youtube-search/models"

	"gorm.io/gorm"
)

var LIMIT = 10

func FetchVideos(publishedDate string) (*[]models.Video, error) {
	var videos []models.Video
	var result *gorm.DB
	if publishedDate == "" {
		result = models.Db.Find(&videos).Limit(LIMIT).Order("published_time desc")
	} else {
		parsedTime, err := time.Parse(time.RFC3339, publishedDate)
		fmt.Println("Parsed time :: ", parsedTime)
		if err != nil {
			return nil, result.Error
		}
		result = models.Db.Where("published_time < ?", parsedTime).Find(&videos).Limit(LIMIT).Order("published_time desc")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &videos, nil
}

func SearchVideos(term string) (*[]models.Video, error) {
	var videos []models.Video
	result := models.Db.Where("MATCH (title) AGAINST (?) IN NATURAL LANGUAGE MODE").Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &videos, nil
}
