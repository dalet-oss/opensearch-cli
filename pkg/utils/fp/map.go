package fp

import "strings"

// Map applies the function f to each element of the slice ts and returns the result as a slice of type R.
func Map[T any, R any](ts []T, f func(t T) R) []R {
	result := make([]R, len(ts))
	for i, t := range ts {
		result[i] = f(t)
	}
	return result
}

// Filter applies a provided filter function to a slice and returns a new slice containing only elements satisfying the condition.
func Filter[T any](ts []T, F func(t T) bool) []T {
	result := make([]T, 0)
	for _, t := range ts {
		if F(t) {
			result = append(result, t)
		}
	}
	return result
}

// GetOrDefault returns the value of `t` if the `condition` is true, otherwise it returns the `defaultValue`.
func GetOrDefault[T any](t T, defaultValue T, condition func(T) bool) T {
	if condition(t) {
		return t
	}
	return defaultValue
}

// NotEmptyString checks if the provided string is not empty or contains only whitespace. Returns true if non-empty, false otherwise.
func NotEmptyString(s string) bool {
	return strings.TrimSpace(s) != ""
}

// AsPointer converts a value of any type to a pointer to that value.
func AsPointer[T any](t T) *T {
	return &t
}

// Ternary returns the value of `trueValue` if the `condition` is true, otherwise it returns the `falseValue`.
func Ternary[T any](trueValue, falseValue T, condition bool) T {
	if condition {
		return trueValue
	}
	return falseValue
}
