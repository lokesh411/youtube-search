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

type VideoPayload struct {
	Title         string         `gorm:"index:,class:FULLTEXT" json:"title"`
	Description   string         `json:"description"`
	ChannelId     string         `json:"channelId"`
	PublishedTime time.Time      `gorm:"index" json:"publishedTime"`
	Thumbnails    datatypes.JSON `json:"thumbnails"`
}

type Video struct {
	gorm.Model
	VideoPayload
}
