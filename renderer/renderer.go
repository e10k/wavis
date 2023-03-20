package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"wav/parser"
)

func ToBlobSvg(wav *parser.Wav, amplitudes []int16, outputWidthPx int, outputHeightPx int, resolution int) string {
	if resolution == 0 {
		resolution = 5
	}

	amplitudesLen := len(amplitudes)

	if resolution > amplitudesLen {
		resolution = amplitudesLen
	}

	chunksCount := amplitudesLen / (int(wav.SampleRate / int32(resolution)))

	samplesPerChunk := amplitudesLen / chunksCount

	var output []int16
	var chunks [][]int16
	for i := 0; i < amplitudesLen; i += samplesPerChunk {
		end := i + samplesPerChunk
		if end > amplitudesLen {
			end = amplitudesLen
		}

		chunks = append(chunks, amplitudes[i:end])
	}

	for _, c := range chunks {
		var maxInChunk int16
		for _, s := range c {
			if s > maxInChunk {
				maxInChunk = s
			}
		}

		//fmt.Printf("maxInChunk: %d\n\n", maxInChunk)

		output = append(output, maxInChunk)
	}

	//fmt.Printf("amplitudesCount: %d, chunksCount: %d, samplesPerChunk: %d, resulted chunks: %d, first chunk size: %d\n",
	//	amplitudesLen, chunksCount, samplesPerChunk, len(chunks), len(chunks[0]))

	var ypoints []int

	for _, v := range output {
		v = v / 2 // cut in half because all these points will be placed in the svg's upper half
		y := outputHeightPx/2 - int(v)
		ypoints = append(ypoints, y)
	}

	/*	// mirror the points and add them to the slice
		var mirrored_ypoints []int
		for _, v := range ypoints {
			mirrored_ypoints = append([]int{outputHeightPx - v}, mirrored_ypoints...)
		}

		ypoints = append(ypoints, mirrored_ypoints...)*/

	type point struct {
		X float64
		Y float64
	}

	var points []point
	var xstep float64

	xstep = float64(outputWidthPx) / float64(chunksCount)

	for i, v := range ypoints {
		points = append(points, point{
			X: float64(i) * xstep,
			Y: float64(v),
		})
	}

	// now mirror the points
	var mirroredPoints []point
	for i := len(points) - 1; i >= 0; i-- {
		p := points[i]
		p.Y = float64(outputHeightPx) - points[i].Y
		mirroredPoints = append(mirroredPoints, p)
	}

	points = append(points, mirroredPoints...)

	// round the points coordinates
	for i := 0; i < len(points); i++ {
		points[i].X = math.Round(points[i].X)
		points[i].Y = math.Round(points[i].Y)
	}

	var pathData bytes.Buffer
	pathData.WriteString(fmt.Sprintf("M %d %d", int(math.Round(points[0].X)), int(math.Round(points[0].Y))))
	for i := 0; i < len(points)-1; i++ {
		xMid := math.Round((points[i].X + points[i+1].X) / 2)
		yMid := math.Round((points[i].Y + points[i+1].Y) / 2)
		cpX1 := math.Round((xMid + points[i].X) / 2)
		cpX2 := math.Round((xMid + points[i+1].X) / 2)

		pathData.WriteString(fmt.Sprintf("Q %d %d %d %d", int(cpX1), int(points[i].Y), int(xMid), int(yMid)))
		pathData.WriteString(fmt.Sprintf("Q %d %d %d %d", int(cpX2), int(points[i+1].Y), int(points[i+1].X), int(points[i+1].Y)))
	}

	type svg struct {
		Width    int
		Height   int
		Points   []point
		PathData string
	}

	svgStruct := svg{
		Width:    outputWidthPx,
		Height:   outputHeightPx,
		Points:   points,
		PathData: pathData.String(),
	}

	svgTemplate := `
<html>
<body>
	<svg width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}" xmlns="http://www.w3.org/2000/svg">
		<path d="{{ .PathData }}" fill="none" stroke="red" stroke-width="1"/>
<!--
		{{range .Points}}<circle cx="{{.X}}" cy="{{.Y}}" r="2"></circle>
		{{end}}
-->
	</svg>
</body></html>`

	var tpl bytes.Buffer
	tmpl, err := template.New("svg").Parse(svgTemplate)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(&tpl, svgStruct)

	return tpl.String()
}

