package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
)

const cacheKey = "my-cache-key"

var cache *bigcache.BigCache

func main() {
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	fmt.Println("running")
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func github(w http.ResponseWriter, r *http.Request) {
	entry, err := cache.Get(cacheKey)
	if err == nil {
		fmt.Println("get from cache")
		json.NewEncoder(w).Encode(string(entry))
		return
	}

	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	cache.Set(cacheKey, []byte(data))
	json.NewEncoder(w).Encode(string(data))
}
