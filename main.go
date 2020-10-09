package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

var Kes *cache.Cache

type Response struct {
	Message string `json:"message"`
}

func github(w http.ResponseWriter, r *http.Request) {
	value, found := Kes.Get("status")
	if found {
		json.NewEncoder(w).Encode(value)
		return
	}

	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var response Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	Kes.Set("status", response, cache.DefaultExpiration)

	json.NewEncoder(w).Encode(string(data))
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	Kes = cache.New(2*time.Minute, 1*time.Minute)
	handleRequests()
}
