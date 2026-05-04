package helper

func AnyMap[K comparable](m map[K]bool) bool {
	for _, v := range m {
		if v {
			return true
		}
	}
	return false
}
func AllMap[K comparable](m map[K]bool) bool {
	for _, v := range m {
		if !v {
			return false
		}
	}
	return true
}

func ValuesMap[K comparable, V any](m map[K]V) []V {
	res := []V{}
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

func FilterMap[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	res := map[K]V{}
	for k, v := range m {
		if f(k, v) {
			res[k] = v
		}
	}
	return res
}
