package api

import "github.com/labstack/echo/v4"

type Shortener interface {
	ShortenURL(c echo.Context) error
	FollowURL(c echo.Context) error
}
