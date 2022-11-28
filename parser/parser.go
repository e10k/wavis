package parser

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"math"
	"os"
)

type Wav struct {
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
	} else if sampleSize == 64 && *audioFormat == 3 { // IEEE_FLOAT
		var sample float64
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

	// TODO can this be improved? get rid of nested loop and break out?
out:
	for {
		var i int16
		for ; i < wav.NumChannels; i++ {
			sample, err := readSample(r, int(wav.BitsPerSample), &wav.AudioFormat)
			if err != nil {
				if err == io.EOF {
					break out
				}

				log.Fatal(err)
			}

			wav.Data[i] = append(wav.Data[i], sample)
		}
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

func scaleToInt16(v interface{}) int16 {
	var input int64
	var inputMin int64
	var inputMax int64

	switch t := v.(type) {
	case uint8:
		input = int64(t)
		inputMin = int64(0)
		inputMax = int64(math.MaxUint8)
	case int16:
		input = int64(t)
		inputMin = int64(math.MinInt16)
		inputMax = int64(math.MaxInt16)
	case int32:
		input = int64(t)
		inputMin = int64(math.MinInt32)
		inputMax = int64(math.MaxInt32)
	case int64:
		input = t
		inputMin = int64(math.MinInt64)
		inputMax = int64(math.MaxInt64)
	case float32:
		input = int64(t)
		inputMin = int64(-1 << 24)
		inputMax = int64(1<<24 - 1)
	case float64:
		input = int64(t)
		inputMin = int64(-1 << 53)
		inputMax = int64(1<<53 - 1)
	default:
		log.Fatal("unsupported type")
	}

	return int16(scale(input, inputMin, inputMax, int64(math.MinInt16), int64(math.MaxInt16)))
}

func scale(input int64, inputMin int64, inputMax int64, outputMin int64, outputMax int64) int64 {
	return (input-inputMin)*(outputMax-outputMin)/(inputMax-inputMin) + outputMin
}
