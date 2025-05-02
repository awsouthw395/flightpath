package cmd

import (
	"flightpath/config"
	"flightpath/domain"
	"flightpath/server"
	"github.com/rs/zerolog"
	"log"
	"os"
)

func Start() {
	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("An error occurred whilst attempting to parse config values: %v", err)
	}

	logLevel, err := zerolog.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Fatalf("An error occurred whilst attempting to parse log level: %v", err)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(logLevel)
	logger.Info().Msg("Starting server")

	server.Start(conf.Server, &logger, domain.ItineraryService{})
}
