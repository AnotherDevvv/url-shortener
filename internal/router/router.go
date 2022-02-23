package router

import (
	stdContext "context"
	"fmt"
	"github.com/labstack/echo/v4"
	"time"
	"urlShortener/internal/shortener"
)

const Port = "1323"

type Router struct {
	shortener *shortener.Shortener
	echo *echo.Echo
}

func NewRouter(sh *shortener.Shortener) *Router {
	e := echo.New()
	e.POST("/shorten", sh.ShortenURL)
	e.GET("/:key", sh.FollowURL)

	return &Router{
		shortener: sh,
		echo: e,
	}
}

func (r *Router) Start(errc chan error) {
	errc <- r.echo.Start(":" + Port)
}

func (r *Router) Close() error {
	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 10*time.Second)
	defer cancel()
	err := r.echo.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to gracefully stop echo server on port %s, reason %w", Port, err)
	}

	return nil
}

func (r *Router) Trace(path string, h echo.HandlerFunc) {
	r.echo.TRACE(path, h)
}