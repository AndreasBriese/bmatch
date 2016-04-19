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

package bmatch

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	DOWNLOAD_URL  = "https://archive.org/download/theciaworldfactb00571gut/571.zip"
	WORLDTEXTFILE = "./theciaworldfactb00571.zip"
	N             = 500 // No of sample patterns
	N_NEG         = 20
)

var (
	pat [][]byte
)

func makeRandomPatterns(lenPatterns int) {
	rand.Seed(time.Now().UnixNano())

	pat = make([][]byte, N+N_NEG) // patterns
	n := len(hay) - 1             // haystack length -1

	// patterns to search
	// N valid patterns
	for i := 0; i < N; i++ {
		// random pattern length >0 && <1024
		m := rand.Intn(lenPatterns) + 1
		si := rand.Intn(n - m)
		pat[i] = hay[si : si+m]
	}
	// 10 unvalid pattern - not in haystack
	for i := N; i < N+N_NEG; i++ {
		// random pattern length >0 && <1024
		m := rand.Intn(lenPatterns) + 1

		si := rand.Intn(n - m)
		pat[i] = make([]byte, m)
		copy(pat[i], hay[si:si+m])
		pat[i][rand.Intn(m)] = byte(252) // 252 is not present in worldText
	}
}

func downloadWFB() (txt []byte) {

	if _, err := os.Lstat(WORLDTEXTFILE); err != nil {
		f, err := os.Create(WORLDTEXTFILE)
		if err != nil {
			log.Fatalf("Create %v failed: %v \n", WORLDTEXTFILE, err)
		}
		defer f.Close()

		resp, err := http.Get(DOWNLOAD_URL)
		if err != nil {
			log.Fatalf("Download from %v failed: %v \n", DOWNLOAD_URL, err)
		}

		zp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Reading download failed: %v \n", err)
		}

		_, err = f.Write(zp)
		if err != nil {
			log.Fatalf("Saving to %v failed: %v \n", WORLDTEXTFILE, err)
		}
		f.Close()

	}

	zr, err := zip.OpenReader(WORLDTEXTFILE)
	if err != nil {
		log.Fatalf("Opening %v failed: %v \n", WORLDTEXTFILE, err)
	}
	defer zr.Close()

	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			log.Fatalf("Opening failed: %v \n", err)
		}
		defer rc.Close()

		txt, err = ioutil.ReadAll(rc)
		if err != nil {
			log.Fatalf("readall failed: %v \n", err)
		}
	}

	return txt
}

var hay []byte

func TestMain(m *testing.M) {

	hay = downloadWFB()

	fmt.Printf("\n###############\nbmatch.go\n")
	// fmt.Println(string(hay[:1000]))

	alpha := 0
	for i := 0; i < 256; i++ {
		needle := []byte{byte(i)}
		if c := mmIndex(&hay, &needle); c != -1 {
			alpha++
		}
	}

	fmt.Printf("Haystack: %v loaded (%v bytes)\nAlphabet size: %v\n\n", WORLDTEXTFILE, len(hay), alpha)

	m.Run()

}

func TestM_Count_BoyerMooreVSBytesIndexVSBmatch(t *testing.T) {

	makeRandomPatterns(1024)

	haystack := string(hay)
	sumBoyerMoore := 0
	sumBytesIndex := 0
	sumBmatch := 0
	for i := range pat {
		needle := string(pat[i])
		r1, _ := strBMCount(&haystack, &needle)
		sumBoyerMoore += r1
		r2, _ := bytesIndexCount(&hay, &(pat[i]))
		sumBytesIndex += r2
		r3, _ := Count(&hay, &(pat[i]))
		sumBmatch += r3
		if r1 != r2 || r2 != r3 {
			fmt.Println("pat", pat[i], len(pat[i]), "BM:", r1, "RK:", r2, "Bmatch:", r3)
			fmt.Printf(">%v<\n", string(pat[i]))
		}
	}

	if sumBmatch != sumBoyerMoore && sumBmatch != sumBytesIndex {
		t.Errorf("FAILED! Total number of found indices = %v; want %v", sumBmatch, sumBoyerMoore)
	}
}

