package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "os"
    "regexp"
)

const (
    AUTH_URL = "https://oauth.vk.com"
    API_URL = "https://api.vk.com"
    API_VERSION = "5.50"
)

type Posts struct {
    Response struct {
        Count int `json:"count"`
        Items []struct {
            Attachments []struct {
                Photo struct {
                    AccessKey string `json:"access_key"`
                    AlbumID   int    `json:"album_id"`
                    Date      int    `json:"date"`
                    Height    int    `json:"height"`
                    ID        int    `json:"id"`
                    OwnerID   int    `json:"owner_id"`
                    Photo130  string `json:"photo_130"`
                    Photo604  string `json:"photo_604"`
                    Photo75   string `json:"photo_75"`
                    Text      string `json:"text"`
                    UserID    int    `json:"user_id"`
                    Width     int    `json:"width"`
                } `json:"photo"`
                Type string `json:"type"`
            } `json:"attachments"`
            Comments struct {
                CanPost int `json:"can_post"`
                Count   int `json:"count"`
            } `json:"comments"`
            Date     int `json:"date"`
            FromID   int `json:"from_id"`
            ID       int `json:"id"`
            IsPinned int `json:"is_pinned"`
            Likes    struct {
                CanLike    int `json:"can_like"`
                CanPublish int `json:"can_publish"`
                Count      int `json:"count"`
                UserLikes  int `json:"user_likes"`
            } `json:"likes"`
            OwnerID    int `json:"owner_id"`
            PostSource struct {
                Type string `json:"type"`
            } `json:"post_source"`
            PostType string `json:"post_type"`
            Reposts  struct {
                Count        int `json:"count"`
                UserReposted int `json:"user_reposted"`
            } `json:"reposts"`
            Text string `json:"text"`
        } `json:"items"`
    } `json:"response"`
}

type Repost struct {
    Response struct {
        LikesCount   int `json:"likes_count"`
        PostID       int `json:"post_id"`
        RepostsCount int `json:"reposts_count"`
        Success      int `json:"success"`
    } `json:"response"`
}

func getData(client *http.Client, url string) []byte {
    response, err := client.Get(url)
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        body, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s\n", err)
            os.Exit(1)
        }
        return body
    }
    return nil
}

func getAccessToken(client *http.Client, clientId, email, pass string) string {
    data := getData(client, AUTH_URL+"/authorize?client_id="+clientId+"&redirect_uri=https://oauth.vk.com/blank.html&display=mobile&scope=wall&v=&response_type=token&v="+API_VERSION)

    r, _ := regexp.Compile("<form method=\"post\" action=\"(.*?)\">")
    match := r.FindStringSubmatch(string(data))
    urlStr := match[1]
    fmt.Println(urlStr)

    r, _ = regexp.Compile("<input type=\"hidden\" name=\"(.*?)\" value=\"(.*?)\" ?/?>")
    matches := r.FindAllStringSubmatch(string(data), -1)
    fmt.Println(matches)

    formData := url.Values{}
    for _, val := range matches {
        formData.Add(val[1], val[2])
    }
    formData.Add("email", email)
    formData.Add("pass", pass)

    response, err := client.PostForm(urlStr, formData)
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    } else {
        r, _ = regexp.Compile("access_token=(.*?)&")
        match = r.FindStringSubmatch(response.Request.URL.String())
        return match[1]
    }

    return ""
}

func getPosts(client *http.Client, domain, count string) *Posts {
    data := getData(client, API_URL+"/method/wall.get?&domain="+domain+"&count="+count+"&filter=all&v="+API_VERSION)

    var posts Posts
    if err := json.Unmarshal(data, &posts); err != nil {
        panic(err)
    }
    return &posts
}

func doRepost(client *http.Client, object, groupId, accessToken string) *Repost {
    data := getData(client, API_URL+"/method/wall.repost?&object="+object+"&group_id="+groupId+"&access_token="+accessToken+"&v="+API_VERSION)

    var repost Repost
    if err := json.Unmarshal(data, &repost); err != nil {
        panic(err)
    }
    return &repost
}

func main() {
    clientId := os.Getenv("CLIENT_ID")
    email := os.Getenv("CLIENT_EMAIL")
    password := os.Getenv("CLIENT_PASSWORD")

    cookieJar, _ := cookiejar.New(nil)
    client := &http.Client{
        Jar: cookieJar,
    }
    accessToken := getAccessToken(client, clientId, email, password)
    fmt.Printf("%s\n", accessToken)

    posts := getPosts(client, "smcat", "10")
    fmt.Println(posts.Response.Items)

    repost := doRepost(client, "wall85635407_3133", "", accessToken)
    fmt.Println(repost.Response.Success)
}