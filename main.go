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

const CACHE_KEY = "github:status"

var c *cache.Cache

func githubResource() (string, error) {
	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		return "", err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	c.Set(string(data), CACHE_KEY, cache.DefaultExpiration)

	return string(data), nil
}

func githubCache() (string, bool) {
	data, found := c.Get(CACHE_KEY)
	fmt.Println(data, found)
	if found {
		return data.(string), true
	}
	return "", false
}

func github(w http.ResponseWriter, r *http.Request) {
	data, found := githubCache()
	if found {
		fmt.Println("from cache")
		json.NewEncoder(w).Encode(data)
		return
	}

	data, err := githubResource()
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	fmt.Println("from resource")
	json.NewEncoder(w).Encode(data)
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	c = cache.New(1*time.Minute, 10*time.Minute)
	handleRequests()
}
