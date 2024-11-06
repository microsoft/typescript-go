package sliceutil

import (
	"iter"
	"slices"
)

func AppendIfUnique[T comparable](array []T, element T) []T {
	if slices.Contains(array, element) {
		return array
	}
	return append(array, element)
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
	result, _ := SameFilter(slice, predicate)
	return result
}

func SameFilter[T any](slice []T, predicate func(T) bool) ([]T, bool) {
	for i, value := range slice {
		if !predicate(value) {
			result := slice[0:i]
			i++
			for i < len(slice) {
				value = slice[i]
				if predicate(value) {
					result = append(result, value)
				}
				i++
			}
			return result, false
		}
	}
	return slice, true
}

func Mapf[T, U any](slice []T, f func(T) U) []U {
	if len(slice) == 0 {
		return nil
	}
	result := make([]U, len(slice))
	for i := range slice {
		result[i] = f(slice[i])
	}
	return result
}

func MapIndex[T, U any](slice []T, f func(T, int) U) []U {
	if len(slice) == 0 {
		return nil
	}
	result := make([]U, len(slice))
	for i := range slice {
		result[i] = f(slice[i], i)
	}
	return result
}

func SameMap[T comparable](slice []T, f func(T) T) []T {
	for i, value := range slice {
		mapped := f(value)
		if mapped != value {
			result := make([]T, len(slice))
			copy(result, slice[:i])
			result[i] = mapped
			for j := i + 1; j < len(slice); j++ {
				result[j] = f(slice[j])
			}
			return result
		}
	}
	return slice
}

func SameMapIndex[T comparable](slice []T, f func(T, int) T) ([]T, bool) {
	for i, value := range slice {
		mapped := f(value, i)
		if mapped != value {
			result := make([]T, len(slice))
			copy(result, slice[:i])
			result[i] = mapped
			for i, value := range slice[i+1:] {
				result[i] = f(value, i)
			}
			return result, false
		}
	}
	return slice, true
}

func Some[T any](array []T, predicate func(T) bool) bool {
	for _, value := range array {
		if predicate(value) {
			return true
		}
	}
	return false
}

func Every[T any](array []T, predicate func(T) bool) bool {
	for _, value := range array {
		if !predicate(value) {
			return false
		}
	}
	return true
}

func InsertSorted[T any](slice []T, element T, cmp func(T, T) int) []T {
	i, _ := slices.BinarySearchFunc(slice, element, cmp)
	return slices.Insert(slice, i, element)
}

func FirstOrNil[T any](slice []T) T {
	if len(slice) != 0 {
		return slice[0]
	}
	return *new(T)
}

func FirstOrNilSeq[T any](seq iter.Seq[T]) T {
	if seq != nil {
		for value := range seq {
			return value
		}
	}
	return *new(T)
}

func LastOrNil[T any](slice []T) T {
	if len(slice) != 0 {
		return slice[len(slice)-1]
	}
	return *new(T)
}

func Find[T any](slice []T, predicate func(T) bool) T {
	for _, value := range slice {
		if predicate(value) {
			return value
		}
	}
	return *new(T)
}

func FindLast[T any](slice []T, predicate func(T) bool) T {
	for i := len(slice) - 1; i >= 0; i-- {
		value := slice[i]
		if predicate(value) {
			return value
		}
	}
	return *new(T)
}

func FindLastIndex[T any](slice []T, predicate func(T) bool) int {
	for i := len(slice) - 1; i >= 0; i-- {
		value := slice[i]
		if predicate(value) {
			return i
		}
	}
	return -1
}

func FindInMap[K comparable, V any](m map[K]V, predicate func(V) bool) V {
	for _, value := range m {
		if predicate(value) {
			return value
		}
	}
	return *new(V)
}

func Concatenate[T any](s1 []T, s2 []T) []T {
	if len(s2) == 0 {
		return s1
	}
	if len(s1) == 0 {
		return s2
	}
	return slices.Concat(s1, s2)
}

func CountWhere[T any](slice []T, predicate func(T) bool) int {
	count := 0
	for _, value := range slice {
		if predicate(value) {
			count++
		}
	}
	return count
}

func ReplaceElement[T any](slice []T, i int, t T) []T {
	result := slices.Clone(slice)
	result[i] = t
	return result
}

func Identical[T any](s1 []T, s2 []T) bool {
	if len(s1) == len(s2) {
		return len(s1) == 0 || &s1[0] == &s2[0]
	}
	return false
}
