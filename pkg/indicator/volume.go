package indicator

func Vwap(period int, closing []float64, volume []float64) []float64 {
	return divide(Sum(period, multiply(closing, volume)), Sum(period, volume))
}
