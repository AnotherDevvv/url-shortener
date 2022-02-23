package db

import (
	"fmt"
	"github.com/boltdb/bolt"
)

const urlBucketName = "shorten"


type URLRepository struct {
	filepath string
	db       *bolt.DB
}

func NewURLRepository(filepath string) *URLRepository {
	return &URLRepository{
		filepath: filepath,
	}
}

func (urlRep *URLRepository) Open() error {
	db, err := bolt.Open(urlRep.filepath, 0600, nil)

	urlRep.db = db
	return err
}

func (urlRep *URLRepository) Close() error {
	err := urlRep.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close db gracefully %w", err)
	}

	return nil
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