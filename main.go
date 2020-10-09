package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/coocood/freecache"
)

var cache *freecache.Cache

func github(w http.ResponseWriter, r *http.Request) {

	cacheVal, err := getCache([]byte("status"))
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	if len(cacheVal) != 0 {
		json.NewEncoder(w).Encode(string(cacheVal))
		return
	}

	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	err = saveCache([]byte("status"), data, 10)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(string(data))
}

func initCache() {
	cache = freecache.NewCache(10 * 1024 * 1024)
}

func saveCache(key []byte, value []byte, expireSeconds int) (err error) {
	fmt.Println("[SAVE CACHE]")
	return cache.Set(key, value, expireSeconds)
}

func getCache(key []byte) ([]byte, error) {
	fmt.Println("[GET CACHE]")
	return cache.Get(key)
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {

	initCache()

	handleRequests()
}

/**
	-. data apa yg dibalikin
		curl --location --request GET 'localhost:10000/github'

		pola data:
		- response
			- "{\"message\":\"GitHub lives! (2020-10-09 02:06:14 -0700) (1)\"}"

		- berubah per 10 detik
		- balikan satu response (general)

	-. cache yang mana
		- freecache

	- keynya by apa
		- public key

	-. kenapa pilih cache itu
		-

**/