func ToSingleLineSvg(wav *parser.Wav, amplitudes []int16, outputWidthPx int, outputHeightPx int, resolution int) string {
	if resolution == 0 {
		resolution = 5
	}

	amplitudesLen := len(amplitudes)

	if resolution > amplitudesLen {
		resolution = amplitudesLen
	}

	chunksCount := amplitudesLen / (int(wav.SampleRate / int32(resolution)))

	samplesPerChunk := amplitudesLen / chunksCount

	var output []int16
	var chunks [][]int16
	for i := 0; i < amplitudesLen; i += samplesPerChunk {
		end := i + samplesPerChunk
		if end > amplitudesLen {
			end = amplitudesLen
		}

		chunks = append(chunks, amplitudes[i:end])
	}

	for _, c := range chunks {
		var maxInChunk int16
		for _, s := range c {
			if s > maxInChunk {
				maxInChunk = s
			}
		}

		//fmt.Printf("maxInChunk: %d\n\n", maxInChunk)

		output = append(output, maxInChunk)
	}

	//fmt.Printf("amplitudesCount: %d, chunksCount: %d, samplesPerChunk: %d, resulted chunks: %d, first chunk size: %d\n",
	//	amplitudesLen, chunksCount, samplesPerChunk, len(chunks), len(chunks[0]))

	var ypoints []int

	for index, v := range output {
		v = v / 2
		modifier := -1
		if index%2 == 0 {
			modifier *= -1
		}
		y := outputHeightPx/2 + int(v)*modifier
		ypoints = append(ypoints, y)
	}

	type point struct {
		X float64
		Y float64
	}

	var points []point
	var xstep float64

	xstep = float64(outputWidthPx) / float64(chunksCount)

	for i, v := range ypoints {
		points = append(points, point{
			X: float64(i) * xstep,
			Y: float64(v),
		})
	}

	// round the points coordinates
	for i := 0; i < len(points); i++ {
		points[i].X = math.Round(points[i].X)
		points[i].Y = math.Round(points[i].Y)
	}

	var pathData bytes.Buffer
	pathData.WriteString(fmt.Sprintf("M %d %d", int(math.Round(points[0].X)), int(math.Round(points[0].Y))))
	for i := 0; i < len(points)-1; i++ {
		xMid := math.Round((points[i].X + points[i+1].X) / 2)
		yMid := math.Round((points[i].Y + points[i+1].Y) / 2)
		cpX1 := math.Round((xMid + points[i].X) / 2)
		cpX2 := math.Round((xMid + points[i+1].X) / 2)

		pathData.WriteString(fmt.Sprintf("Q %d %d %d %d", int(cpX1), int(points[i].Y), int(xMid), int(yMid)))
		pathData.WriteString(fmt.Sprintf("Q %d %d %d %d", int(cpX2), int(points[i+1].Y), int(points[i+1].X), int(points[i+1].Y)))
	}

	type svg struct {
		Width    int
		Height   int
		Points   []point
		PathData string
	}

	svgStruct := svg{
		Width:    outputWidthPx,
		Height:   outputHeightPx,
		Points:   points,
		PathData: pathData.String(),
	}

	svgTemplate := `
<html>
<body>
	<svg width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}" xmlns="http://www.w3.org/2000/svg">
		<path d="{{ .PathData }}" fill="none" stroke="red" stroke-width="1"/>
<!--
		{{range .Points}}<circle cx="{{.X}}" cy="{{.Y}}" r="2"></circle>
		{{end}}
-->
	</svg>
</body></html>`

	var tpl bytes.Buffer
	tmpl, err := template.New("svg").Parse(svgTemplate)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(&tpl, svgStruct)

	return tpl.String()
}

