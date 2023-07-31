// Copyright © 2023 Axoflow
// All rights reserved.

package bytes

import (
	"golang.org/x/exp/constraints"
)

// CommonPrefixLen returns the length of its two parameters' common prefix
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
