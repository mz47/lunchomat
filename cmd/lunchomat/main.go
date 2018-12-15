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
	//updateDatabase()
	startServer()
}

func startServer() {
	router := chi.NewRouter()
	router.Get("/", handleIndex)
	router.Get("/refresh", handleRefresh)
	router.Get("/add", handleAddGet)
	router.Post("/add", handleAddPost)	
	router.Get("/visited/{restaurantId}", handleVisited)
	router.Get("/preferred/{restaurantId}", handlePreferred)
	router.Get("/ignored/{restaurantId}", handleIgnored)
	http.ListenAndServe(":8080", router)
}

func renderTemplate(w http.ResponseWriter) {
	restaurantes := lunchdb.ReceiveAll()
	template, err := template.ParseFiles("../../web/index.html")
	if err != nil {
		log.Fatal("error while parsing template", err)
	}
	template.Execute(w, restaurantes)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("Requesting", r.RequestURI)
	renderTemplate(w)
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	log.Println("Requesting", r.RequestURI)
	renderTemplate(w)
}

func handleVisited(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantId")
	if id != "" {
		log.Println("Requesting", r.RequestURI)
		lunchdb.UpdateBeenThere(id)
	}
	renderTemplate(w)
}

func handleIgnored(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantId")
	if id != "" {
		log.Println("Requesting", r.RequestURI)
		lunchdb.ToggleIgnored(id)
	}
	renderTemplate(w)
}

func handlePreferred(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantId")
	if id != "" {
		log.Println("Requesting", r.RequestURI)
		lunchdb.TogglePreferred(id)
	}
	renderTemplate(w)
}

func handleAddGet(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("../../web/add.html")
	if err != nil {
		log.Fatal("error while parsing template", err)
	}
	template.Execute(w, "")
}

func handleAddPost(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("../../web/add.html")
	if err != nil {
		log.Fatal("error while parsing template", err)
	}
	template.Execute(w, "")
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
		log.Fatal("error while fetching data from api", err)
	}

	req.Header.Set("Authorization", "Bearer "+_apikey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("error while authenticating", err)
	}

	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error while parsing api response", err)
	}

	results := gjson.Get(string(payload), "businesses").Array()
	for _, value := range results {
		lunchdb.Save(
			restaurant.NewRestaurant(
				value.Get("id").String(),
				value.Get("name").String(),
				value.Get("distance").Float(),
				value.Get("rating").Float(),
				value.Get("url").String()))
	}
}
