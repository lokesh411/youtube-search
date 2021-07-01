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
	models.Init()
	cronService := cron.New()
	// run cron every 5 mins
	controller.ScrapeVideos()
	cronService.AddFunc("*/1 * * * *", controller.ScrapeVideos)
	cronService.Start()
	e := echo.New()
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.Logger())
	e.GET("/health-check", func(context echo.Context) error {
		return context.String(http.StatusOK, "Hello, the server is running")
	})
	go func() {
		if err := e.Start(":5000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

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