func ToRadialSvg(wav *parser.Wav, amplitudes []int16, outputWidthPx int, outputHeightPx int, radius int, resolution int) string {
	if resolution == 0 {
		resolution = 5
	}

	amplitudesLen := len(amplitudes)

	if resolution > amplitudesLen {
		resolution = amplitudesLen
	}

	chunksCount := amplitudesLen / (int(wav.SampleRate / int32(resolution)))

	samplesPerChunk := amplitudesLen / chunksCount

	var output []int16
	var chunks [][]int16
	for i := 0; i < amplitudesLen; i += samplesPerChunk {
		end := i + samplesPerChunk
		if end > amplitudesLen {
			end = amplitudesLen
		}

		chunks = append(chunks, amplitudes[i:end])
	}

	for _, c := range chunks {
		var maxInChunk int16
		for _, s := range c {
			if s > maxInChunk {
				maxInChunk = s
			}
		}

		//fmt.Printf("maxInChunk: %d\n\n", maxInChunk)

		output = append(output, maxInChunk)
	}

	//fmt.Printf("amplitudesCount: %d, chunksCount: %d, samplesPerChunk: %d, resulted chunks: %d, first chunk size: %d\n",
	//	amplitudesLen, chunksCount, samplesPerChunk, len(chunks), len(chunks[0]))

	var lengths []int

	for _, v := range output {
		y := int(v)
		lengths = append(lengths, 29+y)
	}

	type point struct {
		X float64
		Y float64
	}

	var points []point
	//var xstep float64

	//xstep = float64(radius) / float64(chunksCount)
	angleIncrement := float64(360) / float64(len(lengths))
	var angle float64 = 270

	for _, l := range lengths {
		points = append(points, point{
			X: float64(l)*math.Cos(math.Pi*float64(angle)/180) + float64(outputWidthPx/2),
			Y: float64(l)*math.Sin(math.Pi*float64(angle)/180) + float64(outputHeightPx/2),
		})

		//fmt.Printf("points: %v, angle: %v, angleIncrement: %v\n", points, angle, angleIncrement)

		angle += angleIncrement
	}

	// round the points coordinates
	for i := 0; i < len(points); i++ {
		points[i].X = math.Round(points[i].X)
		points[i].Y = math.Round(points[i].Y)
	}

	/*	var pathData bytes.Buffer
		pathData.WriteString(fmt.Sprintf("M %d %d", int(math.Round(points[0].X)), int(math.Round(points[0].Y))))
		for i := 0; i < len(points); i++ {
			//xMid := math.Round((points[i].X + points[i+1].X) / 2)
			//yMid := math.Round((points[i].Y + points[i+1].Y) / 2)
			//cpX1 := math.Round((xMid + points[i].X) / 2)
			//cpX2 := math.Round((xMid + points[i+1].X) / 2)

			//pathData.WriteString(fmt.Sprintf("Q %d %d %d %d", int(cpX1), int(points[i].Y), int(xMid), int(yMid)))
			//pathData.WriteString(fmt.Sprintf("Q %d %d %d %d", int(cpX2), int(points[i+1].Y), int(points[i+1].X), int(points[i+1].Y)))
			pathData.WriteString(fmt.Sprintf("L %d %d ", int(points[i].X), int(points[i].Y)))
		}
	*/
	type svg struct {
		Width    int
		Height   int
		Points   []point
		PathData string
	}

	svgStruct := svg{
		Width:  outputWidthPx,
		Height: outputHeightPx,
		Points: points,
		//PathData: pathData.String(),
	}

	svgTemplate := `
<html>
<body>
	<svg width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}" xmlns="http://www.w3.org/2000/svg">
		<!--<path d="{{ .PathData }}" fill="none" stroke="red" stroke-width="1"/> -->
		<!--{{range .Points}}<circle cx="{{.X}}" cy="{{.Y}}" r="2"></circle>
		{{end}}-->
		{{range .Points}}<line x1="250" y1="250" x2="{{.X}}" y2="{{.Y}}" stroke="gray" stroke-width="1"></line>
		{{end}}
		<circle cx="250" cy="250" r="28" fill="white"></circle>
	</svg>
</body></html>`

	var tpl bytes.Buffer
	tmpl, err := template.New("svg").Parse(svgTemplate)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(&tpl, svgStruct)

	return tpl.String()
}
