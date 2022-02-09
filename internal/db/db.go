package db

import (
	"fmt"
	"github.com/boltdb/bolt"
)

//go:generate moq -pkg mock -out mock/repository.go . Repository

type Repository interface {
	Insert(key []byte, url string) error
	Get(key []byte) (string, error)
	Close() error
}

type UrlRepository struct {
	db      * bolt.DB
}

func NewUrlRepository() *UrlRepository {
	urlRep := new(UrlRepository)

	db, err := bolt.Open("shortener.db", 0600, nil)

	if err != nil {
		fmt.Printf("Unable to create embedded db db due to %s", err.Error())
	}

	urlRep.db = db

	return urlRep
}


func (urlRep *UrlRepository) Close() error {
	return urlRep.db.Close()
}

func (urlRep *UrlRepository) Insert(key []byte, url string) error {
	err := urlRep.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("shortener"))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		if b.Put(key, []byte(url)) != nil {
			return err
		}
		return nil
	})

	return err
}

func (urlRep *UrlRepository) Get(key []byte) (string, error) {
	v := []byte{}

	err := urlRep.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("shortener"))

		k := b.Get(key)

		if k != nil {
			v = make([]byte, len(k))
			copy(v, k)
		}

		return nil
	})

	return string(v), err
}
