package indexer

import (
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
)

type Indexer struct {
	db *bolt.DB
}

func NewIndexer() *Indexer {
	db, err := bolt.Open("wiki-search-demo", 0600, nil)
	if err != nil {
		return nil
	}
	return &Indexer{
		db: db,
	}
}

func addToBucket(b *bolt.Bucket, word string, url string) {

	existingVal := b.Get([]byte(word))
	if existingVal == nil {
		b.Put([]byte(word), []byte(url))
		return
	}

	valSet := strings.Split(string(existingVal), ",")

	valSet = append(valSet, url)
	resSet := strings.Join(valSet, ",")

	b.Put([]byte(word), []byte(resSet))
	return
}

func (i *Indexer) Index(data string, url string) error {
	err := i.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("search-data"))
		if err != nil {
			return err
		}
		addToBucket(b, data, url)
		fmt.Println(data, " added to bucket at url ", url)
		return nil
	})

	return err
}

func (i *Indexer) ShowDB() {
	i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("search-data"))
		if b == nil {
			return fmt.Errorf("Error getting bucket: %s", "search-data")
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("Key: %s Value: %s\n", k, v)
		}
		return nil
	})
}
