// go package bhsearch
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

/*
 * 'nos esse quasi nanos gigantum umeris insidentes' (Bernhard von Chartres, 1120)
 * The giants in this respect:
 * This is a modification of the Hash3 algorithm published by Lecroc, 2007
 * LECROQ, T. 2007. Fast exact string matching algorithms. Inf. Process. Lett. 102, 6, 229–235.
 * consisting of two modifications:
 *   a different hash function reducing hash by one shift operation
 *   and a different candidate comparison logic (bitwise computation comparison)
 */

package bh2search

import "runtime"

func findFI(haystack, pattern *[]byte) int {

	var (
		hay       = *haystack
		needle    = *pattern
		n         = len(hay) - 1
		m         = len(needle)
		mm1       = m - 1
		lim       = (m + (1 - mm1&1)) >> 1
		lchr      = needle[mm1]
		jmpMap    = make([]int, ALPHABET)
		i, j, jmp int
	)

	// preprocessing

	for ; i < ALPHABET; i++ {
		jmpMap[i] = mm1
	}

	// h = needle[0] + needle[1]<<2
	jmpMap[uint8(needle[0]+needle[1]<<2)] = m - 2

	for i = 2; i < m; i++ {
		// h = needle[i-1] + needle[i]<<2
		jmpMap[uint8(needle[i-1]+needle[i]<<2)] = mm1 - i
	}

	i = mm1

	for i < n+1 {
		j = 1
		for j != 0 {
			// h = hay[i-1] + hay[i]<<2
			j = jmpMap[uint8(hay[i-1]+hay[i]<<2)]
			i += j
			if i < n {
				continue
			}
			break
		}
		// drive forward
		jmp = i + 1
		// check candidate
		if i < n && 0 == ((hay[i]^needle[mm1])|(hay[i-mm1]^needle[0])) {
			// compare frontmost inner & lastmost inner
			for j = 1; j < lim; j++ {
				if 0 == ((hay[i-mm1+j] ^ needle[j]) | (hay[i-j] ^ needle[mm1-j])) {
					continue
				}
				if hay[i-mm1+j] != needle[j] {
					jmp = i + j
				}
				break
			}
			if j == lim {
				return i - mm1
				jmp += jmpMap[hay[jmp]]
			}
		}
		i = jmp
	}
	// eventually check last character
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
		lim       = (m + (1 - mm1&1)) >> 1
		lchr      = needle[mm1]
		jmpMap    = make([]int, ALPHABET)
		i, j, jmp int
	)

	if m < 2 {
		return found
	}

	// preprocessing

	buflen := 100 + (len(hay)/(1+len(needle)))>>8
	found = make([]int, 0, buflen)

	for ; i < ALPHABET; i++ {
		jmpMap[i] = mm1
	}

	// h = needle[0] + needle[1]<<2
	jmpMap[uint8(needle[0]+needle[1]<<2)] = m - 2

	for i = 2; i < m; i++ {
		// h = needle[i-1] + needle[i]<<2
		jmpMap[uint8(needle[i-1]+needle[i]<<2)] = mm1 - i
	}

	i = mm1

	for i < n+1 {
		j = 1
		for j != 0 {
			// h = hay[i-1] + hay[i]<<2
			j = jmpMap[uint8(hay[i-1]+hay[i]<<2)]
			i += j
			if i < n {
				continue
			}
			break
		}
		// drive forward
		jmp = i + 1
		// check candidate
		if i < n && 0 == ((hay[i]^needle[mm1])|(hay[i-mm1]^needle[0])) {
			// compare frontmost inner & lastmost inner
			for j = 1; j < lim; j++ {
				if 0 == ((hay[i-mm1+j] ^ needle[j]) | (hay[i-j] ^ needle[mm1-j])) {
					continue
				}
				if hay[i-mm1+j] != needle[j] {
					jmp = i + j
				}
				break
			}
			if j == lim {
				found = append(found, i-mm1)
				jmp += jmpMap[hay[jmp]]
			}
		}
		i = jmp
	}
	// eventually check last character
	if i == n && 0 == ((hay[i]^lchr)|(hay[i-mm1]^needle[0])) {
		for j = 1; j < mm1; j++ {
			if hay[i-mm1+j] != needle[j] {
				break
			}
		}
		if j == mm1 {
			found = append(found, i-mm1)
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
		lim       = (m + (1 - mm1&1)) >> 1
		lchr      = needle[mm1]
		jmpMap    = make([]int, ALPHABET)
		i, j, jmp int
	)

	if m < 2 {
		return count
	}

	// preprocessing

	for ; i < ALPHABET; i++ {
		jmpMap[i] = mm1
	}

	// h = needle[0] + needle[1]<<2
	jmpMap[uint8(needle[0]+needle[1]<<2)] = m - 2

	for i = 2; i < m; i++ {
		// h = needle[i-1] + needle[i]<<2
		jmpMap[uint8(needle[i-1]+needle[i]<<2)] = mm1 - i
	}

	i = mm1

	for i < n+1 {
		j = 1
		for j != 0 {
			// h = hay[i-1] + hay[i]<<2
			j = jmpMap[uint8(hay[i-1]+hay[i]<<2)]
			i += j
			if i < n {
				continue
			}
			break
		}
		jmp = i + 1
		if i < n && 0 == ((hay[i]^needle[mm1])|(hay[i-mm1]^needle[0])) {
			// compare frontmost inner & lastmost inner
			for j = 1; j < lim; j++ {
				if 0 == ((hay[i-mm1+j] ^ needle[j]) | (hay[i-j] ^ needle[mm1-j])) {
					continue
				}
				if hay[i-mm1+j] != needle[j] {
					jmp = i + j
				}
				break
			}
			if j == lim {
				jmp += jmpMap[hay[jmp]]
				count++
			}
		}
		i = jmp
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
		lim       = (m + (1 - mm1&1)) >> 1
		lchr      = needle[mm1]
		jmpMap    = make([]int, ALPHABET)
		i, j, jmp int
	)

	// preprocessing

	for ; i < ALPHABET; i++ {
		jmpMap[i] = mm1
	}

	// h = needle[0] + needle[1]<<2
	jmpMap[uint8(needle[0]+needle[1]<<2)] = m - 2

	for i = 2; i < m; i++ {
		// h = needle[i-1] + needle[i]<<2
		jmpMap[uint8(needle[i-1]+needle[i]<<2)] = mm1 - i
	}

	i = mm1

	for i < n {
		j = 1
		for j != 0 {
			// h = hay[i-1] + hay[i]<<2
			j = jmpMap[uint8(hay[i-1]+hay[i]<<2)]
			i += j
			if i < n {
				continue
			}
			break
		}
		// drive forward
		jmp = i + 1
		// check candidate
		if i < n && 0 == ((hay[i]^needle[mm1])|(hay[i-mm1]^needle[0])) {
			// compare frontmost inner & lastmost inner
			for j = 1; j < lim; j++ {
				if 0 == ((hay[i-mm1+j] ^ needle[j]) | (hay[i-j] ^ needle[mm1-j])) {
					continue
				}
				if hay[i-mm1+j] != needle[j] {
					jmp = i + j
				}
				break
			}
			if j == lim {
				return si + i - mm1
				jmp += jmpMap[hay[jmp]]
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
	}
	// eventually check last character
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
		lim       = (m + (1 - mm1&1)) >> 1
		lchr      = needle[mm1]
		jmpMap    = make([]int, ALPHABET)
		i, j, jmp int
	)

	if m < 2 {
		return &found
	}

	// preprocessing

	buflen := 100 + (len(hay)/(1+len(needle)))>>8
	found = make([]int, 0, buflen)

	for ; i < ALPHABET; i++ {
		jmpMap[i] = mm1
	}

	// h = needle[0] + needle[1]<<2
	jmpMap[uint8(needle[0]+needle[1]<<2)] = m - 2

	for i = 2; i < m; i++ {
		// h = needle[i-1] + needle[i]<<2
		jmpMap[uint8(needle[i-1]+needle[i]<<2)] = mm1 - i
	}

	i = mm1

	for i < n {
		j = 1
		for j != 0 {
			// h = hay[i-1] + hay[i]<<2
			j = jmpMap[uint8(hay[i-1]+hay[i]<<2)]
			i += j
			if i < n {
				continue
			}
			break
		}
		// drive forward
		jmp = i + 1
		// check candidate
		if i < n && 0 == ((hay[i]^needle[mm1])|(hay[i-mm1]^needle[0])) {
			// compare frontmost inner & lastmost inner
			for j = 1; j < lim; j++ {
				if 0 == ((hay[i-mm1+j] ^ needle[j]) | (hay[i-j] ^ needle[mm1-j])) {
					continue
				}
				if hay[i-mm1+j] != needle[j] {
					jmp = i + j
				}
				break
			}
			if j == lim {
				found = append(found, si+i-mm1)
				jmp += jmpMap[hay[jmp]]
			}
		}
		i = jmp
	}
	// eventually check last character
	if i == n && 0 == ((hay[i]^lchr)|(hay[i-mm1]^needle[0])) {
		for j = 1; j < mm1; j++ {
			if hay[i-mm1+j] != needle[j] {
				break
			}
		}
		if j == mm1 {
			found = append(found, si+i-mm1)
		}
	}

	return &found
}
