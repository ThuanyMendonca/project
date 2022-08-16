package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySqlDb(host string, port int64, user string, dbName string, password string, timeZone string) (*gorm.DB, error) {
	sn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s",
		user, password, host, port, dbName, timeZone)

	db, err := gorm.Open(mysql.Open(sn), &gorm.Config{})
	if err != nil {
		msg := fmt.Sprintf("Database connection error: %s", err)
		fmt.Println(msg)
		return nil, err
	}

	return db, nil
}
