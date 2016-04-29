// package bmatch
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
 * func memchr(haystack, patttern)
 * This native Go 64-bit algorithm derives
 * from insights into the 32bit implementation in c at
 * http://www.stdlib.net/~colmmacc/strlen.c.html
 * Copyright (C) 1991, 1993, 1997, 2000, 2003 Free Software Foundation, Inc.
 * part of the GNU C Library.
 *    Written by Torbjorn Granlund (tege@sics.se),
 *    with help from Dan Sahlin (dan@sics.se);
 *    commentary by Jim Blandy (jimb@ai.mit.edu).
 */

// .. and another issue: if package "unsafe" is removed from Go in future
// these functions for 1 byte long patterns must become replaced by the (slower)
// bytes.Index() function.

package bmatch

import (
	"bytes"
	// "runtime"
	// "unsafe"
)

/*
 * func mmIndex(haystack, needle *[]byte)
 * returns first index of neddle[0] in haystack
 * This native Go 64-bit algorithm derives
 * from insides into the 32bit implementation in c at
 * http://www.stdlib.net/~colmmacc/strlen.c.html
 */
func mmIndex(haystack, needle *[]byte) int {

	var (
		pat         = *needle
		lastCharIdx = len(pat) - 1
		hayst       = *haystack
		n           = len(hayst)
	)

	if n < len(pat) {
		return -1
	}

	return bytes.Index(hayst, pat)

}

/*
 * func mmFindALL(haystack, needle *[]byte)
 * returns []int containing all indices of neddle[0] in haystack
 * This native Go 64-bit algorithm derives
 * from insides into the 32bit implementation in c at
 * http://www.stdlib.net/~colmmacc/strlen.c.html
 */
func mmFindALL(haystack, needle *[]byte) (found []int) {

	var (
		pat         = *needle
		lastCharIdx = len(pat) - 1
		hayst       = *haystack
		n           = len(hayst)
	)

	if n < len(pat) {
		return found
	}

	buflen := 10 + (n/(1+lastCharIdx))>>3
	found = make([]int, 0, buflen)

	var idx, lastIdx int

	for {
		idx = bytes.Index(hayst, pat)
		if idx == -1 {
			break
		}
		found = append(found, lastIdx+idx)
		lastIdx += idx + 1
		hayst = hayst[idx+1:]
	}

	return found
}

/*
 * func mmCount(haystack, needle *[]byte)
 * returns the total number of neddle[0] found in haystack
 * This native Go 64-bit algorithm derives
 * from insides into the 32bit implementation in c at
 * http://www.stdlib.net/~colmmacc/strlen.c.html
 */
func mmCount(haystack, needle *[]byte) (count int) {

	var (
		pat         = *needle
		lastCharIdx = len(pat) - 1
		hayst       = *haystack
		n           = len(hayst)
	)

	if n < len(pat) {
		return count
	}

	idx := 0
	for {
		idx = bytes.Index(hayst, pat)
		if idx == -1 {
			break
		}
		count++
		hayst = hayst[idx+1:]
	}

	return count
}
