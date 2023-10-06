package utils

func ListHead[K any](s []K) (K, []K) {
	if len(s) == 0 {
		return *(new(K)), []K{}
	}
	if len(s) == 1 {
		return s[0], []K{}
	}
	return s[0], s[1:]
}
