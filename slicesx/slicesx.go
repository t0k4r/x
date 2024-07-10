package slicesx

func Map[T1, T2 any](in []T1, fun func(T1) T2) []T2 {
	out := make([]T2, len(in))
	for i, v := range in {
		out[i] = fun(v)
	}

	return out
}

func Filter[T any](in []T, fun func(T) bool) []T {
	var out []T
	for _, v := range in {
		if fun(v) {
			out = append(out, v)
		}
	}
	return out
}

func MapErr[T1, T2 any](in []T1, fun func(T1) (T2, error)) ([]T2, error) {
	out := make([]T2, len(in))
	for i, v := range in {
		v, err := fun(v)
		if err != nil {
			return out, err
		}
		out[i] = v
	}
	return out, nil
}

func FilterErr[T any](in []T, fun func(T) (bool, error)) ([]T, error) {
	var out []T
	for _, v := range in {
		ok, err := fun(v)
		if err != nil {
			return out, err
		}
		if ok {
			out = append(out, v)
		}
	}

	return out, nil
}
