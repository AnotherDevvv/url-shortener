package api

import (
	stdContext "context"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)


const port = "1323"

type Router struct {
	shortener *Shortener
	echo *echo.Echo
}

func NewRouter(sh *Shortener) *Router {
	e := echo.New()
	e.POST("/shorten", sh.ShortenURL)
	e.GET("/:key", sh.FollowURL)

	return &Router{
		shortener: sh,
		echo: e,
	}
}

func (r *Router) Start() {
	err := r.echo.Start(":" + port)

	switch err {
	case http.ErrServerClosed:
		log.Infof("Server on %s port has been closed", port)
	default:
		log.Fatalf("Unable to start echo server on port %s, reason %s", port, err)
		panic(err)
	}
}

func (r *Router) Close() error {
	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 10*time.Second)
	defer cancel()
	err := r.echo.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("Failed to gracefully stop echo server on port %s, reason %w", port, err)
	}

	return nil
}




