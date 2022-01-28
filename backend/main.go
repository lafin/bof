package main

import (
	"errors"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	api "github.com/lafin/vk"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

func existRepostByID(dbConnect *gorm.DB, postID string) bool {
	result := dbConnect.Where("post = ?", postID).Take(&Post{})
	return result.Error == nil
}

func existRepostByFiles(dbConnect *gorm.DB, files [][]byte) bool {
	records := []Post{}
	result := dbConnect.Where("files IS NOT NULL").Order("created_at").Find(&records)
	if result.Error != nil {
		return false
	}

	for _, record := range records {
		for _, storedFile := range record.Files {
			if len(storedFile) != 0 {
				for _, file := range files {
					if len(file) != 0 {
						percent, err := Compare(storedFile, file)
						if err != nil {
							taskErrors.With(prometheus.Labels{"error": "exist_repost_by_files"}).Inc()
							pushMetrics()
							l.Logf("FATAL existRepostByFiles, %v", err)
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

func doRepost(attachments []string, item *api.Post, group *Group) (bool, error) {
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

func doRemoveDogs(dbConnect *gorm.DB, groupID int) {
	start := 0
	offset := 1000
	for {
		users, err := api.GetListUsersofGroup(groupID, start, offset)
		if err != nil {
			taskErrors.With(prometheus.Labels{"error": "get_list_usersof_group"}).Inc()
			pushMetrics()
			l.Logf("FATAL api.GetListUsersofGroup, %v", err)
		}
		for _, user := range users.Response.Items {
			if !(user.IsDeleted() || user.IsBanned()) {
				continue
			}
			dog := Dog{}
			result := dbConnect.Where("source_id = ? AND user_id = ?", groupID, user.ID).First(&dog)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				result = dbConnect.Create(&Dog{
					SourceID:  groupID,
					UserID:    user.ID,
					CheckedAt: time.Now(),
				})
			} else {
				dog.CheckedAt = time.Now()
				result = dbConnect.Where("source_id = ? AND user_id = ?", groupID, user.ID).Updates(&dog)
			}
			if result.Error != nil {
				taskErrors.With(prometheus.Labels{"error": "do_remove_dogs"}).Inc()
				pushMetrics()
				l.Logf("FATAL doRemoveDogs, %v", result.Error)
			}
		}
		if len(users.Response.Items) < offset {
			break
		} else {
			start += offset
		}
	}

	dbConnect.Unscoped().Where("checked_at < NOW() - INTERVAL '1 day'").Delete(&Dog{})

	dogs := []Dog{}
	result := dbConnect.Where("created_at < NOW() - INTERVAL '1 month'").Find(&dogs)
	if result.Error != nil {
		taskErrors.With(prometheus.Labels{"error": "find_dogs"}).Inc()
		pushMetrics()
		l.Logf("FATAL findDogs, %v", result.Error)
	}
	for _, dog := range dogs {
		status, err := api.RemoveUserFromGroup(dog.SourceID, dog.UserID)
		if err != nil {
			taskErrors.With(prometheus.Labels{"error": "remove_user_from_group"}).Inc()
			pushMetrics()
			l.Logf("FATAL api.RemoveUserFromGroup, %v", err)
		}
		if status.Response != 1 {
			break
		}
		l.Logf("INFO doRemoveDogs group_id: %d user_id: %d", dog.SourceID, dog.UserID)
		dbConnect.Unscoped().Where("source_id = ? AND user_id = ?", dog.SourceID, dog.UserID).Delete(&Dog{})
	}
}

func shouldDoRepost(dbConnect *gorm.DB, info api.Group, group Group, item api.Post, countCheckIn *int) bool {
	postID := getPostID(&info, &item)
	from := info.ID
	to := group.SourceID
	record := &Post{
		Post:  postID,
		Files: nil,
		From:  from,
		To:    to,
	}
	result := dbConnect.Create(record)
	if result.Error != nil {
		taskErrors.With(prometheus.Labels{"error": "insert"}).Inc()
		pushMetrics()
		l.Logf("FATAL post.Insert, %v", result.Error)
	}

	files, attachments := item.GetUniqueFiles()
	if files != nil && existRepostByFiles(dbConnect, files) {
		files = nil
		attachments = nil
	}

	reposted := false
	if attachments != nil {
		if len(attachments) == 0 {
			taskErrors.With(prometheus.Labels{"error": "get_unique_files"}).Inc()
			pushMetrics()
			l.Logf("FATAL item.GetUniqueFiles, empty attachments")
		}
		var err error
		reposted, err = doRepost(attachments, &item, &group)
		if err == nil {
			record.Files = files
			result = dbConnect.Where("post = ?", postID).Updates(&record)
			if result.Error != nil {
				taskErrors.With(prometheus.Labels{"error": "update"}).Inc()
				pushMetrics()
				l.Logf("FATAL post.Update, %v", result.Error)
			}
		} else {
			result = dbConnect.Unscoped().Where("post = ?", postID).Delete(&Post{})
			err = result.Error
		}

		if err != nil {
			taskErrors.With(prometheus.Labels{"error": "do_repost"}).Inc()
			pushMetrics()
			l.Logf("FATAL doRepost, %v post_id: %s", err, postID)
		}
	}

	if reposted {
		*countCheckIn++
		l.Logf("INFO reposted post_id: %s", postID)
	} else {
		l.Logf("INFO skipped post_id: %s", postID)
	}
	if *countCheckIn == maxCountCheckInOneTime {
		l.Logf("INFO interrupted post_id: %s", postID)
		return false
	}

	return true
}

func checkSources(dbConnect *gorm.DB, info api.Group, group Group, countCheckIn *int) {
	posts, err := api.GetPosts(strconv.Itoa(info.ID), "50")
	if err != nil {
		taskErrors.With(prometheus.Labels{"error": "get_posts"}).Inc()
		pushMetrics()
		l.Logf("FATAL api.GetPosts, %v group_id: %d", err, info.ID)
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
					return
				}
			}
		}
	}
}

func checkDestination(dbConnect *gorm.DB, group Group, countCheckIn *int) {
	groupsInfo, err := api.GetGroupsInfo(strconv.Itoa(group.SourceID), "links")
	if err != nil {
		taskErrors.With(prometheus.Labels{"error": "get_groups_info"}).Inc()
		pushMetrics()
		l.Logf("FATAL api.GetGroupsInfo, %v group_id: %d", err, group.SourceID)
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
		taskErrors.With(prometheus.Labels{"error": "get_groups_info"}).Inc()
		pushMetrics()
		l.Logf("FATAL api.GetGroupsInfo, %v", err)
	}

	for _, info := range groupsInfo.Response {
		checkSources(dbConnect, info, group, countCheckIn)
	}
}

const maxCountCheckInOneTime int = 2

func main() {
	_ = godotenv.Load()
	countCheckIn := 0
	metrics(os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_JOB"))

	_, err := api.GetAccessToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_EMAIL"), os.Getenv("CLIENT_PASSWORD"))
	if err != nil {
		taskErrors.With(prometheus.Labels{"error": "get_access_token"}).Inc()
		pushMetrics()
		l.Logf("FATAL api.GetAccessToken, %v", err)
	}

	dbConnect := Connect(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	dbConnect.Unscoped().Where("created_at < NOW() - INTERVAL '1 month'").Delete(&Post{})
	groups, err := GetGroups(dbConnect)
	if err != nil {
		taskErrors.With(prometheus.Labels{"error": "get_groups"}).Inc()
		pushMetrics()
		l.Logf("FATAL db.GetGroups error, %v", err)
	}
	for _, group := range groups {
		doRemoveDogs(dbConnect, group.SourceID)
		checkDestination(dbConnect, group, &countCheckIn)
	}

	pushMetrics()
}
