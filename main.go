package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/henrysusiloo/cache-in-go/cache"
)

var c cache.CacheItf

func github(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://api.github.com/status")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var response = struct {
		Message string `json:"message"`
	}{}

	if err := json.Unmarshal(data, &response); err != nil {
		json.NewEncoder(w).Encode(err)
	}

	if err := c.Set("status", response); err != nil {
		json.NewEncoder(w).Encode(err)
	}

	json.NewEncoder(w).Encode(string(data))

	foo, err := c.Get("status")
	if err == nil {
		fmt.Println(foo)
	}
}

func handleRequests() {
	http.HandleFunc("/github", github)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	c = cache.NewCache("go-cache")

	handleRequests()
}
