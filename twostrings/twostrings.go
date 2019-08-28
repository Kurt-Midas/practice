package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello Twostrings!")
	fmt.Println(equalsWhenOneCharRemoved("x", "y"))        //false
	fmt.Println(equalsWhenOneCharRemoved("x", "XX"))       //false
	fmt.Println(equalsWhenOneCharRemoved("yy", "yx"))      //false
	fmt.Println(equalsWhenOneCharRemoved("abcd", "abxcd")) //true
	fmt.Println(equalsWhenOneCharRemoved("xyz", "xz"))     //true

	// Custom tests
	fmt.Println(equalsWhenOneCharRemoved("abcdefghijklmnopqrstuvwxyz", "bcdefghijklmnopqrstuvwxyz")) //true, missing head
	fmt.Println(equalsWhenOneCharRemoved("abcdefghijklmnopqrstuvwxyz", "abcdefghijklmnopqrstuvwxy")) //true, missing tail
	fmt.Println(equalsWhenOneCharRemoved("abcdefghijklmnopqrstuvwxyz", "bcdefghijklmnopqrstuvwxy"))  //false, length
	fmt.Println(equalsWhenOneCharRemoved("", ""))                                                    //false
	fmt.Println(equalsWhenOneCharRemoved("a", ""))                                                   //true
	fmt.Println(equalsWhenOneCharRemoved("", "a"))                                                   //true
}

// Iterate over both strings at the same time and compare characters at their respective counters.
// At the FIRST mismatch, increment the longer string's counter by 1.
// At the SECOND mismatch, return false.
// Otherwise return true.
func equalsWhenOneCharRemoved(s1, s2 string) bool {
	if len(s1)+1 == len(s2) {
		// s2 is exactly 1 longer, swap them
		s1, s2 = s2, s1
	} else if len(s1) != len(s2)+1 {
		// s1 is not exactly 1 longer. This case is impossible.
		return false
	}

	var c1, c2 int     // respective counters
	for c2 < len(s2) { // until the shorter finishes
		if s1[c1] != s2[c2] { // mismatch
			if c1 > c2 { // prior mismatch detected
				return false
			}
			c1++
			continue // check the match again
		}
		c1++
		c2++
	}
	// if the counters are the same then the mismatch is just the last character of s1
	return true
}
