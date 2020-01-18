package usescases

import (
	"errors"
	"testing"

	"github.com/andrewwebber/tinyurl/pkg/entities"
)

type FaultInjector func() error

const KVSTORE_KEYEXISTS = "key already exists"
const KVSTORE_KEYNOTFOUND = "key not found"

type kvStore struct {
	m             map[string]entities.ShortURL
	faultInjector FaultInjector
}

func (s *kvStore) Insert(key string, shortURL entities.ShortURL) error {
	if s.faultInjector != nil {
		if err := s.faultInjector(); err != nil {
			return err
		}
	}

	if _, ok := s.m[key]; ok {
		return errors.New(KVSTORE_KEYEXISTS)
	}

	s.m[key] = shortURL

	return nil
}

func (s *kvStore) Get(key string) (entities.ShortURL, error) {
	if s.faultInjector != nil {
		if err := s.faultInjector(); err != nil {
			return entities.ShortURL{}, err
		}
	}

	v, ok := s.m[key]
	if !ok {
		return entities.ShortURL{}, errors.New(KVSTORE_KEYNOTFOUND)
	}

	return v, nil
}

func (s *kvStore) IsShortURLExistsError(err error) bool {
	return err.Error() == KVSTORE_KEYEXISTS
}

func TestShortenURL(t *testing.T) {
	url := "https://www.urbandictionary.com/author.php?author=assdkf%3Basdlkfj"
	var injectFault bool
	var fault error
	f := func() error {
		if injectFault {
			return fault
		}

		return nil
	}

	r := &kvStore{m: make(map[string]entities.ShortURL), faultInjector: f}
	tinyURL := NewTinyURL(r, XIDURLShortener)
	shortURL, err := tinyURL.ShortenURL(url)
	if err != nil {
		t.Fatal(err)
	}

	maxLength := 20
	if len(shortURL.Short) > maxLength {
		t.Fatalf("expected URL '%s' to be of length %d, not %d", shortURL.Short, maxLength, len(shortURL.Short))
	}

	shortURL2, err := tinyURL.ShortenURL(url)
	if err != nil {
		t.Fatal(err)
	}

	if shortURL.Short == shortURL2.Short {
		t.Fatalf("short urls cannot be the same - %s : %s", shortURL.Short, shortURL2.Short)
	}
}

func TestShortenURLWithRetries(t *testing.T) {
	url := "https://www.urbandictionary.com/author.php?author=assdkf%3Basdlkfj"

	fault := errors.New(KVSTORE_KEYEXISTS)
	var count int
	f := func() error {
		if count > 3 {
			return nil
		}

		count++
		return fault
	}

	r := &kvStore{m: make(map[string]entities.ShortURL), faultInjector: f}
	tinyURL := NewTinyURL(r, XIDURLShortener)
	if _, err := tinyURL.ShortenURL(url); err != nil {
		t.Fatal(err)
	}
}

func TestShortenURLWithMaxedRetries(t *testing.T) {
	url := "https://www.urbandictionary.com/author.php?author=assdkf%3Basdlkfj"

	fault := errors.New(KVSTORE_KEYEXISTS)
	f := func() error {
		return fault
	}

	r := &kvStore{m: make(map[string]entities.ShortURL), faultInjector: f}
	tinyURL := NewTinyURL(r, XIDURLShortener)
	if _, err := tinyURL.ShortenURL(url); err == nil {
		t.Fatal("expected maxed retries error case")
	}
}

func TestShortenURLWithInfraFailure(t *testing.T) {
	url := "https://www.urbandictionary.com/author.php?author=assdkf%3Basdlkfj"

	fault := errors.New("infra failure")
	f := func() error {
		return fault
	}

	r := &kvStore{m: make(map[string]entities.ShortURL), faultInjector: f}
	tinyURL := NewTinyURL(r, XIDURLShortener)
	_, err := tinyURL.ShortenURL(url)
	if err == nil {
		t.Fatal("expected infra error")
	}

	if err != fault {
		t.Fatal(err)
	}
}
