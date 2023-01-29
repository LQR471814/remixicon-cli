package common

func Clamp(n, floor, ceil int) int {
	if n < floor {
		return floor
	}
	if n > ceil {
		return ceil
	}
	return n
}
