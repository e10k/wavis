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
	width        *int
	height       *int
	padding      *int
	circleRadius *int
	chars        *string
	border       *bool
	resolution   *int
	format       *int
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

var options Options

func init() {
	options.width = flag.Int("width", 0, "output width")
	options.height = flag.Int("height", 0, "output height")
	options.padding = flag.Int("padding", 0, "output padding")
	options.circleRadius = flag.Int("circle-radius", 0, "inner circle radius for radial svg")
	options.chars = flag.String("chars", "* ", "characters to use for the ascii representation")
	options.border = flag.Bool("border", false, "whether the ascii representation should have a border")
	options.resolution = flag.Int("resolution", 0, "data points per second")
	options.format = flag.Int("format", 0, "output format")
}

func main() {
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
		fmt.Println(getSingleLineSvg(wav, &options))
	case 3:
		fmt.Println(getRadialSvg(wav, &options))
	case 4:
		fmt.Println(getAscii(wav, &options))
	default:
		*options.width = 80
		*options.height = 18
		*options.padding = 0
		*options.border = true

		fmt.Println(getInfo(wav, getAscii(wav, &options)))
	}

}

func getBlobSvg(wav *parser.Wav, options *Options) string {
	const (
		defaultWidth      = 800
		defaultHeight     = 300
		defaultResolution = 5
	)

	monoSamples := wav.GetMonoSamples()

	width := *options.width
	if width == 0 {
		width = defaultWidth
	}

	height := *options.height
	if height == 0 {
		height = defaultHeight
	}

	padding := *options.padding

	resolution := *options.resolution
	if resolution == 0 {
		resolution = defaultResolution
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(height-padding))

	svg := renderer.ToBlobSvg(wav, scaledSamples, width, height, resolution)

	return svg
}

func getSingleLineSvg(wav *parser.Wav, options *Options) string {
	const (
		defaultWidth      = 800
		defaultHeight     = 300
		defaultResolution = 10
	)

	monoSamples := wav.GetMonoSamples()

	width := *options.width
	if width == 0 {
		width = defaultWidth
	}

	height := *options.height
	if height == 0 {
		height = defaultHeight
	}

	padding := *options.padding

	resolution := *options.resolution
	if resolution == 0 {
		resolution = defaultResolution
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(height-padding))

	svg := renderer.ToSingleLineSvg(wav, scaledSamples, width, height, resolution)

	return svg
}

func getRadialSvg(wav *parser.Wav, options *Options) string {
	const (
		defaultWidth        = 500
		defaultHeight       = 500
		defaultCircleRadius = 50
		defaultResolution   = 20
	)

	monoSamples := wav.GetMonoSamples()

	width := *options.width
	if width == 0 {
		width = defaultWidth
	}

	height := *options.height
	if height == 0 {
		height = defaultHeight
	}

	padding := *options.padding

	// points per second
	resolution := *options.resolution
	if resolution == 0 {
		resolution = defaultResolution
	}

	circleRadius := *options.circleRadius
	if circleRadius == 0 {
		circleRadius = defaultCircleRadius
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(math.Min(float64(width), float64(height))/2-float64(padding)-float64(circleRadius)))

	svg := renderer.ToRadialSvg(wav, scaledSamples, width, height, circleRadius, resolution)

	return svg
}

func getAscii(wav *parser.Wav, options *Options) string {
	const (
		defaultWidth      = 80
		defaultHeight     = 15
		defaultResolution = 2
	)
	monoSamples := wav.GetMonoSamples()

	width := *options.width
	if width == 0 {
		width = defaultWidth
	}

	height := *options.height
	if height == 0 {
		height = defaultHeight
	}
	if height%2 == 0 {
		height++ // increase it to make it odd so that we can have a middle line
	}

	padding := *options.padding

	border := *options.border

	// points per second
	resolution := *options.resolution
	if resolution == 0 {
		resolution = defaultResolution
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(height/2-padding))

	output := renderer.ToAscii(scaledSamples, width, height, resolution, options.getChars(), border)

	return output
}

func getInfo(wav *parser.Wav, waveform string) string {
	return renderer.ToInfo(wav, waveform)
}
