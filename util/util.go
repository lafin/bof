package util

import (
	"crypto/sha1"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var once sync.Once
var client *http.Client

// Client - get instance of http client
func Client() *http.Client {
	once.Do(func() {
		transport := &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		}
		client = &http.Client{Transport: transport}
	})

	return client
}

// GetSha1 - return sha1 hashsum for file by url
func GetSha1(url string) []byte {
	client := Client()
	res, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	hasher := sha1.New()
	hasher.Write([]byte(body))

	return hasher.Sum(nil)
}