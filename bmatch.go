// go package bmatch
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
 * bmatch is an exact string matching package consisting of a selection of fast string searching
 * algos from the 21th century modified, adjusted & implemented in Go by the author.
 * For pattern strings length > 3 it is about twice as fast as Go/Golangs Boyer-Moore
 * implementation in search.go (part of the standard string library) and about
 * 8-10 times faster than Go/Golangs standard string.Index function based on Rabin-Karp-Algorithm
 * with addition-hash.
 * Go's search for one byte in haystack (bytes.Index, referring to an asm routine in sys) is excelled
 * by about 10-20% by the unsafeMEMCHR single byte search functions on 64bit words.
 *
 * In case you want to search over small alphabets (<30) you might adjust the switch points between the algorithm.
 * Example: The commented switch-points in the Count() function below were found optimal for 23 amino-acid
 * protein sequences.
 */

package bmatch

import (
	"errors"

	// bcj "github.com/AndreasBriese/bmatch/bcjsearch"
	bh2 "github.com/AndreasBriese/bmatch/bh2search"
	bh "github.com/AndreasBriese/bmatch/bhsearch"
	bsf "github.com/AndreasBriese/bmatch/bs_fsbndm"
)

var ALPHABET = 256

func Index(haystack, needle *[]byte) (found int, e error) {

	if len(*needle) < 1 {
		return -1, errors.New("length of needle is smaller 1")
	}

	var matchFn func(haystack, needle *[]byte) (int, error)
	switch {
	case len(*needle) < 2:
		return mmIndex(haystack, needle), nil
	case len(*needle) < 50:
		matchFn = bsf.Index
	case len(*needle) < 12000:
		matchFn = bh.Index
	case len(*needle) < 350000:
		matchFn = bh2.Index
	case len(*needle) < 1<<22:
		matchFn = bsf.Index
	default:
		matchFn = bsf.Index
	}

	return matchFn(haystack, needle)
}

func FindAll(haystack, needle *[]byte) (found []int, e error) {

	if len(*needle) < 1 {
		return found, errors.New("length of needle is smaller 1")
	}

	var matchFn func(haystack, needle *[]byte) ([]int, error)
	switch {
	case len(*needle) < 2:
		return mmFindALL(haystack, needle), nil
	case len(*needle) < 50:
		matchFn = bsf.FindAll
	case len(*needle) < 12000:
		matchFn = bh.FindAll
	case len(*needle) < 350000:
		matchFn = bh2.FindAll
	case len(*needle) < 1<<22:
		matchFn = bsf.FindAll
	default:
		matchFn = bsf.FindAll
	}

	return matchFn(haystack, needle)

}

func Count(haystack, needle *[]byte) (found int, e error) {

	if len(*needle) < 1 {
		return -1, errors.New("length of needle is smaller 1")
	}

	var matchFn func(haystack, needle *[]byte) (int, error)

	switch {
	case len(*needle) < 2:
		return mmCount(haystack, needle), nil
	case len(*needle) < 50:
		//case len(*needle) < 3000: // use for small alphabets
		matchFn = bsf.Count
	case len(*needle) < 12000:
		// case len(*needle) < 16000: // use for small alphabets
		matchFn = bh.Count
	case len(*needle) < 350000:
		matchFn = bh2.Count
	case len(*needle) < 1<<22:
		//matchFn = bcj.Count // use for small alphabets
		matchFn = bsf.Count
	default:
		matchFn = bsf.Count
	}

	return matchFn(haystack, needle)
}
