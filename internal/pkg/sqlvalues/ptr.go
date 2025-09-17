package sqlvalues

func PtrToValue[T any](ptr *T) T {
	var empty T

	if ptr == nil {
		return empty
	}

	return *ptr
}
