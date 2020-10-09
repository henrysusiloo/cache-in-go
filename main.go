package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

const cacheKey = "github"

var c *cache.Cache

func github(w http.ResponseWriter, r *http.Request) {
	dataItf := getCache()
	if dataItf != nil {
		dataStr := fmt.Sprintf("%v", dataItf)
		json.NewEncoder(w).Encode(dataStr)
		fmt.Println("get cache")
		return
	}

	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.NewEncoder(w).Encode(string(data))
	setCache(string(data))
	fmt.Println("get api")
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	initCache()
	handleRequests()
	flushCache()
}

func initCache() {
	// defaultExpiration is 5 minutes
	// cleanupInterval is 10 minutes, expired items will be deleted every 10 minutes
	c = cache.New(5*time.Minute, 10*time.Minute)
}

func setCache(value interface{}) {
	// duration is 0, defaultExpiration will be used
	c.Set(cacheKey, value, 0)
}

func getCache() interface{} {
	value, ok := c.Get(cacheKey)
	if ok {
		return value
	}

	return nil
}

func flushCache() {
	c.Flush()
}
