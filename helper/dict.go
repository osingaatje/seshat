package helper

func Keys[K comparable, V any](m map[K]V) []K {
	set := map[K]bool{}
	for k, _ := range m {
		set[k] = true
	}
	res := make([]K, len(m))
	i := 0
	for k, _ := range set {
		res[i] = k
		i++
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
