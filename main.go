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

var (
	c *cache.Cache
)

func github(w http.ResponseWriter, r *http.Request) {
	status, ok := getCache("github_status")
	if ok {
		fmt.Println("Lewat cache")
		json.NewEncoder(w).Encode(status)
		return
	}

	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	setCache("github_status", string(data))

	fmt.Println("Ga lewat cache")
	json.NewEncoder(w).Encode(string(data))
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	initCache()
	handleRequests()
}

func initCache() {
	c = cache.New(time.Minute*5, time.Minute*5)
}

func setCache(cacheKey string, value interface{}) {
	c.Set(cacheKey, value, 0)
}

func getCache(cacheKey string) (interface{}, bool) {
	return c.Get(cacheKey)
}
