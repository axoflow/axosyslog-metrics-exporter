package bytes

import "golang.org/x/exp/constraints"

func CommonPrefixLen(a, b []byte) int {
	l := min(len(a), len(b))
	a, b = a[:l], b[:l]
	for i := range a {
		if a[i] != b[i] {
			return i
		}
	}
	return l
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
