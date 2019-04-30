package service

import (
	"gangbu/exception"
	"gangbu/module/game/db"
	"github.com/name5566/leaf/log"
	"time"
)

func GetNotice() *db.DBNotice {
	return db.GetCurrentNotice()
}

func PubicNotice(title, context string, pubTime string) {
	log.Debug(pubTime)
	var publicTime time.Time
	if pubTime != "" {
		t, err :=  time.ParseInLocation("2006-01-02 15:04:05", pubTime, time.Local)
		publicTime = t
		if err != nil {
			exception.GameException(exception.TimeParseError)
		}
	} else {
		publicTime = time.Now()
	}
	db.CreateNotice(title, context, publicTime)
}
