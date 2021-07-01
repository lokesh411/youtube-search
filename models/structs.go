package models

import (
	"time"

	datatypes "gorm.io/datatypes"
	"gorm.io/gorm"
)

type Thumbnail struct {
	Url    string `json:"url"`
	Width  uint64 `json:"width"`
	Height uint64 `json:"height"`
}

type Video struct {
	gorm.Model
	Title         string `gorm:"index:,class:FULLTEXT"`
	Description   string
	PublishedTime time.Time      `gorm:"index"`
	Thumbnails    datatypes.JSON
}
