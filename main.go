package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cache-in-go/cache"
)

func github(w http.ResponseWriter, r *http.Request) {
	// get cache
	cacheData, found := cache.GetCache(cache.GithubMessageCacheKey)
	if found && cacheData != "" {
		json.NewEncoder(w).Encode(cacheData)
		return
	}
	fmt.Println("GET GITHUB STATUS")
	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println("DATA: ", string(data))
	// set mapcache
	cache.SetCache(cache.GithubMessageCacheKey, string(data))
	json.NewEncoder(w).Encode(string(data))
}
func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
func main() {
	handleRequests()
}
