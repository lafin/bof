// Package db handle work with db
package db

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

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

// GetPosts - ...
func GetPosts(dbConnect *gorm.DB) ([]Post, error) {
	var posts []Post
	result := dbConnect.Where("created_at > NOW() - INTERVAL '1 month'").Find(&posts)
	return posts, result.Error
}

// GetGroups - ...
func GetGroups(dbConnect *gorm.DB) ([]Group, error) {
	var groups []Group
	result := dbConnect.Find(&groups)
	return groups, result.Error
}
