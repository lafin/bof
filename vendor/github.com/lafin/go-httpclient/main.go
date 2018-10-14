package httpclient

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"
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
			Dial: (&net.Dialer{
				Timeout: time.Second * 10,
			}).Dial,
			TLSHandshakeTimeout: time.Second * 5,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
			DisableKeepAlives:   true,
		}
		cookieJar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}
		client = &http.Client{
			Timeout:   time.Second * 300,
			Transport: transport,
			Jar:       cookieJar,
		}
	})

	return client
}
