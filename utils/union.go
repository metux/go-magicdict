package utils

func UnionSlice[K comparable](a []K, b []K) []K {
	m := make(map[K]bool)
	for _, v := range a {
		m[v] = true
	}
	for _, v := range b {
		m[v] = true
	}
	return MapKeys(m)
}
