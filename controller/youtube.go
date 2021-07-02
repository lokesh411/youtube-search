package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"youtube-search/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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

var ctx = context.Background()
var redisPrefix = "consumedUnits_"

func getKey(maxUnits uint64) (string, bool, error) {
	// split all the strings and check the units consumed
	accessTokens := strings.Split(os.Getenv("youtube_access_tokens"), ",")
	for _, token := range accessTokens {
		val, err := models.Rdb.Get(ctx, redisPrefix+token).Uint64()
		if err != nil {
			if err == redis.Nil {
				return token, true, nil
			}
			return "", false, err
		}
		if val < maxUnits {
			return token, false, nil
		}
	}
	return "", false, errors.New("all the tokens have consumed max amount of units")
}

func incrementUsage(token string, units int64, new bool) error {
	if new {
		loc, e := time.LoadLocation("America/Los_Angeles")
		if e != nil {
			return e
		}
		timeInPST := time.Now().In(loc).UTC()
		year, month, day := time.Now().In(loc).AddDate(0, 0, 1).Date()
		expiryTime := time.Date(year, month, day, 0, 0, 0, 0, loc).UTC().Unix() - timeInPST.Unix()
		fmt.Println("expiry time :: ", expiryTime)
		err := models.Rdb.SetEX(ctx, redisPrefix+token, units, time.Duration(expiryTime*int64(time.Second))).Err()
		if err != nil {
			return e
		}
	} else {
		err := models.Rdb.IncrBy(ctx, redisPrefix+token, units).Err()
		return err
	}
	return nil
}

func scrapeVideosHelper(publishedDate string) error {
	fmt.Println("Search term :: ", os.Getenv("searchTerm"))
	accessToken, isNewToken, err := getKey(10000)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?type=video&order=date&publishedAfter=%s&key=%s&q=%s&part=snippet&maxResults=25", publishedDate, accessToken, os.Getenv("searchTerm"))
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
	err = incrementUsage(accessToken, 100, isNewToken)
	if err != nil {
		return err
	}
	var videoList []models.Video
	for _, videoData := range parsedResponse.Items {
		parsedTime, _ := time.Parse(time.RFC3339, videoData.Snippet.PublishedAt)
		item := models.Video{
			VideoPayload: models.VideoPayload{
				Title:         videoData.Snippet.Title,
				Description:   videoData.Snippet.Description,
				PublishedTime: parsedTime,
				Thumbnails:    videoData.Snippet.Thumbnails,
				ChannelId:     videoData.Snippet.ChannelId,
			},
		}
		videoList = append(videoList, item)
	}
	result := models.Db.Create(&videoList)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ScrapeVideos() {
	var fromPublishedDate string
	result := &models.Video{}
	response := models.Db.Order("published_time desc").Find(result)
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			fromPublishedDate = time.Now().AddDate(0, -1, 0).Format("2021-06-01T20:50:01.170Z")
		} else {
			fmt.Println("Error in fetching the published record :: ", response)
			return
		}
	} else {
		fromPublishedDate = result.PublishedTime.AddDate(0, 0, 1).Format(time.RFC3339)
	}
	error := scrapeVideosHelper(fromPublishedDate)
	if error != nil {
		fmt.Println("Error in scrapping videos :: ", error)
	} else {
		fmt.Println("Scrapped the videos")
	}
}
