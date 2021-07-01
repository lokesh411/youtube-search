package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"youtube-search/models"

	"gorm.io/datatypes"
)

var searchTerm = os.Getenv("searchTerm")

type IDResponse struct {
	Kind    string `json:"kind"`
	VideoId string `json:"videoId"`
}

type SnippetResponse struct {
	PublishedAt string         `json:"publishedAt"`
	ChannelId   string         `json:"channelId"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Thumbnails  datatypes.JSON `json:"thumbnails"`
}

type Items struct {
	Kind    string          `json:"kind"`
	Id      IDResponse      `json:"id"`
	Snippet SnippetResponse `json:"snippet"`
}

type Response struct {
	Kind  string  `json:"kind"`
	Items []Items `json:"items"`
}

func scrapeVideosHelper(publishedDate string) error {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?type=video&order=date&publishedAfter=%s&key=%s&q=%s&part=snippet", publishedDate, os.Getenv("youtube_access_token"), searchTerm)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var parsedResponse Response
	error := json.Unmarshal(responseData, &parsedResponse)
	if error != nil {
		return error
	}
	var videoList []models.Video
	for _, videoData := range parsedResponse.Items {
		parsedTime, _ := time.Parse(time.RFC3339, videoData.Snippet.PublishedAt)
		item := models.Video{Title: videoData.Snippet.Title, Description: videoData.Snippet.Description, PublishedTime: parsedTime, Thumbnails: videoData.Snippet.Thumbnails}
		videoList = append(videoList, item)
	}
	result := models.Db.Create(&videoList)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ScrapeVideos() {
	error := scrapeVideosHelper("2021-06-01T20:50:01.170Z")
	if error != nil {
		fmt.Println("Error in scrapping videos :: ", error)
	}
}
