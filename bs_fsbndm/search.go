// go package bs_fsbndm
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
 * This is a modification of the Forward Semplified BNDM algorithm published by
 * S. Faro and T. Lecroq (2008):
 * Efficient Variants of the Backward-Oracle-Matching Algorithm.
 * Proceedings of the Prague Stringology Conference 2008, pp.146--160, Czech Technical University in Prague, Czech Republic, (2008).
 * Lizence of the authors C-implementation: GNU General Public License V.3 as published by the Free Software Foundation
 */

package bs_fsbndm

import (
	"bytes"
	// "fmt"
	"runtime"
)

func findFI(haystack, pattern *[]byte) int {

	var (
		hay                              = *haystack
		needle                           = *pattern
		n                                = len(hay)
		m                                = len(needle)
		p                                = m // len Pat
		longPat                          = m > 63
		bitPat                           = make([]uint64, ALPHABET)
		bits                             uint64
		i, lastCharIdx, suffIdx, backstp int
	)

	// preprocessing

	if longPat {
		p = 62
	}

	for i = 0; i < ALPHABET; i++ {
		bitPat[i] = 1
	}

	// here comes the magic!!!
	// for each character create a word bloomfilter holding its position(s)!
	// for logPat use the last p bytes of needle for pat check & shift
	suffIdx = m - p
	for i = 0; i < p; i++ {
		bitPat[needle[suffIdx+i]] |= (1 << uint(p-i))
	}

	// search
	if bytes.Equal(hay, needle) {
		return 1
	}

	if longPat { // search for the suffix length p of m
		for i = m; i < n-1; {
			// check character pair at windows right edge
			bits = (bitPat[hay[i+1]] << 1) & bitPat[hay[i]]
			switch bits { // at least hay[i] at window edge chars is found xxx10
			default:
				// run backwards over the candidate; check & shift
				lastCharIdx = i
				backstp = 1
				bits = (bits << 1) & bitPat[hay[i-backstp]]
				for bits != 0 {
					backstp++
					bits = (bits << 1) & bitPat[hay[i-backstp]]
				}
				i += p - backstp
				if i == lastCharIdx {
					if bytes.Equal(needle, hay[lastCharIdx-m+1:lastCharIdx+1]) {
						return lastCharIdx - m + 1
					}
					i++
				}
			case 0: // bits didn't match with any bytes in pat -> shift by p
				i += p
			}
		}
		return -1
	}

	for i = m; i < n-1; {
		// check character pair at windows right edge
		bits = (bitPat[hay[i+1]] << 1) & bitPat[hay[i]]
		switch bits {
		default:
			// run backwards over the candidate; check & shift
			lastCharIdx = i
			backstp = 1
			bits = (bits << 1) & bitPat[hay[i-backstp]]
			for bits != 0 {
				backstp++
				bits = (bits << 1) & bitPat[hay[i-backstp]]
			}
			i += m - backstp
			if backstp == m {
				return lastCharIdx - m + 1
			}
		case 0: // bits didn't match with any bytes in needle -> shift by m
			i += m
		}
	}
	return -1
}

func findALL(haystack, pattern *[]byte) (found []int) {

	var (
		hay                              = *haystack
		needle                           = *pattern
		n                                = len(hay)
		m                                = len(needle)
		p                                = m // len Pat
		longPat                          = m > 63
		bitPat                           = make([]uint64, ALPHABET)
		bits                             uint64
		i, lastCharIdx, suffIdx, backstp int
	)

	buflen := 100 + (len(hay)/(1+len(needle)))>>8
	found = make([]int, 0, buflen)

	// preprocessing

	if longPat {
		p = 62
	}

	for i = 0; i < ALPHABET; i++ {
		bitPat[i] = 1
	}

	// here comes the magic!!!
	// for each character create a word bloomfilter holding its position(s)!
	// for logPat use the last p bytes of needle for pat check & shift
	suffIdx = m - p
	for i = 0; i < p; i++ {
		bitPat[needle[suffIdx+i]] |= (1 << uint(p-i))
	}

	// search
	if bytes.Equal(hay, needle) {
		found = append(found, 0)
	}

	if bytes.Equal(hay[n-m:], needle) {
		found = append(found, n-m)
	}

	if longPat { // search for the suffix length p of m
		for i = m; i < n-1; {
			// check character pair at windows right edge
			bits = (bitPat[hay[i+1]] << 1) & bitPat[hay[i]]
			switch bits { // at least hay[i] at window edge chars is found xxx10
			default:
				// run backwards over the candidate; check & shift
				lastCharIdx = i
				backstp = 1
				bits = (bits << 1) & bitPat[hay[i-backstp]]
				for bits != 0 {
					backstp++
					bits = (bits << 1) & bitPat[hay[i-backstp]]
				}
				i += p - backstp
				if i == lastCharIdx {
					if bytes.Equal(needle, hay[lastCharIdx-m+1:lastCharIdx+1]) {
						found = append(found, lastCharIdx-m+1)
						i = lastCharIdx + 1
					}
					i++
				}
			case 0: // bits didn't match with any bytes in pat -> shift by p
				i += p
			}
		}

		return found
	}

	for i = m; i < n-1; {
		// check character pair at windows right edge
		bits = (bitPat[hay[i+1]] << 1) & bitPat[hay[i]]
		switch bits {
		default:
			// run backwards over the candidate; check & shift
			lastCharIdx = i
			backstp = 1
			bits = (bits << 1) & bitPat[hay[i-backstp]]
			for bits != 0 {
				backstp++
				bits = (bits << 1) & bitPat[hay[i-backstp]]
			}
			i += m - backstp
			if backstp == m {
				found = append(found, lastCharIdx-m+1)
				i = lastCharIdx + 1
			}
		case 0: // bits didn't match with any bytes in needle -> shift by m
			i += m
		}
	}

	return found

}

