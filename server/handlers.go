package server

import (
	"encoding/json"
	"flightpath/domain"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CalculateItineraryRequest [][]string

func (s Server) CalculateItinerary(c echo.Context) error {

	var request CalculateItineraryRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		s.logger.Debug().Err(err).Msg("invalid json")
		return c.String(http.StatusBadRequest, "invalid request body")
	}

	var inputRoutes []domain.Route
	for _, r := range request {
		inputRoutes = append(inputRoutes, domain.Route{Source: r[0], Destination: r[1]})
	}

	sortedRoutes, err := s.ItineraryService.CalculatePath(inputRoutes)
	if err != nil {
		s.logger.Debug().Err(err).Msg("calculating path failed")
		return c.String(http.StatusBadRequest, fmt.Sprintf("invalid route input detected: %v", err.Error()))
	}

	var responseArray []string
	for _, r := range sortedRoutes {
		responseArray = append(responseArray, r.Source)
	}
	responseArray = append(responseArray, sortedRoutes[len(sortedRoutes)-1].Destination)
	return c.JSON(http.StatusOK, responseArray)
}
