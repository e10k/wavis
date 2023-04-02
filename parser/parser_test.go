package parser

import (
	"bytes"
	"testing"
)

func TestReading24BitSamples(t *testing.T) {
	// little-endian order
	data := []byte{
		0b01111111, 0b00000000, 0b00000000, //127
		0b10000001, 0b11111111, 0b11111111, // -127
		0b00000000, 0b00000000, 0b00000000, // 0
		0b00000000, 0b10000000, 0b00000000, // 32768
		0b00000000, 0b10000000, 0b11111111, // -32768
		0b10000000, 0b00000000, 0b00000000, // 128
		0b00000000, 0b00000000, 0b01000000, // 4194304
		0b11111111, 0b11111111, 0b10111111, // -4194304 - 1
		0b00000000, 0b00000000, 0b11000000, // -4194304
		0b11111111, 0b11111111, 0b11111111, // -1
	}

	r := bytes.NewReader(data)

	expected := []int32{127, -127, 0, 32768, -32768, 128, 4194304, -4194304 - 1, -4194304, -1}

	for i, _ := range expected {
		sample, _ := read24BitSample(r)
		if expected[i] != sample {
			t.Errorf("expected %d does not equal given %d", expected[i], sample)
		}
	}
}
