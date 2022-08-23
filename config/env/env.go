package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DataBase struct {
	Host     string
	Port     int64
	User     string
	Name     string
	Password string
	TimeZone string
}

var (
	GinMode   string
	ProjectDb DataBase
)

var (
	AuthorizatorUrl string
)

func Load() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Running application without env file")
	}

	GinMode = os.Getenv("GIN_MODE")

	ProjectDb = DataBase{}
	ProjectDb.Host = os.Getenv("DB_HOST")
	ProjectDb.User = os.Getenv("DB_USER")
	ProjectDb.Name = os.Getenv("DB_DB_NAME")
	ProjectDb.Port, _ = strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)
	ProjectDb.Password = os.Getenv("DB_PASSWORD")
	ProjectDb.TimeZone = os.Getenv("DB_TIME_ZONE")
	AuthorizatorUrl = os.Getenv("AUTHORIZATOR_URL")

}
