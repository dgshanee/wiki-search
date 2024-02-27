package indexer

import (
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/dgshanee/search-engine-demo/crawler"
)

type Indexer struct {
	db *bolt.DB
}

func NewIndexer(db *bolt.DB) *Indexer {
	return &Indexer{
		db: db,
	}
}

func addToBucket(b *bolt.Bucket, set crawler.WordData) {
	setWord := set.Word
	setUrl := set.Url

	existingVal := b.Get([]byte(setWord))
	if existingVal == nil {
		b.Put([]byte(setWord), []byte(setUrl))
		return
	}

	valSet := strings.Split(string(existingVal), ",")

	valSet = append(valSet, setUrl)
	resSet := strings.Join(valSet, ",")

	b.Put([]byte(setWord), []byte(resSet))
	return
}

func (i *Indexer) Index(data []crawler.WordData) {
	i.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("search-data"))
		if err != nil {
			return err
		}

		for _, wordData := range data {
			addToBucket(b, wordData)
		}
		return nil
	})
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
