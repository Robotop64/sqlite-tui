package utils

func RemoveIdx[T any](arr []T, index int) []T {
	if index < 0 || index >= len(arr) {
		return arr // Return the original array if index is out of bounds
	}
	return append(arr[:index], arr[index+1:]...)
}
