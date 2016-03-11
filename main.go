package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"
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

func getAccessToken() string {
    clientId := os.Getenv("CLIENT_ID")
    clientSecret := os.Getenv("CLIENT_SECRET")
    data := getData("https://oauth.vk.com/access_token?client_id="+clientId+"&client_secret="+clientSecret+"&grant_type=client_credentials")
    return data["access_token"].(string)
}

func getPosts() map[string]interface{} {
    return getData("https://api.vk.com/method/wall.get?&domain=smcat&count=10&filter=all")
}

func main() {
    accessToken := getAccessToken()
    fmt.Printf("%s\n", accessToken)
    fmt.Println(getPosts())
}