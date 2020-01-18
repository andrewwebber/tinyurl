package couchbase

import (
	"errors"

	"github.com/andrewwebber/tinyurl/pkg/entities"
	"github.com/couchbase/gocb"
)

type EntitiesRepository struct {
	db *Db
}

func NewEntitiesRepository(db *Db) *EntitiesRepository {
	return &EntitiesRepository{db: db}
}

func (e *EntitiesRepository) Insert(key string, shortURL entities.ShortURL) error {
	return e.db.Insert(key, &shortURL)
}

func (e *EntitiesRepository) Get(key string) (entities.ShortURL, error) {

	var result entities.ShortURL
	err := e.db.Get(key, &result)
	return result, err
}

func (e *EntitiesRepository) IsShortURLExistsError(err error) bool {
	return errors.Is(err, gocb.ErrKeyExists)
}

type Db struct {
	cluster *gocb.Cluster
	bucket  *gocb.Bucket
}

func New(clusterURL string, clusterUsername string, clusterPassword string, bucketName string) (*Db, error) {
	cluster, err := gocb.Connect(clusterURL)
	if err != nil {
		return nil, err
	}

	if err = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: clusterUsername,
		Password: clusterPassword,
	}); err != nil {
		return nil, err
	}

	bucket, err := cluster.OpenBucket(bucketName, "")
	if err != nil {
		return nil, err
	}

	return &Db{
		cluster: cluster,
		bucket:  bucket,
	}, nil
}

func (db *Db) Insert(key string, obj interface{}) error {
	if _, err := db.bucket.Insert(key, obj, 0); err != nil {
		return err
	}

	return nil
}

func (db *Db) Get(key string, obj interface{}) error {
	_, err := db.bucket.Get(key, obj)
	if err != nil {
		if _, err = db.bucket.GetReplica(key, obj, 0); err != nil {
			return err
		}
	}

	return nil
}

func (db *Db) Remove(key string) error {
	if _, err := db.bucket.Remove(key, 0); err != nil {
		return err
	}

	return nil
}
