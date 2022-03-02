package dbConnect

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//dbConnect
func DbConnect(Dbname string) (db *gorm.DB) {
	//配置MySQL
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := "127.0.0.1"
	port := 3306

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("連線失敗, error=" + err.Error())
	}

	return db
}
