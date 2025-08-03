package utils

func RemoveIdx[T any](arr []T, index int) []T {
	if index < 0 || index >= len(arr) {
		return arr // Return the original array if index is out of bounds
	}
	return append(arr[:index], arr[index+1:]...)
}

func RemoveItem[T comparable](arr []T, value T) []T {
	for i := 0; i < len(arr); i++ {
		if arr[i] == value {
			return RemoveIdx(arr, i)
		}
	}
	return arr
}

func IndexOf[T comparable](arr []T, value T) int {
	for i, v := range arr {
		if v == value {
			return i
		}
	}
	return -1 // Return -1 if the value is not found
}
