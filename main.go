package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var gocache *cache.Cache

func github(w http.ResponseWriter, r *http.Request) {
	res, found := gocache.Get("status")
	fmt.Println("res, found", res, found)
	if found {
		json.NewEncoder(w).Encode(res)
		return
	}

	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println("SET CACHE", string(data))
	gocache.Set("status", string(data), 10*time.Second)
	json.NewEncoder(w).Encode(string(data))
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	gocache = cache.New(5*time.Second, 10*time.Second)
	fmt.Println("RUNNING")
	handleRequests()
}