func TestM_FindAll_BoyerMooreVSBytesIndexVSBmatch(t *testing.T) {

	makeRandomPatterns(1024)

	haystack := string(hay) // for "strings" method
	equal := true
	errCnt := 0
	var r1, r2, r3 []int

	for i := range pat {
		needle := string(pat[i]) // for "strings" method
		r1, _ = strBMFindAll(&haystack, &needle)
		r2, _ = bytesIndexFindAll(&hay, &(pat[i]))
		r3, _ = FindAll(&hay, &(pat[i]))
		if len(r1) != len(r2) || len(r2) != len(r3) {
			errCnt++
			continue
		}
		for i := range r1 {
			if r1[i] != r2[i] || r2[i] != r3[i] {
				equal = false
				errCnt++
			}
		}
	}

	for i := range pat {
		r, _ := Index(&hay, &(pat[i]))
		if len(r3) > 0 && r != r3[0] {
			equal = false
			errCnt++
		}
	}

	if !equal {
		t.Errorf("FAILED! At least %v different index\n", errCnt)
	}

}

func BenchmarkM_Bmatch_10_C(b *testing.B) {
	makeRandomPatterns(10)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			Count(&hay, &(pat[i]))
		}
	}
}

func BenchmarkM_BoyerMoore_10_C(b *testing.B) {
	makeRandomPatterns(10)
	haystack := string(hay)
	needles := make([]string, N+N_NEG)
	for i := 0; i < N+N_NEG; i++ {
		needles[i] = string(pat[i])
	}
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			strBMCount(&haystack, &(needles[i]))
		}
	}
}

func BenchmarkM_BytesIndex_10_C(b *testing.B) {
	makeRandomPatterns(10)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			bytesIndexCount(&hay, &(pat[i]))
		}
	}
}

func BenchmarkM_Bmatch_30_C(b *testing.B) {
	makeRandomPatterns(30)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			Count(&hay, &(pat[i]))
		}
	}
}

func BenchmarkM_BoyerMoore_30_C(b *testing.B) {
	makeRandomPatterns(30)
	haystack := string(hay)
	needles := make([]string, N+N_NEG)
	for i := 0; i < N+N_NEG; i++ {
		needles[i] = string(pat[i])
	}
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			strBMCount(&haystack, &(needles[i]))
		}
	}
}

func BenchmarkM_BytesIndex_30_C(b *testing.B) {
	makeRandomPatterns(30)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			bytesIndexCount(&hay, &(pat[i]))
		}
	}
}

func BenchmarkM_Bmatch_1024_C(b *testing.B) {
	makeRandomPatterns(1024)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			Count(&hay, &(pat[i]))
		}
	}
}

func BenchmarkM_BoyerMoore_1024_C(b *testing.B) {
	makeRandomPatterns(1024)
	haystack := string(hay)
	needles := make([]string, N+N_NEG)
	for i := 0; i < N+N_NEG; i++ {
		needles[i] = string(pat[i])
	}
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			strBMCount(&haystack, &(needles[i]))
		}
	}
}

func BenchmarkM_BytesIndex_1024_C(b *testing.B) {
	makeRandomPatterns(1024)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			bytesIndexCount(&hay, &(pat[i]))
		}
	}
}

func BenchmarkM_Bmatch_30_FI(b *testing.B) {
	makeRandomPatterns(30)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			Index(&hay, &(pat[i]))
		}
	}
}

func BenchmarkM_BoyerMoore_30_FI(b *testing.B) {
	makeRandomPatterns(30)
	haystack := string(hay)
	needles := make([]string, N+N_NEG)
	for i := 0; i < N+N_NEG; i++ {
		needles[i] = string(pat[i])
	}
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			strBMFI(&haystack, &(needles[i]))
		}
	}
}

func BenchmarkM_StringsIndex_30_FI(b *testing.B) {
	makeRandomPatterns(30)
	haystack := string(hay)
	needles := make([]string, N+N_NEG)
	for i := 0; i < N+N_NEG; i++ {
		needles[i] = string(pat[i])
	}
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			strings.Index(haystack, needles[i])
		}
	}
}

