package main

import (
	"testing"
)

func TestSimple(t *testing.T) {
	// Simples Beispiel analog zum Skript
	testGraphs := []Graph{
		{
			ID:    0,
			Start: Point{X: -6, Y: 1},
			End:   Point{X: 5, Y: 0},
		},
		{
			ID:    1,
			Start: Point{X: -7, Y: -1},
			End:   Point{X: 3, Y: 6},
		},
		{
			ID:    2,
			Start: Point{X: -1, Y: 6},
			End:   Point{X: 4, Y: -2},
		},
	}
	result := lineSweep(testGraphs)

	expected := 3
	if result != expected {
		t.Errorf("Line Sweep Algorithm returned %d intersections, expected %d", result, expected)
	}
}
