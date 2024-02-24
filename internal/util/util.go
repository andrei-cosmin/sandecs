package util

func NextPowerOfTwo(value uint) uint {
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	value++

	return value
}
