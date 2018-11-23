package main

import (
	"crypto/sha1"
	"encoding/base64"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mz47/lunchomat/database"
	"github.com/tidwall/gjson"
)

var apiChannel = make(chan string)
var dbChannel = make(chan []string)
var testTemplate *template.Template

const (
	_apikey     = "cn2BthaLGbLT4MGwhXUdAycXtHXmhxWUXmI68TCMyHnx9cHeH66AH9RYQ-IsJSd3Bs_IEhCuIGHnTPvza6J0DLeE_2PQG1lOX2n-0rsrWhHRxwvekLKadG8Ae1LbW3Yx"
	_lon        = "10.016290"
	_lat        = "53.554920"
	_radius     = "1500"
	_categories = "lunch"
)

func main() {
	log.Println("starting application")
	database.Connect("restaurantes.db")
	defer database.Disconnect()
	startServer()
}

func startServer() {
	http.HandleFunc("/", handleIndex)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("New Request registered ...")
	w.Header().Set("Content-Type", "text/json")

	updateDatabase()
	database.ReceiveAll()
	db := database.ReceiveAll()

	for _, value := range db {
		w.Write([]byte(value))
		w.Write([]byte("\n"))
	}
}

func updateDatabase() {
	url := "https://api.yelp.com/v3/businesses/search?latitude=" +
		_lat + "&longitude=" +
		_lon + "&radius=" +
		_radius + "&categories=" +
		_categories

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+_apikey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	results := gjson.Get(string(payload), "businesses.#.name").Array()
	for _, value := range results {
		key := generateHash(value.String())
		database.Save(key, value.String())
	}
}

func generateHash(key string) string {
	hasher := sha1.New()
	hasher.Write([]byte(key))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return string(sha)
}
