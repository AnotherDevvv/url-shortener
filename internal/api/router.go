package api

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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
	if err != nil {
		log.Fatalf("Unable to start echo server on port %s, reason %s", port, err)
		panic(err)
	}
}




