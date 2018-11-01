package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/tidwall/gjson"
	"go.etcd.io/bbolt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
)

var database *bolt.DB
var receiverChannel = make(chan string)
var hash = fnv.New32()

const APIKEY = "cn2BthaLGbLT4MGwhXUdAycXtHXmhxWUXmI68TCMyHnx9cHeH66AH9RYQ-IsJSd3Bs_IEhCuIGHnTPvza6J0DLeE_2PQG1lOX2n-0rsrWhHRxwvekLKadG8Ae1LbW3Yx"
const LON = "10.016290"
const LAT = "53.554920"
const RADIUS = "1500"

func main() {
	fmt.Println("starting application")
	Connect("restaurantes.db")
	defer Disconnect()

	startServer()
}

func startServer() {
	log.Println("starting http server on port 8080")
	http.HandleFunc("/", handleIndex)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("accepted request")
	go receiveData()
	payload := <-receiverChannel
	results := gjson.Get(payload, "businesses.#.name").Array()
	for _, value := range results {
		key := generateHash(value.String())
		Save(key, value.String())
		w.Write([]byte(value.String() + "\n"))
	}

}

func receiveData() {
	url := "https://api.yelp.com/v3/businesses/search?latitude=" + LAT + "&longitude=" + LON + "&radius=1000&categories=food"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+APIKEY)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	receiverChannel <- string(payload)
}

func Connect(file string) {
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to database: %s \n", file)
	database = db
}

func Disconnect() {
	if database != nil {
		database.Close()
		log.Println("connection to database closed")
	}
}

func Save(key string, value string) {
	database.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("bucket"))
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = bucket.Put([]byte(key), []byte(value))
		if err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	})
	log.Println("persisted value:", value)
}

func ReceiveAll() []string {
	database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("bucket"))
		bucket.ForEach(func(key, value []byte) error {
			log.Printf("received key: %s, value: %s \n", key, value)
			return nil
		})
		return nil
	})
	return nil
}

func generateHash(key string) string {
	hasher := sha1.New()
	hasher.Write([]byte(key))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return string(sha)
}
