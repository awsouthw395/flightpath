package server

import (
	"flightpath/domain"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"net/http"
	"time"
)

type itineraryService interface {
	CalculatePath(edges []domain.Route) ([]domain.Route, error)
}

type Server struct {
	ItineraryService itineraryService
	logger           *zerolog.Logger
}

func Start(conf Config, logger *zerolog.Logger, itineraryService itineraryService) {

	s := Server{
		ItineraryService: itineraryService,
		logger:           logger,
	}

	e := echo.New()
	e.HideBanner = true
	e.POST("/itinerary", s.CalculateItinerary)

	server := http.Server{
		Addr:        fmt.Sprintf(":%s", conf.Port),
		Handler:     e,
		ReadTimeout: 30 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal().Err(err).Msg("Server experienced a fatal error")
	}
	return
}

type Config struct {
	Port string
}

func (c *Config) ValidateOrDefault() error {
	if c.Port == "" {
		c.Port = "8000"
	}

	return nil
}
