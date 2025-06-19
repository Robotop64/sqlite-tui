package utils

func Map[In any, Out any](input []In, fn func(int, In) Out) []Out {
	result := make([]Out, len(input))
	for i, v := range input {
		result[i] = fn(i, v)
	}
	return result
}

func MapAt[T any](input []T, locs []bool, fn func(int, T) T) []T {
	if len(locs) == 0 {
		return Map(input, fn)
	}
	result := make([]T, len(input))

	for i, v := range input {
		if locs[i] {
			result[i] = fn(i, v)
		} else {
			result[i] = v
		}
	}

	return result
}

func Mask(locs []int, size int) []bool {
	result := make([]bool, size)
	for _, loc := range locs {
		if loc >= 0 && loc < size {
			result[loc] = true
		}
	}
	return result
}
