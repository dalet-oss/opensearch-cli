package fp

import "strings"

func Map[T any, R any](ts []T, f func(t T) R) []R {
	result := make([]R, len(ts))
	for i, t := range ts {
		result[i] = f(t)
	}
	return result
}

func GetOrDefault[T any](t T, defaultValue T, condition func(T) bool) T {
	if condition(t) {
		return t
	}
	return defaultValue
}

func NotEmptyString(s string) bool {
	return strings.TrimSpace(s) != ""
}
