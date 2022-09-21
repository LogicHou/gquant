package utils

import (
	"math"
	"strconv"
	"time"
)

func MsToTime(ms int64) string {
	tm := time.Unix(0, ms*int64(time.Millisecond))
	return tm.Format("2006-01-02 15:04:05.000")
}

func MsToDateTime(ms int64) string {
	tm := time.Unix(0, ms*int64(time.Millisecond))
	return tm.Format("2006-01-02 15:04:05")
}

func F64ToStr(f float64) string {
	s := strconv.FormatFloat(f, 'f', 3, 64)
	return s
}

func StrToF64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func StrToI64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func SliceFind(slice []int64, val int64) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func FRound(x float64) float64 {
	return math.Round(x*1000) / 1000
}

func FRound2(x float64) float64 {
	return math.Round(x*100) / 100
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	}
	return ""
}
