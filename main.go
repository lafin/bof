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
	group, err := db.GroupQuery(session)
	if err != nil {
		log.Fatal(err)
		return
	}
	group.Insert(&db.Group{id, typeGroup, nameGroup})
}

func getGroups(session *mgo.Session) []db.Group {
	group, err := db.GroupQuery(session)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	records := []db.Group{}

	err = group.Find(nil).All(&records)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return records
}

func tryDoRepost(session *mgo.Session, client *http.Client, postId string, from, to int, accessToken string) {
	post, err := db.PostQuery(session)
	if err != nil {
		log.Fatal(err)
		return
	}

	record := db.Post{}
	err = post.Find(bson.M{"post": postId}).One(&record)
	if err != nil {
		err = post.Insert(&db.Post{postId, from, to, time.Now()})
		if err != nil {
			log.Fatal(err)
			return
		}
		repost, err := api.DoRepost(client, postId, to, accessToken)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println(repost.Response.Success)
	}
}

func getMaxCountLikes(posts *api.Post) float32 {
	max := 0
	items := posts.Response.Items
	for _, val := range items {
		if val.Likes.Count > max && val.IsPinned == 0 {
			max = val.Likes.Count
		}
	}
	return float32(max)
}

func main() {
	clientId := os.Getenv("CLIENT_ID")
	email := os.Getenv("CLIENT_EMAIL")
	password := os.Getenv("CLIENT_PASSWORD")
	dbServerAddress := os.Getenv("DB_SERVER")

	client := api.Client()
	accessToken, err := api.GetAccessToken(client, clientId, email, password)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%s\n", accessToken)

	session, err := db.Connect(dbServerAddress)
	if err != nil {
		log.Fatal(err)
		return
	}
	groups := getGroups(session)

	for _, group := range groups {
		posts, err := api.GetPosts(client, group.Name, "50")
		if err != nil {
			log.Fatal(err)
			return
		}

		border := int(getMaxCountLikes(posts) / 2.0 * 1.6)
		items := posts.Response.Items
		for _, val := range items {
			if val.IsPinned == 0 && val.Likes.Count > border {
				tryDoRepost(session, client, "wall-"+strconv.Itoa(group.Id)+"_"+strconv.Itoa(val.ID), group.Id, 117456732, accessToken)
			}
		}
	}

	defer session.Close()
}
