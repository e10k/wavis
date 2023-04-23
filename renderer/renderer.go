package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"path/filepath"
	"wav/parser"
)

func ToBlobSvg(wav *parser.Wav, amplitudes []int16, width int, height int, resolution int) (string, error) {
	if resolution == 0 {
		resolution = 5
	}

	amplitudesLen := len(amplitudes)

	if resolution > amplitudesLen {
		resolution = amplitudesLen
	}

	chunksCount := amplitudesLen / (int(wav.SampleRate / int32(resolution)))

	samplesPerChunk := amplitudesLen / chunksCount
	if samplesPerChunk == 0 {
		return "", fmt.Errorf("not enough samples")
	}

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

		output = append(output, maxInChunk)
	}

	var ypoints []int

	for _, v := range output {
		v = v / 2 // cut in half because all these points will be placed in the svg's upper half
		y := height/2 - int(v)
		ypoints = append(ypoints, y)
	}

	type point struct {
		X float64
		Y float64
	}

	var points []point
	var xstep float64

	xstep = float64(width) / float64(chunksCount)

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
		p.Y = float64(height) - points[i].Y
		mirroredPoints = append(mirroredPoints, p)
	}

	points = append(points, mirroredPoints...)

	// round the points coordinates
	for i := 0; i < len(points); i++ {
		points[i].X = math.Round(points[i].X)
		points[i].Y = math.Round(points[i].Y)
	}

	var pathData bytes.Buffer
	pathData.WriteString(fmt.Sprintf("M %d %d", int(math.Round(points[0].X)), height/2))

	loopLimit := len(points) - 1
	for i := 0; i < loopLimit; i++ {
		xMid := math.Round((points[i].X + points[i+1].X) / 2)
		yMid := math.Round((points[i].Y + points[i+1].Y) / 2)
		cpX1 := math.Round((xMid + points[i].X) / 2)
		cpX2 := math.Round((xMid + points[i+1].X) / 2)

		pathData.WriteString(fmt.Sprintf("Q %d %d %d %d", int(cpX1), int(points[i].Y), int(xMid), int(yMid)))

		lastY := int(points[i+1].Y)
		if i == loopLimit-1 {
			lastY = height / 2
		}
		pathData.WriteString(fmt.Sprintf("Q %d %d %d %d", int(cpX2), int(points[i+1].Y), int(points[i+1].X), lastY))
	}

	type svg struct {
		Width    int
		Height   int
		Points   []point
		PathData string
	}

	svgStruct := svg{
		Width:    width,
		Height:   height,
		Points:   points,
		PathData: pathData.String(),
	}

	svgTemplate := `<svg width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}" xmlns="http://www.w3.org/2000/svg">
	<path d="{{ .PathData }} Z" fill="none" stroke="red" stroke-width="1"/>
</svg>`

	return getStringFromSvgTemplate(svgTemplate, svgStruct)
}

func ToSingleLineSvg(wav *parser.Wav, amplitudes []int16, width int, height int, resolution int) (string, error) {
	if resolution == 0 {
		resolution = 5
	}

	amplitudesLen := len(amplitudes)

	if resolution > amplitudesLen {
		resolution = amplitudesLen
	}

	chunksCount := amplitudesLen / (int(wav.SampleRate / int32(resolution)))

	samplesPerChunk := amplitudesLen / chunksCount
	if samplesPerChunk == 0 {
		return "", fmt.Errorf("not enough samples")
	}

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

		output = append(output, maxInChunk)
	}

	var ypoints []int

	for index, v := range output {
		v = v / 2
		modifier := -1
		if index%2 == 0 {
			modifier *= -1
		}
		y := height/2 + int(v)*modifier
		ypoints = append(ypoints, y)
	}

	type point struct {
		X float64
		Y float64
	}

	var points []point
	var xstep float64

	xstep = float64(width) / float64(chunksCount)

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
		Width:    width,
		Height:   height,
		Points:   points,
		PathData: pathData.String(),
	}

	svgTemplate := `<svg width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}" xmlns="http://www.w3.org/2000/svg">
	<path d="{{ .PathData }}" fill="none" stroke="red" stroke-width="1"/>
</svg>`

	return getStringFromSvgTemplate(svgTemplate, svgStruct)
}

