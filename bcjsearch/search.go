// go package bcjsearch
//
// The MIT License (MIT)
// Copyright (c) 2016 Andreas Briese, eduToolbox@Bri-C GmbH, Sarstedt

// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package bcjsearch

import "runtime"

// import (
// "bytes"
// "fmt"
// )

// helper to speed up the loops
func maximize(a, j int) int {
	if a > j {
		return a
	}
	return j
}

func findFI(haystack, pattern *[]byte) int {

	var (
		hay       = *haystack
		needle    = *pattern
		n         = len(hay) - 1
		m         = len(needle)
		mm1       = m - 1
		z         = 1 - mm1&1
		lim       = (m + z) >> 1
		lchr      = needle[mm1]
		jmpMap    = make([]int, ALPHABET)
		i, j, jmp int
	)

	// preprocessing

	for ; i < ALPHABET; i++ {
		jmpMap[i] = m
	}

	for i = 0; i < m; i++ {
		jmpMap[needle[i]] = mm1 - i
	}

	i = mm1

	for {
		// jump candidate from next char in y
		jmp = i + 1 + jmpMap[hay[i+1]]
		if 0 == ((hay[i] ^ lchr) | (hay[i-mm1] ^ needle[0])) {
			for j = 1; j < lim; j++ {
				if 0 == ((hay[i-mm1+j] ^ needle[j]) | (hay[i-j] ^ needle[mm1-j])) {
					continue
				}
			}
			if j == lim {
				return i - mm1
			}
		}
		for jmp < n {
			i = jmp
			jmp += jmpMap[hay[i]]
			if jmp > i {
				continue
			}
			break
		}
		i = jmp
		if i < n {
			continue
		}
		break
	}

	if i == n && 0 == ((hay[i]^lchr)|(hay[i-mm1]^needle[0])) {
		for j = 1; j < mm1; j++ {
			if hay[i-mm1+j] != needle[j] {
				break
			}
		}
		if j == mm1 {
			return i - mm1
		}
	}

	return -1
}

func findALL(haystack, pattern *[]byte) (found []int) {

	var (
		hay       = *haystack
		needle    = *pattern
		n         = len(hay) - 1
		m         = len(needle)
		mm1       = m - 1
		z         = 1 - mm1&1
		lim       = (m + z) >> 1
		lchr      = needle[mm1]
		jmpMap    = make([]int, ALPHABET)
		i, j, jmp int
	)

	// preprocessing

	for ; i < ALPHABET; i++ {
		jmpMap[i] = m
	}

	for i = 0; i < m; i++ {
		jmpMap[needle[i]] = mm1 - i
	}

	i = mm1

	for {
		// jump candidate from next char in y
		jmp = i + 1 + jmpMap[hay[i+1]]
		if 0 == ((hay[i] ^ lchr) | (hay[i-mm1] ^ needle[0])) {
			for j = 1; j < lim; j++ {
				if 0 == ((hay[i-mm1+j] ^ needle[j]) | (hay[i-j] ^ needle[mm1-j])) {
					continue
				}
			}
			if j == lim {
				found = append(found, i)
			}
		}
		for jmp < n {
			i = jmp
			jmp += jmpMap[hay[i]]
			if jmp > i {
				continue
			}
			break
		}
		i = jmp
		if i < n {
			continue
		}
		break
	}

	if i == n && 0 == ((hay[i]^lchr)|(hay[i-mm1]^needle[0])) {
		for j = 1; j < mm1; j++ {
			if hay[i-mm1+j] != needle[j] {
				break
			}
		}
		if j == mm1 {
			found = append(found, i)
		}
	}

	return found
}

func count(haystack, pattern *[]byte) (count int) {

	var (
		hay       = *haystack
		needle    = *pattern
		n         = len(hay) - 1
		m         = len(needle)
		mm1       = m - 1
		z         = 1 - mm1&1
		lim       = (m + z) >> 1
		lchr      = needle[mm1]
		jmpMap    = make([]int, ALPHABET)
		i, j, jmp int
	)

	// preprocessing

	for ; i < ALPHABET; i++ {
		jmpMap[i] = m
	}

	for i = 0; i < m; i++ {
		jmpMap[needle[i]] = mm1 - i
	}

	i = mm1

	for {
		// jump candidate from next char in y
		jmp = i + 1 + jmpMap[hay[i+1]]
		if 0 == ((hay[i] ^ lchr) | (hay[i-mm1] ^ needle[0])) {
			for j = 1; j < lim; j++ {
				if 0 == ((hay[i-mm1+j] ^ needle[j]) | (hay[i-j] ^ needle[mm1-j])) {
					continue
				}
			}
			if j == lim {
				count++
			}
		}
		for jmp < n {
			i = jmp
			jmp += jmpMap[hay[i]]
			if jmp > i {
				continue
			}
			break
		}
		i = jmp
		if i < n {
			continue
		}
		break
	}

	if i == n && 0 == ((hay[i]^lchr)|(hay[i-mm1]^needle[0])) {
		for j = 1; j < mm1; j++ {
			if hay[i-mm1+j] != needle[j] {
				break
			}
		}
		if j == mm1 {
			count++
		}
	}

	return count
}

