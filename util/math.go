package util

// Clamp restricts value between min and max
func Clamp(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}

// Map maps an input value from interval [inFrom, inTo] to [outFrom, outTo].
// Input value is clamped to input interval.
func Map(value, inFrom, inTo, outFrom, outTo float64) float64 {
	// restrict value to input interval
	value = Clamp(value, inFrom, inTo)
	return outTo*(value-inFrom)/(inTo-inFrom) + outFrom
}
