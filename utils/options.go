package utils

import (
	"flag"
	"fmt"
	"strings"
)

type Options struct {
	Width        *int
	Height       *int
	Padding      *int
	CircleRadius *int
	Chars        *string
	Border       *bool
	Resolution   *int
	Format       *int
}

func (o *Options) GetChars() []string {
	chars := strings.Split(*o.Chars, "")

	l := len(chars)

	if l >= 2 {
		return chars[0:2]
	}

	if l == 1 {
		return []string{chars[0], " "}
	}

	return []string{"•", " "}
}

func (o *Options) Usage(flagSet *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage:\n")
		order := []string{"format", "resolution", "width", "height", "padding", "circle-radius", "chars", "border"}
		for _, name := range order {
			f := flagSet.Lookup(name)
			fmt.Printf("  -%s\n", f.Name)
			fmt.Printf("\t%s\n", f.Usage)
		}
	}
}
