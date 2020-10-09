package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pmylund/go-cache"
)

var (
	c = cache.New(5*time.Second, 10*time.Second)
)

var (
	rds   redis.Conn
	err   error
	reply interface{}
)

func init() {
	rds, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
}

func github(w http.ResponseWriter, r *http.Request) {

	// get from memcache
	memcacheData, found := c.Get("github-status")
	if found {
		json.NewEncoder(w).Encode(fmt.Sprintf("%v", memcacheData))
		return
	}

	// get from redis
	if redisData, err := redis.String(rds.Do("GET", "github-status")); err == nil {
		c.Set("github-status", redisData, cache.DefaultExpiration)
		json.NewEncoder(w).Encode(redisData)
		return
	}

	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)

	// set cache
	c.Set("github-status", string(data), cache.DefaultExpiration)
	_, err = rds.Do("SETEX", "github-status", 60, string(data))
	if err != nil {
		fmt.Println(err)
	}

	resp.Body.Close()
	json.NewEncoder(w).Encode(string(data))
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
