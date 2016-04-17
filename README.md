---
__bmatch.go__
=============
---

bmatch is faster fixed pattern search for Go.

Go provides different search mechanisms to find the indices of a fixed string or byte pattern, that are backed by different search algorithms:

* strings.Index switches by the pattern length 
	* 1: generic assembler function strings•IndexByte(..) (i.e. for OSX/darwin see /usr/local/go/src/runtime/asm_amd64.s)
	* 1-31: generic assembler function strings•indexShortString(..) (i.e. for OSX/darwin see /usr/local/go/src/runtime/asm_amd64.s)
	* >31 calling a Rabin-Karp search algorithm using a rolling addition hash and use string comparison to prove the result 
* strings.Replace for single string replacing invokes a Boyer-Moore search over a string in /usr/local/go/src/strings/search.go implementing a stringFinder type
* bytes.Index using generic assembler code (bytes•IndexByte(..)) to find the index of the first element of the pattern (i.e. for OSX/darwin see /usr/local/go/src/runtime/asm_amd64.s) and then compares the following sequence using another assembler function (Equal(..)).

The before mentioned assembler routines compare each byte of the haystack one by one. See unsafeMEMCHR.go for a faster approach. 
All but the Boyer-Moore search are relatively slow - even in comparison to python str.find function (implemented in C).
bmatch's underlying algorithms outperform all of Go's search functions.


__Usage__

Install bmatch by the usual

    go get github.com/AndreasBriese/bmatch

In your go code use `import "github.com/AndreasBriese/bmatch"` and apply it on **[]byte** types of the haystack (byte sequence to search in) and needle (pattern to search for).

`index, err := bmatch.Index(&haystack, &needle)` gives the first (left) index or -1 if not present,

`indices, err := bmatch.FindAll(&haystack, &needle)` to get an []int array with indices of (overlapping!) occurrences, or

`count, err := bmatch.Count(&haystack, &needle)` to get the number of (overlapping!) occurences of needle in haystack.

__Benchmarks__ (`go test -bench . cpu=1`)

	 ###############
	 bmatch.go
	 Haystack: ./theciaworldfactb00571.zip loaded (3013205 bytes)
	 Alphabet size: 93
	 
	 PASS
	 BenchmarkM_Bmatch_10_C          	       1	1002289595 ns/op 
	 BenchmarkM_BoyerMoore_10_C       	       1	2683892806 ns/op
	 BenchmarkM_BytesIndex_10_C       	       1	2630128009 ns/op
	 BenchmarkM_Bmatch_30_C           	       2	 528011588 ns/op
	 BenchmarkM_BoyerMoore_30_C       	       1	1402345246 ns/op
	 BenchmarkM_BytesIndex_30_C       	       1	2654263863 ns/op
	 BenchmarkM_Bmatch_1024_C         	      20	 105107250 ns/op
	 BenchmarkM_BoyerMoore_1024_C     	       5	 214846485 ns/op
	 BenchmarkM_BytesIndex_1024_C     	       1	2592599288 ns/op
	 BenchmarkM_Bmatch_30_FI        	      10	 104080809 ns/op
	 BenchmarkM_BoyerMoore_30_FI    	       5	 245957698 ns/op
	 BenchmarkM_StringsIndex_30_FI  	       3	 392142229 ns/op
	 BenchmarkM_BytesIndex_30_FI    	       2	 776670641 ns/op
	 BenchmarkM_Bmatch_1024_FI      	      30	  43820654 ns/op
	 BenchmarkM_BoyerMoore_1024_FI  	      10	 103821315 ns/op
	 BenchmarkM_StringsIndex_1024_FI	       1	1306085958 ns/op
	 BenchmarkM_BytesIndex_1024_FI  	       1	1276277951 ns/op
 
 on a MacBookPro 2013 with i7 and 8GB Ram searching for 500 random patterns plus 20 patterns that are not present in the "1995 CIA World Factbook" (~3MB english natural text). Benchmark naming: .._searchFunction_patternMaximumLength_C=count|FI=first left Index
 
 `go test` will try to download the test corpus from http://archive.org if it is not present in the folder. 
 
 __License__   
 bmatch.go (C)opyright 2016 Andreas Briese, eduToolbox@Bri-C GmbH, Sarstedt, with MIT license - see the headers in the code in the subfolders of the various search algorithms for details and reference to their predecessors (C-code mostly taken from the SMART tool http://www.dmi.unict.it/~faro/smart/ v.13.02). 
 
