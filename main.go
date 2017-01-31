package main

import (
	"github.com/lafin/bof/api"
	"github.com/lafin/bof/db"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func existRepostByID(info *api.Group, item *api.Post) bool {
	postID := getPostID(info, item)
	post, err := db.PostQuery()
	if err != nil {
		log.Fatal(err)
		return false
	}
	record := db.Post{}
	err = post.Find(bson.M{"post": postID}).One(&record)
	if err != nil {
		return false
	}
	return true
}

func existRepostByFiles(files [][]byte) bool {
	post, err := db.PostQuery()
	if err != nil {
		log.Fatal(err)
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

func getPostID(info *api.Group, item *api.Post) string {
	return "wall-" + strconv.Itoa(info.ID) + "_" + strconv.Itoa(item.ID)
}

func doRepost(files [][]byte, attachments []string, item *api.Post, info *api.Group, group *db.Group) (bool, error) {
	postID := getPostID(info, item)
	post, err := db.PostQuery()
	if err != nil {
		return false, err
	}

	from := info.ID
	to := group.SourceID
	message := group.Message

	if attachments == nil {
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
		if len(item.Text) > 0 {
			message = item.Text + " " + message
		}
		repost, err := api.DoPost(to, strings.Join(attachments, ","), url.QueryEscape(message))
		if err != nil {
			return false, err
		}
		if repost.Response.PostID > 0 {
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

func doRemoveDogs(groupID int) {
	start := 0
	offset := 1000
	maxRemoveDogs := 3

	for {
		users, err := api.GetListUsersofGroup(groupID, start, offset)
		if err != nil {
			log.Fatal(err)
			return
		}
		for _, user := range users.Response.Items {
			if (user.IsBanned() && maxRemoveDogs > 0) || user.IsDeleted() {
				if user.IsBanned() {
					maxRemoveDogs--
				}
				status, err := api.RemoveUserFromGroup(groupID, user.ID)
				if err != nil {
					log.Fatal(err)
					return
				}
				if status.Response != 1 {
					break
				}
			}
		}
		if len(users.Response.Items) < offset {
			break
		} else {
			start += offset
		}
	}
}

func main() {
	clientID := os.Getenv("CLIENT_ID")
	email := os.Getenv("CLIENT_EMAIL")
	password := os.Getenv("CLIENT_PASSWORD")
	dbServerAddress := os.Getenv("DB_SERVER")

	maxCountCheckInOneTime, err := strconv.ParseInt(os.Getenv("MAX_COUNT_CHECK_IN_ONE_TIME"), 10, 32)
	if err != nil || maxCountCheckInOneTime == 0 {
		maxCountCheckInOneTime = 3
	}

	log.Println("start")
	_, err = api.GetAccessToken(clientID, email, password)
	if err != nil {
		log.Fatal(err)
		return
	}

	session, err := db.Connect(dbServerAddress)
	if err != nil {
		log.Fatal(err)
		return
	}

	groups := db.GetGroups()
	for _, group := range groups {
		go doRemoveDogs(group.SourceID)

		groupsInfo, err := api.GetGroupsInfo(strconv.Itoa(group.SourceID), "links")
		if err != nil {
			log.Fatal(err)
			return
		}

		var ids []string
		for _, info := range groupsInfo.Response {
			r, _ := regexp.Compile("https://vk.com/(.*?)$")
			for _, link := range info.Links {
				ids = append(ids, r.FindStringSubmatch(link.URL)[1])
			}
		}

		groupsInfo, err = api.GetGroupsInfo(strings.Join(ids, ","), "")
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, info := range groupsInfo.Response {
			posts, err := api.GetPosts(strconv.Itoa(info.ID), "50")
			if err != nil {
				log.Fatal(err)
				return
			}

			border := int(posts.GetMaxCountLikes() * group.Border)
			for _, item := range posts.Response.Items {
				if item.IsPinned == 0 && item.Likes.Count > border {
					if !existRepostByID(&info, &item) {
						files, attachments := item.GetUniqueFiles()
						if files != nil && existRepostByFiles(files) {
							files = nil
							attachments = nil
						}
						reposted, err := doRepost(files, attachments, &item, &info, &group)
						if err == nil {
							if reposted {
								log.Println("Reposted")
							} else {
								log.Println("Skipped")
							}
							maxCountCheckInOneTime--
							if maxCountCheckInOneTime == 0 {
								log.Println("interrupted")
								defer session.Close()
								return
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
