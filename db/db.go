package db

import (
	"gopkg.in/mgo.v2"
	"time"
)

type Post struct {
	Post  string
	Group int
	Date  time.Time
}

type Group struct {
	Id   int
	Type string
	Name string
}

func Connect(dbServerAddress string) *mgo.Session {
	session, err := mgo.Dial(dbServerAddress)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	return session
}

func MessageQuery(session *mgo.Session) *mgo.Collection {
	connect := session.DB("bof").C("message")

	duration, _ := time.ParseDuration("30d")
	index := mgo.Index{
		Key:         []string{"Post"},
		Unique:      true,
		ExpireAfter: duration,
	}
	err := connect.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return connect
}

func GroupQuery(session *mgo.Session) *mgo.Collection {
	connect := session.DB("bof").C("group")

	index := mgo.Index{
		Key:    []string{"Type", "Name"},
		Unique: true,
	}
	err := connect.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return connect
}
