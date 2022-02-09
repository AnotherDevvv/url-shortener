package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"urlShortener/internal/db"
)

type Shortener struct {
	urlRepository db.Repository
	address       string
}

func NewShortener(address string, r db.Repository) *Shortener {
	return &Shortener {
		address: address,
		urlRepository: r,
	}
}

func (sh *Shortener) ShortenUrl(c echo.Context) error {
	url := c.QueryParam("url")

	fmt.Println(url)

	e := validate(url)
	if e != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("URL cannot be shortened because %s", e.Error()))
	}

	hash := sha256.Sum256([]byte(url))

	err := sh.urlRepository.Insert(hash[:4], url)

	if err != nil {
		return c.String(http.StatusServiceUnavailable, e.Error())
	}

	return c.String(http.StatusOK, sh.address + hex.EncodeToString(hash[:4]))
}

func (sh *Shortener) FollowUrl(c echo.Context) error  {
	key := c.Param("key")

	hexKey, err := hex.DecodeString(key)

	url, err := sh.urlRepository.Get(hexKey)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if len(url) > 0 {
		return c.Redirect(http.StatusMovedPermanently, url)
	} else {
		return c.String(http.StatusBadRequest, fmt.Sprintf("short link with %s", key))
	}
}

func validate(url string) error {
	if !strings.HasPrefix(url, "http") || !strings.HasPrefix(url, "https") {
		return errors.New(fmt.Sprintf("url %s does not start with http(s)", url))
	} else {
		return nil
	}
}
