package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"urlShortener/internal/api"
	"urlShortener/internal/db"
)

func main() {
	initLogger()

	rep := db.NewURLRepository()
	defer rep.Close()

	sh := api.NewShortener(rep)

	api.NewRouter(sh).Start()
}

func initLogger()  {
	logfile, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic("Unable to create log file")
	}

	log.SetOutput(logfile)
}