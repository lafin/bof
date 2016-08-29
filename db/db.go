package db

import (
	"time"

	"gopkg.in/mgo.v2"
)

var session *mgo.Session

// Connect - start session to the db
func Connect(dbServerAddress string) (*mgo.Session, error) {
	var err error
	session, err = mgo.Dial(dbServerAddress)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	return session, nil
}

// PostQuery - get connection for the post collection
func PostQuery() (*mgo.Collection, error) {
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

	index = mgo.Index{
		Key: []string{"fingerprints"},
	}
	err = connect.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	return connect, nil
}

// GroupQuery - get connection for the group collection
func GroupQuery() (*mgo.Collection, error) {
	connect := session.DB("bof").C("group")

	index := mgo.Index{
		Key:        []string{"SourceID"},
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := connect.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	return connect, nil
}
