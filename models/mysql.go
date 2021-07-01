package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB
var err error

func Init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
			os.Getenv("db_username"),
			os.Getenv("db_password"),
			os.Getenv("db_host"),
			os.Getenv("db_port"),
			os.Getenv("db_name"),
		),
	}), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Println("Failed to connect mysql :: ", err)
	} else {
		log.Println("Connected to mysql")
	}
	Db.AutoMigrate(&Video{})
}
