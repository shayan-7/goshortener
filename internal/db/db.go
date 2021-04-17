package db

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       1,
	})
}

// GlobalDB a global db object will be used across different packages
var GlobalDB *gorm.DB

// InitDatabase creates a sqlite db
func InitDatabase() (err error) {
	dsn := url.URL{
		User:     url.UserPassword("postgres", "postgres"),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s", "localhost"),
		Path:     "goshortener_dev",
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	url := dsn.String()
	if os.Getenv("GOSHORTENER_DB_URL") == "UNIX_SOCKET" {
		url = "postgres:///goshortener?host=/var/run/postgresql/"
	}

	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  url,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal(err)
		return err
	}

	GlobalDB = db
	return
}
