package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClientAndMigrate(user string, pass string, host string, port string, dbname string) (*Client, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)+"?parseTime=True&loc=Asia%2FTokyo&charset=utf8mb4"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&ScheduledTask{}, &Topic{}); err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c Client) Close() {
	sqlDB, err := c.db.DB()
	if err != nil {
		log.Panicf("Cannot Collect DB: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Errorf("Cannot Close DB connection: %v", err)
	}
}
