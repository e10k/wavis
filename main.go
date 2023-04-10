package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"wav/parser"
	"wav/renderer"
	"wav/utils"
)

var wav *parser.Wav

var options utils.Options

func init() {
	options.Width = flag.Int("width", 0, "output width")
	options.Height = flag.Int("height", 0, "output height")
	options.Padding = flag.Int("padding", 0, "output padding")
	options.CircleRadius = flag.Int("circle-radius", 0, "inner circle radius for radial svg")
	options.Chars = flag.String("chars", "* ", "characters to use for the ascii representation")
	options.Border = flag.Bool("border", false, "whether the ascii representation should have a border")
	options.Resolution = flag.Int("resolution", 0, "data points per second")
	options.Format = flag.Int("format", 0, "output format")

	flag.Usage = options.Usage(flag.CommandLine)
}

func main() {
	flag.Parse()

	filename := flag.Arg(0)
	ext := strings.ToLower(filepath.Ext(filename))
	if len(filename) < 1 || ext != ".wav" {
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

	switch *options.Format {
	case 1:
		fmt.Println(getBlobSvg(wav, &options))
	case 2:
		fmt.Println(getSingleLineSvg(wav, &options))
	case 3:
		fmt.Println(getRadialSvg(wav, &options))
	case 4:
		fmt.Println(getAscii(wav, &options))
	default:
		*options.Width = 80
		*options.Height = 18
		*options.Padding = 0
		*options.Border = true

		fmt.Println(getInfo(wav, getAscii(wav, &options)))
	}

}

func getBlobSvg(wav *parser.Wav, options *utils.Options) string {
	const (
		defaultWidth      = 800
		defaultHeight     = 300
		defaultResolution = 5
	)

	monoSamples := wav.GetMonoSamples()

	width := *options.Width
	if width == 0 {
		width = defaultWidth
	}

	height := *options.Height
	if height == 0 {
		height = defaultHeight
	}

	padding := *options.Padding

	resolution := *options.Resolution
	if resolution == 0 {
		resolution = defaultResolution
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(height-padding))

	svg := renderer.ToBlobSvg(wav, scaledSamples, width, height, resolution)

	return svg
}

func getSingleLineSvg(wav *parser.Wav, options *utils.Options) string {
	const (
		defaultWidth      = 800
		defaultHeight     = 300
		defaultResolution = 10
	)

	monoSamples := wav.GetMonoSamples()

	width := *options.Width
	if width == 0 {
		width = defaultWidth
	}

	height := *options.Height
	if height == 0 {
		height = defaultHeight
	}

	padding := *options.Padding

	resolution := *options.Resolution
	if resolution == 0 {
		resolution = defaultResolution
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(height-padding))

	svg := renderer.ToSingleLineSvg(wav, scaledSamples, width, height, resolution)

	return svg
}

func getRadialSvg(wav *parser.Wav, options *utils.Options) string {
	const (
		defaultWidth        = 500
		defaultHeight       = 500
		defaultCircleRadius = 50
		defaultResolution   = 20
	)

	monoSamples := wav.GetMonoSamples()

	width := *options.Width
	if width == 0 {
		width = defaultWidth
	}

	height := *options.Height
	if height == 0 {
		height = defaultHeight
	}

	padding := *options.Padding

	// points per second
	resolution := *options.Resolution
	if resolution == 0 {
		resolution = defaultResolution
	}

	circleRadius := *options.CircleRadius
	if circleRadius == 0 {
		circleRadius = defaultCircleRadius
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(math.Min(float64(width), float64(height))/2-float64(padding)-float64(circleRadius)))

	svg := renderer.ToRadialSvg(wav, scaledSamples, width, height, circleRadius, resolution)

	return svg
}

func getAscii(wav *parser.Wav, options *utils.Options) string {
	const (
		defaultWidth      = 80
		defaultHeight     = 15
		defaultResolution = 2
	)
	monoSamples := wav.GetMonoSamples()

	width := *options.Width
	if width == 0 {
		width = defaultWidth
	}

	height := *options.Height
	if height == 0 {
		height = defaultHeight
	}
	if height%2 == 0 {
		height++ // increase it to make it odd so that we can have a middle line
	}

	padding := *options.Padding

	border := *options.Border

	// points per second
	resolution := *options.Resolution
	if resolution == 0 {
		resolution = defaultResolution
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(height/2-padding))

	output := renderer.ToAscii(scaledSamples, width, height, resolution, options.GetChars(), border)

	return output
}

func getInfo(wav *parser.Wav, waveform string) string {
	return renderer.ToInfo(wav, waveform)
}
