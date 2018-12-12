package main

import (
	"errors"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lafin/bof/db"
	"github.com/lafin/bof/utils"
	copier "github.com/lafin/copier"
	api "github.com/lafin/vk"
	"gopkg.in/mgo.v2/bson"
)

func existRepostByID(info *api.Group, item *api.Post) bool {
	postID := getPostID(info, item)
	post, err := db.PostQuery()
	if err != nil {
		log.Fatalf("[existRepostByID] error: %s post_id: %s", err, postID)
		return false
	}
	record := db.Post{}
	err = post.Find(bson.M{"post": postID}).One(&record)
	return err == nil
}

func existRepostByFiles(files [][]byte) bool {
	post, err := db.PostQuery()
	if err != nil {
		log.Fatalf("[existRepostByFiles] error: %s", err)
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
						percent, err := utils.Compare(storedFile, file)
						if err != nil {
							log.Fatalf("[existRepostByFiles] error: %s", err)
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

func doRepost(attachments []string, item *api.Post, group *db.Group) (bool, error) {
	to := group.SourceID
	message := group.Message

	if len(item.Text) > 0 {
		r := regexp.MustCompile(`(\n+|\s+)?#(\p{L}|\p{P})+(\n+|\s+)?`)
		item.Text = r.ReplaceAllString(item.Text, " ")
		r = regexp.MustCompile(`(\n|\s)+`)
		item.Text = r.ReplaceAllString(item.Text, " ")
		message = strings.Trim(item.Text, " ") + " " + message
	}
	repost, err := api.DoPost(to, strings.Join(attachments, ","), url.QueryEscape(message))
	if err != nil {
		return false, err
	}

	if repost.Error.ErrorCode > 0 {
		err = errors.New(repost.Error.ErrorMsg)
	}
	return repost.Response.PostID > 0, err
}

func doRemoveDogs(groupID int) {
	start := 0
	offset := 1000

	totalUsers := 0
	var usersList []int

	for {
		users, err := api.GetListUsersofGroup(groupID, start, offset)
		if err != nil {
			log.Fatalf("[doRemoveDogs] error: %s", err)
			return
		}

		totalUsers += len(users.Response.Items)
		for _, user := range users.Response.Items {
			if user.IsDeleted() || user.IsBanned() {
				usersList = append(usersList, user.ID)
			}
		}
		if len(users.Response.Items) < offset {
			break
		} else {
			start += offset
		}
	}

	percentBadUsers := float32(len(usersList)) / float32(totalUsers)
	if percentBadUsers > minPercentOfBadUsers {
		count := int(float32(totalUsers) * percentBadUsers)
		for index := 0; index < count; index++ {
			userID := usersList[index]
			status, err := api.RemoveUserFromGroup(groupID, userID)
			if err != nil {
				log.Fatalf("[doRemoveDogs] error: %s", err)
				return
			}
			if status.Response != 1 {
				break
			}
			log.Printf("[doRemoveDogs] user_id: %d group_id: %d", userID, groupID)
		}
	}
}

func shouldDoRepost(info api.Group, group db.Group, item api.Post, countCheckIn *int) bool {
	postID := getPostID(&info, &item)
	post, err := db.PostQuery()
	if err != nil {
		log.Fatalf("[main:db.PostQuery] error: %s", err)
		return false
	}

	from := info.ID
	to := group.SourceID
	record := &db.Post{
		Post:  postID,
		Files: nil,
		From:  from,
		To:    to,
		Date:  time.Now()}
	err = post.Insert(record)
	if err != nil {
		log.Fatalf("[main:post.Insert] error: %s", err)
		return false
	}

	files, attachments := item.GetUniqueFiles()
	if files != nil && existRepostByFiles(files) {
		files = nil
		attachments = nil
	}

	reposted := false
	if attachments != nil {
		copier.UploadFiles(attachments, to)

		reposted, err = doRepost(attachments, &item, &group)
		if err == nil {
			record.Files = files
			err = post.Update(bson.M{"post": postID}, record)
			if err != nil {
				log.Fatalf("[main:post.Update] error: %s", err)
				return false
			}
		} else {
			err = post.Remove(bson.M{"post": postID})
		}

		if err != nil {
			log.Fatalf("[main:doRepost] error: %s post_id: %s", err, postID)
			return false
		}
	}

	if reposted {
		*countCheckIn++
		log.Printf("[main] reposted: %s", postID)
	} else {
		log.Printf("[main] skipped: %s", postID)
	}
	if *countCheckIn == maxCountCheckInOneTime {
		log.Printf("[main] interrupted: %s", postID)
		return false
	}

	return true
}

func checkSources(info api.Group, group db.Group, countCheckIn *int) {
	posts, err := api.GetPosts(strconv.Itoa(info.ID), "50")
	if err != nil {
		log.Fatalf("[main:api.GetPosts] error: %s group_id: %d", err, info.ID)
		return
	}

	border := int(posts.GetMaxCountLikes() * group.Border)
	for _, item := range posts.Response.Items {
		// skip posts with links
		r := regexp.MustCompile(`.*\[club\d+\|.*`)
		if r.MatchString(item.Text) {
			continue
		}
		if item.IsPinned == 0 && item.Likes.Count > border {
			if !existRepostByID(&info, &item) {
				repostedSuccess := shouldDoRepost(info, group, item, countCheckIn)
				if !repostedSuccess {
					os.Exit(0)
				}
			}
		}
	}
}

func checkDestination(group db.Group, countCheckIn *int) {
	groupsInfo, err := api.GetGroupsInfo(strconv.Itoa(group.SourceID), "links")
	if err != nil {
		log.Fatalf("[main:api.GetGroupsInfo] error: %s group_id: %d", err, group.SourceID)
		return
	}

	var ids []string
	for _, info := range groupsInfo.Response {
		r := regexp.MustCompile("https://vk.com/(.*?)$")
		for _, link := range info.Links {
			ids = append(ids, r.FindStringSubmatch(link.URL)[1])
		}
	}

	groupsInfo, err = api.GetGroupsInfo(strings.Join(ids, ","), "")
	if err != nil {
		log.Fatalf("[main:api.GetGroupsInfo] error: %s", err)
		return
	}

	for _, info := range groupsInfo.Response {
		checkSources(info, group, countCheckIn)
	}
}

const maxCountCheckInOneTime int = 2
const minPercentOfBadUsers float32 = 0.01

func main() {
	clientID := os.Getenv("CLIENT_ID")
	email := os.Getenv("CLIENT_EMAIL")
	password := os.Getenv("CLIENT_PASSWORD")
	dbServerAddress := os.Getenv("DB_SERVER")
	countCheckIn := 0

	log.Println("start")
	_, err := api.GetAccessToken(clientID, email, password)
	if err != nil {
		log.Fatalf("[main:api.GetAccessToken] error: %s", err)
		return
	}

	session, err := db.Connect(dbServerAddress)
	if err != nil {
		log.Fatalf("[main:db.Connect] error: %s", err)
		return
	}
	defer session.Close()

	groups, err := db.GetGroups()
	if err != nil {
		log.Fatalf("[main:db.GetGroups] error: %s", err)
		return
	}
	for _, group := range groups {
		go doRemoveDogs(group.SourceID)
		checkDestination(group, &countCheckIn)
	}

	log.Println("done")
}
