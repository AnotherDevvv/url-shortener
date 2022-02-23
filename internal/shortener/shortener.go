package shortener

import (
	"crypto/md5"
	"encoding/base32"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"urlShortener/internal/api"
)

type Shortener struct {
	urlRepository api.Repository
}

func NewShortener(repository api.Repository) *Shortener {
	return &Shortener{
		urlRepository: repository,
	}
}

func (sh *Shortener) ShortenURL(c echo.Context) error {
	url := c.QueryParam("url")

	log.Debugf("Shortening %s", url)

	err := validate(url)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("URL cannot be shortened because %s", err))
	}

	key := Encode(url)
	value, err := sh.urlRepository.Get(key)
	if err != nil {
		log.Errorf("Failed to check key existence for url %s. Cause: %s", url, err)
		return c.NoContent(http.StatusServiceUnavailable)
	}

	if len(value) > 0 {
		if url == value {
			return c.String(http.StatusOK, key)
		}

		key, err = sh.regenerate(key)
	}
	if err != nil {
		log.Errorf("Failed to regenerate key for url %s. Cause: %s", url, err)
		return c.NoContent(http.StatusServiceUnavailable)
	}

	err = sh.urlRepository.Insert(key, url)
	if err != nil {
		log.Errorf("Failed to insert url %s with key %s. Cause: %s", url, key, err)
		return c.NoContent(http.StatusServiceUnavailable)
	}

	log.Infof("URL %s has been encoded to %s key", url, key)

	return c.String(http.StatusOK, key)
}

func (sh *Shortener) FollowURL(c echo.Context) error  {
	key := c.Param("key")

	url, err := sh.urlRepository.Get(key)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if len(url) > 0 {
		return c.Redirect(http.StatusMovedPermanently, url)
	}

	return c.String(http.StatusNotFound, fmt.Sprintf("short link with %s key not found", key))
}

func (sh *Shortener) regenerate(key string) (string, error) {
	newKey := key
	for {
		newKey = Encode(newKey)
		value, err := sh.urlRepository.Get(newKey)
		if len(value) == 0 || err != nil {
			return newKey, err
		}
	}
}

// md5 produces 128bit. hash[:5] is a slice of first 5 bytes (40 bits).
// base32 encodes symbols by 5 bits and produces exactly (40/5)=8 symbols for short link
// md5 is chosen because its less cpu intensive than sha. base32 doesnt generate into "/" and "+" as base64 does
func Encode(url string) string {
	hash := md5.Sum([]byte(url))
	return base32.StdEncoding.EncodeToString(hash[:5])
}

func validate(url string) error {
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		return errors.New(fmt.Sprintf("url %s does not start with http(s)", url))
	}

	return nil
}