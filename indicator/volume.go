package indicator

func Vwap(period int, values []float64, volume []float64) []float64 {
	return divide(Sum(period, multiply(values, volume)), Sum(period, volume))
}
