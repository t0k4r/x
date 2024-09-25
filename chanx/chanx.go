package chanx

import (
	"iter"

	"github.com/t0k4r/x/iterx"
)

func All[T any](c chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range c {
			if !yield(item) {
				break
			}
		}
	}
}
func Some[T any](c chan T, filter func(T) bool) iter.Seq[T] {
	return iterx.Filter(All(c), filter)
}
func Transform[In, Out any](c chan In, mapf func(In) Out) iter.Seq[Out] {
	return iterx.Map(All(c), mapf)
}
func Uniq[T comparable](c chan T) iter.Seq[T] {
	return iterx.Uniq(All(c))
}
func Collect[T any](it iter.Seq[T]) chan T {
	c := make(chan T)
	go func() {
		for item := range it {
			c <- item
		}
		close(c)
	}()
	return c
}
