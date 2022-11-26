package main

import (
	"testing"
	"wav/parser"
	"wav/utils"
)

func TestToMonoSamples(t *testing.T) {
	w := &parser.Wav{
		Data: [][]int32{
			{0, 0, 5, 5, 12345, 12345, -12345},
			{0, 1, 7, 8, 67890, -12345, 0},
		},
	}

	given := w.GetMonoSamples()
	expected := []int32{0, 0, 6, 6, 40117, 0, -6172}

	var i int
	for ; i < len(expected); i++ {
		if expected[i] != given[i] {
			t.Errorf("expected %d does not equal given %d", expected[i], given[i])
		}
	}
}

func TestScaleBetween(t *testing.T) {
	given := utils.ScaleBetween([]int32{-4, 0, 5, 6, 9}, 0, 100)
	expected := []int32{44, 0, 55, 66, 100}

	for i, _ := range expected {
		if expected[i] != given[i] {
			t.Errorf("expected %d does not equal given %d", expected[i], given[i])
		}
	}
}
