package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"judgeMore/pkg/utils"
	"time"
)

var db *gorm.DB

func Init() {
	var err error
	dsn, err := utils.GetMysqlDSN()
	if err != nil {
		panic(err)
	}
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(20 * time.Second)
}
