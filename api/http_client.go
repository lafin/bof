package api

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"sync"
)

var once sync.Once
var client *http.Client

// GetData - get data by an url
func GetData(url string) ([]byte, error) {
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
		transport := &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		}
		cookieJar, _ := cookiejar.New(nil)
		client = &http.Client{Transport: transport, Jar: cookieJar}
	})

	return client
}
