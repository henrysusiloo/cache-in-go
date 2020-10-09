package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var c *cache.Cache

func github(w http.ResponseWriter, r *http.Request) {
	cacheData, err := getFromCache()
	if err != nil {
		log.Println("Error when get from cache, err: ", err.Error())
	} else if err == nil {
		json.NewEncoder(w).Encode(string(cacheData))
		return
	}

	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)

	setRespToMemcache("githubresp", data)

	resp.Body.Close()
	json.NewEncoder(w).Encode(string(data))
}

func getFromCache() ([]byte, error) {
	data, found := c.Get("githubresp")
	if !found {
		return []byte{}, errors.New("Memcache key not found")
	}

	dataByte, ok := data.([]byte)
	if !ok {
		return []byte{}, errors.New("Memcache key not found")
	}

	return dataByte, nil
}

func setRespToMemcache(key string, data []byte) {
	c.Set(key, string(data), cache.DefaultExpiration)
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func initCache() {
	c = cache.New(5*time.Minute, 10*time.Minute)
}

func main() {
	initCache()
	handleRequests()
}
