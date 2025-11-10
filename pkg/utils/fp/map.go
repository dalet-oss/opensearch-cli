package fp

import "strings"

func Map[T any, R any](ts []T, f func(t T) R) []R {
	result := make([]R, len(ts))
	for i, t := range ts {
		result[i] = f(t)
	}
	return result
}

func Filter[T any](ts []T, F func(t T) bool) []T {
	result := make([]T, 0)
	for _, t := range ts {
		if F(t) {
			result = append(result, t)
		}
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

func AsPointer[T any](t T) *T {
	return &t
}

func FromPointer[T any](t *T) T {
	return *t
}

func Ternary[T any](trueValue, falseValue T, condition bool) T {
	if condition {
		return trueValue
	}
	return falseValue
}
