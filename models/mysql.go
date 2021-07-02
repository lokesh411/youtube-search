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

func MysqlInit() {
	// initalize the logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	// open mysql connection
	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
			os.Getenv("mysql_username"),
			os.Getenv("mysql_password"),
			os.Getenv("mysql_host"),
			os.Getenv("mysql_port"),
			os.Getenv("mysql_db_name"),
		),
	}), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Println("Failed to connect mysql :: ", err)
	} else {
		log.Println("Connected to mysql")
	}
	Db.AutoMigrate(&Video{})
}
