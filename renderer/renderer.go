package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

func ToSvg(amplitudes []int32, outputWidthPx int, outputHeightPx int, chunksCount int) string {
	if chunksCount == 0 {
		chunksCount = 1
	}

	amplitudesLen := len(amplitudes)

	if chunksCount > amplitudesLen {
		chunksCount = amplitudesLen
	}

	samplesPerChunk := amplitudesLen / chunksCount

	var output []int32
	var chunks [][]int32
	for i := 0; i < amplitudesLen; i += samplesPerChunk {
		end := i + samplesPerChunk
		if end > amplitudesLen {
			end = amplitudesLen
		}

		chunks = append(chunks, amplitudes[i:end])
	}

	for _, c := range chunks {
		var maxInChunk int32
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
		X float32
		Y int
	}

	var points []point
	var xstep float32

	xstep = float32(outputWidthPx) / float32(chunksCount)

	for i, v := range ypoints {
		points = append(points, point{
			X: float32(i) * xstep,
			Y: v,
		})
	}

	// now mirror the points
	var mirroredPoints []point
	for i := len(points) - 1; i >= 0; i-- {
		p := points[i]
		p.Y = outputHeightPx - points[i].Y
		mirroredPoints = append(mirroredPoints, p)
	}

	points = append(points, mirroredPoints...)

	var pathData bytes.Buffer
	pathData.WriteString(fmt.Sprintf("M %f %d", points[0].X, points[0].Y))
	for i := 0; i < len(points) - 1; i++ {
		xMid := (points[i].X + points[i+1].X) / 2
		yMid := (points[i].Y + points[i+1].Y) / 2
		cpX1 := (xMid + points[i].X) / 2
		cpX2 := (xMid + points[i+1].X) / 2

		pathData.WriteString(fmt.Sprintf("Q %f %d %f %d", cpX1, points[i].Y, xMid, yMid))
		pathData.WriteString(fmt.Sprintf("Q %f %d %f %d", cpX2, points[i+1].Y, points[i+1].X, points[i+1].Y))
	}


	type svg struct {
		Width  int
		Height int
		Points []point
		PathData string
	}

	svgStruct := svg{
		Width:  outputWidthPx,
		Height: outputHeightPx,
		Points: points,
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
