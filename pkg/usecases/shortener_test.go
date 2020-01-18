package usescases

import (
	"errors"
	"log"
	"testing"

	"github.com/andrewwebber/tinyurl/pkg/entities"
	"github.com/andrewwebber/tinyurl/pkg/internal"
)

const baseURL = "http://localhost:8080"

func TestShortenURL(t *testing.T) {
	log.SetFlags(log.Llongfile)
	url := "https://www.urbandictionary.com/author.php?author=assdkf%3Basdlkfj"
	var injectFault bool
	var fault error
	f := func() error {
		if injectFault {
			return fault
		}

		return nil
	}

	r := &internal.KVStore{M: make(map[string]entities.ShortURL), FaultInjector: f}
	tinyURL := NewTinyURL(baseURL, r, XIDURLShortener)
	shortURL, err := tinyURL.ShortenURL(url)
	if err != nil {
		t.Fatal(err)
	}

	if len(shortURL.Short[len(baseURL):]) > maxLength {
		t.Fatalf("expected URL '%s' to be of length %d, not %d", shortURL.Short[len(baseURL):], maxLength, len(shortURL.Short[len(baseURL):]))
	}

	shortURL2, err := tinyURL.ShortenURL(url)
	if err != nil {
		t.Fatal(err)
	}

	if shortURL.Short == shortURL2.Short {
		t.Fatalf("short urls cannot be the same - %s : %s", shortURL.Short, shortURL2.Short)
	}

	t.Logf("lookup url %s", shortURL2.Short)

	shortURL2Get, err := tinyURL.URL(shortURL2.Short)
	if err != nil {
		t.Fatal(err)
	}

	if shortURL2.URL != shortURL2Get {
		t.Fatalf("unexpected url, expected %s - found %s", shortURL2.URL, shortURL2Get)
	}
}

func TestShortenURLWithRetries(t *testing.T) {
	url := "https://www.urbandictionary.com/author.php?author=assdkf%3Basdlkfj"

	fault := errors.New(internal.KVSTORE_KEYEXISTS)
	var count int
	f := func() error {
		if count > 3 {
			return nil
		}

		count++
		return fault
	}

	r := &internal.KVStore{M: make(map[string]entities.ShortURL), FaultInjector: f}
	tinyURL := NewTinyURL(baseURL, r, XIDURLShortener)
	if _, err := tinyURL.ShortenURL(url); err != nil {
		t.Fatal(err)
	}
}

func TestShortenURLWithMaxedRetries(t *testing.T) {
	url := "https://www.urbandictionary.com/author.php?author=assdkf%3Basdlkfj"

	fault := errors.New(internal.KVSTORE_KEYEXISTS)
	f := func() error {
		return fault
	}

	r := &internal.KVStore{M: make(map[string]entities.ShortURL), FaultInjector: f}
	tinyURL := NewTinyURL(baseURL, r, XIDURLShortener)
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

	r := &internal.KVStore{M: make(map[string]entities.ShortURL), FaultInjector: f}
	tinyURL := NewTinyURL(baseURL, r, XIDURLShortener)
	_, err := tinyURL.ShortenURL(url)
	if err == nil {
		t.Fatal("expected infra error")
	}

	if err != fault {
		t.Fatal(err)
	}
}