func ToRadialSvg(wav *parser.Wav, amplitudes []int16, width int, height int, CircleRadius int, resolution int) (string, error) {
	const (
		defaultResolution = 5
	)

	if resolution == 0 {
		resolution = defaultResolution
	}

	amplitudesLen := len(amplitudes)

	if resolution > amplitudesLen {
		resolution = amplitudesLen
	}

	chunksCount := amplitudesLen / (int(wav.SampleRate / int32(resolution)))

	samplesPerChunk := amplitudesLen / chunksCount
	if samplesPerChunk == 0 {
		return "", fmt.Errorf("not enough samples")
	}

	var chunks [][]int16
	for i := 0; i < amplitudesLen; i += samplesPerChunk {
		end := i + samplesPerChunk
		if end > amplitudesLen {
			end = amplitudesLen
		}

		chunks = append(chunks, amplitudes[i:end])
	}

	var output []int16
	for _, c := range chunks {
		var maxInChunk int16
		for _, s := range c {
			if s > maxInChunk {
				maxInChunk = s
			}
		}

		output = append(output, maxInChunk)
	}

	var lengths []int

	for _, v := range output {
		lengths = append(lengths, CircleRadius+int(v))
	}

	type point struct {
		X float64
		Y float64
	}

	var points []point

	angleIncrement := float64(360) / float64(len(lengths))
	var angle float64 = 270

	for _, l := range lengths {
		points = append(points, point{
			X: float64(l)*math.Cos(math.Pi*float64(angle)/180) + float64(width/2),
			Y: float64(l)*math.Sin(math.Pi*float64(angle)/180) + float64(height/2),
		})

		angle += angleIncrement
	}

	// round the points coordinates
	for i := 0; i < len(points); i++ {
		points[i].X = math.Round(points[i].X)
		points[i].Y = math.Round(points[i].Y)
	}

	type svg struct {
		Width        int
		Height       int
		CenterX      int
		CenterY      int
		CircleRadius int
		Points       []point
	}

	svgStruct := svg{
		Width:        width,
		Height:       height,
		CenterX:      width / 2,
		CenterY:      height / 2,
		CircleRadius: CircleRadius,
		Points:       points,
	}

	svgTemplate := `<svg width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}" xmlns="http://www.w3.org/2000/svg">
	{{range .Points}}<line x1="{{$.CenterX}}" y1="{{$.CenterY}}" x2="{{.X}}" y2="{{.Y}}" stroke="red" stroke-width="1"></line>
	{{end}}<circle cx="{{.CenterX}}" cy="{{.CenterY}}" r="{{.CircleRadius}}" fill="white"></circle>
</svg>`

	return getStringFromSvgTemplate(svgTemplate, svgStruct)
}

func ToAscii(amplitudes []int16, width int, height int, chars []string, border bool) (string, error) {
	amplitudesLen := len(amplitudes)

	samplesPerChunk := amplitudesLen / width
	if samplesPerChunk == 0 {
		return "", fmt.Errorf("not enough samples")
	}

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

		output = append(output, maxInChunk)
	}

	var lengths []int

	for _, v := range output {
		lengths = append(lengths, int(v))
	}

	m := height/2 + 1
	var b bytes.Buffer

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			penDown := y >= m-lengths[x]-1 && y < m+lengths[x]

			if border {
				if x == 0 && y == 0 {
					b.WriteRune('╭')
				} else if x == 0 && y == height-1 {
					b.WriteRune('╰')
				} else if x == width-1 && y == 0 {
					b.WriteRune('╮')
				} else if x == width-1 && y == height-1 {
					b.WriteRune('╯')
				} else if y == 0 || y == height-1 {
					b.WriteRune('─')
				} else if x == 0 || x == width-1 {
					b.WriteRune('│')
				} else if penDown {
					b.WriteString(chars[0])
				} else {
					b.WriteString(chars[1])
				}
			} else if penDown {
				b.WriteString(chars[0])
			} else {
				b.WriteString(chars[1])
			}
		}

		if y < height-1 {
			b.WriteByte('\n')
		}
	}

	return b.String(), nil
}

func ToInfo(wav *parser.Wav, waveform string) string {
	var b bytes.Buffer

	b.WriteByte('\n')
	b.WriteString(fmt.Sprintf("File:\t\t%s\n", filepath.Base(wav.Name)))
	b.WriteString(fmt.Sprintf("Channels:\t%d\n", wav.NumChannels))
	b.WriteString(fmt.Sprintf("Sample Rate:\t%d\n", wav.SampleRate))
	b.WriteString(fmt.Sprintf("Precision:\t%d-bit\n", wav.BitsPerSample))
	b.WriteString(fmt.Sprintf("Byte Rate:\t%d\n", wav.ByteRate))
	b.WriteString(fmt.Sprintf("Duration:\t%s\n", wav.GetFormattedDuration()))
	b.WriteString(fmt.Sprintf("File Size:\t%d", wav.GetFileSize()))

	if len(waveform) > 0 {
		b.WriteString("\n\n")
		b.WriteString(waveform)
	}

	return b.String()
}

func getStringFromSvgTemplate(svgTemplate string, svgStruct interface{}) (string, error) {
	var tpl bytes.Buffer
	tmpl, err := template.New("svg").Parse(svgTemplate)
	if err != nil {
		return "", fmt.Errorf("template error: %v", err)
	}
	if err = tmpl.Execute(&tpl, svgStruct); err != nil {
		return "", fmt.Errorf("template error: %v", err)
	}

	return tpl.String(), nil
}
