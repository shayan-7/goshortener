package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type Stats struct {
	ID        string `gorm:"unique"`
	Time      uint64
	IsBrowser bool
	Count     uint
}

func Subscribe(l *log.Logger, cache *redis.Client, db *gorm.DB) {
	fmt.Println("Processing started on redis queue:", "goshortener_stats")
	psNewMessage := cache.Subscribe("goshortener_stats")
	for {
		msg, _ := psNewMessage.ReceiveMessage()
		splitedMsg := strings.Split(msg.Payload, ":")
		id, isBrowser := splitedMsg[0], splitedMsg[1]
		l.Println("\033[34;1m:>Received message:", msg.Payload, "\033[0m")

		now := time.Now()
		currentTime := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			0,
			0,
			0,
			time.UTC,
		)

		var statExists bool
		err := db.Raw(
			"SELECT exists(SELECT * FROM stats WHERE id = ? AND time = ?)",
			id,
			currentTime.Unix(),
		).Row().Scan(&statExists)
		if err != nil {
			l.Println(err)
		}

		if !statExists {
			db.Exec(
				"INSERT INTO stats (id, time, is_browser, count) VALUES (?, ?, ?, ?)",
				id,
				currentTime.Unix(),
				isBrowser,
				1,
			)
		} else {
			db.Exec(
				"UPDATE stats SET count = count + 1 WHERE id = ? AND time = ? AND is_browser = ?",
				id,
				currentTime.Unix(),
				isBrowser,
			)
		}
	}
}
