package lunchdb

import (
	"encoding/json"
	"log"

	"github.com/mz47/lunchomat/internal/restaurant"
	"go.etcd.io/bbolt"
)

var lunchdb *bolt.DB

// Connect establishes a database connection
func Connect() {
	db, err := bolt.Open("database.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to database: %s \n", "database.db")
	lunchdb = db
}

// Disconnect terminates database connection and closes file
func Disconnect() {
	if lunchdb != nil {
		lunchdb.Close()
		log.Println("connection to database closed")
	}
}

// Save key and value in database
func Save(r restaurant.Restaurant) {
	lunchdb.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("bucket"))
		if err != nil {
			log.Fatal(err)
			return err
		}

		rjson, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = bucket.Put([]byte(r.Name), rjson)
		if err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	})
}

// ReceiveAll gets all values from database and return restaurant array
func ReceiveAll() []restaurant.Restaurant {
	restaurantes := []restaurant.Restaurant{}
	var item restaurant.Restaurant
	lunchdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("bucket"))
		bucket.ForEach(func(key, value []byte) error {
			err := json.Unmarshal(value, &item)
			if err != nil {
				log.Fatal(err)
			}
			restaurantes = append(restaurantes, item)
			return nil
		})
		return nil
	})
	return restaurantes
}
