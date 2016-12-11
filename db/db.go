package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"time"
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

	duration, err := time.ParseDuration("720h")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	index := mgo.Index{
		Key:         []string{"post"},
		Unique:      true,
		ExpireAfter: duration,
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

// GetGroups - get list of groups
func GetGroups() []Group {
	group, err := GroupQuery()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	records := []Group{}
	err = group.Find(nil).All(&records)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return records
}
