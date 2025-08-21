package sugar

func Default[T any](value T, err error) T {
	var defaultValue T
	if err != nil {
		return defaultValue
	}
	return value
}

func Value[T any](pointer *T) T {
	var defaultValue T
	if pointer == nil {
		return defaultValue
	}
	return *pointer
}

func And[T comparable](a T, b T) T {
	var defaultValue T
	if a == defaultValue {
		return a
	}
	return b
}

func Or[T comparable](a T, b T) T {
	var defaultValue T
	if a == defaultValue {
		return b
	}
	return a
}

func If[T any](condition bool, ifValue T, elseValue T) T {
	if condition {
		return ifValue
	} else {
		return elseValue
	}
}
