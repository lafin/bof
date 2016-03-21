package main

import (
	"bof/api"
	"bof/db"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"strconv"
	"time"
)

func addGroup(session *mgo.Session, typeGroup string, nameGroup string, id int) {
	group := db.GroupQuery(session)
	group.Insert(&db.Group{id, typeGroup, nameGroup})
}

func getGroups(session *mgo.Session) []db.Group {
	group := db.GroupQuery(session)
	records := []db.Group{}

	err := group.Find(nil).All(&records)
	if err != nil {
		log.Fatal(err)
	}
	return records
}

func tryDoRepost(session *mgo.Session, post string, group int) {
	message := db.MessageQuery(session)
	record := db.Post{}
	err := message.Find(bson.M{"post": post}).One(&record)
	if err != nil {
		err = message.Insert(&db.Post{post, group, time.Now()})
		if err != nil {
			log.Fatal(err)
		}
		// repost := api.DoRepost(client, "wall1_148582", "117456732", accessToken)
		// fmt.Println(repost.Response.Success)
	}
}

func main() {
	clientId := os.Getenv("CLIENT_ID")
	email := os.Getenv("CLIENT_EMAIL")
	password := os.Getenv("CLIENT_PASSWORD")
	dbServerAddress := os.Getenv("DB_SERVER")

	client := api.GetClient()
	accessToken := api.GetAccessToken(client, clientId, email, password)
	fmt.Printf("%s\n", accessToken)

	session := db.Connect(dbServerAddress)
	groups := getGroups(session)

	for _, group := range groups {
		posts := api.GetPosts(client, group.Name, "50")
		items := posts.Response.Items

		for _, val := range items {
			if val.IsPinned == 0 {
				fmt.Println("wall"+"1"+"_"+strconv.Itoa(val.ID), val.Date, val.Likes.Count, val.Text)
			}
		}
	}

	defer session.Close()
}
