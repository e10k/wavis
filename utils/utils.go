package utils

func ScaleBetween(numbers []int32, scaledMin, scaledMax int32) []int32 {
	// first make all numbers positive
	for i, v := range numbers {
		if v < 0 {
			numbers[i] = -v
		}
	}

	var min, max int32

	for _, v := range numbers {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	var scaledSamples []int32

	for _, v := range numbers {
		scaledValue := (scaledMax-scaledMin)*(v-min)/(max-min) + scaledMin
		scaledSamples = append(scaledSamples, scaledValue)
	}

	return scaledSamples
}
