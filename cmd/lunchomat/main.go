package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/mz47/lunchomat/internal/lunchdb"
	"github.com/mz47/lunchomat/internal/restaurant"
	"github.com/tidwall/gjson"
)

var apiChannel = make(chan string)
var dbChannel = make(chan []string)

const (
	_apikey     = "cn2BthaLGbLT4MGwhXUdAycXtHXmhxWUXmI68TCMyHnx9cHeH66AH9RYQ-IsJSd3Bs_IEhCuIGHnTPvza6J0DLeE_2PQG1lOX2n-0rsrWhHRxwvekLKadG8Ae1LbW3Yx"
	_lon        = "10.016290"
	_lat        = "53.554920"
	_radius     = "1500"
	_categories = "lunch"
	_attributes = "GoodForMeal.lunch"
)

func main() {
	log.Println("starting application")
	lunchdb.Connect()
	defer lunchdb.Disconnect()
	updateDatabase()
	startServer()
}

func startServer() {
	router := chi.NewRouter()
	router.Get("/", handleIndex)
	router.Get("/refresh", handleRefresh)
	router.Get("/visited/{restaurantId}", handleVisited)
	router.Get("/preferred/{restaurantId}", handlePreferred)
	router.Get("/ignored/{restaurantId}", handleIgnored)

	http.ListenAndServe(":8080", router)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("Requesting", r.RequestURI)
	restaurantes := lunchdb.ReceiveAll()
	template, _ := template.ParseFiles("../../web/index.html")
	template.Execute(w, restaurantes)
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	log.Println("Requesting", r.RequestURI)
	restaurantes := lunchdb.ReceiveAll()
	template, _ := template.ParseFiles("../../web/index.html")
	template.Execute(w, restaurantes)
}

func handleVisited(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantId")
	if id != "" {
		log.Println("Requesting", r.RequestURI)
		lunchdb.UpdateBeenThere(id)
	}
	http.Redirect(w, r, "/", 302)
}

func handleIgnored(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantId")
	if id != "" {
		log.Println("Requesting", r.RequestURI)
		lunchdb.ToggleIgnored(id)
	}
	http.Redirect(w, r, "/", 302)
}

func handlePreferred(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantId")
	if id != "" {
		log.Println("Requesting", r.RequestURI)
		lunchdb.TogglePreferred(id)
	}
	http.Redirect(w, r, "/", 302)
}

func updateDatabase() {
	url := "https://api.yelp.com/v3/businesses/search" +
		"?latitude=" + _lat +
		"&longitude=" + _lon +
		"&radius=" + _radius +
		"&categories=" + _categories +
		"&attrs=" + _attributes

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

	results := gjson.Get(string(payload), "businesses").Array()
	for _, value := range results {
		//if !lunchdb.Exists(value.Get("id").String()) {
		lunchdb.Save(
			restaurant.NewRestaurant(
				value.Get("id").String(),
				value.Get("name").String(),
				value.Get("distance").Float(),
				value.Get("rating").Float(),
				value.Get("url").String()))
		//}
	}
}
