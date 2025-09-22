package fp

func Map[T any, R any](ts []T, f func(t T) R) []R {
	result := make([]R, len(ts))
	for i, t := range ts {
		result[i] = f(t)
	}
	return result
}
