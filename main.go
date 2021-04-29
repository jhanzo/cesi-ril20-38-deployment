package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Server without framework example")

	if err := run(); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
}

func run() error {
	srv := newServer()

	config := config()

	if config.Debug {
		log.Info().Msg("debug mode")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if config.DB.Hostname == "" {
		return errors.New("missing db hostname")
	}

	if config.DB.Username == "" {
		return errors.New("missing db username")
	}

	if config.DB.Password == "" {
		return errors.New("missing db password")
	}

	if config.DB.Name == "" {
		return errors.New("missing db dbName")
	}

	srv.store = &dbStore{}
	err := srv.store.Open()
	if err != nil {
		return err
	}
	defer srv.store.Close()

	http.HandleFunc("/", srv.servHTTP)
	log.Info().Msg("Serving HTTP on port 9000")
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		return err
	}

	return nil
}
