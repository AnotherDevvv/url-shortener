package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"urlShortener/internal/api"
	"urlShortener/internal/db"
)

const port = "1323"

func main() {
	e := echo.New()
	rep := db.NewUrlRepository()

	sh := api.NewShortener(fmt.Sprintf("http://localhost:%s/", port), rep)

	e.POST("/shorten", sh.ShortenUrl)
	e.GET("/:key", sh.FollowUrl)

	err := e.Start(":" + port)
	if err != nil {
		fmt.Printf("Unable to start echo server on port %s, reason %s", port, err.Error())
		return
	}
}
