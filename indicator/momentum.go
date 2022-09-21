package indicator

func RsiPeriod(period int, closing []float64) ([]float64, []float64) {
	gains := make([]float64, len(closing))
	losses := make([]float64, len(closing))

	for i := 1; i < len(closing); i++ {
		difference := closing[i] - closing[i-1]

		if difference > 0 {
			gains[i] = difference
			losses[i] = 0
		} else {
			losses[i] = -difference
			gains[i] = 0
		}
	}

	meanGains := Rma(period, gains)
	meanLosses := Rma(period, losses)

	rsi := make([]float64, len(closing))
	rs := make([]float64, len(closing))

	for i := 0; i < len(rsi); i++ {
		rs[i] = meanGains[i] / meanLosses[i]
		rsi[i] = 100 - (100 / (1 + rs[i]))
	}

	return rs, rsi
}
