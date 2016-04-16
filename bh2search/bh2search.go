// go package bh2search
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
 *   and a different candidate comparison logic (zigzag)
 */

package bh2search

import (
	"errors"
)

// Alphabet & Errors
var (
	ALPHABET    = 256
	NEEDLESHORT = errors.New("Length needle is < 3")
	NEEDLELONG  = errors.New("Length needle > length haystack")
)

func Index(haystack, needle *[]byte) (int, error) {

	// check length needle
	if len(*haystack) < len(*needle) {
		return -1, NEEDLELONG
	}
	if len(*needle) < 3 {
		return -1, NEEDLESHORT
	}

	return findFI(haystack, needle), nil
}

func Count(haystack, needle *[]byte) (int, error) {

	// check length needle
	if len(*haystack) < len(*needle) {
		return -1, NEEDLELONG
	}
	if len(*needle) < 2 {
		return -1, NEEDLESHORT
	}

	return count(haystack, needle), nil
}

func FindAll(haystack, needle *[]byte) (found []int, e error) {

	// check length needle
	if len(*haystack) < len(*needle) {
		return found, NEEDLELONG
	}
	if len(*needle) < 2 {
		return found, NEEDLESHORT
	}

	return findALL(haystack, needle), nil
}

func FindAllCC(haystack, needle *[]byte, startIdx, partLen, bufLen int, threads chan bool) *[]int {

	found := []int{}

	// check length needle
	if len(*haystack) < len(*needle) {
		return &found
	}
	if len(*needle) < 2 {
		return &found
	}

	return findALL_CC(haystack, needle, startIdx, partLen, bufLen, threads)
}

func FindIndexCC(haystack, needle *[]byte, startIdx, partLen int, breaker chan int, threads chan bool) int {

	// check length needle
	if len(*haystack) < len(*needle) {
		return -1
	}
	if len(*needle) < 2 {
		return -1
	}

	return findFI_CC(haystack, needle, startIdx, partLen, breaker, threads)
}
