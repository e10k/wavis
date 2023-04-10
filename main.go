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
		log.Fatal("no .wav file provided")
	}

	f, err := os.Open(filename)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal("failed closing the file")
		}
	}(f)
	if err != nil {
		log.Fatalf("failed opening the file: %v", err)
	}

	wav, err = parser.Parse(f)
	if err != nil {
		log.Fatalf("error parsing the file: %v", err)
	}

	switch *options.Format {
	case 1:
		if s, err := getBlobSvg(wav, &options); err != nil {
			log.Fatalf("error creating blob svg: %v", err)
		} else {
			fmt.Println(s)
		}
	case 2:
		if s, err := getSingleLineSvg(wav, &options); err != nil {
			log.Fatalf("error creating single line svg: %v", err)
		} else {
			fmt.Println(s)
		}
	case 3:
		if s, err := getRadialSvg(wav, &options); err != nil {
			log.Fatalf("error creating radial svg: %v", err)
		} else {
			fmt.Println(s)
		}
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

func getBlobSvg(wav *parser.Wav, options *utils.Options) (string, error) {
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

	svg, err := renderer.ToBlobSvg(wav, scaledSamples, width, height, resolution)

	return svg, err
}

func getSingleLineSvg(wav *parser.Wav, options *utils.Options) (string, error) {
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

	svg, err := renderer.ToSingleLineSvg(wav, scaledSamples, width, height, resolution)

	return svg, err
}

func getRadialSvg(wav *parser.Wav, options *utils.Options) (string, error) {
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

	svg, err := renderer.ToRadialSvg(wav, scaledSamples, width, height, circleRadius, resolution)

	return svg, err
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
