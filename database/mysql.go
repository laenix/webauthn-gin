package database

import (
	"fmt"
	"webauthn/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	//加载config中的参数
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	datebase := viper.GetString("mysql.database")
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	charset := viper.GetString("mysql.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		username,
		password,
		host,
		port,
		datebase,
		charset,
	)
	//连接数据库
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}
	//自动建立表
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Certificate{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
