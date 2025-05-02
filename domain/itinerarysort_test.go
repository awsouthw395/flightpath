package domain_test

import (
	"flightpath/domain"
	"slices"
	"testing"
)

// [["LAX","DXB"], ["JFK","LAX"], ["SFO","SJC"], ["DXB","SFO"]]
func TestBasicPath(t *testing.T) {
	routes := []domain.Route{
		{"LAX", "DXB"},
		{"JFK", "LAX"},
		{"SFO", "SJC"},
		{"DXB", "SFO"},
	}
	expectedResult := []domain.Route{
		{"JFK", "LAX"},
		{"LAX", "DXB"},
		{"DXB", "SFO"},
		{"SFO", "SJC"},
	}
	result, err := domain.ItineraryService{}.CalculatePath(routes)
	if err != nil {
		t.Errorf("TestBasicPath experienced an error: %v", err)
	}

	if !slices.Equal(result, expectedResult) {
		t.Errorf("TestBasicPath experienced an unexpected result: %v != %v", result, expectedResult)
	}
}

func TestMultipleTerminations(t *testing.T) {
	routes := []domain.Route{
		{"LAX", "DXB"},
		{"JFK", "LAX"},
		{"SFO", "SJC"},
		{"DXB", "SFO"},
		{"SFO", "ABC"},
	}
	_, err := domain.ItineraryService{}.CalculatePath(routes)
	if err.Error() != "multiple terminations found" {
		t.Errorf("TestMultipleTerminations experienced an unexpected result: %v", err)
	}
}

func TestMultipleOrigins(t *testing.T) {
	routes := []domain.Route{
		{"LAX", "DXB"},
		{"JFK", "LAX"},
		{"SFO", "SJC"},
		{"DXB", "SFO"},
		{"ABC", "DXB"},
	}
	_, err := domain.ItineraryService{}.CalculatePath(routes)
	if err.Error() != "multiple origins found" {
		t.Errorf("TestMultipleOrigins experienced an unexpected result: %v", err)
	}
}

func TestCompleteCircle(t *testing.T) {
	routes := []domain.Route{
		{"LAX", "DXB"},
		{"JFK", "LAX"},
		{"SFO", "SJC"},
		{"DXB", "SFO"},
		{"SJC", "JFK"},
	}
	_, err := domain.ItineraryService{}.CalculatePath(routes)
	if err.Error() != "no start or end node exists, all nodes belong to a circle" {
		t.Errorf("TestCompleteCircle experienced an unexpected result: %v", err)
	}
}

func TestPathEndsInACircle(t *testing.T) {
	routes := []domain.Route{
		{"LAX", "DXB"},
		{"JFK", "LAX"},
		{"SFO", "SJC"},
		{"DXB", "SFO"},
		{"SJC", "LAX"},
	}
	expectedResult := []domain.Route{
		{"JFK", "LAX"},
		{"LAX", "DXB"},
		{"DXB", "SFO"},
		{"SFO", "SJC"},
		{"SJC", "LAX"},
	}
	result, err := domain.ItineraryService{}.CalculatePath(routes)
	if err != nil {
		t.Errorf("TestBasicPath experienced an error: %v", err)
	}

	if !slices.Equal(result, expectedResult) {
		t.Errorf("TestBasicPath experienced an unexpected result: %v != %v", result, expectedResult)
	}
}

func TestPathStartsInACircle(t *testing.T) {
	routes := []domain.Route{
		{"SFO", "JFK"},
		{"LAX", "DXB"},
		{"JFK", "LAX"},
		{"SFO", "SJC"},
		{"DXB", "SFO"},
	}
	expectedResult := []domain.Route{
		{"SFO", "JFK"},
		{"JFK", "LAX"},
		{"LAX", "DXB"},
		{"DXB", "SFO"},
		{"SFO", "SJC"},
	}
	result, err := domain.ItineraryService{}.CalculatePath(routes)
	if err != nil {
		t.Errorf("TestBasicPath experienced an error: %v", err)
	}

	if !slices.Equal(result, expectedResult) {
		t.Errorf("TestBasicPath experienced an unexpected result: %v != %v", result, expectedResult)
	}
}

// [["LAX","DXB"], ["JFK","LAX"], ["SFO","SJC"], ["DXB","SFO"], ["SJC","JFK"], ["CDG","JFK"], ["JFK","LHR"]]
func TestPathContainsACircle(t *testing.T) {
	routes := []domain.Route{
		{"LAX", "DXB"},
		{"JFK", "LAX"},
		{"SFO", "SJC"},
		{"DXB", "SFO"},
		{"SJC", "JFK"},
		{"CDG", "JFK"},
		{"JFK", "LHR"},
	}
	expectedResult := []domain.Route{
		{"CDG", "JFK"},
		{"JFK", "LAX"},
		{"LAX", "DXB"},
		{"DXB", "SFO"},
		{"SFO", "SJC"},
		{"SJC", "JFK"},
		{"JFK", "LHR"},
	}
	result, err := domain.ItineraryService{}.CalculatePath(routes)
	if err != nil {
		t.Errorf("TestBasicPath experienced an error: %v", err)
	}

	if !slices.Equal(result, expectedResult) {
		t.Errorf("TestBasicPath experienced an unexpected result: %v != %v", result, expectedResult)
	}
}

func TestPathContainsTwoCirclesOnOneNode(t *testing.T) {
	routes := []domain.Route{
		{"LAX", "DXB"},
		{"JFK", "LAX"},
		{"SFO", "SJC"},
		{"DXB", "SFO"},
		{"SJC", "JFK"},
		{"CDG", "JFK"},
		{"JFK", "LHR"},
		{"JFK", "SFO"},
		{"SFO", "JFK"},
	}

	_, err := domain.ItineraryService{}.CalculatePath(routes)
	if err.Error() != "path cannot be determined, unresolvable node paths" {
		t.Errorf("TestCompleteCircle experienced an unexpected result: %v", err)
	}
}
