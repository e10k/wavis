package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"wav/parser"
	"wav/renderer"
	"wav/utils"
)

var wav *parser.Wav

type Options struct {
	width       *int
	height      *int
	padding     *int
	innerRadius *int
	chars       *string
	resolution  *int
	format      *int
}

func (o *Options) getChars() []string {
	chars := strings.Split(*o.chars, "")

	l := len(chars)

	if l >= 2 {
		return chars[0:2]
	}

	if l == 1 {
		return []string{chars[0], " "}
	}

	return []string{"*", " "}
}

func main() {
	var options Options
	options.width = flag.Int("width", 1000, "output width")
	options.height = flag.Int("height", 400, "output height")
	options.padding = flag.Int("padding", 40, "output vertical padding")
	options.innerRadius = flag.Int("inner-radius", 40, "inner radius for radial svg")
	options.chars = flag.String("chars", "* ", "characters to use for the ascii representation")
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
		fmt.Println(getBlobSvg(wav, &options))
	case 2:
		fmt.Println("true shape (not symmetrical), not implemented; probably useless? tbd")
	case 3:
		fmt.Println(getSingleLineSvg(wav, &options))
	case 4:
		fmt.Println(getRadialSvg(wav, &options))
	case 5:
		fmt.Println("city skyline, not implemented")
	case 6:
		fmt.Println(getAscii(wav, &options))
	default:
		fmt.Println("output waveform information, no visualisation, not implemented")
	}

}

func getBlobSvg(wav *parser.Wav, options *Options) string {
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

	svg := renderer.ToBlobSvg(wav, scaledSamples, width, height, resolution)

	return svg
}

func getSingleLineSvg(wav *parser.Wav, options *Options) string {
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

	svg := renderer.ToSingleLineSvg(wav, scaledSamples, width, height, resolution)

	return svg
}

func getRadialSvg(wav *parser.Wav, options *Options) string {
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

	innerRadius := *options.innerRadius
	if innerRadius == 0 {
		innerRadius = 29
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(math.Min(float64(width), float64(height))/2-float64(padding)-float64(innerRadius)))

	svg := renderer.ToRadialSvg(wav, scaledSamples, width, height, innerRadius, resolution)

	return svg
}

func getAscii(wav *parser.Wav, options *Options) string {
	monoSamples := wav.GetMonoSamples()

	width := *options.width
	if width == 0 {
		width = 80
	}

	height := *options.height
	if height == 0 {
		height = 15
	}
	if height%2 == 0 {
		height++ // make it odd so that we can have a middle line
	}

	padding := *options.padding

	// points per second
	resolution := *options.resolution
	if resolution == 0 {
		resolution = 2
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(height/2-padding))

	svg := renderer.ToAscii(scaledSamples, width, height, resolution, options.getChars())

	return svg
}
