package helper

func Find[E any](arr []E, f func(*E) bool) (*E, bool) {
	for _, e := range arr {
		if f(&e) {
			return &e, true
		}
	}
	return nil, false
}
