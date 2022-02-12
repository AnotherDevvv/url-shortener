package db

import (
	"fmt"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

const urlBucketName = "shorten"

//go:generate moq -pkg mock -out mock/repository.go . Repository
type Repository interface {
	Insert(key string, value string) error
	Get(key string) (string, error)
	Close()
}

type URLRepository struct {
	db *bolt.DB
}

func NewURLRepository() *URLRepository {
	return &URLRepository{
		db: func() *bolt.DB {
			db, err := bolt.Open("shortener.db", 0600, nil)
			if err != nil {
				log.Fatalf("Unable to create embedded db due to %s", err)
				panic(err)
			}

			return db
		}(),
	}
}

func (urlRep *URLRepository) Close() {
	err := urlRep.db.Close()
	if err != nil{
		log.Error("Failed to close db gracefully %s", err)
	}
}

func (urlRep *URLRepository) Insert(key string, url string) error {
	return urlRep.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(urlBucketName))
		if err != nil {
			return fmt.Errorf("create bucket %s: %w", urlBucketName, err)
		}

		err = bucket.Put([]byte(key), []byte(url))
		if  err != nil {
			return fmt.Errorf("put into bucket %s: %w", urlBucketName, err)
		}

		return nil
	})
}

func (urlRep *URLRepository) Get(key string) (string, error) {
	var value []byte

	err := urlRep.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(urlBucketName))

		key := bucket.Get([]byte(key))
		if key != nil {
			value = make([]byte, len(key))
			copy(value, key)
		}

		return nil
	})

	return string(value), err
}
