package main

import (
	"net/http"
	"os"

	"github.com/Peikkin/postgres-golang/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	router := router.Router()
	log.Info().Msg("Запуск сервера на порту :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal().Err(err).Msg("Ошибка запуска сервера")
	}
}
