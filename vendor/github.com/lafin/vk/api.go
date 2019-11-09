package vk

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"

	"github.com/lafin/go-httpclient"
)

var accessToken string

// GetAccessToken - get access toket for authorize on the vk.com
func GetAccessToken(clientID, email, pass string) (string, error) {
	if len(accessToken) > 0 {
		return accessToken, nil
	}

	client := httpclient.Client()
	data, err := httpclient.GetData(AuthURL + "/authorize?client_id=" + clientID + "&redirect_uri=https://oauth.vk.com/blank.html&display=mobile&scope=" + Permissions + "&response_type=token&v=" + APIVersion)
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile("<form method=\"post\" action=\"(.*?)\">")
	match := r.FindStringSubmatch(string(data))
	urlStr := match[1]

	r = regexp.MustCompile("<input type=\"hidden\" name=\"(.*?)\" value=\"(.*?)\" ?/?>")
	matches := r.FindAllStringSubmatch(string(data), -1)

	formData := url.Values{}
	for _, val := range matches {
		formData.Add(val[1], val[2])
	}
	formData.Add("email", email)
	formData.Add("pass", pass)

	response, err := client.PostForm(urlStr, formData)
	if err != nil {
		return "", err
	}

	r = regexp.MustCompile("__q_hash=.*?")
	if r.MatchString(response.Request.URL.String()) {
		data, err := httpclient.GetData(response.Request.URL.String())
		if err != nil {
			return "", err
		}

		r := regexp.MustCompile("<form method=\"post\" action=\"(.*?)\">")
		match := r.FindStringSubmatch(string(data))
		response, err = client.PostForm(match[1], url.Values{})
		if err != nil {
			return "", err
		}
	}

	r = regexp.MustCompile("access_token=(.*?)&")
	match = r.FindStringSubmatch(response.Request.URL.String())
	if len(match) > 0 {
		accessToken = match[1]
		return accessToken, nil
	}

	return "", errors.New("can't find the access_token")
}

// GetPosts - get list of posts
func GetPosts(groupID, count string) (*Posts, error) {
	data, err := httpclient.GetData(APIURL + "/method/wall.get?owner_id=-" + groupID + "&count=" + count + "&filter=all&access_token=" + accessToken + "&v=" + APIVersion)
	if err != nil {
		return nil, err
	}

	var posts Posts
	if err := json.Unmarshal(data, &posts); err != nil {
		return nil, err
	}
	return &posts, nil
}

// GetGroupsInfo - get group info
func GetGroupsInfo(groupIDs, fields string) (*Groups, error) {
	data, err := httpclient.GetData(APIURL + "/method/groups.getById?group_ids=" + groupIDs + "&fields=" + fields + "&access_token=" + accessToken + "&v=" + APIVersion)
	if err != nil {
		return nil, err
	}

	var groups Groups
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, err
	}
	return &groups, nil
}

// DoRepost - do repost the post
func DoRepost(object string, groupID int, message string) (*ResponseRepost, error) {
	data, err := httpclient.GetData(APIURL + "/method/wall.repost?object=" + object + "&group_id=" + strconv.Itoa(groupID) + "&message=" + message + "&access_token=" + accessToken + "&v=" + APIVersion)
	if err != nil {
		return nil, err
	}

	var repost ResponseRepost
	if err := json.Unmarshal(data, &repost); err != nil {
		return nil, err
	}
	return &repost, nil
}

// DoPost - do post the post
func DoPost(groupID int, attachments, message string) (*ResponsePost, error) {
	data, err := httpclient.GetData(APIURL + "/method/wall.post?owner_id=-" + strconv.Itoa(groupID) + "&from_group=1&mark_as_ads=0&attachments=" + attachments + "&message=" + message + "&access_token=" + accessToken + "&v=" + APIVersion)
	if err != nil {
		return nil, err
	}

	var post ResponsePost
	if err := json.Unmarshal(data, &post); err != nil {
		return nil, err
	}
	return &post, nil
}

// GetMaxCountLikes - return max likes in list of posts
func (p *Posts) GetMaxCountLikes() float32 {
	max := 0
	for _, item := range p.Response.Items {
		if item.Likes.Count > max && item.IsPinned == 0 {
			max = item.Likes.Count
		}
	}
	return float32(max)
}

// GetUniqueFiles - get lists of files
func (p *Post) GetUniqueFiles() ([][]byte, []string) {
	var attachments = make([]string, 0, len(p.Attachments))
	var files [][]byte

	for _, item := range p.Attachments {
		var err error
		var file []byte
		var attachment string
		var photoURL string

		switch item.Type {
		case "photo":
			for _, size := range item.Photo.Sizes {
				if size.Type == "s" {
					photoURL = size.URL
					attachment = item.Type + strconv.Itoa(item.Photo.OwnerID) + "_" + strconv.Itoa(item.Photo.ID) + "_" + item.Photo.AccessKey
					break
				}
			}
		case "doc":
			if len(item.Doc.URL) > 0 {
				photoURL = item.Doc.Preview.GetSmallPreview()
				attachment = item.Type + strconv.Itoa(item.Doc.OwnerID) + "_" + strconv.Itoa(item.Doc.ID) + "_" + item.Doc.AccessKey
			}
		case "video":
			attachment = item.Type + strconv.Itoa(item.Video.OwnerID) + "_" + strconv.Itoa(item.Video.ID) + "_" + item.Video.AccessKey
		}
		if photoURL != "" {
			file, err = httpclient.GetData(photoURL)
		}
		if err != nil {
			return nil, nil
		}
		if file != nil {
			files = append(files, file)
		}
		if len(attachment) > 0 {
			attachments = append(attachments, attachment)
		}
	}
	return files, attachments
}

// GetListUsersofGroup - get list deactivated users
func GetListUsersofGroup(groupID int, offset int, count int) (*ResponseUsersOfGroup, error) {
	data, err := httpclient.GetData(APIURL + "/method/groups.getMembers?group_id=" + strconv.Itoa(groupID) + "&offset=" + strconv.Itoa(offset) + "&count=" + strconv.Itoa(count) + "&fields=last_seen&access_token=" + accessToken + "&v=" + APIVersion)
	if err != nil {
		return nil, err
	}

	var users ResponseUsersOfGroup
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	return &users, nil
}

// RemoveUserFromGroup - remove user from group
func RemoveUserFromGroup(groupID int, userID int) (*ResponseRemoveUser, error) {
	data, err := httpclient.GetData(APIURL + "/method/groups.removeUser?group_id=" + strconv.Itoa(groupID) + "&user_id=" + strconv.Itoa(userID) + "&access_token=" + accessToken + "&v=" + APIVersion)
	if err != nil {
		return nil, err
	}

	var status ResponseRemoveUser
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, err
	}
	return &status, nil
}
