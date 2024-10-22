package iterx

import (
	"iter"
	"slices"
)

func Map[In, Out any](it iter.Seq[In], mapf func(In) Out) iter.Seq[Out] {
	return func(yield func(Out) bool) {
		for item := range it {
			if !yield(mapf(item)) {
				break
			}
		}
	}
}
func Map2[InK, InV, OutK, OutV any](it iter.Seq2[InK, InV], mapf func(InK, InV) (OutK, OutV)) iter.Seq2[OutK, OutV] {
	return func(yield func(OutK, OutV) bool) {
		for k, v := range it {
			if !yield(mapf(k, v)) {
				break
			}
		}
	}
}
func Filter[T any](it iter.Seq[T], filter func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range it {
			if !filter(item) {
				continue
			}
			if !yield(item) {
				break
			}
		}
	}
}
func Filter2[K, V any](it iter.Seq2[K, V], filter func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range it {
			if !filter(k, v) {
				continue
			}
			if !yield(k, v) {
				break
			}
		}
	}
}
func Uniq[T comparable](it iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var items []T
		for item := range it {
			if slices.Contains(items, item) {
				continue
			}
			items = append(items, item)
			if !yield(item) {
				break
			}
		}
	}
}
