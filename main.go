package main

import (
	"bof/api"
	"bof/db"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
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

func tryDoRepost(session *mgo.Session, client *http.Client, postId string, from, to int, accessToken string) {
	post := db.PostQuery(session)
	record := db.Post{}
	err := post.Find(bson.M{"post": postId}).One(&record)
	fmt.Println(err, postId)
	if err != nil {
		err = post.Insert(&db.Post{postId, from, to, time.Now()})
		if err != nil {
			log.Fatal(err)
		}
		repost := api.DoRepost(client, postId, to, accessToken)
		fmt.Println(repost.Response.Success)
	}
}

func main() {
	clientId := os.Getenv("CLIENT_ID")
	email := os.Getenv("CLIENT_EMAIL")
	password := os.Getenv("CLIENT_PASSWORD")
	dbServerAddress := os.Getenv("DB_SERVER")

	client := api.Client()
	accessToken := api.GetAccessToken(client, clientId, email, password)
	fmt.Printf("%s\n", accessToken)

	session := db.Connect(dbServerAddress)
	groups := getGroups(session)

	for _, group := range groups {
		posts := api.GetPosts(client, group.Name, "50")
		items := posts.Response.Items

		for _, val := range items {
			if val.IsPinned == 0 && val.Likes.Count > 1300 {
				tryDoRepost(session, client, "wall-"+strconv.Itoa(group.Id)+"_"+strconv.Itoa(val.ID), group.Id, 117456732, accessToken)
			}
		}
	}

	defer session.Close()
}
