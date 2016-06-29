package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lafin/bof/api"
	"github.com/lafin/bof/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func addGroup(session *mgo.Session, typeGroup string, sourceName string, sourceID int, destinationID int, border float32) {
	group, err := db.GroupQuery(session)
	if err != nil {
		log.Fatal(err)
		return
	}
	group.Insert(&db.Group{SourceID: sourceID, Border: border})
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

func existRepost(session *mgo.Session, postID string) bool {
	post, err := db.PostQuery(session)
	if err != nil {
		log.Fatal(err)
		return true
	}

	record := db.Post{}
	err = post.Find(bson.M{"post": postID}).One(&record)
	if err != nil {
		return false
	}
	return true
}

func doRepost(session *mgo.Session, client *http.Client, postID string, from, to int, message, accessToken string) int {
	post, err := db.PostQuery(session)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	repost, err := api.DoRepost(client, postID, to, message, accessToken)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	fmt.Println(repost.Response.Success)
	if repost.Response.Success == 1 {
		err = post.Insert(&db.Post{postID, from, to, time.Now()})
		if err != nil {
			log.Fatal(err)
			return 0
		}
		return repost.Response.PostID
	}
	return 0
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
	clientID := os.Getenv("CLIENT_ID")
	email := os.Getenv("CLIENT_EMAIL")
	password := os.Getenv("CLIENT_PASSWORD")
	dbServerAddress := os.Getenv("DB_SERVER")

	client := api.Client()
	accessToken, err := api.GetAccessToken(client, clientID, email, password)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(accessToken)

	session, err := db.Connect(dbServerAddress)
	if err != nil {
		log.Fatal(err)
		return
	}

	groups := getGroups(session)
	for _, group := range groups {
		groupInfo, err := api.GetGroupsInfo(client, strconv.Itoa(group.SourceID), "links")
		if err != nil {
			log.Fatal(err)
			return
		}

		links := groupInfo.Response[0].Links
		r, _ := regexp.Compile("https://vk.com/(.*?)$")
		var ids []string
		for _, link := range links {
			ids = append(ids, r.FindStringSubmatch(link.URL)[1])
		}

		groupsInfo, err := api.GetGroupsInfo(client, strings.Join(ids, ","), "")
		if err != nil {
			log.Fatal(err)
			return
		}

		infos := groupsInfo.Response
		for _, info := range infos {
			posts, err := api.GetPosts(client, strconv.Itoa(info.ID), "50")
			if err != nil {
				log.Fatal(err)
				return
			}

			border := int(getMaxCountLikes(posts) * group.Border)
			items := posts.Response.Items
			var repostID int
			var postID string
			var exist bool
			for _, val := range items {
				if val.IsPinned == 0 && val.Likes.Count > border {
					postID = "wall-" + strconv.Itoa(info.ID) + "_" + strconv.Itoa(val.ID)
					exist = existRepost(session, postID)
					if exist == false {
						repostID = doRepost(session, client, postID, info.ID, group.SourceID, group.Message, accessToken)
						if repostID == 0 {
							fmt.Println("Unsuccess try do repost")
							return
						}
					}
				}
			}
		}
	}

	defer session.Close()
}
