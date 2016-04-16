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
 * from insides into the 32bit implementation in c at
 * http://www.stdlib.net/~colmmacc/strlen.c.html
 * Copyright (C) 1991, 1993, 1997, 2000, 2003 Free Software Foundation, Inc.
 * part of the GNU C Library.
 *    Written by Torbjorn Granlund (tege@sics.se),
 *    with help from Dan Sahlin (dan@sics.se);
 *    commentary by Jim Blandy (jimb@ai.mit.edu).
 */

package bmatch

import (
	"unsafe"
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
		pat            = *needle
		lastCharIdx    = len(pat) - 1
		char           = pat[lastCharIdx]
		hayst          = *haystack
		n              = len(hayst)
		readIdx        = uintptr(unsafe.Pointer(&hayst[0]))
		maxIdx         = readIdx + uintptr(len(hayst))
		limIdx         = uintptr(maxIdx >> 3 << 3)
		hay, NOT_hay   uint64
		needleMask     uint64
		idx, uint64Idx int
		magicBITS      = uint64(0x7efefefefefefeff)
		NOT_magicBITS  = uint64(0xffffffffffffffff) ^ magicBITS // go has no bitwise '~' operator
	)

	if n < len(pat) {
		return -1
	}

	// prepare index mask &
	// prepare needleMask mask
	for i := uint8(0); i < 7; i++ {
		needleMask |= uint64(char)
		needleMask <<= 8
	}
	needleMask |= uint64(char)

	// run over haystack
	uint64Idx = 0
	for {
		hay = *(*uint64)(unsafe.Pointer(readIdx))
		hay ^= needleMask
		NOT_hay = uint64(0xffffffffffffffff) ^ hay // go has no bitwise '~' operator
		if (((hay + magicBITS) ^ NOT_hay) & NOT_magicBITS) != 0 {
			for idx = 0; idx < 8; idx++ {
				if hay<<uint((7-idx)<<3)>>56 == 0 {
					return uint64Idx + idx
				}
			}
		}
		readIdx += 8
		uint64Idx += 8
		if readIdx < limIdx {
			continue
		}
		break
	}
	if n%8 != 0 {
		uint64Idx -= 8
		for uint64Idx < len(hayst) {
			if hayst[uint64Idx] == char {
				return uint64Idx
			}
			uint64Idx++
		}
	}

	return -1

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
		pat            = *needle
		lastCharIdx    = len(pat) - 1
		char           = pat[lastCharIdx]
		hayst          = *haystack
		n              = len(hayst)
		readIdx        = uintptr(unsafe.Pointer(&hayst[0]))
		maxIdx         = readIdx + uintptr(len(hayst))
		limIdx         = uintptr(maxIdx >> 3 << 3)
		hay, NOT_hay   uint64
		needleMask     uint64
		idx, uint64Idx int
		magicBITS      = uint64(0x7efefefefefefeff)
		NOT_magicBITS  = uint64(0xffffffffffffffff) ^ magicBITS // go has no bitwise '~' operator
	)

	if n < len(pat) {
		return found
	}

	buflen := 10 + (n/(1+lastCharIdx))>>3
	found = make([]int, 0, buflen)

	// prepare index mask &
	// prepare needleMask mask
	for i := uint8(0); i < 7; i++ {
		needleMask |= uint64(char)
		needleMask <<= 8
	}
	needleMask |= uint64(char)

	// run over haystack
	uint64Idx = 0
	for {
		hay = *(*uint64)(unsafe.Pointer(readIdx))
		hay ^= needleMask
		NOT_hay = uint64(0xffffffffffffffff) ^ hay // go has no bitwise '~' operator
		if (((hay + magicBITS) ^ NOT_hay) & NOT_magicBITS) != 0 {
			for idx = 0; idx < 8; idx++ {
				if hay<<uint((7-idx)<<3)>>56 == 0 {
					found = append(found, uint64Idx+idx)
				}
			}
		}
		readIdx += 8
		uint64Idx += 8
		if readIdx < limIdx {
			continue
		}
		break
	}
	if n%8 != 0 {
		uint64Idx -= 8
		for uint64Idx < len(hayst) {
			if hayst[uint64Idx] == char {
				found = append(found, uint64Idx)
			}
			uint64Idx++
		}
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
		pat            = *needle
		lastCharIdx    = len(pat) - 1
		char           = pat[lastCharIdx]
		hayst          = *haystack
		n              = len(hayst)
		readIdx        = uintptr(unsafe.Pointer(&hayst[0]))
		maxIdx         = readIdx + uintptr(len(hayst))
		limIdx         = uintptr(maxIdx >> 3 << 3)
		hay, NOT_hay   uint64
		needleMask     uint64
		idx, uint64Idx int
		magicBITS      = uint64(0x7efefefefefefeff)
		NOT_magicBITS  = uint64(0xffffffffffffffff) ^ magicBITS // go has no bitwise '~' operator
	)

	if n < len(pat) {
		return count
	}

	// prepare index mask &
	// prepare needleMask mask
	for i := uint8(0); i < 7; i++ {
		needleMask |= uint64(char)
		needleMask <<= 8
	}
	needleMask |= uint64(char)

	// run over haystack
	uint64Idx = 0
	for {
		hay = *(*uint64)(unsafe.Pointer(readIdx))
		hay ^= needleMask
		NOT_hay = uint64(0xffffffffffffffff) ^ hay // go has no bitwise '~' operator
		if (((hay + magicBITS) ^ NOT_hay) & NOT_magicBITS) != 0 {
			for idx = 0; idx < 8; idx++ {
				if hay<<uint((7-idx)<<3)>>56 == 0 {
					count++
				}
			}
		}
		readIdx += 8
		uint64Idx += 8
		if readIdx < limIdx {
			continue
		}
		break
	}

	if n%8 != 0 {
		uint64Idx = n - n%8
		for uint64Idx < len(hayst) {
			if hayst[uint64Idx] == char {
				count++
			}
			uint64Idx++
		}
	}

	return count
}
