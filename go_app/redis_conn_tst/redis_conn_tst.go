package main

import "github.com/garyburd/redigo/redis"
import (
        "fmt"
        "strconv"
)

const REDIS_SERVER_URL = "localhost"
const REDIS_SERVER_PORT = 9999
const REDIS_PASSWORD = "redis_p@ssw0rd"

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

        fmt.Printf("Connected!", response)

         //Set the message in Message Queue
        conn.Do("SET", "message1", "Hello World")

        //Get the message from Message Queue
        message, err := redis.String(conn.Do("GET", "message1"))
        if err != nil {
                fmt.Println("Key Not found")
        }
        fmt.Printf("message1", message)

        defer conn.Close()
}

