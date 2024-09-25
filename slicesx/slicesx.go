package slicesx

import (
	"iter"
	"slices"

	"github.com/t0k4r/x/iterx"
)

func All[T any](s []T) iter.Seq[T] {
	return iterx.DropK(slices.All(s))
}

func Some[T any](s []T, filter func(T) bool) iter.Seq[T] {
	return iterx.Filter(All(s), filter)
}

func Transform[In, Out any](s []In, mapf func(In) Out) iter.Seq[Out] {
	return iterx.Map(All(s), mapf)
}
