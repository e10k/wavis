package utils

import (
	"strings"
	"testing"
)

func TestOptions_GetChars(t *testing.T) {
	given := []string{"! ", "✨・", "abcd"}

	var expected [][]string
	expected = append(expected, []string{"!", " "}, []string{"✨", "・"}, []string{"a", "b"})

	var options Options
	var chars []string
	for i, _ := range given {
		options.Chars = &given[i]
		chars = options.GetChars()

		if len(chars) != 2 {
			t.Errorf("the chars length should be 2")
		}

		givenSplit := strings.Split(given[i], "")
		if expected[i][0] != givenSplit[0] || expected[i][1] != givenSplit[1] {
			t.Errorf("expected {%s, %s} do not equal the given {%s, %s}", expected[i][0], expected[i][1], givenSplit[0], givenSplit[1])
		}
	}

	var testChars string

	testChars = ""
	options.Chars = &testChars
	chars = options.GetChars()
	if "*" != chars[0] || " " != chars[1] {
		t.Errorf("expected {*, } do not equal the given {%s, %s}", chars[0], chars[1])
	}

	testChars = "@"
	options.Chars = &testChars
	chars = options.GetChars()
	if "@" != chars[0] || " " != chars[1] {
		t.Errorf("expected {@, } do not equal the given {%s, %s}", chars[0], chars[1])
	}
}
