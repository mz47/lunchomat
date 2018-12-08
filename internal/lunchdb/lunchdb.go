package lunchdb

import (
	"encoding/json"
	"log"

	"github.com/mz47/lunchomat/internal/restaurant"
	"go.etcd.io/bbolt"
)

var lunchdb *bolt.DB

const _file = "database.db"
const _bucket = "bucket"

// Connect establishes a database connection
func Connect() {
	db, err := bolt.Open(_file, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to database: %s \n", _file)
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
		bucket, err := tx.CreateBucketIfNotExists([]byte(_bucket))
		if err != nil {
			log.Fatal(err)
			return err
		}

		rjson, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = bucket.Put([]byte(r.Id), rjson)
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
		bucket := tx.Bucket([]byte(_bucket))
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

// UpdateBeenThere updates a restaurant by the given id
func UpdateBeenThere(id string) {
	var item restaurant.Restaurant
	lunchdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_bucket))
		value := bucket.Get([]byte(id))
		if value != nil {
			error := json.Unmarshal(value, &item)
			if error != nil {
				log.Fatal(error)
				return nil
			}
		}
		return nil
	})
	if (restaurant.Restaurant{}) != item {
		log.Println("updated restaurant", item.Name)
		item.TimesVisited++
		Save(item)
		log.Println("saved restaurant", item.Name)
	} else {
		log.Println("No result with id", id, "found")
	}
}

// TogglePreferred updates a restaurant by the given id
func TogglePreferred(id string) {
	var item restaurant.Restaurant
	lunchdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_bucket))
		value := bucket.Get([]byte(id))
		if value != nil {
			error := json.Unmarshal(value, &item)
			if error != nil {
				log.Fatal(error)
				return nil
			}
		}
		return nil
	})
	if (restaurant.Restaurant{}) != item {
		log.Println("updated restaurant", item.Name)
		item.Preferred = !item.Preferred
		Save(item)
		log.Println("saved restaurant", item.Name)
	} else {
		log.Println("No result with id", id, "found")
	}
}

// ToggleIgnored updates a restaurant by the given id
func ToggleIgnored(id string) {
	var item restaurant.Restaurant
	lunchdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(_bucket))
		value := bucket.Get([]byte(id))
		if value != nil {
			error := json.Unmarshal(value, &item)
			if error != nil {
				log.Fatal(error)
				return nil
			}
		}
		return nil
	})
	if (restaurant.Restaurant{}) != item {
		log.Println("updated restaurant", item.Name)
		item.Ignored = !item.Ignored
		Save(item)
		log.Println("saved restaurant", item.Name)
	} else {
		log.Println("No result with id", id, "found")
	}
}

// Exists checks the existance of an id
func Exists(id string) bool {
	exists := false
	lunchdb.View(func(tx *bolt.Tx) error {
		bucket, _ := tx.CreateBucketIfNotExists([]byte(_bucket))
		value := bucket.Get([]byte(id))
		if value != nil {
			exists = true
			return nil
		}
		return nil
	})
	return exists
}
