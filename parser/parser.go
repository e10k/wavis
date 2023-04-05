package parser

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

type Wav struct {
	Name          string
	ChunkID       [4]byte
	ChunkSize     int32
	Format        [4]byte
	Subchunk1ID   [4]byte
	Subchunk1Size int32
	AudioFormat   int16
	NumChannels   int16
	SampleRate    int32
	ByteRate      int32
	BlockAlign    int16
	BitsPerSample int16
	Subchunk2ID   [4]byte
	Subchunk2Size int32
	Data          [][]int16
}

func readSample(r io.Reader, sampleSize int, audioFormat *int16) (int16, error) {
	if sampleSize == 8 {
		var sample uint8
		err := binary.Read(r, binary.LittleEndian, &sample)
		if err != nil {
			return int16(0), err
		}

		return scaleToInt16(sample), nil
	} else if sampleSize == 16 {
		var sample int16
		err := binary.Read(r, binary.LittleEndian, &sample)
		if err != nil {
			return int16(0), err
		}

		return scaleToInt16(sample), nil
	} else if sampleSize == 24 {
		sample, err := read24BitSample(r)
		if err != nil {
			return int16(0), err
		}

		return scaleToInt16(sample), nil
	} else if sampleSize == 32 && *audioFormat == 1 { // PCM
		var sample int32
		err := binary.Read(r, binary.LittleEndian, &sample)
		if err != nil {
			return int16(0), err
		}

		return scaleToInt16(sample), nil
	} else if sampleSize == 32 && *audioFormat == 3 { // IEEE_FLOAT
		var sample float32
		err := binary.Read(r, binary.LittleEndian, &sample)
		if err != nil {
			return int16(0), err
		}

		return scaleToInt16(sample), nil
	} else {
		return int16(0), errors.New("invalid sample size")
	}
}

func Parse(f *os.File) *Wav {
	var wav Wav
	wav.Name = f.Name()

	r := bufio.NewReader(f)

	if err := binary.Read(r, binary.BigEndian, &wav.ChunkID); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.ChunkSize); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.BigEndian, &wav.Format); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.BigEndian, &wav.Subchunk1ID); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.Subchunk1Size); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.AudioFormat); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.NumChannels); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.SampleRate); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.ByteRate); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.BlockAlign); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.BitsPerSample); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.BigEndian, &wav.Subchunk2ID); err != nil {
		log.Fatal(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &wav.Subchunk2Size); err != nil {
		log.Fatal(err)
	}

	wav.Data = make([][]int16, wav.NumChannels)

	numSamples := wav.GetNumSamples()

	if numSamples == 0 {
		// this can happen when the wav format is unsupported
		// (e.g. is it 64-bit float and contains a "fact" chunk that messes the parsing above);
		// TODO find better solution for handling this
		panic("could not get the number of samples")
	}

	var i int16
	for s := int32(0); s < numSamples; s++ {
		for ; i < wav.NumChannels; i++ {
			sample, err := readSample(r, int(wav.BitsPerSample), &wav.AudioFormat)
			if err != nil {
				panic(err)
			}

			wav.Data[i] = append(wav.Data[i], sample)
		}

		i = 0
	}

	// all channels should have the same number of channels, but it's not always the case
	// so make sure to trim the longer ones
	minSamples := len(wav.Data[0])
	for i := int16(0); i < wav.NumChannels; i++ {
		l := len(wav.Data[0])
		if l < minSamples {
			minSamples = l
		}
	}

	for i := int16(0); i < wav.NumChannels; i++ {
		wav.Data[i] = wav.Data[i][:minSamples]
	}

	return &wav
}

func (w *Wav) GetMonoSamples() []int16 {
	var monoSamples []int16
	length := len(w.Data[0])

	var i int
	for ; i < length; i++ {
		var sum int32
		var j int
		numChannels := len(w.Data)
		for ; j < numChannels; j++ {
			sum += int32(w.Data[j][i])
		}
		mean := int16(sum / int32(numChannels))

		monoSamples = append(monoSamples, mean)
	}

	return monoSamples
}

func (w *Wav) GetFileSize() int32 {
	return w.ChunkSize + 8
}

func (w *Wav) GetNumSamples() int32 {
	return w.Subchunk2Size / int32(w.NumChannels*w.BitsPerSample/8)
}

func (w *Wav) GetDuration() (float64, int32) {
	numSamples := float64(w.Subchunk2Size / int32(w.NumChannels*w.BitsPerSample/8))

	return numSamples / float64(w.SampleRate), int32(numSamples)
}

func (w *Wav) GetFormattedDuration() string {
	duration, samples := w.GetDuration()

	d := int(duration)
	milliseconds := int((duration - float64(d)) * 1000)

	hours := d / 3600
	minutes := d % 3600 / 60
	seconds := d % 3600 % 60

	return fmt.Sprintf("%02d:%02d:%02d.%03d = %d samples", hours, minutes, seconds, milliseconds, samples)
}

func scaleToInt16(v interface{}) int16 {
	var input float64
	var inputMin float64
	var inputMax float64

	switch t := v.(type) {
	case uint8:
		input = float64(t)
		inputMin = float64(0)
		inputMax = float64(math.MaxUint8)
	case int16:
		input = float64(t)
		inputMin = float64(math.MinInt16)
		inputMax = float64(math.MaxInt16)
	case int32:
		input = float64(t)
		inputMin = float64(math.MinInt32)
		inputMax = float64(math.MaxInt32)
	case int64:
		input = float64(t)
		inputMin = float64(math.MinInt64)
		inputMax = float64(math.MaxInt64)
	case float32:
		input = float64(t)
		inputMin = float64(-1)
		inputMax = float64(1)
	default:
		log.Fatal("unsupported type")
	}

	return int16(scale(input, inputMin, inputMax, float64(math.MinInt16), float64(math.MaxInt16)))
}

func scale(input float64, inputMin float64, inputMax float64, outputMin float64, outputMax float64) float64 {
	return (input-inputMin)*(outputMax-outputMin)/(inputMax-inputMin) + outputMin
}

func read24BitSample(r io.Reader) (int32, error) {
	var buf [3]byte

	err := binary.Read(r, binary.LittleEndian, &buf)
	if err != nil {
		return 0, err
	}

	sample := int32(buf[0]) | (int32(buf[1]) << 8) | (int32(buf[2]) << 16)
	if (sample & (1 << 23)) != 0 {
		sample |= ^((1 << 24) - 1)
	}

	return sample, nil
}
