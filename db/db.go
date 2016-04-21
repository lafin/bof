package db

import (
	"gopkg.in/mgo.v2"
	"time"
)

type Post struct {
	Post string
	From int
	To   int
	Date time.Time
}

type Group struct {
	Id   int
	Type string
	Name string
}

func Connect(dbServerAddress string) (*mgo.Session, error) {
	session, err := mgo.Dial(dbServerAddress)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	return session, nil
}

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
