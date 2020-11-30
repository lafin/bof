package main

import (
	"errors"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lafin/bof/db"
	"github.com/lafin/bof/misc"
	api "github.com/lafin/vk"
	"gorm.io/gorm"
)

func existRepostByID(dbConnect *gorm.DB, postID string) bool {
	result := dbConnect.Where("post = ?", postID).Take(&db.Post{})
	return result.Error == nil
}

func existRepostByFiles(dbConnect *gorm.DB, files [][]byte) bool {
	records := []db.Post{}
	result := dbConnect.Where("files IS NOT NULL").Order("created_at").Find(&records)
	if result.Error != nil {
		return false
	}

	for _, record := range records {
		for _, storedFile := range record.Files {
			if len(storedFile) != 0 {
				for _, file := range files {
					if len(file) != 0 {
						percent, err := misc.Compare(storedFile, file)
						if err != nil {
							log.Fatalf("[existRepostByFiles] error: %s", err)
						} else if percent < 0.05 {
							return true
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

func shouldDoRepost(dbConnect *gorm.DB, info api.Group, group db.Group, item api.Post, countCheckIn *int) bool {
	postID := getPostID(&info, &item)
	from := info.ID
	to := group.SourceID
	record := &db.Post{
		Post:  postID,
		Files: nil,
		From:  from,
		To:    to,
	}
	result := dbConnect.Create(record)
	if result.Error != nil {
		log.Fatalf("[main:post.Insert] error: %s", result.Error)
		return false
	}

	files, attachments := item.GetUniqueFiles()
	if files != nil && existRepostByFiles(dbConnect, files) {
		files = nil
		attachments = nil
	}

	reposted := false
	if attachments != nil {
		if len(attachments) == 0 {
			log.Fatalf("[main:item.GetUniqueFiles] error: empty attachments")
			return false
		}
		var err error
		reposted, err = doRepost(attachments, &item, &group)
		if err == nil {
			record.Files = files
			result = dbConnect.Where("post = ?", postID).Updates(&record)
			if result.Error != nil {
				log.Fatalf("[main:post.Update] error: %s", result.Error)
				return false
			}
		} else {
			result = dbConnect.Unscoped().Where("post = ?", postID).Delete(&db.Post{})
			err = result.Error
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

func checkSources(dbConnect *gorm.DB, info api.Group, group db.Group, countCheckIn *int) {
	posts, err := api.GetPosts(strconv.Itoa(info.ID), "50")
	if err != nil {
		log.Fatalf("[main:api.GetPosts] error: %s group_id: %d", err, info.ID)
		return
	}

	border := int(posts.GetMaxCountLikes() * group.Border)
	for i, item := range posts.Response.Items {
		// skip posts with CopyHistory
		if len(item.CopyHistory) > 0 {
			continue
		}
		// skip posts with links
		r := regexp.MustCompile(`.*\[club\d+\|.*`)
		if r.MatchString(item.Text) {
			continue
		}
		postID := getPostID(&info, &posts.Response.Items[i])
		if item.IsPinned == 0 && item.Likes.Count > border {
			if !existRepostByID(dbConnect, postID) {
				repostedSuccess := shouldDoRepost(dbConnect, info, group, item, countCheckIn)
				if !repostedSuccess {
					os.Exit(0)
				}
			}
		}
	}
}

func checkDestination(dbConnect *gorm.DB, group db.Group, countCheckIn *int) {
	groupsInfo, err := api.GetGroupsInfo(strconv.Itoa(group.SourceID), "links")
	if err != nil {
		log.Fatalf("[main:api.GetGroupsInfo] error: %s group_id: %d", err, group.SourceID)
		return
	}

	var ids []string
	for _, info := range groupsInfo.Response {
		r := regexp.MustCompile(`https://vk\.com/(.*?)$`)
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
		checkSources(dbConnect, info, group, countCheckIn)
	}
}

const maxCountCheckInOneTime int = 2
const minPercentOfBadUsers float32 = 0.01

func main() {
	_ = godotenv.Load()
	countCheckIn := 0

	log.Println("start")
	_, err := api.GetAccessToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_EMAIL"), os.Getenv("CLIENT_PASSWORD"))
	if err != nil {
		log.Fatalf("[main:api.GetAccessToken] error: %s", err)
		return
	}

	dbConnect := db.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	groups, err := db.GetGroups(dbConnect)
	if err != nil {
		log.Fatalf("[main:db.GetGroups] error: %s", err)
		return
	}
	for _, group := range groups {
		go doRemoveDogs(group.SourceID)
		checkDestination(dbConnect, group, &countCheckIn)
	}

	log.Println("done")
}