func BenchmarkM_BytesIndex_30_FI(b *testing.B) {
	makeRandomPatterns(30)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			bytesIndexFI(&hay, &(pat[i]))
		}
	}
}

func BenchmarkM_Bmatch_1024_FI(b *testing.B) {
	makeRandomPatterns(1024)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			Index(&hay, &(pat[i]))
		}
	}
}

func BenchmarkM_BoyerMoore_1024_FI(b *testing.B) {
	makeRandomPatterns(1024)
	haystack := string(hay)
	needles := make([]string, N+N_NEG)
	for i := 0; i < N+N_NEG; i++ {
		needles[i] = string(pat[i])
	}
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			strBMFI(&haystack, &(needles[i]))
		}
	}
}

func BenchmarkM_StringsIndex_1024_FI(b *testing.B) {
	makeRandomPatterns(1024)
	haystack := string(hay)
	needles := make([]string, N+N_NEG)
	for i := 0; i < N+N_NEG; i++ {
		needles[i] = string(pat[i])
	}
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			strings.Index(haystack, needles[i])
		}
	}
}

func BenchmarkM_BytesIndex_1024_FI(b *testing.B) {
	makeRandomPatterns(1024)
	b.ResetTimer()
	for r := 0; r < b.N; r++ {
		for i := 0; i < N+N_NEG; i++ {
			bytesIndexFI(&hay, &(pat[i]))
		}
	}
}

// stringFinder efficiently finds strings in a source text. It's implemented
// using the Boyer-Moore string search algorithm:
// http://en.wikipedia.org/wiki/Boyer-Moore_string_search_algorithm
// http://www.cs.utexas.edu/~moore/publications/fstrpos.pdf (note: this aged
// document uses 1-based indexing)
//
type stringFinder struct {
	// pattern is the string that we are searching for in the text.
	pattern string

	// badCharSkip[b] contains the distance between the last byte of pattern
	// and the rightmost occurrence of b in pattern. If b is not in pattern,
	// badCharSkip[b] is len(pattern).
	//
	// Whenever a mismatch is found with byte b in the text, we can safely
	// shift the matching frame at least badCharSkip[b] until the next time
	// the matching char could be in alignment.
	badCharSkip [256]int

	// goodSuffixSkip[i] defines how far we can shift the matching frame given
	// that the suffix pattern[i+1:] matches, but the byte pattern[i] does
	// not. There are two cases to consider:
	//
	// 1. The matched suffix occurs elsewhere in pattern (with a different
	// byte preceding it that we might possibly match). In this case, we can
	// shift the matching frame to align with the next suffix chunk. For
	// example, the pattern "mississi" has the suffix "issi" next occurring
	// (in right-to-left order) at index 1, so goodSuffixSkip[3] ==
	// shift+len(suffix) == 3+4 == 7.
	//
	// 2. If the matched suffix does not occur elsewhere in pattern, then the
	// matching frame may share part of its prefix with the end of the
	// matching suffix. In this case, goodSuffixSkip[i] will contain how far
	// to shift the frame to align this portion of the prefix to the
	// suffix. For example, in the pattern "abcxxxabc", when the first
	// mismatch from the back is found to be in position 3, the matching
	// suffix "xxabc" is not found elsewhere in the pattern. However, its
	// rightmost "abc" (at position 6) is a prefix of the whole pattern, so
	// goodSuffixSkip[3] == shift+len(suffix) == 6+5 == 11.
	goodSuffixSkip []int
}

