package utils

// Clamp constrains v to the inclusive [minVal, maxVal] range.
//
// Parameters:
//   - v: value to clamp.
//   - minVal: lower bound.
//   - maxVal: upper bound.
//
// Returns:
//   - clamped value within the provided bounds.
func Clamp(v, minVal, maxVal float64) float64 {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}

// Max returns the larger of a and b.
//
// Parameters:
//   - a: first integer.
//   - b: second integer.
//
// Returns:
//   - the greater of a or b.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
