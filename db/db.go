package db

import (
	"time"

	"gopkg.in/mgo.v2"
)

var session *mgo.Session

type (
	// Collection - ...
	Collection = mgo.Collection
	// Session - ...
	Session = mgo.Session
)

// Connect - start session to the db
func Connect(dbServerAddress string) (*Session, error) {
	var err error
	session, err = mgo.Dial(dbServerAddress)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

// PostQuery - get connection for the post collection
func PostQuery() (*Collection, error) {
	connect := session.DB("bof").C("post")

	duration, err := time.ParseDuration("720h")
	if err != nil {
		return nil, err
	}
	index := mgo.Index{
		Key:         []string{"date"},
		ExpireAfter: duration,
	}
	err = connect.EnsureIndex(index)
	if err != nil {
		return nil, err
	}
	return connect, nil
}

// GroupQuery - get connection for the group collection
func GroupQuery() (*Collection, error) {
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

// GetGroups - get list of groups
func GetGroups() ([]Group, error) {
	group, err := GroupQuery()
	if err != nil {
		return nil, err
	}
	records := []Group{}
	err = group.Find(nil).All(&records)
	if err != nil {
		return nil, err
	}
	return records, nil
}
