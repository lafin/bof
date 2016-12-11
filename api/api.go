package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"sync"
)

var once sync.Once
var client *http.Client
var accessToken string

func getData(url string) ([]byte, error) {
	client := Client()
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Client - get instance of http client
func Client() *http.Client {
	once.Do(func() {
		cookieJar, _ := cookiejar.New(nil)
		client = &http.Client{
			Jar: cookieJar,
		}
	})

	return client
}

// GetAccessToken - get access toket for authorize on the vk.com
func GetAccessToken(clientID, email, pass string) (string, error) {
	if len(accessToken) > 0 {
		return accessToken, nil
	}

	client := Client()
	data, err := getData(AuthURL + "/authorize?client_id=" + clientID + "&redirect_uri=https://oauth.vk.com/blank.html&display=mobile&scope=wall&v=&response_type=token&v=" + APIVersion)
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
	data, err := getData(APIURL + "/method/wall.get?&owner_id=-" + groupID + "&count=" + count + "&filter=all&v=" + APIVersion)
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
	data, err := getData(APIURL + "/method/groups.getById?&group_ids=" + groupIDs + "&fields=" + fields + "&v=" + APIVersion)
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
	data, err := getData(APIURL + "/method/wall.repost?&object=" + object + "&group_id=" + strconv.Itoa(groupID) + "&message=" + message + "&access_token=" + accessToken + "&v=" + APIVersion)
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
	data, err := getData(APIURL + "/method/wall.post?&owner_id=-" + strconv.Itoa(groupID) + "&from_group=1&mark_as_ads=0&attachments=" + attachments + "&message=" + message + "&access_token=" + accessToken + "&v=" + APIVersion)
	if err != nil {
		return nil, err
	}

	var post ResponsePost
	if err := json.Unmarshal(data, &post); err != nil {
		return nil, err
	}
	return &post, nil
}
