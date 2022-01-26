package db

import (
	"fmt"

	"github.com/go-pkgz/lgr"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var l = lgr.New(lgr.Msec, lgr.Debug, lgr.CallerFile, lgr.CallerFunc)

// Connect - connection to a db
func Connect(dbHost, dbUser, dbPassword, dbName string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", dbHost, dbUser, dbPassword, dbName)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		l.Logf("FATAL failed to connect database, %v", err)
	}
	err = db.AutoMigrate(&Group{}, &Post{}, &Dog{})
	if err != nil {
		l.Logf("FATAL db migration, %v", err)
	}
	return db
}
