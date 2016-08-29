package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lafin/bof/api"
	"github.com/lafin/bof/db"
	"github.com/lafin/bof/util"
	"gopkg.in/mgo.v2/bson"
)

func addGroup(typeGroup string, sourceName string, sourceID int, destinationID int, border float32) {
	group, err := db.GroupQuery()
	if err != nil {
		log.Fatal(err)
		return
	}
	group.Insert(&db.Group{SourceID: sourceID, Border: border})
}

func getGroups() []db.Group {
	group, err := db.GroupQuery()
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

func existRepost(postID string) bool {
	post, err := db.PostQuery()
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

func checkOnUniqueness(post api.PostItem) bool {
	for _, val := range post.Attachments {
		if len(val.Photo.Photo75) > 0 {
			fmt.Println("Photo75")
			fmt.Println(util.GetSha1(val.Photo.Photo75))
		}
	}
	return false
}

func doRepost(postID string, fingerprints []string, from, to int, message string) int {
	post, err := db.PostQuery()
	if err != nil {
		log.Fatal(err)
		return 0
	}

	repost, err := api.DoRepost(postID, to, message)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	fmt.Println(repost.Response.Success)
	if repost.Response.Success == 1 {
		err = post.Insert(&db.Post{postID, fingerprints, from, to, time.Now()})
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

	_, err := api.GetAccessToken(clientID, email, password)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("connected")

	session, err := db.Connect(dbServerAddress)
	if err != nil {
		log.Fatal(err)
		return
	}

	groups := getGroups()
	for _, group := range groups {
		groupInfo, err := api.GetGroupsInfo(strconv.Itoa(group.SourceID), "links")
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

		groupsInfo, err := api.GetGroupsInfo(strings.Join(ids, ","), "")
		if err != nil {
			log.Fatal(err)
			return
		}

		infos := groupsInfo.Response
		for _, info := range infos {
			posts, err := api.GetPosts(strconv.Itoa(info.ID), "50")
			if err != nil {
				log.Fatal(err)
				return
			}

			border := int(getMaxCountLikes(posts) * group.Border)
			items := posts.Response.Items
			// var repostID int
			var postID string
			for _, val := range items {
				if val.IsPinned == 0 && val.Likes.Count > border {
					postID = "wall-" + strconv.Itoa(info.ID) + "_" + strconv.Itoa(val.ID)
					if !existRepost(postID) && checkOnUniqueness(val) {
						// repostID = doRepost(postID, info.ID, group.SourceID, group.Message)
						// if repostID == 0 {
						// 	fmt.Println("Unsuccess try do repost")
						// 	return
						// }
					}
					checkOnUniqueness(val)
				}
			}
		}
	}

	defer session.Close()
}
