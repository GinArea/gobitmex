package bitmexapi

func identity[T any](v T) (T, error) {
	return v, nil
}
