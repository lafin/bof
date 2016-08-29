package db

import "time"

// Post - define struct of posts collection
type Post struct {
	Post         string
	Fingerprints []string
	From         int
	To           int
	Date         time.Time
}

// Group - define struct of groups collection
type Group struct {
	SourceID int
	Border   float32
	Message  string
}
