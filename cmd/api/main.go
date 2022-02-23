package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"urlShortener/internal/api"
	"urlShortener/internal/db"
)

func main() {
	var err error

	if err = initLogger(); err != nil {
		fmt.Printf("Cannot init logger %s", err)
		os.Exit(1)
	}

	repository := db.NewURLRepository("shortener.db")
	if err = repository.Open(); err != nil {
		log.Errorf("Unable to create embedded db due to %s", err)
		os.Exit(1)
	}

	shortener := api.NewShortener(repository)

	router := api.NewRouter(shortener)
	routerc := make(chan error, 1)
	go router.Start(routerc)

	if err = awaitTermination(routerc, router, repository); err != nil {
		os.Exit(1)
	}
}

func initLogger() error {
	logfile, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("unable to create log file %w", err)
	}

	log.SetOutput(logfile)
	return nil
}

func awaitTermination(routerc chan error, closers ...io.Closer) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig:= <-signals:
		log.Infof("Received %s signal", sig)
		for _, c := range closers {
			err := c.Close()
			if err != nil {
				log.Errorf("Unable to close: %s", err)
			}
		}
	case err:= <-routerc:
		switch err {
		case http.ErrServerClosed:
			log.Infof("Server on %s port has been closed", api.Port)
		default:
			log.Errorf("Server on %s port failed to start", api.Port)
			return err
		}
	}

	log.Infof("Exiting application")
	return nil
}