package controller

import (
	"fmt"
	"time"
	models "youtube-search/models"

	"gorm.io/gorm"
)

var LIMIT = 10

// method to query mysql and fetch videos
func FetchVideos(publishedDate string) (*[]models.VideoPayload, error) {
	var videos []models.VideoPayload
	var result *gorm.DB
	if publishedDate == "" {
		result = models.Db.Limit(LIMIT).Order("published_time desc").Table("videos").Find(&videos)
	} else {
		parsedTime, err := time.Parse(time.RFC3339, publishedDate)
		fmt.Println("Parsed time :: ", parsedTime)
		if err != nil {
			return nil, result.Error
		}
		result = models.Db.Table("videos").Select("title", "description", "published_time", "thumbnails").Where("published_time < ?", parsedTime).Limit(LIMIT).Order("published_time desc").Find(&videos)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &videos, nil
}

// method to search videos - Full text search
func SearchVideos(term string) (*[]models.VideoPayload, error) {
	var videos []models.VideoPayload
	result := models.Db.Table("videos").Select("title", "description", "published_time", "thumbnails").Where("MATCH (title) AGAINST (? IN NATURAL LANGUAGE MODE)", term).Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}
	return &videos, nil
}