func count(haystack, pattern *[]byte) (count int) {

	var (
		hay                              = *haystack
		needle                           = *pattern
		n                                = len(hay)
		m                                = len(needle)
		p                                = m // len Pat
		longPat                          = m > 63
		bitPat                           = make([]uint64, ALPHABET)
		bits                             uint64
		i, lastCharIdx, suffIdx, backstp int
	)

	// preprocessing

	if longPat {
		p = 62
	}

	for i = 0; i < ALPHABET; i++ {
		bitPat[i] = 1
	}

	// here comes the magic!!!
	// for each character create a word bloomfilter holding its position(s)!
	// for logPat use the last p bytes of needle for pat check & shift
	suffIdx = m - p
	for i = 0; i < p; i++ {
		bitPat[needle[suffIdx+i]] |= (1 << uint(p-i))
	}

	// search
	if bytes.Equal(hay, needle) {
		count++
	}

	if bytes.Equal(hay[n-m:], needle) {
		count++
	}

	if longPat { // search for the suffix length p of m
		for i = m; i < n-1; {
			// check character pair at windows right edge
			bits = (bitPat[hay[i+1]] << 1) & bitPat[hay[i]]
			switch bits { // at least hay[i] at window edge chars is found xxx10
			default:
				// run backwards over the candidate; check & shift
				lastCharIdx = i
				backstp = 1
				bits = (bits << 1) & bitPat[hay[i-backstp]]
				for bits != 0 {
					backstp++
					bits = (bits << 1) & bitPat[hay[i-backstp]]
				}
				i += p - backstp
				if i == lastCharIdx {
					if bytes.Equal(needle, hay[lastCharIdx-m+1:lastCharIdx+1]) {
						count++
						i = lastCharIdx + 1
					}
					i++
				}
			case 0: // bits didn't match with any bytes in pat -> shift by p
				i += p
			}
		}

		return count
	}

	for i = m; i < n-1; {
		// check character pair at windows right edge
		bits = (bitPat[hay[i+1]] << 1) & bitPat[hay[i]]
		switch bits {
		default:
			// run backwards over the candidate; check & shift
			lastCharIdx = i
			backstp = 1
			bits = (bits << 1) & bitPat[hay[i-backstp]]
			for bits != 0 {
				backstp++
				bits = (bits << 1) & bitPat[hay[i-backstp]]
			}
			i += m - backstp
			if backstp == m {
				count++
				i = lastCharIdx + 1
			}
		case 0: // bits didn't match with any bytes in needle -> shift by m
			i += m
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
		hay                              = (*haystack)[si : startIdx+partLen]
		needle                           = *pattern
		n                                = len(hay)
		m                                = len(needle)
		p                                = m // len Pat
		longPat                          = m > 63
		bitPat                           = make([]uint64, ALPHABET)
		bits                             uint64
		i, lastCharIdx, suffIdx, backstp int
	)

	// preprocessing

	if longPat {
		p = 62
	}

	for i = 0; i < ALPHABET; i++ {
		bitPat[i] = 1
	}

	// here comes the magic!!!
	// for each character create a word bloomfilter holding its position(s)!
	// for logPat use the last p bytes of needle for pat check & shift
	suffIdx = m - p
	for i = 0; i < p; i++ {
		bitPat[needle[suffIdx+i]] |= (1 << uint(p-i))
	}

	// search
	if bytes.Equal(hay, needle) {
		return si
	}

	if longPat { // search for the suffix length p of m
		for i = m; i < n-1; {
			// check character pair at windows right edge
			bits = (bitPat[hay[i+1]] << 1) & bitPat[hay[i]]
			switch bits { // at least hay[i] at window edge chars is found xxx10
			default:
				// run backwards over the candidate; check & shift
				lastCharIdx = i
				backstp = 1
				bits = (bits << 1) & bitPat[hay[i-backstp]]
				for bits != 0 {
					backstp++
					bits = (bits << 1) & bitPat[hay[i-backstp]]
				}
				i += p - backstp
				if i == lastCharIdx {
					if bytes.Equal(needle, hay[lastCharIdx-m+1:lastCharIdx+1]) {
						return si + lastCharIdx - m + 1
					}
					i++
				}
			case 0: // bits didn't match with any bytes in pat -> shift by p
				i += p
			}
			select {
			case broken := <-breaker:
				if broken < si {
					return len(*haystack)
				}
			default:
			}
		}
		return -1
	}

	for i = m; i < n-1; {
		// check character pair at windows right edge
		bits = (bitPat[hay[i+1]] << 1) & bitPat[hay[i]]
		switch bits {
		default:
			// run backwards over the candidate; check & shift
			lastCharIdx = i
			backstp = 1
			bits = (bits << 1) & bitPat[hay[i-backstp]]
			for bits != 0 {
				backstp++
				bits = (bits << 1) & bitPat[hay[i-backstp]]
			}
			i += m - backstp
			if backstp == m {
				return si + lastCharIdx - m + 1
			}
		case 0: // bits didn't match with any bytes in needle -> shift by m
			i += m
		}
		select {
		case broken := <-breaker:
			if broken < si {
				return len(*haystack)
			}
		default:
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
		found                            = make([]int, 0, bufLen)
		hay                              = (*haystack)[si : startIdx+partLen]
		needle                           = *pattern
		n                                = len(hay)
		m                                = len(needle)
		p                                = m // len Pat
		longPat                          = m > 63
		bitPat                           = make([]uint64, ALPHABET)
		bits                             uint64
		i, lastCharIdx, suffIdx, backstp int
	)

	// fmt.Println("::", si, n)

	// preprocessing

	if longPat {
		p = 62
	}

	for i = 0; i < ALPHABET; i++ {
		bitPat[i] = 1
	}

	// here comes the magic!!!
	// for each character create a word bloomfilter holding its position(s)!
	// for logPat use the last p bytes of needle for pat check & shift
	suffIdx = m - p
	for i = 0; i < p; i++ {
		bitPat[needle[suffIdx+i]] |= (1 << uint(p-i))
	}

	// search
	if bytes.Equal(hay, needle) {
		found = append(found, 0)
	}

	if longPat { // search for the suffix length p of m
		for i = m; i < n-1; {
			// check character pair at windows right edge
			bits = (bitPat[hay[i+1]] << 1) & bitPat[hay[i]]
			switch bits { // at least hay[i] at window edge chars is found xxx10
			default:
				// run backwards over the candidate; check & shift
				lastCharIdx = i
				backstp = 1
				bits = (bits << 1) & bitPat[hay[i-backstp]]
				for bits != 0 {
					backstp++
					bits = (bits << 1) & bitPat[hay[i-backstp]]
				}
				i += p - backstp
				if i == lastCharIdx {
					if bytes.Equal(needle, hay[lastCharIdx-m+1:lastCharIdx+1]) {
						found = append(found, si+lastCharIdx-m+1)
						i = lastCharIdx + 1
					}
					i++
				}
			case 0: // bits didn't match with any bytes in pat -> shift by p
				i += p
			}
		}

		return &found
	}

	for i = m; i < n-1; {
		// check character pair at windows right edge
		bits = (bitPat[hay[i+1]] << 1) & bitPat[hay[i]]
		switch bits {
		default:
			// run backwards over the candidate; check & shift
			lastCharIdx = i
			backstp = 1
			bits = (bits << 1) & bitPat[hay[i-backstp]]
			for bits != 0 {
				backstp++
				bits = (bits << 1) & bitPat[hay[i-backstp]]
			}
			i += m - backstp
			if backstp == m {
				found = append(found, si+lastCharIdx-m+1)
				i = lastCharIdx + 1
			}
		case 0: // bits didn't match with any bytes in needle -> shift by m
			i += m
		}
	}

	return &found

}
