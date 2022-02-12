package api

import (
	"crypto/md5"
	"encoding/base32"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"urlShortener/internal/db"
)

type Shortener struct {
	urlRepository db.Repository
}

func NewShortener(repository db.Repository) *Shortener {
	return &Shortener {
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

	key := encode(url)

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

func validate(url string) error {
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		return errors.New(fmt.Sprintf("url %s does not start with http(s)", url))
	}

	return nil
}

// md5 produces 128bit. hash[:5] is a slice of first 5 bytes (40 bits).
// base32 encodes symbols by 5 bits and produces exactly (40/5)=8 symbols for short link
// md5 is chosen because its less cpu intensive than sha. base32 doesnt encode into "/" and "+" as base64 does
func encode(url string) string {
	hash := md5.Sum([]byte(url))

	return base32.StdEncoding.EncodeToString(hash[:5])
}