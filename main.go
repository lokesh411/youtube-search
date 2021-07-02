package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"youtube-search/controller"
	"youtube-search/handlers"
	"youtube-search/models"

	"github.com/joho/godotenv"
	cron "github.com/robfig/cron/v3"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		fmt.Println("Successfully loaded .env file")
	}
	fmt.Println("Sleeping for 20 seconds to wait for starting mysql")
	time.Sleep(20 * time.Second)
	fmt.Println("Sleep done")
	// initialize mysql and redis
	models.Init()
	// Running to scrape videos at the start of the service
	controller.ScrapeVideos()
	cronService := cron.New()
	// run cron every 15 mins
	cronService.AddFunc("*/15 * * * *", controller.ScrapeVideos)
	cronService.Start()
	e := echo.New()
	// Initialising all the middlewares
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.Logger())
	e.GET("/health-check", func(context echo.Context) error {
		return context.String(http.StatusOK, "Hello, the server is running")
	})
	e.GET("/videos", handlers.GetVideos)
	e.GET("/videos/search", handlers.SearchVideos)
	go func() {
		if err := e.Start(":5000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server", err)
		}
	}()
	// having a graceful shutdown
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
