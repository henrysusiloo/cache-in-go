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
	cacheSystem *cache.Cache
)

const (
	CacheKeyGithubResponse = "github_response"

	CacheTTLGithubResponse = time.Minute * 3
)

func github(w http.ResponseWriter, r *http.Request) {
	cacheData, isFound := getCache(CacheKeyGithubResponse)
	if isFound {
		json.NewEncoder(w).Encode(cacheData)
		return
	}

	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	setCache(CacheKeyGithubResponse, string(data))
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
	cacheSystem = cache.New(5*time.Minute, 10*time.Minute)
	if cacheSystem == nil {
		log.Fatal("Failed to init cache")
	}
}

func setCache(cacheKey string, data string) {
	cacheSystem.Set(cacheKey, data, CacheTTLGithubResponse)
	log.Println(fmt.Sprintf("Set cache success for key %s", cacheKey)) // log if needed
}

func getCache(cacheKey string) (data interface{}, isFound bool) {
	data, isFound = cacheSystem.Get(cacheKey)

	return data, isFound
}
