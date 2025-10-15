package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"log"
)

func InitDatabase() *gorm.DB {
	dbPsswPath := os.Getenv("MYSQL_ROOT_PASSWORD_FILE")
	dbPswd, fileErr := os.ReadFile(dbPsswPath)
	if fileErr != nil {
		log.Panicf("Password file not found: %s", fileErr)
	}

	dsn := fmt.Sprintf("root:%s@tcp(chat_database:3306)/chat_api?charset=utf8&parseTime=true", dbPswd)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect to database: %s", err)
	}

	return db
}
