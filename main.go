package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"wav/parser"
	"wav/renderer"
	"wav/utils"
)

var wav *parser.Wav

type Options struct {
	width      *int
	height     *int
	resolution *int
}

func main() {
	var options Options
	options.width = flag.Int("width", 500, "output width")
	options.height = flag.Int("height", 100, "output height")
	options.resolution = flag.Int("resolution", 5, "data points per second")

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

	/*monoSamples := wav.GetMonoSamples()

	outputWidthPx := 200
	outputHeightPx := 100
	outputSlicePx := 5

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int32(outputHeightPx)/2)

	//fmt.Printf("scaledSamples: %#v", scaledSamples)

	renderer.ToSvg(scaledSamples, outputWidthPx, outputHeightPx, outputSlicePx)*/

	http.HandleFunc("/", getSvg)

	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getSvg(w http.ResponseWriter, r *http.Request) {
	log.Printf("%#v", r.URL.Query().Get("a"))

	monoSamples := wav.GetMonoSamples()

	width, _ := strconv.Atoi(r.URL.Query().Get("width"))
	if width == 0 {
		width = 200
	}
	height, _ := strconv.Atoi(r.URL.Query().Get("height"))
	if height == 0 {
		height = 100
	}

	// points per second
	resolution, _ := strconv.Atoi(r.URL.Query().Get("resolution"))
	if resolution == 0 {
		resolution = 2
	}

	scaledSamples := utils.ScaleBetween(monoSamples, 0, int16(height)/2)

	svg := renderer.ToSvg(wav, scaledSamples, width, height, resolution)

	io.WriteString(w, svg)
}
