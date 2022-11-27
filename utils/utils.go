package utils

func ScaleBetween(numbers []int16, scaledMin, scaledMax int16) []int16 {
	// first make all numbers positive
	for i, v := range numbers {
		if v < 0 {
			numbers[i] = -v
		}
	}

	var inputMin, inputMax int32

	var numbersI32 []int32

	for _, v := range numbers {
		numbersI32 = append(numbersI32, int32(v))
	}

	for _, v := range numbersI32 {
		if v < inputMin {
			inputMin = v
		}
		if v > inputMax {
			inputMax = v
		}
	}

	var scaledSamples []int16

	for _, v := range numbersI32 {
		scaledValue := (int32(scaledMax)-int32(scaledMin))*(v-inputMin)/(inputMax-inputMin) + int32(scaledMin)
		scaledSamples = append(scaledSamples, int16(scaledValue))
	}

	return scaledSamples
}
