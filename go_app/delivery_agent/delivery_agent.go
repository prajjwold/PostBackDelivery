package main

import "github.com/garyburd/redigo/redis"
import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const REDIS_SERVER_URL = "localhost"
const REDIS_SERVER_PORT = 9999
const REDIS_PASSWORD = "redis_p@ssw0rd"
const REDIS_MESSAGE_QUEUE = "request"
const RETRY_LIMIT = 5

var REDIS_CON_STRING = REDIS_SERVER_URL + ":" + strconv.Itoa(REDIS_SERVER_PORT)

func main() {
	//Connect the Redis Server
	conn, err := redis.Dial("tcp", REDIS_CON_STRING)
	if err != nil {
		panic(err)
	}

	response, err := conn.Do("Auth", REDIS_PASSWORD)

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Redis Server", response)

	//for {
	//Get the message from Message Queue
	message, err := redis.String(conn.Do("RPOP", REDIS_MESSAGE_QUEUE))
	if err != nil {
		fmt.Println("Key Not found")
	}
	fmt.Println("Got Request:\n", message)

	var request map[string]interface{}
	if err := json.Unmarshal([]byte(message), &request); err != nil {
		panic(err)
	}
	fmt.Println(request)

	endpoint := request["endpoint"].(map[string]interface{})
	url := endpoint["url"].(string)
	method := endpoint["method"].(string)

	data := request["data"].(map[string]interface{})

	fmt.Println(url)
	fmt.Println(method)
	fmt.Println(data)

	var dataMap map[string]string
	dataMap = make(map[string]string)
	for key, value := range data {
		dataMap[key] = value.(string)
	}
	//Construct formatted url replacing placeholder with values from dataMap
	formattedUrl := ConstructFormattedUrl(url, dataMap)
	fmt.Println(formattedUrl)

	//Send Response
	if strings.ToUpper(method) == "GET" {
		SendHttpGetRequest(url)
	} else if strings.ToUpper(method) == "POST" {
		SendHttpPostRequest(url, dataMap)
	} else {
		fmt.Printf("HTTP METHOD %s not supported for URL %s.\n", method, url)
	}

	//}
	defer conn.Close()
}

func ConstructFormattedUrl(url string, dataMap map[string]string) string {
	re := regexp.MustCompile("{([a-zA-Z0-9_@]+)}")
	matches := re.FindAllString(url, -1)
	for _, match := range matches {
		fmt.Println(match, dataMap[match])
		r := strings.NewReplacer("{", "", "}", "")
		var dataMapkey = r.Replace(match)
		url = strings.Replace(url, match, dataMap[dataMapkey], -1)
	}
	return url
}

func SendHttpGetRequest(url string) {
	response, err := http.Get(url)
	for i := 0; i < RETRY_LIMIT; i++ {
		if err != nil {
			response, err = http.Get(url)
			fmt.Printf("Get Request:", response)
		} else {
			break
		}
	}
}

func SendHttpPostRequest(postUrl string, dataMap map[string]string) {
	data := url.Values{}
	for key, value := range dataMap {
		data.Add(key, value)
	}
	response, err := http.PostForm(postUrl, data)
	for i := 0; i < RETRY_LIMIT; i++ {
		if err != nil {
			response, err = http.PostForm(postUrl, data)
			fmt.Printf("Post Request:", response)
		} else {
			break
		}
	}
}