func makeStringFinder(pattern string) *stringFinder {
	f := &stringFinder{
		pattern:        pattern,
		goodSuffixSkip: make([]int, len(pattern)),
	}
	// last is the index of the last character in the pattern.
	last := len(pattern) - 1

	// Build bad character table.
	// Bytes not in the pattern can skip one pattern's length.
	for i := range f.badCharSkip {
		f.badCharSkip[i] = len(pattern)
	}
	// The loop condition is < instead of <= so that the last byte does not
	// have a zero distance to itself. Finding this byte out of place implies
	// that it is not in the last position.
	for i := 0; i < last; i++ {
		f.badCharSkip[pattern[i]] = last - i
	}

	// Build good suffix table.
	// First pass: set each value to the next index which starts a prefix of
	// pattern.
	lastPrefix := last
	for i := last; i >= 0; i-- {
		if strings.HasPrefix(pattern, pattern[i+1:]) {
			lastPrefix = i + 1
		}
		// lastPrefix is the shift, and (last-i) is len(suffix).
		f.goodSuffixSkip[i] = lastPrefix + last - i
	}
	// Second pass: find repeats of pattern's suffix starting from the front.
	for i := 0; i < last; i++ {
		lenSuffix := longestCommonSuffix(pattern, pattern[1:i+1])
		if pattern[i-lenSuffix] != pattern[last-lenSuffix] {
			// (last-i) is the shift, and lenSuffix is len(suffix).
			f.goodSuffixSkip[last-lenSuffix] = lenSuffix + last - i
		}
	}

	return f
}

func longestCommonSuffix(a, b string) (i int) {
	for ; i < len(a) && i < len(b); i++ {
		if a[len(a)-1-i] != b[len(b)-1-i] {
			break
		}
	}
	return
}

// next returns the index in text of the first occurrence of the pattern. If
// the pattern is not found, it returns -1.
func (f *stringFinder) next(text string) int {
	i := len(f.pattern) - 1
	for i < len(text) {
		// Compare backwards from the end until the first unmatching character.
		j := len(f.pattern) - 1
		for j >= 0 && text[i] == f.pattern[j] {
			i--
			j--
		}
		if j < 0 {
			return i + 1 // match
		}
		i += max(f.badCharSkip[text[i]], f.goodSuffixSkip[j])
	}
	return -1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func strBMCount(haystack, needle *string) (found int, e error) {

	s := *haystack
	sep := *needle

	strFind := makeStringFinder(sep)
	idx := 0
	for len(s) > len(sep) {
		idx = strFind.next(s)
		if idx == -1 {
			break
		}
		// fmt.Printf("<%v> <%v>\n", sep, s[idx:idx+len(sep)])
		found++
		s = s[idx+1:]
	}
	return found, nil
}

func strBMFI(haystack, needle *string) (found int, e error) {

	s := *haystack
	sep := *needle
	strFind := makeStringFinder(sep)

	found = strFind.next(s)

	return found, nil
}

func strBMFindAll(haystack, needle *string) (found []int, e error) {

	s := *haystack
	sep := *needle

	strFind := makeStringFinder(sep)
	idx := 0
	lastIdx := 0

	for {
		idx = strFind.next(s)
		if idx == -1 {
			break
		}
		// fmt.Printf("<%v> <%v>\n", sep, s[idx:idx+len(sep)])
		found = append(found, lastIdx+idx)
		lastIdx += idx + 1
		s = s[idx+1:]
	}

	return found, nil
}

func bytesIndexCount(haystack, needle *[]byte) (found int, e error) {
	s := *haystack
	sep := *needle

	idx := 0
	for {
		idx = bytes.Index(s, sep)
		if idx == -1 {
			break
		}
		// fmt.Printf("<%v> <%v>\n", sep, s[idx:idx+len(sep)])
		found++
		s = s[idx+1:]
	}
	return found, nil
}

func bytesIndexFI(haystack, needle *[]byte) (found int, e error) {
	s := *haystack
	sep := *needle
	found = bytes.Index(s, sep)
	return found, nil
}

func bytesIndexFindAll(haystack, needle *[]byte) (found []int, e error) {
	s := *haystack
	sep := *needle

	idx := 0
	lastIdx := 0

	for {
		idx = bytes.Index(s, sep)
		if idx == -1 {
			break
		}
		// fmt.Printf("<%v> <%v>\n", sep, s[idx:idx+len(sep)])
		found = append(found, lastIdx+idx)
		lastIdx += idx + 1
		s = s[idx+1:]
	}

	return found, nil
}
