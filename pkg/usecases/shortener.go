package usescases

import (
	"github.com/andrewwebber/tinyurl/pkg/entities"
	"github.com/rs/xid"
)

type ShortURLRepository interface {
	Insert(key string, shortURL entities.ShortURL) error
	Get(key string) (entities.ShortURL, error)
	IsShortURLExistsError(err error) bool
}

type URLShortener func(url string) string

func XIDURLShortener(url string) string {
	return xid.New().String()
}

func NewTinyURL(r ShortURLRepository, s URLShortener) TinyURL {
	return TinyURL{
		repository: r,
		shortener:  s,
	}
}

type TinyURL struct {
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
	}

	return result, err
}
