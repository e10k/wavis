package parser

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"log"
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
	Data          [][]int32
}

func readSample(r io.Reader, sampleSize int) (int32, error) {
	if sampleSize == 8 {
		var sample uint8
		err := binary.Read(r, binary.LittleEndian, &sample)
		if err != nil {
			return int32(0), err
		}

		return int32(sample), nil
	} else if sampleSize == 16 {
		var sample int16
		err := binary.Read(r, binary.LittleEndian, &sample)
		if err != nil {
			return int32(0), err
		}

		return int32(sample), nil
	} else if sampleSize == 32 {
		var sample int32
		err := binary.Read(r, binary.LittleEndian, &sample)
		if err != nil {
			return int32(0), err
		}

		return int32(sample), nil
	} else {
		return int32(0), errors.New("invalid sample size")
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

	wav.Data = make([][]int32, wav.NumChannels)

	// TODO can this be improved? get rid of nested loop and break out?
out:
	for {
		var i int16
		for ; i < wav.NumChannels; i++ {
			sample, err := readSample(r, int(wav.BitsPerSample))
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

func (w *Wav) GetMonoSamples() []int32 {
	var monoSamples []int32
	length := len(w.Data[0])

	var i int
	for ; i < length; i++ {
		var sum int32
		var j int
		numChannels := len(w.Data)
		for ; j < numChannels; j++ {
			sum += w.Data[j][i]
		}
		mean := sum / int32(numChannels)

		monoSamples = append(monoSamples, mean)
	}

	return monoSamples
}
