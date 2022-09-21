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

func Sma(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := float64(0)

	for i, value := range values {
		count := i + 1
		sum += value

		if i >= period {
			sum -= values[i-period]
			count = period
		}

		result[i] = sum / float64(count)
	}

	return result
}

func Rma(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := float64(0)

	for i, value := range values {
		count := i + 1

		if i < period {
			sum += value
		} else {
			sum = (result[i-1] * float64(period-1)) + value
			count = period
		}

		result[i] = sum / float64(count)
	}

	return result
}

func Kdj(rPeriod, kPeriod, dPeriod int, high, low, closing []float64) (k, d, j []float64) {
	clen := len(closing)
	rsv := make([]float64, clen)
	j = make([]float64, clen)

	if rsv[0] == 0 {
		rsv[0] = 50.0
	}
	for i := 1; i < clen; i++ {
		m := i + 1 - rPeriod
		if m < 0 {
			m = 0
		}
		h := maxHigh(high[m : i+1])
		l := minLow(low[m : i+1])
		rsv[i] = (closing[i] - l) * 100.0 / (h - l)
	}

	k = kdjMa(float64(kPeriod), rsv)
	d = kdjMa(float64(dPeriod), k)

	for i := 0; i < clen; i++ {
		j[i] = 3.0*k[i] - 2.0*d[i]
	}
	return
}

func kdjMa(n float64, values []float64) (r []float64) {
	r = make([]float64, len(values))
	for i := 0; i < len(values); i++ {
		if i == 0 {
			r[i] = values[i]
		} else {
			r[i] = (1.0*values[i] + (n-1.0)*r[i-1]) / n
		}
	}
	return
}

func maxHigh(values []float64) (h float64) {
	h = values[0]
	for i := 0; i < len(values); i++ {
		if values[i] > h {
			h = values[i]
		}
	}
	return
}

func minLow(values []float64) (l float64) {
	l = values[0]
	for i := 0; i < len(values); i++ {
		if values[i] < l {
			l = values[i]
		}
	}
	return
}
