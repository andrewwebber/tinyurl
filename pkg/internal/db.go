package internal

import (
	"errors"
	"fmt"
	"log"

	"github.com/andrewwebber/tinyurl/pkg/entities"
)

type FaultInjector func() error

const KVSTORE_KEYEXISTS = "key already exists"
const KVSTORE_KEYNOTFOUND = "key not found - %s"

type KVStore struct {
	M             map[string]entities.ShortURL
	FaultInjector FaultInjector
}

func (s *KVStore) Insert(key string, shortURL entities.ShortURL) error {
	if s.FaultInjector != nil {
		if err := s.FaultInjector(); err != nil {
			return err
		}
	}

	if _, ok := s.M[key]; ok {
		return errors.New(KVSTORE_KEYEXISTS)
	}

	s.M[key] = shortURL
	log.Printf("insert %s", key)

	return nil
}

func (s *KVStore) Get(key string) (entities.ShortURL, error) {
	if s.FaultInjector != nil {
		if err := s.FaultInjector(); err != nil {
			return entities.ShortURL{}, err
		}
	}

	v, ok := s.M[key]
	if !ok {
		return entities.ShortURL{}, fmt.Errorf(KVSTORE_KEYNOTFOUND, key)
	}

	log.Printf("get %s", key)
	return v, nil
}

func (s *KVStore) IsShortURLExistsError(err error) bool {
	return err.Error() == KVSTORE_KEYEXISTS
}
