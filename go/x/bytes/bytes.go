// Copyright © 2023 Axoflow
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
