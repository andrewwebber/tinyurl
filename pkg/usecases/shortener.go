package usescases

import (
	"fmt"
	"strings"

	"github.com/andrewwebber/tinyurl/pkg/entities"
	"github.com/rs/xid"
)

const maxLength = 21

type ShortURLRepository interface {
	Insert(key string, shortURL entities.ShortURL) error
	Get(key string) (entities.ShortURL, error)
	IsShortURLExistsError(err error) bool
}

type URLShortener func(url string) string

func XIDURLShortener(url string) string {
	return xid.New().String()
}

func NewTinyURL(baseURL string, r ShortURLRepository, s URLShortener) TinyURL {
	return TinyURL{
		baseURL:    baseURL,
		repository: r,
		shortener:  s,
	}
}

type TinyURL struct {
	baseURL    string
	shortener  URLShortener
	repository ShortURLRepository
}

func (t *TinyURL) ShortenURL(url string) (entities.ShortURL, error) {
	result := entities.ShortURL{URL: url}

	var err error
	for retryCount := 0; retryCount < 5; retryCount++ {
		short := t.shortener(url)
		result.Short = short
		result.URL = url
		if err = t.repository.Insert(short, result); err != nil {
			if t.repository.IsShortURLExistsError(err) {
				continue
			}

			return result, err
		}

		result.Short = fmt.Sprintf("%s/%s", t.baseURL, result.Short)
		break
	}

	return result, err
}

func (t *TinyURL) URL(shortURL string) (string, error) {
	var result string
	shortURL = strings.Replace(shortURL, t.baseURL+"/", "", 1)
	if len(shortURL) > maxLength {
		return result, fmt.Errorf("invalid shorturl length %d", len(shortURL))
	}

	e, err := t.repository.Get(shortURL)
	if err != nil {
		return result, err
	}

	result = e.URL
	return result, nil
}
