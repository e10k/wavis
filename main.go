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
	padding    *int
	resolution *int
	format     *int
}

func main() {
	var options Options
	options.width = flag.Int("width", 1000, "output width")
	options.height = flag.Int("height", 400, "output height")
	options.padding = flag.Int("padding", 40, "output vertical padding")
	options.resolution = flag.Int("resolution", 10, "data points per second")
	options.format = flag.Int("format", 1, "output format") // 1 is symmetrical svg

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

	switch *options.format {
	case 1:
		fmt.Println(getSvg(wav, &options))
	case 2:
		fmt.Println("true shape (not symmetrical), not implemented")
	case 3:
		fmt.Println("single line svg, not implemented")
	case 4:
		fmt.Println("radial svg, not implemented")
	case 5:
		fmt.Println("city skyline, not implemented")
	case 6:
		fmt.Println("ascii, not implemented")
	default:
		fmt.Println("output waveform information, no visualisation, not implemented")
	}

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

	padding := *options.padding

	// points per second
	resolution := *options.resolution
	if resolution == 0 {
		resolution = 2
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(height-padding))

	svg := renderer.ToSvg(wav, scaledSamples, width, height, resolution)

	return svg
}
