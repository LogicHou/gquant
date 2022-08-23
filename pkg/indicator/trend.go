package indicator

func Sum(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := 0.0

	for i := 0; i < len(values); i++ {
		sum += values[i]

		if i >= period {
			sum -= values[i-period]
		}

		result[i] = sum
	}

	return result
}

func multiply(values1, values2 []float64) []float64 {
	checkSameSize(values1, values2)

	result := make([]float64, len(values1))

	for i := 0; i < len(result); i++ {
		result[i] = values1[i] * values2[i]
	}

	return result
}
