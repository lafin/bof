package main

import (
	"github.com/lafin/bof/api"
	"github.com/lafin/bof/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func getPostQuery() (*mgo.Collection, error) {
	post, err := db.PostQuery()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return post, nil
}

func existRepostByID(postID string) bool {
	post, err := getPostQuery()
	if err != nil {
		return false
	}
	record := db.Post{}
	err = post.Find(bson.M{"post": postID}).One(&record)
	if err != nil {
		return false
	}
	return true
}

func existRepostByFiles(files [][]byte, postContext api.Post) bool {
	post, err := getPostQuery()
	if err != nil {
		return false
	}
	records := []db.Post{}
	err = post.Find(bson.M{"files": bson.M{"$exists": true}}).Sort("date").All(&records)
	if err != nil {
		return false
	}

	for _, record := range records {
		for _, storedFile := range record.Files {
			if len(storedFile) != 0 {
				for _, file := range files {
					if len(file) != 0 {
						percent, err := Compare(storedFile, file)
						if err != nil {
							log.Fatal(err)
						} else {
							if percent < 0.05 {
								postID := "wall" + strconv.Itoa(postContext.OwnerID) + "_" + strconv.Itoa(postContext.ID)
								log.Println("filtered", record.Post, postID, percent)
								return true
							}
						}
					}
				}
			}
		}
	}

	return false
}

func getUniqueFiles(post api.Post) [][]byte {
	files := make([][]byte, 5)
	var file []byte

	for _, val := range post.Attachments {
		switch val.Type {
		case "photo":
			if len(val.Photo.Photo75) > 0 {
				file = GetData(val.Photo.Photo75)
			}
		case "doc":
			if len(val.Doc.URL) > 0 {
				file = GetData(val.Doc.URL)
			}
		}
		files = append(files, file)
	}
	found := existRepostByFiles(files, post)
	if found {
		files = nil
	}
	return files
}

func getPostID(info *api.Group, item *api.Post) string {
	return "wall-" + strconv.Itoa(info.ID) + "_" + strconv.Itoa(item.ID)
}

func doRepost(postID string, files [][]byte, info *api.Group, group *db.Group) (bool, error) {
	post, err := db.PostQuery()
	if err != nil {
		return false, err
	}

	from := info.ID
	to := group.SourceID
	message := group.Message

	if files == nil {
		err = post.Insert(&db.Post{
			Post:  postID,
			Files: files,
			From:  from,
			To:    to,
			Date:  time.Now()})
		if err != nil {
			return false, err
		}
	} else {
		repost, err := api.DoRepost(postID, to, message)
		if err != nil {
			return false, err
		}
		if repost.Response.Success == 1 {
			err = post.Insert(&db.Post{
				Post:  postID,
				Files: files,
				From:  from,
				To:    to,
				Date:  time.Now()})
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func getMaxCountLikes(posts *api.Posts) float32 {
	max := 0
	items := posts.Response.Items
	for _, item := range items {
		if item.Likes.Count > max && item.IsPinned == 0 {
			max = item.Likes.Count
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
	log.Println("start")

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
			var reposted bool
			var postID string
			for _, item := range items {
				if item.IsPinned == 0 && item.Likes.Count > border {
					postID = getPostID(&info, &item)
					if !existRepostByID(postID) {
						files := getUniqueFiles(item)
						reposted, err = doRepost(postID, files, &info, &group)
						if err == nil {
							if reposted {
								log.Println("Reposted")
							} else {
								log.Println("Skipped")
							}
						} else {
							log.Fatal(err)
							defer session.Close()
							return
						}
					}
				}
			}
		}
	}

	log.Println("done")
	defer session.Close()
}
