package helper

func Keys[K comparable, V any](m map[K]V) []K {
	res := make([]K, len(m))
	for k, _ := range m {
		res = append(res, k)
	}
	return res
}

func Vals[K comparable, V any](m map[K]V) []V {
	res := make([]V, len(m))
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

func KeysVals[K comparable, V any](m map[K]V) ([]K, []V) {
	resK := make([]K, len(m))
	resV := make([]V, len(m))
	for k, v := range m {
		resK = append(resK, k)
		resV = append(resV, v)
	}
	return resK, resV
}
