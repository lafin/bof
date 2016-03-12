package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"
    "errors"
)

func parseJson(body []byte) map[string]interface{} {
    var data map[string]interface{}
    if err := json.Unmarshal(body, &data); err != nil {
        panic(err)
    }
    return data
}

func getData(url string) map[string]interface{} {
    var result map[string]interface{}

    response, err := http.Get(url)
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
        result = parseJson(body)
    }
    return result
}

func getAccessToken(clientId, clientSecret string) string {
    data := getData("https://oauth.vk.com/access_token?client_id="+clientId+"&client_secret="+clientSecret+"&grant_type=client_credentials")
    return data["access_token"].(string)
}

func getPosts(args ...interface{}) map[string]interface{} {
    var err error
    var domain string = ""
    var count string = "10"

    for i, p := range args {
        switch i {
            case 0:
                val, ok := p.(string)
                if ok {
                    domain = val
                } else {
                    err = errors.New("Parameter not type string.")
                }
            case 1:
                val, ok := p.(string)
                if ok {
                    count = val
                } else {
                    err = errors.New("Parameter not type string.")
                }
            default:
                err = errors.New("Wrong parameter count.")
        }
    }

    if err != nil {
        fmt.Println(err)
        return nil
    }

    return getData("https://api.vk.com/method/wall.get?&domain="+domain+"&count="+count+"&filter=all")
}

func main() {
    clientId := os.Getenv("CLIENT_ID")
    clientSecret := os.Getenv("CLIENT_SECRET")
    accessToken := getAccessToken(clientId, clientSecret)
    fmt.Printf("%s\n", accessToken)

    posts := getPosts("smcat")
    fmt.Println(posts)
}