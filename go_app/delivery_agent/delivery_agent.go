package main

import "github.com/garyburd/redigo/redis"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"postbackdelivery/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const REDIS_SERVER_URL = "localhost"
const REDIS_SERVER_PORT = 9999
const REDIS_PASSWORD = "redis_p@ssw0rd"
const REDIS_MESSAGE_QUEUE = "request"
const RETRY_LIMIT = 5

var REDIS_CON_STRING = REDIS_SERVER_URL + ":" + strconv.Itoa(REDIS_SERVER_PORT)

func main() {
	//Setup Logger
	utils.SetupLogger()

	//Connect the Redis Server
	conn, err := redis.Dial("tcp", REDIS_CON_STRING)
	if err != nil {
		panic(err)
		utils.Info.Println(err)
	}

	response, err := conn.Do("Auth", REDIS_PASSWORD)
	if err != nil {
		panic(err)
		utils.Info.Println(err)
	}
	fmt.Println("Connected to Redis Server")
	utils.Info.Println("\nConnected to Redis Server", response, "\n")

	for {
		//Get the message from Message Queue
		message, _ := redis.String(conn.Do("RPOP", REDIS_MESSAGE_QUEUE))
		//message,_ := redis.String(conn.Do("BRPOP", REDIS_MESSAGE_QUEUE, 0))

		if message != "" {
			utils.Info.Println("\nGot PostBack Obect Request:\n")

			var request map[string]interface{}
			if err := json.Unmarshal([]byte(message), &request); err != nil {
				panic(err)
			}
			utils.PostBack.Println("\nRequest Time:", time.Now(), "\nPostBack Request:\n", message, "\n")

			endpoint := request["endpoint"].(map[string]interface{})
			url := endpoint["url"].(string)
			method := endpoint["method"].(string)
			data := request["data"].(map[string]interface{})

			var dataMap map[string]string
			dataMap = make(map[string]string)
			for key, value := range data {
				dataMap[key] = value.(string)
			}
			//Construct formatted url replacing placeholder with values from dataMap
			formattedUrl := ConstructFormattedUrl(url, dataMap)

			//Send Response
			if strings.ToUpper(method) == "GET" {
				utils.Info.Println("\nSending HTTP GET request:\n", formattedUrl, "\n")
				SendHttpGetRequest(formattedUrl)
			} else if strings.ToUpper(method) == "POST" {
				utils.Info.Println("\nSending HTTP POST request:\n", formattedUrl, dataMap, "\n")
				SendHttpPostRequest(formattedUrl, dataMap)
			} else {
				utils.Info.Println("\nHTTP Method Not Supported\n", method, formattedUrl, "\n")
			}
		}

	}

	defer conn.Close()
}

func ConstructFormattedUrl(url string, dataMap map[string]string) string {
	re := regexp.MustCompile("{([a-zA-Z0-9_@]+)}")
	matches := re.FindAllString(url, -1)
	for _, match := range matches {
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
		} else {
			break
		}
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	utils.PostBack.Println("\nResponse Time:", time.Now(), "\nResponse Status Code:", response.Status, "\n Response Body:\n", string(body), "\n")
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
			//fmt.Printf("Post Request:", response)
		} else {
			break
		}
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	//utils.PostBack.Println("\nResponse Time:", time.Now(), "\nResponse Status Code:", response.Status,"\n")
	utils.PostBack.Println("\nResponse Time:", time.Now(), "\nResponse Status Code:", response.Status, "\n Response Body:\n", string(body), "\n")
}
