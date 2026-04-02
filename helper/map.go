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
