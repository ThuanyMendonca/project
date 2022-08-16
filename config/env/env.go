package env

import (
	"fmt"
	"os"

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

func Load() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Running application without env file")
	}

	GinMode = os.Getenv("GIN_MODE")

	ProjectDb = DataBase{}
	ProjectDb.Host = os.Getenv("DB_HOST")
	ProjectDb.User = os.Getenv("DB_USER")
	ProjectDb.Name = os.Getenv("DB_DB_NAME")
	ProjectDb.Password = os.Getenv("DB_PASSWORD")
	ProjectDb.TimeZone = os.Getenv("DB_TIME_ZONE")
	os.Getenv("DB_PORT")

}
