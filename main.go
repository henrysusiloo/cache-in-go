package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
)

var cache *bigcache.BigCache

func github(w http.ResponseWriter, r *http.Request) {
	cacheResp, err := getFromCache()
	if err == nil {
		json.NewEncoder(w).Encode(string(cacheResp))
	}

	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)

	// setToCache(data)
	resp.Body.Close()
	json.NewEncoder(w).Encode(string(data))
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func getFromCache() ([]byte, error) {
	data, err := cache.Get("github")
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func setToCache(data []byte) {
	cache.Set("github", data)
}

func initCache() {
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
}

func main() {
	initCache()
	handleRequests()
}
