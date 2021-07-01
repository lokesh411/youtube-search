package handlers

import (
	"fmt"
	"net/http"
	"youtube-search/controller"
	"youtube-search/utils"

	"github.com/labstack/echo/v4"
)

func GetVideos(context echo.Context) error {
	publishedDate := context.QueryParams().Get("publishedDate")
	videos, err := controller.FetchVideos(publishedDate)
	if err != nil {
		fmt.Println("Error in fetching videos :: ", err)
		return context.JSON(http.StatusInternalServerError, utils.StdResponse{Success: false})
	}
	return context.JSON(http.StatusOK, utils.StdResponse{Success: true, Data: videos, Message: "Success"})
}

func SearchVideos(context echo.Context) error {
	term := context.QueryParams().Get("term")
	videos, err := controller.SearchVideos(term)
	if err != nil {
		fmt.Println("Error in searching for videos :: ", err)
		return context.JSON(http.StatusInternalServerError, utils.StdResponse{Success: false})
	}
	return context.JSON(http.StatusOK, utils.StdResponse{Success: true, Data: videos, Message: "Success"})
}
