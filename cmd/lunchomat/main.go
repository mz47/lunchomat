package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

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
	startServer()
}

func startServer() {
	router := httprouter.New()
	router.GET("/", handleIndex)
	router.GET("/beenthere/:id", handleBeenThere)
	router.GET("/megusta/:id", handleMeGusta)
	router.GET("/nomegusta/:id", handleNoMeGusta)
	router.GET("/ignore/:id", handleIgnore)
	http.ListenAndServe(":8080", router)
}

func handleIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Requesting /")
	updateDatabase()
	restaurantes := lunchdb.ReceiveAll()
	template, _ := template.ParseFiles("../../web/index.html")
	template.Execute(w, restaurantes)
}

func handleBeenThere(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Requesting /beenthere for id: ", p.ByName("id"))
	restaurantes := lunchdb.ReceiveAll()
	template, _ := template.ParseFiles("../../web/index.html")
	template.Execute(w, restaurantes)
}
func handleMeGusta(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Requesting /megusta for id: ", p.ByName("id"))
	restaurantes := lunchdb.ReceiveAll()
	template, _ := template.ParseFiles("../../web/index.html")
	template.Execute(w, restaurantes)
}
func handleNoMeGusta(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Requesting /nomegusta for id: ", p.ByName("id"))
	restaurantes := lunchdb.ReceiveAll()
	template, _ := template.ParseFiles("../../web/index.html")
	template.Execute(w, restaurantes)
}
func handleIgnore(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Requesting /ignore for id: ", p.ByName("id"))
	restaurantes := lunchdb.ReceiveAll()
	template, _ := template.ParseFiles("../../web/index.html")
	template.Execute(w, restaurantes)
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
		lunchdb.Save(restaurant.NewRestaurant(value.Get("id").String(), value.Get("name").String(), value.Get("distance").Float(), value.Get("rating").Float(), value.Get("url").String()))
	}
}
