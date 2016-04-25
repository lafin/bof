package db

import (
	"time"

	"gopkg.in/mgo.v2"
)

// Post - define struct of post repost collection
type Post struct {
	Post string
	From int
	To   int
	Date time.Time
}

// Group - define struct of group collection
type Group struct {
	ID   int
	Type string
	Name string
}

// Connect - start session to the db
func Connect(dbServerAddress string) (*mgo.Session, error) {
	session, err := mgo.Dial(dbServerAddress)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	return session, nil
}

// PostQuery - get connection for the post collection
func PostQuery(session *mgo.Session) (*mgo.Collection, error) {
	connect := session.DB("bof").C("post")

	duration, _ := time.ParseDuration("30d")
	index := mgo.Index{
		Key:         []string{"post"},
		Unique:      true,
		ExpireAfter: duration,
	}
	err := connect.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	return connect, nil
}

// GroupQuery - get connection for the group collection
func GroupQuery(session *mgo.Session) (*mgo.Collection, error) {
	connect := session.DB("bof").C("group")

	index := mgo.Index{
		Key:    []string{"Type", "Name"},
		Unique: true,
	}
	err := connect.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	return connect, nil
}
