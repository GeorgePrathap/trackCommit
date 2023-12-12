package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"trackCommit/model"

	"github.com/boltdb/bolt"
)

func ConnectDB() *bolt.DB {
	// Open or create BoltDB database
	db, err := bolt.Open("mydb.db", 0600, nil)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}

	return db
}

func CloseDatabase(db *bolt.DB) {
	db.Close()
}

func CreateBucket(db *bolt.DB, bucketName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
}

func DeleteBucket(db *bolt.DB, bucketName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
}

func SaveCommitDetails(db *bolt.DB, bucketName string, details model.CommitDetails) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}

		// Encode CommitDetails struct to binary
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(details); err != nil {
			return fmt.Errorf("encode commit details: %v", err)
		}

		// Save to BoltDB
		if err := bucket.Put([]byte(details.Hash), buf.Bytes()); err != nil {
			return fmt.Errorf("put commit details: %v", err)
		}

		return nil
	})
}

func RetrieveCommitDetailsByHash(db *bolt.DB, bucketName string, hash string) (model.CommitDetails, error) {
	var details model.CommitDetails

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		data := bucket.Get([]byte(hash))
		if data == nil {
			return fmt.Errorf("commit details not found for hash %s", hash)
		}

		// Decode binary data into CommitDetails struct
		dec := gob.NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&details); err != nil {
			return fmt.Errorf("decode commit details: %v", err)
		}

		return nil
	})

	return details, err
}