func findFI_CC(haystack, pattern *[]byte,
	startIdx, partLen int,
	breaker chan int,
	threads chan bool) int {

	runtime.LockOSThread()
	defer func() {
		threads <- true
		runtime.UnlockOSThread()
	}()

	si := startIdx - (len(*pattern) - 1)
	if si < 0 {
		si = 0
	}

	var (
		hay       = (*haystack)[si : startIdx+partLen]
		needle    = *pattern
		n         = len(hay) - 1
		m         = len(needle)
		mm1       = m - 1
		z         = 1 - mm1&1
		lim       = (m + z) >> 1
		lchr      = needle[mm1]
		jmpMap    = make([]int, ALPHABET)
		i, j, jmp int
	)

	// preprocessing

	for ; i < ALPHABET; i++ {
		jmpMap[i] = m
	}

	for i = 0; i < m; i++ {
		jmpMap[needle[i]] = mm1 - i
	}

	i = mm1

	for {
		// jump candidate from next char in y
		jmp = i + 1 + jmpMap[hay[i+1]]
		if 0 == ((hay[i] ^ lchr) | (hay[i-mm1] ^ needle[0])) {
			for j = 1; j < lim; j++ {
				if 0 == ((hay[i-mm1+j] ^ needle[j]) | (hay[i-j] ^ needle[mm1-j])) {
					continue
				}
				if hay[i-j] != needle[mm1-j] {
					//jmp = i-mm1+j //i-j+jmpMap[hay[i-j]], jmp)
					for jmp < n {
						i = jmp
						jmp += jmpMap[hay[i]]
						if jmp > i {
							continue
						}
						break
					}
					break
				}
				if hay[i-mm1+j] != needle[j] {
					jmp = i + j //  maximize(i+j, jmp)
					for jmp < n {
						i = jmp
						jmp += jmpMap[hay[i]]
						if jmp > i {
							continue
						}
						break
					}
					break
				}
			}
			if j == lim {
				return si + i - mm1
			}
		}
		i = jmp
		select {
		case broken := <-breaker:
			if broken < si {
				return len(*haystack)
			}
		default:
		}
		if i < n {
			continue
		}
		break
	}

	if i == n && 0 == ((hay[i]^lchr)|(hay[i-mm1]^needle[0])) {
		for j = 1; j < mm1; j++ {
			if hay[i-mm1+j] != needle[j] {
				break
			}
		}
		if j == mm1 {
			return si + i - mm1
		}
	}

	return -1
}

func findALL_CC(haystack, pattern *[]byte,
	startIdx, partLen, bufLen int,
	threads chan bool) *[]int {

	runtime.LockOSThread()
	defer func() {
		threads <- true
		runtime.UnlockOSThread()
	}()

	si := startIdx - (len(*pattern) - 1)
	if si < 0 {
		si = 0
	}

	var (
		found     = make([]int, 0, bufLen)
		hay       = (*haystack)[si : startIdx+partLen]
		needle    = *pattern
		n         = len(hay) - 1
		m         = len(needle)
		mm1       = m - 1
		z         = 1 - mm1&1
		lim       = (m + z) >> 1
		lchr      = needle[mm1]
		jmpMap    = make([]int, ALPHABET)
		i, j, jmp int
	)

	// preprocessing

	for ; i < ALPHABET; i++ {
		jmpMap[i] = m
	}

	for i = 0; i < m; i++ {
		jmpMap[needle[i]] = mm1 - i
	}

	i = mm1

	switch m {
	case 2:
		for {
			jmp = i + 1 + jmpMap[hay[i+1]]
			if 0 == ((hay[i] ^ lchr) | (hay[i-mm1] ^ needle[0])) {
				found = append(found, si+i)
			}
			i = jmp
			if i < n {
				continue
			}
			break
		}
	case 3:
		for {
			jmp = i + 1 + jmpMap[hay[i+1]]
			if 0 == ((hay[i] ^ lchr) | (hay[i-mm1] ^ needle[0]) | (hay[i-1] ^ needle[1])) {
				found = append(found, si+i)
			}
			i = jmp
			if i < n {
				continue
			}
			break
		}
	default:
		for {
			// jump candidate from next char in y
			jmp = i + 1 + jmpMap[hay[i+1]]
			if 0 == ((hay[i] ^ lchr) | (hay[i-mm1] ^ needle[0])) {
				for j = 1; j < lim; j++ {
					if 0 == ((hay[i-mm1+j] ^ needle[j]) | (hay[i-j] ^ needle[mm1-j])) {
						continue
					}
					break
				}
				if j == lim {
					found = append(found, si+i)
				}
			}
			i = jmp
			for jmp < n {
				i = jmp
				jmp += jmpMap[hay[i]]
				if jmp > i {
					continue
				}
				break
			}
			if i < n {
				continue
			}
			break
		}
	}

	if i == n && hay[i] == lchr {
		for j = 0; j < mm1; j++ {
			if hay[i-mm1+j] != needle[j] {
				break
			}
		}
		if j == mm1 {
			found = append(found, si+i)
		}
	}

	return &found
}
