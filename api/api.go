package api

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strconv"
)

var accessToken string

// GetAccessToken - get access toket for authorize on the vk.com
func GetAccessToken(clientID, email, pass string) (string, error) {
	if len(accessToken) > 0 {
		return accessToken, nil
	}

	client := Client()
	data, err := GetData(AuthURL + "/authorize?client_id=" + clientID + "&redirect_uri=https://oauth.vk.com/blank.html&display=mobile&scope=wall&v=&response_type=token&v=" + APIVersion)
	if err != nil {
		return "", err
	}

	r, _ := regexp.Compile("<form method=\"post\" action=\"(.*?)\">")
	match := r.FindStringSubmatch(string(data))
	urlStr := match[1]

	r, _ = regexp.Compile("<input type=\"hidden\" name=\"(.*?)\" value=\"(.*?)\" ?/?>")
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

	r, _ = regexp.Compile("access_token=(.*?)&")
	match = r.FindStringSubmatch(response.Request.URL.String())
	if len(match) > 0 {
		accessToken = match[1]
		return accessToken, nil
	}

	return "", errors.New("can't find the access_token")
}

// GetPosts - get list of posts
func GetPosts(groupID, count string) (*Posts, error) {
	data, err := GetData(APIURL + "/method/wall.get?&owner_id=-" + groupID + "&count=" + count + "&filter=all&v=" + APIVersion)
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
	data, err := GetData(APIURL + "/method/groups.getById?&group_ids=" + groupIDs + "&fields=" + fields + "&v=" + APIVersion)
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
	data, err := GetData(APIURL + "/method/wall.repost?&object=" + object + "&group_id=" + strconv.Itoa(groupID) + "&message=" + message + "&access_token=" + accessToken + "&v=" + APIVersion)
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
	data, err := GetData(APIURL + "/method/wall.post?&owner_id=-" + strconv.Itoa(groupID) + "&from_group=1&mark_as_ads=0&attachments=" + attachments + "&message=" + message + "&access_token=" + accessToken + "&v=" + APIVersion)
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
	var attachments []string
	var attachment string
	var files [][]byte
	var file []byte
	var err error

	for _, item := range p.Attachments {
		switch item.Type {
		case "photo":
			if len(item.Photo.Photo75) > 0 {
				file, err = GetData(item.Photo.Photo75)
				attachment = item.Type + strconv.Itoa(item.Photo.OwnerID) + "_" + strconv.Itoa(item.Photo.ID)
			}
		case "doc":
			if len(item.Doc.URL) > 0 {
				file, err = GetData(item.Doc.URL)
				attachment = item.Type + strconv.Itoa(item.Doc.OwnerID) + "_" + strconv.Itoa(item.Doc.ID)
			}
		}
		if err != nil {
			return nil, nil
		}
		files = append(files, file)
		attachments = append(attachments, attachment)
	}
	return files, attachments
}
