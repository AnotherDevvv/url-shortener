package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"syscall"
	"urlShortener/internal/api"
	"urlShortener/internal/db"
)

func main() {
	initLogger()

	repository := db.NewURLRepository()

	shortener := api.NewShortener(repository)

	router := api.NewRouter(shortener)
	go router.Start()

	awaitTermination(router, repository)
}

func initLogger()  {
	logfile, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic("Unable to create log file")
	}

	log.SetOutput(logfile)
}

func awaitTermination(closers ...io.Closer) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	terminated := make(chan bool, 1)

	go func() {
		sig := <-signals
		log.Infof("Received %s signal", sig)
		for _, c := range closers {
			err := c.Close()
			if err != nil {
				log.Errorf("Unable to close: %s", err)
			}
		}
		terminated <- true
	}()

	<-terminated
	log.Infof("Exiting application")
}