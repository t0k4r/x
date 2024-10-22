package chanx

import (
	"iter"
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
