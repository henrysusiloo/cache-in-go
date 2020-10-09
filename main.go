package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

func github(w http.ResponseWriter, r *http.Request) {
	value, found := c.Get("data")
	if found {
		json.NewEncoder(w).Encode(value)
		return
	}
	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	c.Set("data", resp, cache.DefaultExpiration)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.NewEncoder(w).Encode(string(data))
}
func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
func main() {
	handleRequests()
}
