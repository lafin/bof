package main

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

// Post is
type Post struct {
	Post      string        `gorm:"primaryKey"`
	Files     pq.ByteaArray `gorm:"type:bytea[]"`
	From      int
	To        int
	CreatedAt time.Time `gorm:"index"`
}

// Group is
type Group struct {
	SourceID int `gorm:"primaryKey"`
	Border   float32
	Message  string
}

// Dog is
type Dog struct {
	SourceID  int       `gorm:"uniqueIndex:idx_dog"`
	UserID    int       `gorm:"uniqueIndex:idx_dog"`
	CheckedAt time.Time `gorm:"index"`
	CreatedAt time.Time `gorm:"index"`
}

// GetGroups - ...
func GetGroups(dbConnect *gorm.DB) ([]Group, error) {
	var groups []Group
	result := dbConnect.Find(&groups)
	return groups, result.Error
}
