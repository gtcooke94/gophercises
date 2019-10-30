package main

import (
	"fmt"
)

// Complete the caesarCipher function below.
func caesarCipher(s string, k int32) string {
	returnStr := ""
	for _, char := range s {
		newChar := rotateChar(char, k)
		returnStr += string(newChar)
	}
	return returnStr
}

func rotateChar(c rune, k int32) rune {
	a := 'a'
	z := 'z'
	A := 'A'
	Z := 'Z'
	newC := c
	if c >= a && c <= z {
		newC = (c-a+k)%26 + a
	} else if c >= A && c <= Z {
		newC = (c-A+k)%26 + A
	}
	return newC
}

func main() {
	var s string
	var length, shift int32
	fmt.Scanf("%d", &length)
	fmt.Scanf("%s", &s)
	fmt.Scanf("%d", &shift)
	fmt.Println(s)
	fmt.Println(caesarCipher(s, shift))
}
