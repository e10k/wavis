package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"wav/parser"
	"wav/renderer"
	"wav/utils"
)

var wav *parser.Wav

type Options struct {
	width      *int
	height     *int
	resolution *int
	normalise  *bool
}

func main() {
	var options Options
	options.width = flag.Int("width", 1000, "output width")
	options.height = flag.Int("height", 400, "output height")
	options.resolution = flag.Int("resolution", 10, "data points per second")
	options.normalise = flag.Bool("normalise", false, "normalise")

	flag.Parse()

	filename := flag.Arg(0)
	if len(filename) < 1 {
		log.Fatal("No .wav file provided.")
	}

	f, err := os.Open(filename)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal("Failed closing the file.")
		}
	}(f)
	if err != nil {
		log.Fatal(err)
	}

	wav = parser.Parse(f)

	fmt.Println(getSvg(wav, &options))
}

func getSvg(wav *parser.Wav, options *Options) string {
	monoSamples := wav.GetMonoSamples()

	width := *options.width
	if width == 0 {
		width = 200
	}

	height := *options.height
	if height == 0 {
		height = 100
	}

	// points per second
	resolution := *options.resolution
	if resolution == 0 {
		resolution = 2
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(height))

	svg := renderer.ToSvg(wav, scaledSamples, width, height, resolution)

	return svg
}
