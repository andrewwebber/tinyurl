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

	er := NewEntitiesRepository(db)
	var err error
	url := entities.ShortURL{URL: "url", Short: "short"}

	if err = er.Insert(insertKey, url); err != nil {
		t.Fatal(err)
	}

	if err = er.Insert(insertKey, url); err == nil {
		t.Fatal("expected 'key already exists' error")
	}

	if err != gocb.ErrKeyExists {
		t.Fatal(err)
	}

	if er.IsShortURLExistsError(err) {
		return
	}

	t.Fatalf("unexpected error - %v", err)
}

func TestGet(t *testing.T) {
	db := dbConnect(t)
	getKey := "TestGetReplica"
	defer func() {
		db.Remove(getKey)
	}()

	er := NewEntitiesRepository(db)
	var err error
	url := entities.ShortURL{URL: "url", Short: "short"}
	if err = er.Insert(getKey, url); err != nil {
		t.Fatal(err)
	}

	var replica entities.ShortURL
	if replica, err = er.Get(getKey); err != nil {
		t.Fatal(err)
	}

	if replica.Short != url.Short {
		t.Fatalf("unexpected short url, found %s - expected %s", replica.Short, url.Short)
	}
}
