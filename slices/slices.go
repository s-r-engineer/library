package librarySlices

func InsertAt[T any](slice []T, value T, index int) []T {
	if index < 0 || index > len(slice) {
		panic("index out of range")
	}
	slice = append(slice, value)
	copy(slice[index+1:], slice[index:])
	slice[index] = value
	return slice
}
