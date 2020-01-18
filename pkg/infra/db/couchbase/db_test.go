package couchbase

import (
	"testing"

	"github.com/andrewwebber/tinyurl/pkg/entities"
	"github.com/couchbase/gocb"
)

const (
	CLUSTERURL       = "couchbase://localhost"
	CLUSTER_USERNAME = "Administrator"
	CLUSTER_PASSWORD = "password"
	BUCKET_NAME      = "tinyurl"
)

func dbConnect(t *testing.T) *Db {
	db, err := New(CLUSTERURL, CLUSTER_USERNAME, CLUSTER_PASSWORD, BUCKET_NAME)
	if err != nil {
		t.Fatal(err)
	}

	if db == nil {
		t.Fatal("did not expect empty database")
	}

	return db
}

func TestInsert(t *testing.T) {
	db := dbConnect(t)

	insertKey := "TestInsert"
	defer func() {
		db.Remove(insertKey)
	}()

	var err error
	url := entities.ShortURL{URL: "url", Short: "short"}
	if err = db.Insert(insertKey, &url); err != nil {
		t.Fatal(err)
	}

	if err = db.Insert(insertKey, &url); err == nil {
		t.Fatal("expected 'key already exists' error")
	}

	if err != gocb.ErrKeyExists {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	db := dbConnect(t)
	getKey := "TestGetReplica"
	defer func() {
		db.Remove(getKey)
	}()

	url := entities.ShortURL{URL: "url", Short: "short"}
	if err := db.Insert(getKey, &url); err != nil {
		t.Fatal(err)
	}

	var replica entities.ShortURL
	if err := db.Get(getKey, &replica); err != nil {
		t.Fatal(err)
	}

}
