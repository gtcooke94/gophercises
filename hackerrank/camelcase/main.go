package main

import "fmt"

func main() {
	var s string
	fmt.Scanf("%s", &s)
	var wordCount int32 = 0
	if len(s) > 0 {
		wordCount = 1
	}
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			wordCount++
		}
	}
	fmt.Println(wordCount)
}
