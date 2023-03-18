package main

import (
	"testing"
	"wav/parser"
	"wav/utils"
)

func TestToMonoSamples(t *testing.T) {
	w := &parser.Wav{
		Data: [][]int16{
			{0, 0, 5, 5, 12345, 12345, -12345},
			{0, 1, 7, 8, 16789, -12345, 0},
		},
	}

	given := w.GetMonoSamples()
	expected := []int16{0, 0, 6, 6, 14567, 0, -6172}

	var i int
	for ; i < len(expected); i++ {
		if expected[i] != given[i] {
			t.Errorf("expected %d does not equal given %d", expected[i], given[i])
		}
	}
}

func TestScaleBetween(t *testing.T) {
	given := utils.ScaleBetween([]int16{-4, 0, 5, 6, 9}, 0, 100)
	expected := []int16{44, 0, 55, 66, 100}

	for i, _ := range expected {
		if expected[i] != given[i] {
			t.Errorf("expected %d does not equal given %d", expected[i], given[i])
		}
	}
}
