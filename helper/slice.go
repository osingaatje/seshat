package helper

func FindValue[K comparable, V any](mp map[K]V, f func(V) bool) (*K, *V, bool) {
	for k, v := range mp {
		if f(v) {
			return &k, &v, true
		}
	}
	return nil, nil, false
}

func Find[E any](arr []E, f func(*E) bool) (*E, bool) {
	for _, e := range arr {
		if f(&e) {
			return &e, true
		}
	}
	return nil, false
}

func Any[E any](arr []E, f func(*E) bool) bool {
	for _, e := range arr {
		if f(&e) {
			return true
		}
	}
	return false
}

func AnyBool(arr []bool) bool {
	for _, e := range arr {
		if e {
			return true
		}
	}
	return false
}
