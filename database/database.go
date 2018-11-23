package database

import (
	"log"

	"go.etcd.io/bbolt"
)

var database *bolt.DB

func main() {
	log.Print("database")
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
}

func ReceiveAll() []string {
	restaurantes := []string{}
	database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("bucket"))
		bucket.ForEach(func(key, value []byte) error {
			//log.Printf("received key: %s, value: %s \n", key, value)
			restaurantes = append(restaurantes, string(value))
			return nil
		})
		return nil
	})
	return restaurantes
}
